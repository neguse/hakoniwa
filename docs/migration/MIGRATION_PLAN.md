# 箱庭諸島 Perl → Go 移植計画

## 概要

本ドキュメントは、箱庭諸島（Hakoniwa）のPerlコードベース（7,400行）をGoに移植するための詳細な計画書です。

### 移植方針

1. **段階的アプローチ**: 逐語的移植 → Go慣用的リファクタリング → 最適化
2. **互換性維持**: データフォーマット・ゲームロジックの完全な互換性を保証
3. **テスト駆動**: 各段階で既存のPerl版との出力比較テストを実施

---

## フェーズ概要

```
Phase 0: 準備・設計（1-2日）
   ↓
Phase 1: 逐語的移植（2-3週間）
   ↓
Phase 2: Go的リファクタリング（1-2週間）
   ↓
Phase 3: 最適化・クリーンアップ（1週間）
```

---

## Phase 0: 準備・設計（1-2日）

### 目標

- 移植の基盤となるドキュメントとリポジトリ構造を整備
- Perl → Go の対応関係を明確化
- テスト戦略を確立

### 成果物

1. **ドキュメント**
   - [x] `MIGRATION_PLAN.md` - 本ドキュメント
   - [ ] `PERL_TO_GO_MAPPING.md` - Perl/Go対応表
   - [ ] `DATA_FORMAT.md` - データフォーマット仕様
   - [ ] `API_REFERENCE.md` - 内部API設計
   - [ ] `TESTING_STRATEGY.md` - テスト戦略

2. **リポジトリ構造**
   - [ ] `go/` ディレクトリ配下の標準レイアウト作成
   - [ ] `perl/` ディレクトリへの既存コード移動
   - [ ] テスト用ディレクトリ構造の準備

3. **開発環境**
   - [ ] `go.mod` 初期化
   - [ ] `Makefile` 作成
   - [ ] CI/CD設定（後日）

---

## Phase 1: 逐語的移植（2-3週間）

### 目標

- Perlコードを機械的にGoへ変換
- 構造・ロジックを可能な限り元のまま維持
- 各モジュールごとに互換性テストを実施

### 基本方針

#### DO ✅

- Perlの構造をそのまま維持
- 関数名を機械的に変換（`snake_case` → `camelCase`）
- グローバル変数をそのまま使用（`variable/variable.go`）
- Perlのロジックを1:1で翻訳
- コメントで元のPerl行番号を記録

#### DON'T ❌

- Phase 1で設計を改善しようとしない
- グローバル変数を排除しない（Phase 2で実施）
- インターフェースを抽出しない（Phase 2で実施）
- エラーハンドリングを改善しない（Phase 2で実施）

### 実装順序

1. **依存関係の少ないモジュールから開始**

```
1. internal/hako/const     (依存: なし)
2. internal/hako/variable  (依存: const)
3. internal/hako/core      (依存: const, variable)
   - core.go               (エントリーポイント)
   - fileio.go             (ファイルI/O)
   - lock.go               (ロック機構)
   - utils.go              (ユーティリティ)
   - template.go           (HTMLテンプレート)
4. internal/hako/turn      (依存: const, variable, core)
   - turn.go               (ターン処理)
   - disaster.go           (災害処理)
   - monster.go            (モンスター処理)
5. internal/hako/mapview   (依存: const, variable, core)
6. internal/hako/top       (依存: const, variable, core)
7. internal/hako/maintenance (依存: const, variable, core)
8. internal/web            (依存: 上記すべて)
   - handler.go            (HTTPハンドラ)
   - middleware.go         (ミドルウェア)
9. cmd/hako-main           (エントリーポイント)
10. cmd/hako-mente         (メンテナンスツール)
```

2. **各モジュールの実装手順**

```
For each module:
  1. Perlコードを読み、関数一覧を作成
  2. Go構造体・型定義を作成
  3. 各関数を逐語的に移植
  4. 単体テストを作成
  5. Perl版との出力比較テストを実行
  6. パスするまで修正
  7. 次のモジュールへ
```

### 命名規則

#### モジュール名

| Perl | Go |
|------|-----|
| `Hako::Main` | `internal/hako/core` |
| `Hako::Const` | `internal/hako/const` |
| `Hako::Variable` | `internal/hako/variable` |
| `Hako::Turn` | `internal/hako/turn` |
| `Hako::Map` | `internal/hako/mapview` |
| `Hako::Top` | `internal/hako/top` |
| `Hako::Maintenance` | `internal/hako/maintenance` |

#### 関数名

| Perl | Go | 可視性 |
|------|-----|--------|
| `sub readIslandsFile` | `func ReadIslandsFile` | 公開 |
| `sub htmlEscape` | `func htmlEscape` | 非公開 |
| `sub hakolock` | `func hakoLock` | 非公開 |
| `sub cgiInput` | `func cgiInput` | 非公開 |

**ルール:**
- 公開関数: 大文字開始（`ReadIslandsFile`）
- 非公開関数: 小文字開始（`htmlEscape`）

#### 変数名

| Perl | Go |
|------|-----|
| `$HunitTime` | `const.UnitTime` |
| `$mode` | `variable.Mode` |
| `@Hislands` | `variable.Islands` |
| `%HidToName` | `variable.IDToName` |

### 型変換パターン

```go
// Perlハッシュ → Goマップ
my %hash = ();          →  hash := make(map[string]string)

// Perl配列 → Goスライス
my @array = ();         →  array := []string{}

// Perlリファレンス → Goポインタ
my $ref = \%hash;       →  ref := &hash

// Perl島データ → Go構造体
my @island = (...);     →  type Island struct { ... }
```

### エラーハンドリング（Phase 1）

```go
// Phase 1: Perlのdieをそのままpanic
sub foo {
    die "error" if !$ok;
}
↓
func foo() {
    if !ok {
        panic("error")  // Phase 1ではこれでOK
    }
}
```

**注**: Phase 2で適切な `error` 型に改善

### テスト方針（Phase 1）

#### 1. 単体テスト

```go
// 各関数の基本動作を確認
func TestHtmlEscape(t *testing.T) {
    tests := []struct {
        input    string
        expected string
    }{
        {"<script>", "&lt;script&gt;"},
        {"&", "&amp;"},
    }
    for _, tt := range tests {
        got := htmlEscape(tt.input)
        assert.Equal(t, tt.expected, got)
    }
}
```

#### 2. Golden Testing（互換性テスト）

```go
// Perl版の出力を「正解」として保存し、Go版と比較
func TestOutputCompatibility(t *testing.T) {
    // Perl版で生成した期待値を読み込み
    expected := readGoldenFile("testdata/golden/top_page.html")

    // Go版で同じ入力を処理
    got := runGoHandler(map[string]string{"mode": "top"})

    // 差分チェック
    if diff := cmp.Diff(expected, got); diff != "" {
        t.Errorf("Output mismatch (-want +got):\n%s", diff)
    }
}
```

#### 3. データ互換性テスト

```go
// Perl版が生成したデータファイルをGo版で読み込み
func TestDataFileCompatibility(t *testing.T) {
    // Perl版でデータ生成（テスト前処理）
    runPerlScript("t/generate_test_data.pl")

    // Go版で読み込み
    islands, err := core.ReadIslandsFile("testdata/hakojima.dat")
    require.NoError(t, err)

    // 内容を検証
    assert.Equal(t, 5, len(islands))
    assert.Equal(t, "テスト島", islands[0].Name)
}
```

### マイルストーン

- **Week 1**: `const`, `variable`, `core` パッケージ完成
- **Week 2**: `turn`, `mapview`, `top`, `maintenance` パッケージ完成
- **Week 3**: Web層とエントリーポイント完成、総合テスト

---

## Phase 2: Go的リファクタリング（1-2週間）

### 目標

- Phase 1で作成した逐語的コードをGo慣用的な設計に改善
- 保守性・テスタビリティの向上
- Phase 1の互換性テストがすべてパスすることを保証

### リファクタリング方針

#### アーキテクチャ改善

1. **グローバル変数の排除**
   - `variable` パッケージのグローバル変数 → 構造体フィールドへ

   ```go
   // Before (Phase 1)
   var Mode string
   var IslandID string

   // After (Phase 2)
   type Request struct {
       Mode     string
       IslandID string
   }
   ```

2. **関数 → メソッド化**
   ```go
   // Before
   func RunMain() { ... }

   // After
   type Game struct {
       config  *Config
       storage Storage
       lockMgr LockManager
   }
   func (g *Game) HandleRequest(ctx context.Context, req *Request) (*Response, error)
   ```

3. **インターフェース抽出**
   ```go
   type Storage interface {
       ReadIslands() ([]Island, error)
       WriteIslands(islands []Island) error
   }

   type FileStorage struct { ... }      // 本番用
   type MemoryStorage struct { ... }    // テスト用
   ```

4. **エラーハンドリング改善**
   ```go
   // Before: panic
   if !ok {
       panic("error")
   }

   // After: error返却
   if !ok {
       return fmt.Errorf("operation failed: %w", err)
   }
   ```

5. **コンテキスト対応**
   ```go
   func (g *Game) HandleRequest(ctx context.Context, req *Request) (*Response, error) {
       if err := g.lockMgr.Lock(ctx); err != nil {
           return nil, err
       }
       defer g.lockMgr.Unlock()
       // ...
   }
   ```

#### コード品質

- `golangci-lint` 準拠
- `Effective Go` 準拠
- テストカバレッジ 80% 以上

### リファクタリング順序

1. インターフェース定義（Storage, LockManager, Logger）
2. Game構造体の作成とメソッド化
3. グローバル変数の削除
4. エラーハンドリング改善
5. コンテキスト対応
6. テストの書き直し

### マイルストーン

- **Week 1**: インターフェース抽出、構造体化
- **Week 2**: エラーハンドリング改善、テスト修正、品質チェック

---

## Phase 3: 最適化・クリーンアップ（1週間）

### 目標

- パフォーマンス最適化
- ドキュメント整備
- デプロイメント準備

### 最適化項目

1. **並行処理の改善**
   - `sync.RWMutex` の適切な使用
   - コンテキストキャンセル対応

2. **メモリ効率化**
   - 不要なアロケーション削減
   - 構造体のメモリレイアウト最適化

3. **I/O最適化**
   - バッファリング
   - ファイルロックの効率化

4. **ベンチマーク**
   ```go
   func BenchmarkTurnProcessing(b *testing.B) {
       // ターン処理のベンチマーク
   }
   ```

### ドキュメント

- GoDoc形式のコメント追加
- 使用方法のREADME更新
- デプロイメントガイド作成

---

## リスク管理

### 既知のリスク

| リスク | 影響 | 対策 |
|-------|------|------|
| データフォーマットの互換性問題 | 高 | Golden Testingで早期検出 |
| Perlの挙動を完全に再現できない | 中 | 差異を文書化、必要に応じてPerl側修正 |
| パフォーマンス低下 | 低 | ベンチマーク計測、プロファイリング |

---

## 成功基準

### Phase 1 完了条件

- [ ] すべてのPerl関数がGoに移植されている
- [ ] 単体テストがすべてパス
- [ ] Golden Testがすべてパス（Perl版と出力が一致）
- [ ] データ互換性テストがパス

### Phase 2 完了条件

- [ ] グローバル変数がすべて排除されている
- [ ] `golangci-lint` がエラーなし
- [ ] テストカバレッジ 80% 以上
- [ ] Phase 1の互換性テストがすべてパス（リファクタ後も）

### Phase 3 完了条件

- [ ] パフォーマンスがPerl版と同等以上
- [ ] GoDocがすべての公開APIに記述されている
- [ ] デプロイメントガイド完成

---

## スケジュール

```
Week 1:  Phase 0 (準備) + Phase 1 開始 (const, variable, core)
Week 2:  Phase 1 継続 (turn, mapview, top, maintenance)
Week 3:  Phase 1 完了 (web層, エントリーポイント) + テスト
Week 4:  Phase 2 (リファクタリング) 開始
Week 5:  Phase 2 完了 + Phase 3 (最適化)
Week 6:  バッファ・ドキュメント整備
```

**合計**: 約 6 週間

---

## 次のステップ

1. [ ] Phase 0の残りドキュメント作成
2. [ ] `go.mod` 初期化
3. [ ] `internal/hako/const` パッケージの実装開始

---

## 参考資料

- [Effective Go](https://go.dev/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Standard Go Project Layout](https://github.com/golang-standards/project-layout)
- [箱庭諸島オリジナルドキュメント](../original/hako-readme.txt)

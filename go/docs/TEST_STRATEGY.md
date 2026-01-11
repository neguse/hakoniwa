# テスト戦略書

## 概要

Phase 1完了後の品質保証として、3層のテストを実施します。

## テスト階層

### 1. ユニットテスト (Unit Tests)

**目的**: 個別関数の正確性を検証

**対象**:
- `core/utils.go`: adjustXY, getAround1, getAround2等の座標計算
- `core/fileio.go`: データのエンコード/デコード
- `turn/turn.go`: 個別コマンド処理関数
- `mapview/mapview.go`: HTML生成関数

**実装方針**:
- Table-driven tests（テーブル駆動テスト）
- 各関数ごとに `*_test.go` ファイル
- エッジケースの網羅（境界値、異常値）

**ファイル配置**:
```
internal/hako/core/utils_test.go
internal/hako/core/fileio_test.go
internal/hako/turn/turn_test.go
```

### 2. データ互換性テスト (Compatibility Tests)

**目的**: Perl版のデータファイル形式との完全互換性を保証

**検証項目**:
- 島データファイル（.cgi）の読み込み
- 島データファイル（.cgi）の書き込み
- ログファイル（.log）の読み書き
- Hex形式エンコーディングの正確性
- コマンド配列のCSV形式

**テストデータ**:
- `testdata/islands/`: サンプル島データ
- `testdata/expected/`: 期待される出力

**実装**:
```
internal/hako/core/fileio_compat_test.go
```

### 3. ゴールデンテスト (Golden Tests)

**目的**: Perl版と完全に同一の出力を生成することを保証

**方法**:
1. 同一の初期データを用意
2. Perl版でターン処理実行 → 出力保存（golden output）
3. Go版でターン処理実行 → 出力保存
4. 両者を比較（diff）

**検証内容**:
- ターン処理後の島データが一致
- ログ出力が一致
- HTML出力が一致（空白・改行除く）

**実装**:
```
internal/hako/turn/golden_test.go
testdata/golden/: golden outputファイル群
```

## テストデータ構造

```
go/
├── internal/hako/
│   ├── core/
│   │   ├── utils_test.go
│   │   ├── fileio_test.go
│   │   └── fileio_compat_test.go
│   ├── turn/
│   │   ├── turn_test.go
│   │   └── golden_test.go
│   └── mapview/
│       └── mapview_test.go
└── testdata/
    ├── islands/          # サンプル島データ
    │   ├── island001.cgi
    │   ├── island002.cgi
    │   └── ...
    ├── expected/         # 期待される出力
    │   ├── island001_after_turn.cgi
    │   └── ...
    └── golden/           # Perl版の出力（golden）
        ├── scenario1/
        │   ├── before/
        │   └── after/
        └── scenario2/
```

## 実装順序

1. **Phase 1: ユニットテスト基盤** ✅ 次のステップ
   - testdata/ディレクトリ作成
   - サンプル島データ生成ヘルパー
   - core/utils_test.go実装

2. **Phase 2: データ互換性テスト**
   - core/fileio_test.go実装
   - Hex形式エンコード/デコードテスト
   - 島データ読み書きテスト

3. **Phase 3: コマンド個別テスト**
   - turn/turn_test.go実装
   - 各コマンド関数の単体テスト

4. **Phase 4: ゴールデンテスト**
   - Perl版との出力比較環境構築
   - シナリオベーステスト

## テスト実行コマンド

```bash
# 全テスト実行
go test ./...

# カバレッジ付き実行
go test -cover ./...

# 詳細出力
go test -v ./...

# 特定パッケージのみ
go test ./internal/hako/core/

# ベンチマーク
go test -bench=. ./...
```

## カバレッジ目標

- **ユニットテスト**: 80%以上
- **データ互換性**: 100%（全フォーマット）
- **ゴールデンテスト**: 主要シナリオ5件以上

## CI/CD統合（将来）

- GitHub Actionsで自動テスト実行
- プルリクエスト時の自動検証
- カバレッジレポート生成

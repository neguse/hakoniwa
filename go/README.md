# 箱庭諸島 Go版

箱庭諸島（Hakoniwa）のGo実装です。Perl版からの移植プロジェクトです。

## 概要

本プロジェクトは、Perl実装の箱庭諸島をGoに移植したものです。

- **Phase 1**: 逐語的移植（Perlの構造をそのまま維持）
- **Phase 2**: Go慣用的なリファクタリング
- **Phase 3**: 最適化・クリーンアップ

詳細は [移植計画](../docs/migration/MIGRATION_PLAN.md) を参照してください。

## 必要要件

- **Go**: 1.21 以上
- **Perl**: 5.x（互換性テスト用）

## インストール

```bash
# リポジトリをクローン
git clone https://github.com/neguse/hakoniwa.git
cd hakoniwa/go

# 依存関係をインストール
go mod tidy
```

## ビルド

```bash
# すべてのバイナリをビルド
make build

# 個別にビルド
go build -o bin/hako-main ./cmd/hako-main
go build -o bin/hako-mente ./cmd/hako-mente
```

## 実行

```bash
# メインサーバーを起動
./bin/hako-main

# または開発モードで実行
make run-main
```

デフォルトで `http://localhost:8080` でアクセス可能になります。

## テスト

```bash
# すべてのテストを実行
make test

# 単体テストのみ
make test-unit

# 互換性テスト（Golden Testing）
make test-compat

# 統合テスト
make test-integration

# カバレッジレポート生成
make coverage
```

## 開発

### ディレクトリ構造

```
go/
├── cmd/                    # エントリーポイント
│   ├── hako-main/         # メインゲームサーバー
│   └── hako-mente/        # メンテナンスツール
├── internal/              # 内部パッケージ
│   ├── hako/             # コアロジック
│   │   ├── const/        # 定数
│   │   ├── variable/     # グローバル変数（Phase 1のみ）
│   │   ├── core/         # コア機能
│   │   ├── turn/         # ターン処理
│   │   ├── mapview/      # 島の表示
│   │   ├── top/          # トップページ
│   │   └── maintenance/  # メンテナンス
│   ├── web/              # Webレイヤー
│   └── testutil/         # テストユーティリティ
├── test/                  # テスト
│   ├── integration/      # 統合テスト
│   └── compatibility/    # 互換性テスト
└── web/                   # 静的ファイル
    └── images/           # ゲーム画像
```

### コーディング規約

- [Effective Go](https://go.dev/doc/effective_go) に準拠
- `golangci-lint` でチェック
- テストカバレッジ 80% 以上を目標

```bash
# コードフォーマット
make fmt

# Linter実行
make lint
```

### Perl版との互換性

データフォーマットはPerl版と完全互換です。

```bash
# Perl版が生成したデータをGo版で読み込み
./bin/hako-main --data-dir ../cgi/data

# Go版が生成したデータをPerl版で読み込み
cd ../perl
plackup app.psgi
```

## ドキュメント

- [移植計画](../docs/migration/MIGRATION_PLAN.md)
- [Perl→Goマッピング](../docs/migration/PERL_TO_GO_MAPPING.md)
- [データフォーマット仕様](../docs/migration/DATA_FORMAT.md)
- [API リファレンス](../docs/migration/API_REFERENCE.md)
- [テスト戦略](../docs/migration/TESTING_STRATEGY.md)

## Phase 1 の実装状況

- [ ] `internal/hako/const` - 定数定義
- [ ] `internal/hako/variable` - グローバル変数
- [ ] `internal/hako/core` - コアロジック
- [ ] `internal/hako/turn` - ターン処理
- [ ] `internal/hako/mapview` - 島の表示
- [ ] `internal/hako/top` - トップページ
- [ ] `internal/hako/maintenance` - メンテナンス
- [ ] `cmd/hako-main` - メインサーバー
- [ ] `cmd/hako-mente` - メンテナンスツール

## ライセンス

オリジナルの箱庭諸島 ver2.3 のライセンスに準拠します。

詳細は [README.md](../docs/original/README.md) を参照してください。

## 貢献

プルリクエストを歓迎します！

1. フォークする
2. フィーチャーブランチを作成 (`git checkout -b feature/amazing-feature`)
3. コミット (`git commit -m 'Add amazing feature'`)
4. プッシュ (`git push origin feature/amazing-feature`)
5. プルリクエストを作成

## 参考資料

- [箱庭諸島オリジナル](http://www.bekkoame.ne.jp/~tokuoka/hakoniwa.html)
- [Perl版実装](../perl/)

# 箱庭諸島

箱庭諸島 ver2.3 の現代的な実装プロジェクトです。

## 概要

本プロジェクトは以下の2つの実装を含みます:

1. **Perl版** (`perl/`) - Plack/PSGIを使用した現代的なPerl実装
2. **Go版** (`go/`) - Perl版からのポート（移植中）

## プロジェクト構成

```
hakoniwa/
├── perl/              # Perl実装（現行版）
│   ├── cgi/          # CGIスクリプト
│   ├── lib/Hako/     # コアモジュール
│   ├── t/            # テスト
│   └── app.psgi      # Plackアプリケーション
├── go/                # Go実装（移植版）
│   ├── cmd/          # エントリーポイント
│   ├── internal/     # 内部パッケージ
│   └── test/         # テスト
└── docs/              # ドキュメント
    ├── migration/    # 移植関連ドキュメント
    └── original/     # オリジナルドキュメント
```

## クイックスタート

### Perl版を実行

```bash
cd perl
plackup app.psgi
# http://localhost:5000 でアクセス
```

### Go版を実行（移植完了後）

```bash
cd go
make build
./bin/hako-main
# http://localhost:8080 でアクセス
```

## Go移植プロジェクト

### 移植計画

詳細は [移植計画ドキュメント](docs/migration/MIGRATION_PLAN.md) を参照してください。

**フェーズ:**
- **Phase 0**: 準備・設計 ✅
- **Phase 1**: 逐語的移植（進行中）
- **Phase 2**: Go慣用的リファクタリング（予定）
- **Phase 3**: 最適化・クリーンアップ（予定）

### ドキュメント

- [移植計画](docs/migration/MIGRATION_PLAN.md) - 全体計画と手順
- [Perl→Goマッピング](docs/migration/PERL_TO_GO_MAPPING.md) - コード変換対応表
- [データフォーマット仕様](docs/migration/DATA_FORMAT.md) - ファイル形式
- [APIリファレンス](docs/migration/API_REFERENCE.md) - 内部API設計
- [テスト戦略](docs/migration/TESTING_STRATEGY.md) - テスト方針

## ライセンス

オリジナルの箱庭諸島 ver2.3 に準じます。
素晴らしいゲームを制作された原作者の方々に敬意を表します。

>  箱庭諸島 ver2.3
>
>    字: 徳岡宏樹
>    絵: 小川克人
>  題字: 稲葉修吾
>  テストプレイ他協力: 井上友博、小澤武史、さかもと、ほえほえ、ありづか
>  箱庭諸島のページ: http://www.bekkoame.ne.jp/~tokuoka/hakoniwa.html

### スクリプトについて

以下、オリジナルのreadme.txtより。

> 箱庭諸島2のスクリプトを改変し、それを他人に譲渡、配布する場合には、
> 以下の制約を課します。
>
> ・無料配布であること。
> ・ゲーム画面のトップに表示される、スクリプトの配布元へのリンクを
>   消すのを禁止すること。また、それ以外の改造は許可すること。
> ・本条件と同等に、改造したものの配布を許可すること。
> ・配布するページにおいて、オリジナルスクリプトの配布元として当サイトへ
>   のリンクを置くこと。

…といっても、オリジナルの作者である[徳岡宏樹さんのWebサイト](http://t.pos.to/hako/)は
現時点でアクセスできない状態になってしまっているため、4つめの制約はあまり意味のないものになってしまっています。

### 画像ファイルについて

以下、[徳岡宏樹さんのWebサイトサイトのアーカイブ](https://web.archive.org/web/20070113153728/http://t.pos.to/hako/)より。

> [Q3] オリジナルに付属していた画像については、再配布や改変は可能ですか？
>
> [A3] 商用利用を除いて、許可するものとします。
> オリジナルスクリプトに付属していた文書では「箱庭諸島以外の用途に使用してはならない」と書いてありました。
> しかし、その後原作者より「商用でない限り、箱庭諸島以外でも配布・改変可」という許可を得ています。
> 従って、商用利用でなければ配布も改変も可能です。
> もちろん何らかの問題が発生したとしても画像の原作者は関知しません。
> 自己責任でお願いします。

## 参考にさせていただいたWebサイト

* 再配布
    * [箱庭諸島の保管庫](http://www.hakoniwa.net/hako/)
    * [箱庭なページ](http://hako.gob.jp/)
    * [Neo-INO](http://neo-sub.sakura.ne.jp/ino/hako/download.html)
* 解説
    * [箱庭解体新書](http://qqmh3psd.web.fc2.com/sadoga/)

## 貢献

プルリクエストを歓迎します！移植プロジェクトへの貢献については、[移植計画](docs/migration/MIGRATION_PLAN.md)を参照してください。

## 開発環境

### Perl版
- Perl 5.x
- Plack/PSGI
- 依存モジュール: `cpanfile` 参照

### Go版
- Go 1.21以上
- 依存モジュール: `go.mod` 参照

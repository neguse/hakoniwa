# テスト戦略

本ドキュメントは、箱庭諸島のPerl→Go移植におけるテスト戦略を定義します。

---

## 目次

1. [テスト方針](#テスト方針)
2. [テストレベル](#テストレベル)
3. [Golden Testing（互換性テスト）](#golden-testing互換性テスト)
4. [単体テスト](#単体テスト)
5. [統合テスト](#統合テスト)
6. [テストツール・ライブラリ](#テストツールライブラリ)
7. [テストデータ管理](#テストデータ管理)
8. [CI/CD統合](#cicd統合)

---

## テスト方針

### 目標

1. **完全な互換性保証**: Perl版とGo版で同一の動作を保証
2. **回帰防止**: リファクタリング後も互換性を維持
3. **高いカバレッジ**: 80%以上のテストカバレッジ
4. **自動化**: すべてのテストを自動実行可能に

### 基本原則

- **Perl版を「正解」とする**: Golden Testingでデータ・出力を比較
- **段階的テスト**: 各モジュール → 統合 → E2E の順で実施
- **テストファースト**: 実装前にテストケースを定義
- **継続的検証**: Phase 2リファクタリング後もPhase 1のテストが通ること

---

## テストレベル

```
┌─────────────────────────────────────┐
│ E2Eテスト（ブラウザテスト）          │  ← 最終的に追加
├─────────────────────────────────────┤
│ 統合テスト（HTTPハンドラ）           │
├─────────────────────────────────────┤
│ Golden Testing（互換性テスト）       │  ← **最重要**
├─────────────────────────────────────┤
│ 単体テスト（関数・メソッド）         │
└─────────────────────────────────────┘
```

### 優先順位

1. **Golden Testing** - Perl版との互換性を保証（Phase 1で必須）
2. **単体テスト** - 各関数の動作を確認
3. **統合テスト** - HTTPリクエスト全体の動作を確認
4. **E2Eテスト** - 実際のブラウザ操作を確認（オプション）

---

## Golden Testing（互換性テスト）

### 概要

Perl版の出力を「Golden（正解）」として保存し、Go版の出力と比較します。

### テストケース

#### 1. データ互換性テスト

Perl版が生成したデータファイルをGo版で読み込み、正しく解釈できることを確認。

```go
// test/compatibility/data_compat_test.go
package compatibility

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestReadPerlGeneratedData(t *testing.T) {
    // Perl版でテストデータ生成
    runPerlScript(t, "scripts/generate_test_data.pl")

    // Go版で読み込み
    islands, err := core.ReadIslandsFile("testdata/perl_generated/hakojima.dat")
    require.NoError(t, err)

    // 内容を検証
    assert.Equal(t, 5, len(islands))
    assert.Equal(t, "テスト島", islands[0].Name)
    assert.Equal(t, "0", islands[0].ID)
    assert.Equal(t, 100, islands[0].Money)
}

func TestWriteGoGeneratedData(t *testing.T) {
    // Go版でデータ生成
    islands := []Island{
        {Name: "Go島", ID: "0", Money: 100},
    }
    err := core.WriteIslandsFile("testdata/go_generated/hakojima.dat", islands)
    require.NoError(t, err)

    // Perl版で読み込み
    output := runPerlScript(t, "scripts/read_go_data.pl", "testdata/go_generated/hakojima.dat")

    // 読み込めたことを確認
    assert.Contains(t, output, "Go島")
}
```

#### 2. 出力互換性テスト（HTML出力）

同じ入力に対して、Perl版とGo版が同じHTML出力を生成することを確認。

```go
// test/compatibility/output_compat_test.go
func TestTopPageOutput(t *testing.T) {
    // Perl版の出力を取得（Golden）
    perlOutput := runPerlCGI(t, map[string]string{"mode": "top"})
    saveGoldenFile(t, "testdata/golden/top_page.html", perlOutput)

    // Go版の出力を取得
    goOutput := runGoHandler(t, map[string]string{"mode": "top"})

    // 比較（厳密なdiffまたは正規化後の比較）
    if diff := cmp.Diff(normalizeHTML(perlOutput), normalizeHTML(goOutput)); diff != "" {
        t.Errorf("Output mismatch (-perl +go):\n%s", diff)
    }
}
```

**正規化処理:**
```go
func normalizeHTML(html string) string {
    // タイムスタンプなど動的な部分を除去
    html = regexp.MustCompile(`\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}`).ReplaceAllString(html, "TIMESTAMP")
    // 空白を正規化
    html = strings.TrimSpace(html)
    return html
}
```

#### 3. ゲームロジック互換性テスト

ターン処理など、ゲームロジックの結果が一致することを確認。

```go
func TestTurnProcessing(t *testing.T) {
    // 初期状態を準備
    setupTestData(t, "testdata/initial_state")

    // Perl版でターン処理
    runPerlScript(t, "scripts/run_turn.pl")
    perlResult := readIslandsFile(t, "testdata/perl_result/hakojima.dat")

    // 初期状態を復元
    setupTestData(t, "testdata/initial_state")

    // Go版でターン処理
    turn.TurnMain()
    core.WriteIslandsFile("testdata/go_result/hakojima.dat")
    goResult := readIslandsFile(t, "testdata/go_result/hakojima.dat")

    // 結果を比較
    assert.Equal(t, perlResult, goResult)
}
```

### Golden ファイルの管理

```
test/compatibility/golden/
├── top_page.html              # トップページ
├── new_island.html            # 新島作成
├── print_island.html          # 島観光
├── owner_page.html            # 開発画面
├── command_result.html        # コマンド登録結果
└── turn_result/               # ターン処理結果
    ├── hakojima.dat
    ├── island.0
    └── hakojima.his
```

**更新方法:**
```bash
# Golden ファイルを更新（Perl版の出力で上書き）
go test ./test/compatibility -update-golden
```

---

## 単体テスト

### ユーティリティ関数のテスト

```go
// internal/hako/core/utils_test.go
package core

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestHtmlEscape(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
    }{
        {
            name:     "Escape <",
            input:    "<script>",
            expected: "&lt;script&gt;",
        },
        {
            name:     "Escape &",
            input:    "A & B",
            expected: "A &amp; B",
        },
        {
            name:     "Escape quote",
            input:    `Say "Hello"`,
            expected: `Say &quot;Hello&quot;`,
        },
        {
            name:     "No escape needed",
            input:    "Hello World",
            expected: "Hello World",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := htmlEscape(tt.input)
            assert.Equal(t, tt.expected, got)
        })
    }
}

func TestRandom(t *testing.T) {
    // 範囲チェック
    for i := 0; i < 100; i++ {
        r := random(10)
        assert.GreaterOrEqual(t, r, 0)
        assert.Less(t, r, 10)
    }
}

func TestExpToLevel(t *testing.T) {
    tests := []struct {
        exp      int
        expected int
    }{
        {0, 1},
        {100, 1},
        {200, 2},
        {500, 3},
    }

    for _, tt := range tests {
        got := expToLevel(tt.exp)
        assert.Equal(t, tt.expected, got)
    }
}
```

### ファイルI/Oのテスト

```go
// internal/hako/core/fileio_test.go
func TestReadWriteIslandsFile(t *testing.T) {
    // 一時ディレクトリ作成
    tmpDir := t.TempDir()
    const.DirName = tmpDir

    // テストデータ
    islands := []Island{
        {
            Name:  "テスト島",
            ID:    "0",
            Money: 100,
            Food:  50,
        },
    }

    // 書き込み
    err := WriteIslandsFile(islands)
    require.NoError(t, err)

    // 読み込み
    loaded, err := ReadIslandsFile("0")
    require.NoError(t, err)

    // 検証
    assert.Equal(t, 1, len(loaded))
    assert.Equal(t, "テスト島", loaded[0].Name)
    assert.Equal(t, 100, loaded[0].Money)
}
```

### テーブル駆動テスト

```go
func TestCheckPassword(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        stored   string
        expected bool
    }{
        {"Correct password", "test", encode("test"), true},
        {"Wrong password", "wrong", encode("test"), false},
        {"Empty password", "", encode("test"), false},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := checkPassword(tt.input, tt.stored)
            assert.Equal(t, tt.expected, got)
        })
    }
}
```

---

## 統合テスト

### HTTPハンドラテスト

```go
// test/integration/main_test.go
package integration

import (
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestTopPageHandler(t *testing.T) {
    // テストサーバー起動
    handler := http.HandlerFunc(mainHandler)
    server := httptest.NewServer(handler)
    defer server.Close()

    // リクエスト送信
    resp, err := http.Get(server.URL + "?mode=top")
    require.NoError(t, err)
    defer resp.Body.Close()

    // ステータスコード確認
    assert.Equal(t, http.StatusOK, resp.StatusCode)

    // Content-Type確認
    assert.Equal(t, "text/html; charset=UTF-8", resp.Header.Get("Content-Type"))

    // ボディ確認
    body, err := io.ReadAll(resp.Body)
    require.NoError(t, err)
    assert.Contains(t, string(body), "箱庭諸島")
}

func TestNewIslandFlow(t *testing.T) {
    // 新島作成フロー
    handler := http.HandlerFunc(mainHandler)
    server := httptest.NewServer(handler)
    defer server.Close()

    // POSTリクエスト
    values := url.Values{}
    values.Set("mode", "new")
    values.Set("ISLANDNAME", "新しい島")
    values.Set("PASSWORD", "pass123")

    resp, err := http.PostForm(server.URL, values)
    require.NoError(t, err)
    defer resp.Body.Close()

    // 成功確認
    assert.Equal(t, http.StatusOK, resp.StatusCode)

    // データファイル確認
    islands, err := core.ReadIslandsFile("0")
    require.NoError(t, err)
    assert.Equal(t, 1, len(islands))
    assert.Equal(t, "新しい島", islands[0].Name)
}
```

---

## テストツール・ライブラリ

### 必須ライブラリ

```go
// go.mod
require (
    github.com/stretchr/testify v1.8.4      // アサーション
    github.com/google/go-cmp v0.6.0         // 差分比較
    golang.org/x/text v0.14.0               // エンコーディング
)
```

### テストヘルパー

```go
// internal/testutil/testutil.go
package testutil

import (
    "os"
    "os/exec"
    "testing"
)

// RunPerlScript はPerlスクリプトを実行
func RunPerlScript(t *testing.T, script string, args ...string) string {
    cmd := exec.Command("perl", append([]string{script}, args...)...)
    output, err := cmd.CombinedOutput()
    if err != nil {
        t.Fatalf("Perl script failed: %v\nOutput: %s", err, output)
    }
    return string(output)
}

// SetupTestData はテストデータをセットアップ
func SetupTestData(t *testing.T, sourceDir string) {
    tmpDir := t.TempDir()
    // sourceDir の内容を tmpDir にコピー
    // ...
    const.DirName = tmpDir
}

// SaveGoldenFile はGoldenファイルを保存
func SaveGoldenFile(t *testing.T, path string, content string) {
    if os.Getenv("UPDATE_GOLDEN") == "1" {
        err := os.WriteFile(path, []byte(content), 0644)
        if err != nil {
            t.Fatalf("Failed to save golden file: %v", err)
        }
    }
}
```

---

## テストデータ管理

### ディレクトリ構造

```
test/
├── compatibility/
│   ├── golden/                 # Goldenファイル（Perl版の出力）
│   │   ├── top_page.html
│   │   └── ...
│   ├── testdata/
│   │   ├── perl_generated/     # Perl版が生成したデータ
│   │   ├── go_generated/       # Go版が生成したデータ
│   │   └── initial_state/      # テスト初期状態
│   └── compat_test.go
├── integration/
│   ├── testdata/
│   └── integration_test.go
└── scripts/                    # テスト用Perlスクリプト
    ├── generate_test_data.pl   # テストデータ生成
    ├── read_go_data.pl         # Go版データ読み込み
    └── run_turn.pl             # ターン処理実行
```

### テストデータ生成スクリプト

```perl
#!/usr/bin/env perl
# test/scripts/generate_test_data.pl

use strict;
use warnings;
use lib '../../perl/lib';
use Hako::Main;
use Hako::Const;

# テストデータ生成
$HdirName = 'test/compatibility/testdata/perl_generated';
mkdir $HdirName unless -d $HdirName;

# 5つの島を作成
for my $i (0..4) {
    my $island = [
        "テスト島$i",  # name
        "$i",          # id
        0,             # prize
        0,             # absent
        "説明$i",      # comment
        encode("pass$i"), # password
        100 + $i * 10, # money
        50,            # food
        # ...
    ];
    push(@Hislands, $island);
}

writeIslandsFile();
```

---

## CI/CD統合

### GitHub Actions 例

```yaml
# .github/workflows/test.yml
name: Test

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'

      - name: Set up Perl
        run: |
          sudo apt-get update
          sudo apt-get install -y perl libplack-perl

      - name: Install Go dependencies
        run: go mod download

      - name: Run unit tests
        run: go test -v ./internal/...

      - name: Run compatibility tests
        run: go test -v ./test/compatibility/...

      - name: Run integration tests
        run: go test -v ./test/integration/...

      - name: Check coverage
        run: |
          go test -coverprofile=coverage.out ./...
          go tool cover -func=coverage.out
```

### Makefile

```makefile
# Makefile
.PHONY: test test-unit test-compat test-integration coverage

test: test-unit test-compat test-integration

test-unit:
	go test -v ./internal/...

test-compat:
	go test -v ./test/compatibility/...

test-integration:
	go test -v ./test/integration/...

coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

update-golden:
	UPDATE_GOLDEN=1 go test ./test/compatibility/...
```

---

## カバレッジ目標

### Phase 1

- **単体テスト**: 60%以上
- **Golden Testing**: 主要機能100%
- **統合テスト**: 主要フロー100%

### Phase 2

- **単体テスト**: 80%以上
- **Golden Testing**: 引き続き100%
- **統合テスト**: 全フロー100%

---

## チェックリスト

テスト実装時の確認項目:

- [ ] すべての公開関数に単体テストがある
- [ ] Golden Testingで主要機能をカバーしている
- [ ] Perl版が生成したデータをGo版で読み込めることを確認
- [ ] Go版が生成したデータをPerl版で読み込めることを確認
- [ ] テストデータ生成スクリプトが動作する
- [ ] CI/CDでテストが自動実行される
- [ ] カバレッジレポートが生成される
- [ ] Phase 2リファクタリング後も全テストがパスする

---

## 参考資料

- [Go Testing Guide](https://go.dev/doc/tutorial/add-a-test)
- [stretchr/testify](https://github.com/stretchr/testify)
- [google/go-cmp](https://github.com/google/go-cmp)

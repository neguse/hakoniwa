# Perl → Go マッピングガイド

本ドキュメントは、箱庭諸島のPerlコードをGoに変換する際の具体的な対応関係を示します。

---

## 目次

1. [モジュール構造のマッピング](#モジュール構造のマッピング)
2. [命名規則](#命名規則)
3. [データ型の変換](#データ型の変換)
4. [構文パターンの変換](#構文パターンの変換)
5. [関数の具体例](#関数の具体例)
6. [特殊なケース](#特殊なケース)

---

## モジュール構造のマッピング

### ディレクトリ・パッケージ対応表

| Perlモジュール | Goパッケージパス | 説明 |
|--------------|----------------|------|
| `Hako::Main` | `internal/hako/core` | コアロジック・ユーティリティ |
| `Hako::Const` | `internal/hako/const` | 定数・設定値 |
| `Hako::Variable` | `internal/hako/variable` | グローバル変数（Phase 1のみ） |
| `Hako::Turn` | `internal/hako/turn` | ターン処理・ゲームロジック |
| `Hako::Map` | `internal/hako/mapview` | 島の表示・インタラクション |
| `Hako::Top` | `internal/hako/top` | トップページ |
| `Hako::Maintenance` | `internal/hako/maintenance` | メンテナンスツール |
| `cgi/hako-main.cgi` | `cmd/hako-main/main.go` | メインエントリーポイント |
| `cgi/hako-mente.cgi` | `cmd/hako-mente/main.go` | メンテナンスエントリーポイント |

### ファイル分割方針

Perlの大きなモジュール（`Main.pm`: 1,175行）は、Goでは責務ごとに複数ファイルに分割します。

**例: `Hako::Main` → `internal/hako/core/`**

```
internal/hako/core/
├── core.go        # エントリーポイント (RunMain)
├── fileio.go      # ファイルI/O (ReadIslandsFile, WriteIslandsFile)
├── lock.go        # ロック機構 (hakoLock, unlock)
├── utils.go       # ユーティリティ (htmlEscape, random, encode)
└── template.go    # HTMLテンプレート (tempHeader, tempFooter)
```

---

## 命名規則

### パッケージ名

- Perlのモジュール名 `Hako::XXX` → Goのパッケージ名 `xxx`（小文字）
- `Hako::Main` → `core`（より明確な名前に変更）
- `Hako::Map` → `mapview`（Goの予約語との衝突を回避）

### 関数名

#### 基本ルール

```perl
# Perl: スネークケース
sub read_islands_file { ... }
sub html_escape { ... }

# Go: キャメルケース + 可視性
func ReadIslandsFile() { ... }  # 公開（先頭大文字）
func htmlEscape() { ... }        # 非公開（先頭小文字）
```

#### 可視性の判断基準

| Perlでの状況 | Go |
|------------|-----|
| `EXPORT` リストに含まれる | 公開（大文字開始） |
| モジュール内部でのみ使用 | 非公開（小文字開始） |

#### 具体例

| Perl関数 | Go関数 | 可視性 |
|---------|--------|--------|
| `sub run_main` | `func RunMain()` | 公開 |
| `sub readIslandsFile` | `func ReadIslandsFile()` | 公開 |
| `sub writeIsland` | `func WriteIsland()` | 公開 |
| `sub hakolock` | `func hakoLock()` | 非公開 |
| `sub unlock` | `func unlock()` | 非公開 |
| `sub htmlEscape` | `func htmlEscape()` | 非公開 |
| `sub cgiInput` | `func cgiInput()` | 非公開 |
| `sub random` | `func random()` | 非公開 |
| `sub tempHeader` | `func tempHeader()` | 非公開 |

### 定数名

```perl
# Perl: $H プレフィックス
$HunitTime = 21600;
$HmaxIsland = 30;
$HinitialMoney = 10000;

# Go: パッケージ名で名前空間を区別
package const

const (
    UnitTime     = 21600
    MaxIsland    = 30
    InitialMoney = 10000
)
```

### 変数名

```perl
# Perl: スカラー($), 配列(@), ハッシュ(%)
$mode           # スカラー
@Hislands       # 配列
%HidToName      # ハッシュ

# Go: 型で表現
variable.Mode          // string
variable.Islands       // []Island
variable.IDToName      // map[string]string
```

---

## データ型の変換

### 基本型

| Perl | Go |
|------|-----|
| `my $str = "";` | `var str string` |
| `my $num = 0;` | `var num int` |
| `my $flag = 0;` | `var flag bool` (0/1 → false/true) |

### 配列・スライス

```perl
# Perl
my @array = ();
push(@array, $value);
my $len = scalar @array;
my $first = $array[0];

# Go
var array []string
array = append(array, value)
len := len(array)
first := array[0]
```

### ハッシュ・マップ

```perl
# Perl
my %hash = ();
$hash{$key} = $value;
my $value = $hash{$key};
if (exists $hash{$key}) { ... }

# Go
hash := make(map[string]string)
hash[key] = value
value := hash[key]
if _, ok := hash[key]; ok { ... }
```

### リファレンス・ポインタ

```perl
# Perl
my $ref = \%hash;
my $value = $ref->{$key};

# Go
ref := &hash
value := (*ref)[key]
// または構造体の場合
ref := &myStruct
value := ref.Field  // 自動デリファレンス
```

### 配列リファレンスの配列（島データ）

```perl
# Perl: 島データは配列のリファレンスの配列
my @islands = ();
my $island = [
    $name,      # [0] 島の名前
    $id,        # [1] 島のID
    $prize,     # [2] 受賞
    # ...
];
push(@islands, $island);
my $name = $islands[0]->[0];

# Go: 構造体のスライス
type Island struct {
    Name  string
    ID    string
    Prize int
    // ...
}
var islands []Island
island := Island{
    Name:  name,
    ID:    id,
    Prize: prize,
}
islands = append(islands, island)
name := islands[0].Name
```

---

## 構文パターンの変換

### 条件分岐

```perl
# Perl
if ($condition) {
    # ...
} elsif ($other) {
    # ...
} else {
    # ...
}

# 後置if
die "error" if !$ok;

# Go
if condition {
    // ...
} else if other {
    // ...
} else {
    // ...
}

// Phase 1: panic（後置if相当）
if !ok {
    panic("error")
}

// Phase 2: エラー返却
if !ok {
    return fmt.Errorf("error")
}
```

### ループ

```perl
# Perl: for配列
for my $item (@array) {
    # ...
}

# Perl: C-style for
for (my $i = 0; $i < $max; $i++) {
    # ...
}

# Perl: while
while (<FILEHANDLE>) {
    # ...
}

# Go: range
for _, item := range array {
    // ...
}

# Go: C-style for
for i := 0; i < max; i++ {
    // ...
}

# Go: bufio.Scanner
scanner := bufio.NewScanner(file)
for scanner.Scan() {
    line := scanner.Text()
    // ...
}
```

### ファイルI/O

```perl
# Perl: 読み込み
open(my $in, '<', $file) or die "Cannot open $file: $!";
while (my $line = <$in>) {
    chomp($line);
    # ...
}
close($in);

# Perl: 書き込み
open(my $out, '>', $file) or die "Cannot open $file: $!";
print $out $data;
close($out);

# Go: 読み込み
file, err := os.Open(filename)
if err != nil {
    return fmt.Errorf("cannot open file: %w", err)
}
defer file.Close()

scanner := bufio.NewScanner(file)
for scanner.Scan() {
    line := scanner.Text()
    // ...
}

# Go: 書き込み
file, err := os.Create(filename)
if err != nil {
    return fmt.Errorf("cannot create file: %w", err)
}
defer file.Close()

writer := bufio.NewWriter(file)
fmt.Fprintln(writer, data)
writer.Flush()
```

### 文字列操作

```perl
# Perl: 連結
my $str = $a . $b . $c;

# Perl: 分割
my @parts = split(/,/, $str);

# Perl: 結合
my $joined = join(',', @parts);

# Perl: 置換
$str =~ s/old/new/g;

# Perl: マッチ
if ($str =~ /pattern/) { ... }

# Go: 連結
str := a + b + c
// または
str := fmt.Sprintf("%s%s%s", a, b, c)

# Go: 分割
parts := strings.Split(str, ",")

# Go: 結合
joined := strings.Join(parts, ",")

# Go: 置換
str = strings.ReplaceAll(str, "old", "new")

# Go: マッチ
if matched, _ := regexp.MatchString("pattern", str); matched {
    // ...
}
```

### エラーハンドリング

```perl
# Perl: die
sub foo {
    die "error message" if !$ok;
    return $value;
}

# Perl: eval (try-catch相当)
eval {
    dangerous_operation();
};
if ($@) {
    print "Error: $@\n";
}

# Go (Phase 1): panic
func foo() string {
    if !ok {
        panic("error message")
    }
    return value
}

# Go (Phase 2): error返却
func foo() (string, error) {
    if !ok {
        return "", fmt.Errorf("error message")
    }
    return value, nil
}

// recover (panic-catch相当)
func wrapper() {
    defer func() {
        if r := recover(); r != nil {
            fmt.Printf("Recovered: %v\n", r)
        }
    }()
    foo()
}
```

---

## 関数の具体例

### 例1: `htmlEscape` (utils.go)

```perl
# lib/Hako/Main.pm (行番号: 約500)
sub htmlEscape {
    my ($str) = @_;
    $str =~ s/&/&amp;/g;
    $str =~ s/</&lt;/g;
    $str =~ s/>/&gt;/g;
    $str =~ s/"/&quot;/g;
    return $str;
}
```

```go
// internal/hako/core/utils.go
package core

import "strings"

// htmlEscape escapes HTML special characters
// Ref: lib/Hako/Main.pm:~500
func htmlEscape(str string) string {
    str = strings.ReplaceAll(str, "&", "&amp;")
    str = strings.ReplaceAll(str, "<", "&lt;")
    str = strings.ReplaceAll(str, ">", "&gt;")
    str = strings.ReplaceAll(str, "\"", "&quot;")
    return str
}
```

### 例2: `random` (utils.go)

```perl
# lib/Hako/Main.pm
sub random {
    my ($m) = @_;
    return int( rand($m) );
}
```

```go
// internal/hako/core/utils.go
package core

import "math/rand"

// random returns a random integer in [0, m)
// Ref: lib/Hako/Main.pm:~520
func random(m int) int {
    return rand.Intn(m)
}
```

### 例3: `ReadIslandsFile` (fileio.go)

```perl
# lib/Hako/Main.pm (簡略版)
sub readIslandsFile {
    my ($id) = @_;

    my $file = "$HdirName/hakojima.dat";
    open(my $in, '<', $file) or return 0;

    # ヘッダー読み込み
    my $header = <$in>;
    chomp($header);
    my ($turn, $islandNumber, $islandLastID, $defaultID) = split(/,/, $header);

    # 島データ読み込み
    while (my $line = <$in>) {
        chomp($line);
        my @data = split(/,/, $line);
        # ... 島データをパース
        push(@Hislands, \@island);
    }
    close($in);
    return 1;
}
```

```go
// internal/hako/core/fileio.go
package core

import (
    "bufio"
    "fmt"
    "os"
    "strings"

    "hakoniwa/internal/hako/const"
    "hakoniwa/internal/hako/variable"
)

// ReadIslandsFile reads the main island data file
// Returns true on success, false on failure
// Ref: lib/Hako/Main.pm:~200
func ReadIslandsFile(id string) bool {
    filename := fmt.Sprintf("%s/hakojima.dat", const.DirName)

    file, err := os.Open(filename)
    if err != nil {
        return false
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)

    // Read header
    if !scanner.Scan() {
        return false
    }
    header := scanner.Text()
    parts := strings.Split(header, ",")
    if len(parts) < 4 {
        return false
    }

    variable.Turn = parts[0]
    variable.IslandNumber = parts[1]
    variable.IslandLastID = parts[2]
    variable.DefaultID = parts[3]

    // Read islands
    variable.Islands = []Island{}
    for scanner.Scan() {
        line := scanner.Text()
        island := parseIslandLine(line)
        variable.Islands = append(variable.Islands, island)
    }

    return true
}
```

### 例4: `RunMain` エントリーポイント (core.go)

```perl
# lib/Hako/Main.pm:65-120
sub run_main {
    # ロックをかける
    if (!hakolock()) {
        tempHeader();
        tempLockFail();
        tempFooter();
        exit(0);
    }

    # 乱数の初期化
    srand(time ^ $$);

    # COOKIE読みこみ
    cookieInput();

    # CGI読みこみ
    cgiInput();

    # 島データの読みこみ
    if (readIslandsFile($HcurrentID || 0) == 0) {
        unlock();
        tempHeader();
        tempNoDataFile();
        tempFooter();
        exit(0);
    }

    # モード別処理
    if ($mode eq 'top') {
        topPageMain();
    } elsif ($mode eq 'print') {
        printIslandMain();
    }
    # ...

    unlock();
}
```

```go
// internal/hako/core/core.go
package core

import (
    "math/rand"
    "time"

    "hakoniwa/internal/hako/mapview"
    "hakoniwa/internal/hako/top"
    "hakoniwa/internal/hako/variable"
)

// RunMain is the main entry point for the game
// Ref: lib/Hako/Main.pm:65
func RunMain() {
    // Acquire lock
    if !hakoLock() {
        tempHeader()
        tempLockFail()
        tempFooter()
        return
    }
    defer unlock()

    // Initialize random seed
    rand.Seed(time.Now().UnixNano())

    // Read cookie
    cookieInput()

    // Parse CGI input
    cgiInput()

    // Read island data
    currentID := variable.CurrentID
    if currentID == "" {
        currentID = "0"
    }
    if !ReadIslandsFile(currentID) {
        tempHeader()
        tempNoDataFile()
        tempFooter()
        return
    }

    // Mode dispatch
    switch variable.Mode {
    case "top":
        top.TopPageMain()
    case "print":
        mapview.PrintIslandMain()
    // ... other modes
    default:
        top.TopPageMain()
    }
}
```

---

## 特殊なケース

### 1. グローバル変数（Phase 1）

```perl
# Perl: パッケージ変数
package Hako::Variable;
our $mode = '';
our @Hislands = ();

# 他のパッケージから使用
use Hako::Variable;
$mode = 'top';
```

```go
// Go (Phase 1): パッケージ変数
package variable

var Mode string
var Islands []Island

// 他のパッケージから使用
import "hakoniwa/internal/hako/variable"

variable.Mode = "top"
```

**注**: Phase 2ではこれを構造体フィールドに変換します。

### 2. 配列インデックスによるデータアクセス

```perl
# Perl: 島データは配列インデックスで管理
my $island = $Hislands[0];
my $name = $island->[0];
my $id = $island->[1];
my $prize = $island->[2];
```

```go
// Go: 構造体フィールドで管理
island := variable.Islands[0]
name := island.Name
id := island.ID
prize := island.Prize
```

### 3. ファイルロック

```perl
# Perl: 複数のロック方式に対応
use Fcntl ':flock';

sub hakolock {
    if ($lockMode == 0) {
        # flock方式
        open(my $lockFile, '>', "$HdirName/hakojima.lock");
        flock($lockFile, LOCK_EX);
    } elsif ($lockMode == 1) {
        # symlink方式
        # ...
    }
}
```

```go
// Go: github.com/gofrs/flock を使用
package core

import (
    "github.com/gofrs/flock"
    "hakoniwa/internal/hako/const"
)

var fileLock *flock.Flock

func hakoLock() bool {
    lockPath := fmt.Sprintf("%s/hakojima.lock", const.DirName)
    fileLock = flock.New(lockPath)

    locked, err := fileLock.TryLock()
    if err != nil || !locked {
        return false
    }
    return true
}

func unlock() {
    if fileLock != nil {
        fileLock.Unlock()
    }
}
```

### 4. CGIパラメータのパース

```perl
# Perl: 手動でパース
sub cgiInput {
    my $buffer;
    if ($ENV{'REQUEST_METHOD'} eq 'POST') {
        read(STDIN, $buffer, $ENV{'CONTENT_LENGTH'});
    } else {
        $buffer = $ENV{'QUERY_STRING'};
    }

    foreach my $pair (split(/&/, $buffer)) {
        my ($name, $value) = split(/=/, $pair);
        $value =~ tr/+/ /;
        $value =~ s/%([0-9a-fA-F]{2})/pack("C", hex($1))/eg;
        # $nameに応じて変数に格納
    }
}
```

```go
// Go: net/http の標準機能を使用
package core

import (
    "net/http"
    "hakoniwa/internal/hako/variable"
)

func cgiInput(r *http.Request) {
    if err := r.ParseForm(); err != nil {
        return
    }

    variable.Mode = r.FormValue("mode")
    variable.IslandID = r.FormValue("ISLANDID")
    variable.IslandName = r.FormValue("ISLANDNAME")
    // ...
}
```

### 5. HTMLテンプレート出力

```perl
# Perl: print文で直接出力
sub tempHeader {
    print "Content-type: text/html; charset=UTF-8\n\n";
    print "<html>\n";
    print "<head><title>箱庭諸島</title></head>\n";
    print "<body>\n";
}
```

```go
// Go (Phase 1): 文字列バッファに蓄積
package core

import (
    "bytes"
    "hakoniwa/internal/hako/variable"
)

func tempHeader() {
    variable.OutputBuffer.WriteString("<html>\n")
    variable.OutputBuffer.WriteString("<head><title>箱庭諸島</title></head>\n")
    variable.OutputBuffer.WriteString("<body>\n")
}

// HTTPハンドラで一括出力
func Handler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html; charset=UTF-8")
    RunMain(r)
    w.Write(variable.OutputBuffer.Bytes())
}
```

**Phase 2**: `html/template` を使用した改善

---

## チェックリスト

Phase 1の逐語的移植時に、以下を確認してください:

- [ ] Perlの関数名を適切にGoの命名規則に変換した
- [ ] 配列・ハッシュを適切なGo型（slice, map, struct）に変換した
- [ ] ファイルI/Oのエラーハンドリングを実装した（Phase 1はpanicでOK）
- [ ] グローバル変数をvariableパッケージに配置した
- [ ] 元のPerl行番号をコメントで記録した
- [ ] 単体テストを作成した

---

## 参考資料

- [Effective Go](https://go.dev/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [箱庭諸島オリジナルコード](../../perl/)

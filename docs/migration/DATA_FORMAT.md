# データフォーマット仕様

本ドキュメントは、箱庭諸島のデータファイル形式を詳細に記述します。Perl版とGo版で**完全な互換性**を維持する必要があります。

---

## 目次

1. [ファイル一覧](#ファイル一覧)
2. [hakojima.dat - メインデータファイル](#hakojimadat---メインデータファイル)
3. [island.{ID} - 個別島データ](#islandid---個別島データ)
4. [hakojima.his - 履歴ファイル](#hakojimahis---履歴ファイル)
5. [hakojima.log{N} - 詳細ログファイル](#hakojimalogn---詳細ログファイル)
6. [ロックファイル](#ロックファイル)
7. [バックアップファイル](#バックアップファイル)
8. [エンコーディング](#エンコーディング)
9. [Go実装の注意点](#go実装の注意点)

---

## ファイル一覧

すべてのデータファイルは `cgi/data/` ディレクトリに格納されます（設定により変更可能）。

| ファイル名 | 説明 | 形式 |
|----------|------|-----|
| `hakojima.dat` | メインデータ（ターン情報、全島の基本情報） | テキスト（カンマ区切り） |
| `island.{ID}` | 個別島データ（地形、コマンド、掲示板） | テキスト（Hex + カンマ区切り） |
| `hakojima.his` | ゲーム履歴ログ | テキスト（カンマ区切り） |
| `hakojima.log{0-9}` | 詳細ログファイル（10ファイルローテーション） | テキスト（カンマ区切り） |
| `hakojima.lock` | ロックファイル | 空ファイル（存在チェックのみ） |
| `data.bak{0-3}` | バックアップディレクトリ | ディレクトリ（ローテーション） |

---

## hakojima.dat - メインデータファイル

### フォーマット

```
<ヘッダー行>
<島データ行1>
<島データ行2>
...
<島データ行N>
```

### ヘッダー行

カンマ区切りで以下の情報を含みます:

```
ターン数,最終更新時間,島の総数,次に割り当てるID
```

| フィールド | 変数名 | 型 | 説明 |
|----------|--------|---|------|
| ターン数 | `$HislandTurn` | 整数 | 現在のターン番号（0から開始） |
| 最終更新時間 | `$HislandLastTime` | UNIX時刻 | 最後にターンが更新された時刻 |
| 島の総数 | `$HislandNumber` | 整数 | 現在存在する島の数 |
| 次に割り当てるID | `$HislandNextID` | 整数 | 新しい島に割り当てるID（連番） |

**例:**
```
123,1699876543,5,6
```
- ターン123
- 最終更新: 2023-11-13 12:34:03 (UNIX時刻)
- 島の総数: 5
- 次ID: 6

### 島データ行

各島ごとに1行。カンマ区切りで以下の情報を含みます:

```
名前,ID,受賞,連続資金繰り数,コメント,暗号化パスワード,資金,食料,人口,広さ,農場,工場,採掘場
```

| # | フィールド | 変数名 | 型 | 説明 |
|--|----------|--------|---|------|
| 0 | 名前 | `$island->[0]` | 文字列 | 島の名前（最大32文字） |
| 1 | ID | `$island->[1]` | 整数 | 島の一意なID |
| 2 | 受賞 | `$island->[2]` | 整数 | 受賞フラグ（ビットマスク） |
| 3 | 連続資金繰り数 | `$island->[3]` | 整数 | 連続で放置されているターン数 |
| 4 | コメント | `$island->[4]` | 文字列 | 島の説明（最大80文字） |
| 5 | 暗号化パスワード | `$island->[5]` | 文字列 | encode()関数で暗号化されたパスワード |
| 6 | 資金 | `$island->[6]` | 整数 | 現在の資金（100億円単位） |
| 7 | 食料 | `$island->[7]` | 整数 | 現在の食料（100トン単位） |
| 8 | 人口 | `$island->[8]` | 整数 | 現在の人口（100人単位） |
| 9 | 広さ | `$island->[9]` | 整数 | 島の面積（100平方km単位） |
| 10 | 農場 | `$island->[10]` | 整数 | 農場の数 |
| 11 | 工場 | `$island->[11]` | 整数 | 工場の数 |
| 12 | 採掘場 | `$island->[12]` | 整数 | 採掘場の数 |

**例:**
```
テスト島,0,0,0,平和な島です,xYz123,100,50,20,121,5,3,0
```

### Go実装例

```go
type Island struct {
    Name          string  // [0] 名前
    ID            string  // [1] ID
    Prize         int     // [2] 受賞
    Absent        int     // [3] 連続資金繰り数
    Comment       string  // [4] コメント
    Password      string  // [5] 暗号化パスワード
    Money         int     // [6] 資金
    Food          int     // [7] 食料
    Pop           int     // [8] 人口
    Area          int     // [9] 広さ
    Farm          int     // [10] 農場
    Factory       int     // [11] 工場
    Mountain      int     // [12] 採掘場

    // island.{ID}ファイルから読み込まれるデータ
    Land          [][]LandInfo  // 地形データ (11x11)
    LandValue     [][]int       // 地形の値 (11x11)
    Commands      []Command     // コマンドキュー
    Lbbs          []LbbsEntry   // ローカル掲示板
}
```

---

## island.{ID} - 個別島データ

各島ごとに `island.0`, `island.1`, ... のファイルが作成されます。

### フォーマット

```
<地形データ（1行、363文字）>
<コマンド行1>
<コマンド行2>
...
<コマンド行N> （最大 HcommandMax 行）
<ローカル掲示板行1>
<ローカル掲示板行2>
...
```

### 地形データ（1行目）

11x11 = 121ヘクスの地形情報を16進数で表現。**363文字**の固定長。

**フォーマット:** 各ヘクスは3文字
```
<land(1桁)><landValue(2桁)>
```

- **land** (1桁): 地形タイプ（0-9, A-F）
- **landValue** (2桁): 地形の値（00-FF、16進数）

**例:**
```
000000000101010202030405060708090A0B0C0D0E0F0...（363文字）
```

| 位置 | 文字 | 意味 |
|-----|------|------|
| 0-2 | `000` | (0, 0) = 地形0, 値0 |
| 3-5 | `000` | (0, 1) = 地形0, 値0 |
| ... | ... | ... |
| 360-362 | `0F0` | (10, 10) = 地形0, 値F(15) |

**読み取り順序:** 左上から右へ、行ごとに下へ（合計121ヘクス）

**地形タイプ一覧:**

| 値 | 定数名 | 地形 |
|----|--------|------|
| 0 | `$HlandSea` | 海 |
| 1 | `$HlandWaste` | 荒地 |
| 2 | `$HlandPlains` | 平地 |
| 3 | `$HlandTown` | 町 |
| 4 | `$HlandForest` | 森 |
| 5 | `$HlandFarm` | 農場 |
| 6 | `$HlandFactory` | 工場 |
| 7 | `$HlandBase` | 基地 |
| 8 | `$HlandSbase` | 海底基地 |
| 9 | `$HlandDefence` | 防衛施設 |
| A | `$HlandMountain` | 山 |
| B | `$HlandMonster` | モンスター |
| C | `$HlandSbase2` | 海底都市 |
| D | `$HlandOil` | 海底油田 |
| E | `$HlandMonument` | 記念碑 |

### コマンド行（2行目以降）

カンマ区切りで5フィールド:

```
kind,target,x,y,arg
```

| フィールド | 説明 | 型 |
|----------|------|---|
| kind | コマンドの種類 | 整数 |
| target | ターゲットID（援助・ミサイルなど） | 整数 |
| x | X座標（0-10） | 整数 |
| y | Y座標（0-10） | 整数 |
| arg | 引数（数量など） | 整数 |

**例:**
```
0,0,5,5,0
1,2,3,4,10
```

**最大行数:** `$HcommandMax`（デフォルト30）

### ローカル掲示板行

実装により異なる（詳細は `$HuseLbbs` による）

---

## hakojima.his - 履歴ファイル

ゲームの公開イベント履歴を保存します。

### フォーマット

```
<ターン数>,<ログメッセージ>
<ターン数>,<ログメッセージ>
...
```

**例:**
```
120,テスト島で大地震が発生！
121,モンスター「いのら」がサンプル島に接近！
122,テスト島がサンプル島に援助を実施
```

**最大行数:** `$HhistoryMax`（デフォルト100）

---

## hakojima.log{N} - 詳細ログファイル

10個のログファイル（`hakojima.log0` ~ `hakojima.log9`）にローテーションで記録されます。

### フォーマット

```
<機密フラグ>,<ターン数>,<id1>,<id2>,<メッセージ>
```

| フィールド | 説明 |
|----------|------|
| 機密フラグ | 0=公開, 1=当事者のみ, 2=機密 |
| ターン数 | イベント発生ターン |
| id1 | 当事者の島ID |
| id2 | 相手の島ID（該当しない場合は空） |
| メッセージ | ログメッセージ |

**例:**
```
0,120,0,,テスト島: 人口が増加
1,121,0,1,テスト島 → サンプル島: ミサイル発射
2,122,1,,サンプル島: 極秘作戦実行
```

---

## ロックファイル

### hakojima.lock

並行アクセスを防ぐためのロックファイル。

**ロック方式:**

| モード | 方式 | 説明 |
|-------|------|------|
| 0 | flock | ファイルロック（推奨） |
| 1 | symlink | シンボリックリンク |
| 2 | mkdir | ディレクトリ作成 |
| 3 | no lock | ロックなし（開発用） |

**Go実装:**
```go
import "github.com/gofrs/flock"

fileLock := flock.New("data/hakojima.lock")
locked, err := fileLock.TryLock()
if err != nil || !locked {
    return fmt.Errorf("failed to acquire lock")
}
defer fileLock.Unlock()
```

---

## バックアップファイル

### data.bak{0-3}

定期的に `data/` ディレクトリ全体をバックアップ。

**ローテーション:**
- `$HbackupTurn` ターンごとにバックアップ（デフォルト: 10ターン）
- `$HbackupTimes` 個のバックアップを保持（デフォルト: 4）

**ディレクトリ構造:**
```
data/           # 現在のデータ
data.bak0/      # 最新のバックアップ
data.bak1/      # 1つ前
data.bak2/      # 2つ前
data.bak3/      # 3つ前（最古）
```

---

## エンコーディング

### 文字コード

- **ファイル:** UTF-8
- **Perl:** `use utf8; use open ':encoding(utf8)';`
- **Go:** すべて UTF-8（標準）

### パスワード暗号化

**Perl実装 (`encode` 関数):**
```perl
sub encode {
    my ($password) = @_;
    if ($cryptOn) {
        return crypt($password, $password);
    } else {
        return $password;
    }
}
```

**Go実装:**
```go
import "golang.org/x/crypto/bcrypt"

func encode(password string) (string, error) {
    if const.CryptOn {
        // crypt() 互換の実装が必要
        // Perl の crypt() は DES ベースなので注意
        return cryptCompat(password, password)
    }
    return password, nil
}
```

**注意:** Perlの `crypt()` は環境依存。互換性のため、同じアルゴリズム（DES）を使用する必要があります。

---

## Go実装の注意点

### 1. ファイルI/Oの順序保証

Perlでは順次処理されるため、Go実装でも**同じ順序**でファイルを読み書きする必要があります。

### 2. 数値のフォーマット

**16進数地形データ:**
```go
// 読み込み
landType := int(hexStr[0])  // '0'-'9', 'A'-'F'
landValue, _ := strconv.ParseInt(hexStr[1:3], 16, 64)

// 書き込み
hexStr := fmt.Sprintf("%X%02X", landType, landValue)
```

### 3. カンマ区切りのパース

```go
parts := strings.Split(line, ",")
if len(parts) != expectedCount {
    return fmt.Errorf("invalid format")
}
```

**注意:** コメントや名前にカンマが含まれる可能性があるため、適切にエスケープ処理が必要。

### 4. ファイルパーミッション

```go
// ディレクトリ
os.MkdirAll(dirPath, 0755)

// ファイル
os.WriteFile(filePath, data, 0644)
```

### 5. アトミック書き込み

```go
// 一時ファイルに書き込み → リネーム
tmpFile := filepath + ".tmp"
if err := os.WriteFile(tmpFile, data, 0644); err != nil {
    return err
}
if err := os.Rename(tmpFile, filepath); err != nil {
    return err
}
```

### 6. テストデータ生成

互換性テスト用に、Perl版と同じデータを生成できるようにする:

```go
func GenerateTestIsland() Island {
    return Island{
        Name:     "テスト島",
        ID:       "0",
        Prize:    0,
        Absent:   0,
        Comment:  "平和な島です",
        Password: "xYz123",  // encode("test")
        Money:    100,
        Food:     50,
        Pop:      20,
        Area:     121,
        Farm:     5,
        Factory:  3,
        Mountain: 0,
    }
}
```

---

## チェックリスト

データフォーマット互換性を確保するために:

- [ ] ファイルのエンコーディングがUTF-8
- [ ] カンマ区切りのフィールド数が正しい
- [ ] 16進数地形データが363文字固定
- [ ] パスワード暗号化がPerl版と一致
- [ ] ファイルロックが正しく動作
- [ ] バックアップローテーションが正しい
- [ ] Perl版が生成したデータをGo版で読み込める
- [ ] Go版が生成したデータをPerl版で読み込める

---

## 参考資料

- [memo.txt](../../perl/memo.txt) - オリジナルのメモ
- [Hako::Main::readIslandsFile](../../perl/lib/Hako/Main.pm) - Perl読み込み実装
- [Hako::Main::writeIslandsFile](../../perl/lib/Hako/Main.pm) - Perl書き込み実装

# API リファレンス

本ドキュメントは、箱庭諸島Go版の内部APIを定義します。Phase 1（逐語的移植）とPhase 2（リファクタリング後）の両方の設計を記載します。

---

## 目次

1. [Phase 1 API（逐語的移植）](#phase-1-api逐語的移植)
2. [Phase 2 API（リファクタリング後）](#phase-2-apiリファクタリング後)
3. [パッケージ一覧](#パッケージ一覧)
4. [共通型定義](#共通型定義)

---

## Phase 1 API（逐語的移植）

Phase 1では、Perlの構造をそのまま維持し、グローバル変数を使用します。

### パッケージ構成

```
internal/hako/
├── const/         # 定数
├── variable/      # グローバル変数
├── core/          # コアロジック
├── turn/          # ターン処理
├── mapview/       # 島の表示
├── top/           # トップページ
└── maintenance/   # メンテナンス
```

---

### const パッケージ

**ファイル:** `internal/hako/const/const.go`

#### 設定定数

```go
package const

// ディレクトリ・パス設定
const (
    BaseDir         = ""
    ImageDir        = "/images"
    DirName         = "data"
    DirMode         = 0755
)

// 認証設定
const (
    MasterPassword  = "yourpassword"
    SpecialPassword = ""
    AdminName       = "管理者"
    Email           = "admin@example.com"
)

// ゲーム設定
const (
    UnitTime       = 21600  // 6時間（秒）
    MaxIsland      = 30     // 最大島数
    TopLogTurn     = 30     // トップページに表示するログのターン数
    LogMax         = 10     // ログファイルの最大数
    BackupTurn     = 10     // バックアップ間隔（ターン）
    BackupTimes    = 4      // バックアップ世代数
    HistoryMax     = 100    // 履歴の最大行数
    GiveupTurn     = 24     // 放棄判定ターン数
    CommandMax     = 30     // コマンドキューの最大数
    IslandSize     = 11     // 島のサイズ（11x11）
)

// 初期値
const (
    InitialMoney   = 100    // 初期資金（100億円単位）
    InitialFood    = 50     // 初期食料（100トン単位）
)

// 単位
const (
    UnitMoney      = 100    // 資金の単位（億円）
    UnitFood       = 100    // 食料の単位（トン）
    UnitPop        = 100    // 人口の単位（人）
    UnitArea       = 100    // 面積の単位（平方km）
)
```

#### 地形定数

```go
// 地形タイプ
const (
    LandSea       = 0   // 海
    LandWaste     = 1   // 荒地
    LandPlains    = 2   // 平地
    LandTown      = 3   // 町
    LandForest    = 4   // 森
    LandFarm      = 5   // 農場
    LandFactory   = 6   // 工場
    LandBase      = 7   // ミサイル基地
    LandSbase     = 8   // 海底基地
    LandDefence   = 9   // 防衛施設
    LandMountain  = 10  // 山
    LandMonster   = 11  // モンスター
    LandSbase2    = 12  // 海底都市
    LandOil       = 13  // 海底油田
    LandMonument  = 14  // 記念碑
)
```

#### 災害・モンスター定数

```go
// 災害確率
const (
    DisEarthquake  = 10   // 地震
    DisTsunami     = 10   // 津波
    DisTyphoon     = 10   // 台風
    DisMeteo       = 5    // 隕石
    DisHugeMeteo   = 10   // 巨大隕石
    DisEruption    = 5    // 噴火
    DisFire        = 20   // 火災
)

// モンスター種類
const (
    MonsterNumber = 9
)

// モンスター名
var MonsterName = []string{
    "いのら",
    "サンジラ",
    "レッドいのら",
    "ダークいのら",
    "いのらゴースト",
    "クジラ",
    "キングいのら",
    "メカいのら",
    "くじらず",
}

// モンスターHP
var MonsterBHP = []int{1, 1, 2, 3, 1, 10, 5, 3, 20}
var MonsterDHP = []int{1, 2, 2, 3, 1, 10, 5, 5, 30}

// モンスター特殊能力
var MonsterSpecial = []int{0, 0, 1, 2, 3, 4, 5, 6, 7}

// モンスター経験値
var MonsterExp = []int{100, 200, 300, 500, 10, 1000, 800, 500, 2000}
```

---

### variable パッケージ

**ファイル:** `internal/hako/variable/variable.go`

#### グローバル変数

```go
package variable

import "bytes"

// リクエストパラメータ
var (
    Mode           string   // 処理モード
    CurrentName    string   // 島の名前（入力）
    CurrentID      string   // 島のID（入力）
    DefaultID      string   // デフォルトID
    OldPassword    string   // 旧パスワード
    InputPassword  string   // 新パスワード
    InputPassword2 string   // 新パスワード（確認）
    Message        string   // メッセージ
)

// コマンド関連
var (
    CommandPlanNumber int     // コマンド番号
    CommandKind       int     // コマンド種類
    CommandTarget     string  // ターゲットID
    CommandX          int     // X座標
    CommandY          int     // Y座標
    CommandArg        int     // 引数
    CommandMode       string  // コマンドモード（write/insert/delete）
)

// ゲーム状態
var (
    Turn            string     // 現在のターン
    IslandNumber    string     // 島の総数
    IslandLastID    string     // 最後のID
    Islands         []Island   // 島データ
    IDToName        map[string]string  // ID → 名前
    IDToNumber      map[string]int     // ID → 番号
)

// ログプール
var (
    LogPool       []string   // 通常ログ
    LogPoolSecure []string   // 機密ログ
)

// 出力バッファ
var (
    OutputBuffer  bytes.Buffer  // HTML出力バッファ
)

// 乱数座標配列
var (
    Rpx [121]int  // ランダムX座標
    Rpy [121]int  // ランダムY座標
)
```

---

### core パッケージ

**ファイル:** `internal/hako/core/`

#### core.go - エントリーポイント

```go
package core

// RunMain はメインエントリーポイント
// Ref: lib/Hako/Main.pm:65
func RunMain(r *http.Request)
```

#### fileio.go - ファイルI/O

```go
// ReadIslandsFile はメインデータファイルを読み込む
// Ref: lib/Hako/Main.pm:~200
func ReadIslandsFile(id string) bool

// ReadIsland は個別島データを読み込む
// Ref: lib/Hako/Main.pm:~300
func ReadIsland(island *Island) bool

// WriteIslandsFile はメインデータファイルに書き込む
// Ref: lib/Hako/Main.pm:~400
func WriteIslandsFile() bool

// WriteIsland は個別島データを書き込む
// Ref: lib/Hako/Main.pm:~450
func WriteIsland(island *Island) bool
```

#### lock.go - ロック機構

```go
// hakoLock はファイルロックを取得
// Ref: lib/Hako/Main.pm:~600
func hakoLock() bool

// unlock はファイルロックを解放
// Ref: lib/Hako/Main.pm:~650
func unlock()
```

#### utils.go - ユーティリティ

```go
// htmlEscape はHTML特殊文字をエスケープ
// Ref: lib/Hako/Main.pm:~500
func htmlEscape(str string) string

// random は乱数を生成 [0, m)
// Ref: lib/Hako/Main.pm:~520
func random(m int) int

// encode はパスワードを暗号化
// Ref: lib/Hako/Main.pm:~540
func encode(password string) string

// checkPassword はパスワードを検証
// Ref: lib/Hako/Main.pm:~560
func checkPassword(input, stored string) bool

// expToLevel は経験値からレベルを計算
// Ref: lib/Hako/Main.pm:~700
func expToLevel(exp int) int

// monsterSpec はモンスターの仕様を取得
// Ref: lib/Hako/Main.pm:~720
func monsterSpec(lv int) (name string, hp int, exp int)
```

#### template.go - HTMLテンプレート

```go
// tempHeader はHTMLヘッダーを出力
// Ref: lib/Hako/Main.pm:~800
func tempHeader()

// tempFooter はHTMLフッターを出力
// Ref: lib/Hako/Main.pm:~850
func tempFooter()

// tempLockFail はロック失敗メッセージを出力
// Ref: lib/Hako/Main.pm:~900
func tempLockFail()

// tempNoDataFile はデータファイルなしメッセージを出力
// Ref: lib/Hako/Main.pm:~920
func tempNoDataFile()

// tempWrongPassword はパスワード誤りメッセージを出力
// Ref: lib/Hako/Main.pm:~940
func tempWrongPassword()
```

---

### turn パッケージ

**ファイル:** `internal/hako/turn/`

#### turn.go

```go
package turn

// NewIslandMain は新しい島を作成
// Ref: lib/Hako/Turn.pm:~50
func NewIslandMain()

// TurnMain はターン処理を実行
// Ref: lib/Hako/Turn.pm:~500
func TurnMain()

// ChangeMain は島の情報を変更
// Ref: lib/Hako/Turn.pm:~200
func ChangeMain()
```

---

### mapview パッケージ

**ファイル:** `internal/hako/mapview/mapview.go`

```go
package mapview

// PrintIslandMain は島を観光表示
// Ref: lib/Hako/Map.pm:~50
func PrintIslandMain()

// OwnerMain は島の開発画面を表示
// Ref: lib/Hako/Map.pm:~200
func OwnerMain()

// CommandMain はコマンドを登録
// Ref: lib/Hako/Map.pm:~400
func CommandMain()

// CommentMain はコメントを更新
// Ref: lib/Hako/Map.pm:~600
func CommentMain()
```

---

### top パッケージ

**ファイル:** `internal/hako/top/top.go`

```go
package top

// TopPageMain はトップページを表示
// Ref: lib/Hako/Top.pm:~50
func TopPageMain()
```

---

### maintenance パッケージ

**ファイル:** `internal/hako/maintenance/maintenance.go`

```go
package maintenance

// RunMaintenance はメンテナンスツールを実行
// Ref: lib/Hako/Maintenance.pm:~50
func RunMaintenance(r *http.Request)
```

---

## Phase 2 API（リファクタリング後）

Phase 2では、グローバル変数を排除し、構造体ベースの設計に移行します。

### 設計方針

- グローバル変数 → 構造体フィールド
- 関数 → メソッド
- インターフェースの導入
- エラーハンドリングの改善
- コンテキスト対応

---

### コアインターフェース

```go
// Storage はデータ永続化のインターフェース
type Storage interface {
    ReadIslands(ctx context.Context) ([]Island, error)
    WriteIslands(ctx context.Context, islands []Island) error
    ReadIsland(ctx context.Context, id string) (*Island, error)
    WriteIsland(ctx context.Context, island *Island) error
    ReadHistory(ctx context.Context) ([]HistoryEntry, error)
    WriteHistory(ctx context.Context, entries []HistoryEntry) error
}

// LockManager はロック管理のインターフェース
type LockManager interface {
    Lock(ctx context.Context) error
    Unlock(ctx context.Context) error
}

// Logger はログ出力のインターフェース
type Logger interface {
    Info(msg string, fields ...Field)
    Error(msg string, err error, fields ...Field)
}
```

---

### Game 構造体（Phase 2のコア）

```go
package game

import (
    "context"
    "hakoniwa/internal/hako/const"
)

// Game はゲームのメインロジックを管理
type Game struct {
    config  *Config
    storage Storage
    lockMgr LockManager
    logger  Logger
    rand    *rand.Rand
}

// NewGame は新しいGameインスタンスを作成
func NewGame(cfg *Config, storage Storage, lockMgr LockManager, logger Logger) *Game {
    return &Game{
        config:  cfg,
        storage: storage,
        lockMgr: lockMgr,
        logger:  logger,
        rand:    rand.New(rand.NewSource(time.Now().UnixNano())),
    }
}

// HandleRequest はHTTPリクエストを処理
func (g *Game) HandleRequest(ctx context.Context, req *Request) (*Response, error) {
    // ロック取得
    if err := g.lockMgr.Lock(ctx); err != nil {
        return nil, fmt.Errorf("failed to acquire lock: %w", err)
    }
    defer g.lockMgr.Unlock(ctx)

    // データ読み込み
    islands, err := g.storage.ReadIslands(ctx)
    if err != nil {
        return nil, fmt.Errorf("failed to read islands: %w", err)
    }

    // モード別処理
    var response *Response
    switch req.Mode {
    case "top":
        response, err = g.handleTopPage(ctx, islands)
    case "print":
        response, err = g.handlePrintIsland(ctx, req, islands)
    case "owner":
        response, err = g.handleOwner(ctx, req, islands)
    // ...
    default:
        response, err = g.handleTopPage(ctx, islands)
    }

    if err != nil {
        return nil, err
    }

    return response, nil
}
```

---

### Request/Response 構造体

```go
// Request はHTTPリクエストを表現
type Request struct {
    Mode           string
    IslandName     string
    IslandID       string
    Password       string
    OldPassword    string
    Message        string
    Command        *Command
}

// Response はHTTPレスポンスを表現
type Response struct {
    StatusCode  int
    ContentType string
    Body        []byte
    Cookies     []*http.Cookie
}

// Command はコマンドを表現
type Command struct {
    Number int
    Kind   int
    Target string
    X      int
    Y      int
    Arg    int
    Mode   string  // write/insert/delete
}
```

---

### Storage 実装

#### FileStorage（本番用）

```go
package storage

type FileStorage struct {
    dataDir string
    mu      sync.RWMutex
}

func NewFileStorage(dataDir string) *FileStorage {
    return &FileStorage{dataDir: dataDir}
}

func (fs *FileStorage) ReadIslands(ctx context.Context) ([]Island, error) {
    fs.mu.RLock()
    defer fs.mu.RUnlock()

    filepath := fmt.Sprintf("%s/hakojima.dat", fs.dataDir)
    // 実装...
}

func (fs *FileStorage) WriteIslands(ctx context.Context, islands []Island) error {
    fs.mu.Lock()
    defer fs.mu.Unlock()

    filepath := fmt.Sprintf("%s/hakojima.dat", fs.dataDir)
    // 実装...
}
```

#### MemoryStorage（テスト用）

```go
package storage

type MemoryStorage struct {
    islands map[string]*Island
    mu      sync.RWMutex
}

func NewMemoryStorage() *MemoryStorage {
    return &MemoryStorage{
        islands: make(map[string]*Island),
    }
}

func (ms *MemoryStorage) ReadIslands(ctx context.Context) ([]Island, error) {
    ms.mu.RLock()
    defer ms.mu.RUnlock()

    result := make([]Island, 0, len(ms.islands))
    for _, island := range ms.islands {
        result = append(result, *island)
    }
    return result, nil
}
```

---

## 共通型定義

すべてのフェーズで使用される型定義。

### Island 構造体

```go
// Island は島のデータを表現
type Island struct {
    // 基本情報
    Name     string `json:"name"`
    ID       string `json:"id"`
    Prize    int    `json:"prize"`
    Absent   int    `json:"absent"`
    Comment  string `json:"comment"`
    Password string `json:"-"`  // JSONには含めない

    // リソース
    Money    int `json:"money"`
    Food     int `json:"food"`
    Pop      int `json:"pop"`
    Area     int `json:"area"`

    // 施設
    Farm     int `json:"farm"`
    Factory  int `json:"factory"`
    Mountain int `json:"mountain"`

    // 地形データ
    Land      [][]int `json:"land"`       // 11x11
    LandValue [][]int `json:"land_value"` // 11x11

    // コマンド・掲示板
    Commands []Command   `json:"commands"`
    Lbbs     []LbbsEntry `json:"lbbs"`
}

// NewIsland は新しい島を作成
func NewIsland(name, id, password string) *Island {
    island := &Island{
        Name:      name,
        ID:        id,
        Password:  password,
        Money:     const.InitialMoney,
        Food:      const.InitialFood,
        Area:      const.IslandSize * const.IslandSize,
        Land:      make([][]int, const.IslandSize),
        LandValue: make([][]int, const.IslandSize),
    }

    // 地形初期化
    for i := 0; i < const.IslandSize; i++ {
        island.Land[i] = make([]int, const.IslandSize)
        island.LandValue[i] = make([]int, const.IslandSize)
    }

    return island
}
```

### LandInfo 構造体

```go
// LandInfo は地形情報を表現
type LandInfo struct {
    Type  int  // 地形タイプ
    Value int  // 地形の値
}
```

### Command 構造体

```go
// Command はコマンドを表現
type Command struct {
    Kind   int    `json:"kind"`
    Target string `json:"target"`
    X      int    `json:"x"`
    Y      int    `json:"y"`
    Arg    int    `json:"arg"`
}
```

### LbbsEntry 構造体

```go
// LbbsEntry はローカル掲示板のエントリーを表現
type LbbsEntry struct {
    Name    string `json:"name"`
    Message string `json:"message"`
    Time    int64  `json:"time"`
}
```

### HistoryEntry 構造体

```go
// HistoryEntry は履歴エントリーを表現
type HistoryEntry struct {
    Turn    int    `json:"turn"`
    Message string `json:"message"`
}
```

### LogEntry 構造体

```go
// LogEntry は詳細ログエントリーを表現
type LogEntry struct {
    SecretFlag int    `json:"secret_flag"`  // 0=公開, 1=当事者のみ, 2=機密
    Turn       int    `json:"turn"`
    ID1        string `json:"id1"`  // 当事者
    ID2        string `json:"id2"`  // 相手
    Message    string `json:"message"`
}
```

---

## 使用例

### Phase 1 の使用例

```go
import (
    "net/http"
    "hakoniwa/internal/hako/core"
)

func handler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html; charset=UTF-8")
    core.RunMain(r)
    w.Write(variable.OutputBuffer.Bytes())
}
```

### Phase 2 の使用例

```go
import (
    "context"
    "net/http"
    "hakoniwa/internal/game"
    "hakoniwa/internal/storage"
)

func handler(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()

    // 依存性注入
    storage := storage.NewFileStorage("data")
    lockMgr := lock.NewFileLock("data/hakojima.lock")
    logger := logger.NewStdLogger()
    cfg := LoadConfig()

    game := game.NewGame(cfg, storage, lockMgr, logger)

    // リクエストパース
    req := parseRequest(r)

    // 処理実行
    resp, err := game.HandleRequest(ctx, req)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // レスポンス送信
    w.Header().Set("Content-Type", resp.ContentType)
    w.WriteHeader(resp.StatusCode)
    w.Write(resp.Body)
}
```

---

## 参考資料

- [PERL_TO_GO_MAPPING.md](./PERL_TO_GO_MAPPING.md) - Perl/Go対応表
- [DATA_FORMAT.md](./DATA_FORMAT.md) - データフォーマット仕様

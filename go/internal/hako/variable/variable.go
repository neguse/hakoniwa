// Package variable defines all global variables for the Hakoniwa game.
// This is a literal translation from Perl lib/Hako/Variable.pm
// Phase 1: These are package-level variables (will be refactored in Phase 2)
//
// Ref: perl/lib/Hako/Variable.pm
package variable

import (
	"bytes"

	"github.com/neguse/hakoniwa/internal/hako/types"
)

//----------------------------------------------------------------------
// ユーザ入力値
//----------------------------------------------------------------------

// CurrentID は操作・表示対象の島ID
// Ref: perl/lib/Hako/Variable.pm:61
var CurrentID string

// CurrentNumber はCurrentIDの島番号(@Hislands内のindex)
// Ref: perl/lib/Hako/Variable.pm:64
var CurrentNumber int

// コマンド関連
// Ref: perl/lib/Hako/Variable.pm:67-73

var (
	CommandArg        int    // コマンド引数
	CommandKind       int    // コマンド種類
	CommandMode       string // コマンドモード（write/insert/delete）
	CommandPlanNumber int    // コマンド計画番号
	CommandTarget     string // コマンドターゲットID
	CommandX          int    // コマンドX座標
	CommandY          int    // コマンドY座標
)

// CurrentName は島新規作成・名前変更時の名前
// Ref: perl/lib/Hako/Variable.pm:76
var CurrentName string

// InputPassword はユーザ入力パスワード
// Ref: perl/lib/Hako/Variable.pm:79
var InputPassword string

// InputPassword2 はユーザ入力パスワード（確認用）
// Ref: perl/lib/Hako/Variable.pm:80
var InputPassword2 string

// LbbsMessage はLBBS投稿メッセージ
// Ref: perl/lib/Hako/Variable.pm:83
var LbbsMessage string

// LbbsName はLBBS投稿者名
// Ref: perl/lib/Hako/Variable.pm:86
var LbbsName string

// Message はコメントに設定するメッセージ
// Ref: perl/lib/Hako/Variable.pm:89
var Message string

// OldPassword は変更前パスワード
// Ref: perl/lib/Hako/Variable.pm:92
var OldPassword string

// デフォルト値（Cookieに保存され、あらかじめフォームなどに設定される）
// Ref: perl/lib/Hako/Variable.pm:95-101

var (
	DefaultID       string // デフォルト島ID
	DefaultTarget   string // デフォルトターゲットID
	DefaultKind     int    // デフォルトコマンド種類
	DefaultName     string // デフォルト名前
	DefaultPassword string // デフォルトパスワード
	DefaultX        int    // デフォルトX座標
	DefaultY        int    // デフォルトY座標
)

//----------------------------------------------------------------------
// 島情報（ファイルから読みこんだもの）
//----------------------------------------------------------------------

// IslandTurn はターン数
// Ref: perl/lib/Hako/Variable.pm:108
var IslandTurn int

// IslandLastTime は最終ターン経過時刻
// Ref: perl/lib/Hako/Variable.pm:111
var IslandLastTime int64

// IslandNumber は島の総数
// Ref: perl/lib/Hako/Variable.pm:114
var IslandNumber int

// IslandNextID は次に割り当てる島ID
// Ref: perl/lib/Hako/Variable.pm:117
var IslandNextID int

// Islands は島データのリスト
// Ref: perl/lib/Hako/Variable.pm:120
var Islands []*types.Island

//----------------------------------------------------------------------
// テンポラリな情報
//----------------------------------------------------------------------

// MainMode はモード（turn, new, print, owner, command, comment, lbbs, change, top）
// Ref: perl/lib/Hako/Variable.pm:127
var MainMode string

// LbbsMode はLBBSモード（0:観光者, 1:島主, 2:削除）
// Ref: perl/lib/Hako/Variable.pm:130
var LbbsMode int

// IslandList は島セレクト用OPTIONコントロール（デフォルト自分）
// Ref: perl/lib/Hako/Variable.pm:133
var IslandList string

// TargetList はTargetの島セレクト用OPTIONコントロール
// Ref: perl/lib/Hako/Variable.pm:136
var TargetList string

// IDToName は島ID->名前のマップ
// Ref: perl/lib/Hako/Variable.pm:139
var IDToName map[string]string

// IDToNumber は島ID->島番号のマップ
// Ref: perl/lib/Hako/Variable.pm:142
var IDToNumber map[string]int

// 島の全座標をランダムにシャッフルしたもの
// Ref: perl/lib/Hako/Variable.pm:145-146

var (
	Rpx []int // ランダムX座標配列
	Rpy []int // ランダムY座標配列
)

// DefenceHex は防衛施設の座標 [島ID][X][Y]
// Ref: perl/lib/Hako/Variable.pm:149
var DefenceHex [][][]bool

// LogPool は通常ログ
// Ref: perl/lib/Hako/Variable.pm:152
var LogPool []string

// LateLogPool は遅延ログ
// Ref: perl/lib/Hako/Variable.pm:155
var LateLogPool []string

// SecretLogPool は機密ログ
// Ref: perl/lib/Hako/Variable.pm:158
var SecretLogPool []string

// LockID はロック用ファイルハンドル（lockMode == 2の時に利用）
// Ref: perl/lib/Hako/Variable.pm:161
var LockID interface{}

// OutputBuffer はHTML出力用バッファ（Phase 1用）
var OutputBuffer bytes.Buffer

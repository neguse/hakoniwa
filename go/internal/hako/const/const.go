// Package const defines all constants and configuration values for the Hakoniwa game.
// This is a literal translation from Perl lib/Hako/Const.pm
//
// Ref: perl/lib/Hako/Const.pm
package const

//----------------------------------------------------------------------
// 各種設定値
//----------------------------------------------------------------------

//----------------------------------------------------------------------
// 必ず設定する部分
//----------------------------------------------------------------------

// BaseDir はCGIファイルを置くディレクトリ
// Ref: perl/lib/Hako/Const.pm:192
var BaseDir = "http://localhost:5000/cgi-bin"

// ImageDir は画像ファイルを置くディレクトリ
// Ref: perl/lib/Hako/Const.pm:196
var ImageDir = "http://localhost:5000/images"

// MasterPassword はすべての島のパスワードを代用できるマスターパスワード
// Ref: perl/lib/Hako/Const.pm:201
var MasterPassword = "yourpassword"

// SpecialPassword は特殊パスワード（名前変更時に資金・食料が最大になる）
// Ref: perl/lib/Hako/Const.pm:206
var SpecialPassword = "yourspecialpassword"

// AdminName は管理者名
// Ref: perl/lib/Hako/Const.pm:209
var AdminName = "管理者の名前"

// Email は管理者のメールアドレス
// Ref: perl/lib/Hako/Const.pm:212
var Email = "管理者@どこか.どこか.どこか"

// BBS は掲示板アドレス
// Ref: perl/lib/Hako/Const.pm:215
var BBS = "http://サーバー/掲示板.cgi"

// TopPage はホームページのアドレス
// Ref: perl/lib/Hako/Const.pm:218
var TopPage = "http://サーバー/ホームページ.html"

// DirMode はディレクトリのパーミッション
// Ref: perl/lib/Hako/Const.pm:222
const DirMode = 0755

// DirName はデータディレクトリの名前
// Ref: perl/lib/Hako/Const.pm:228
var DirName = "data"

// LockMode はロックの方式
// 1: ディレクトリ, 2: システムコール(flock), 3: シンボリックリンク, 4: 通常ファイル
// Ref: perl/lib/Hako/Const.pm:237
const LockMode = 2

//----------------------------------------------------------------------
// ゲームの進行やファイルなど
//----------------------------------------------------------------------

// UnitTime は1ターンが何秒か (デフォルト: 6時間)
// Ref: perl/lib/Hako/Const.pm:254
const UnitTime = 21600

// UnlockTime は異常終了基準時間（ロック後何秒で強制解除するか）
// Ref: perl/lib/Hako/Const.pm:258
const UnlockTime = 120

// MaxIsland は島の最大数
// Ref: perl/lib/Hako/Const.pm:261
const MaxIsland = 30

// TopLogTurn はトップページに表示するログのターン数
// Ref: perl/lib/Hako/Const.pm:264
const TopLogTurn = 1

// LogMax はログファイル保持ターン数
// Ref: perl/lib/Hako/Const.pm:267
const LogMax = 8

// BackupTurn はバックアップを何ターンおきに取るか
// Ref: perl/lib/Hako/Const.pm:270
const BackupTurn = 12

// BackupTimes はバックアップを何回分残すか
// Ref: perl/lib/Hako/Const.pm:273
const BackupTimes = 4

// HistoryMax は発見ログ保持行数
// Ref: perl/lib/Hako/Const.pm:276
const HistoryMax = 10

// GiveupTurn は放棄コマンド自動入力ターン数
// Ref: perl/lib/Hako/Const.pm:279
const GiveupTurn = 28

// CommandMax はコマンド入力限界数
// Ref: perl/lib/Hako/Const.pm:283
const CommandMax = 20

// UseLbbs はローカル掲示板を使用するかどうか (0:使用しない、1:使用する)
// Ref: perl/lib/Hako/Const.pm:286
const UseLbbs = 0

// LbbsMax はローカル掲示板行数
// Ref: perl/lib/Hako/Const.pm:289
const LbbsMax = 10

// IslandSize は島の大きさ
// Ref: perl/lib/Hako/Const.pm:293
const IslandSize = 12

// HideMoneyMode は他人から資金を見えなくするか (0: 見えない, 1: 見える, 2: 100の位で四捨五入)
// Ref: perl/lib/Hako/Const.pm:299
const HideMoneyMode = 2

// CryptOn はパスワードの暗号化 (false: 暗号化しない, true: 暗号化する)
// Ref: perl/lib/Hako/Const.pm:302
const CryptOn = true

// Debug はデバッグモード (trueだと、「ターンを進める」ボタンが使用できる)
// Ref: perl/lib/Hako/Const.pm:305
const Debug = false

//----------------------------------------------------------------------
// 資金、食料などの設定値と単位
//----------------------------------------------------------------------

// InitialMoney は初期資金
// Ref: perl/lib/Hako/Const.pm:311
const InitialMoney = 100

// InitialFood は初期食料
// Ref: perl/lib/Hako/Const.pm:314
const InitialFood = 100

// UnitMoney はお金の単位
// Ref: perl/lib/Hako/Const.pm:317
const UnitMoney = "億円"

// UnitFood は食料の単位
// Ref: perl/lib/Hako/Const.pm:320
const UnitFood = "00トン"

// UnitPop は人口の単位
// Ref: perl/lib/Hako/Const.pm:323
const UnitPop = "00人"

// UnitArea は広さの単位
// Ref: perl/lib/Hako/Const.pm:326
const UnitArea = "00万坪"

// UnitTree は木の数の単位
// Ref: perl/lib/Hako/Const.pm:329
const UnitTree = "00本"

// TreeValue は木の単位当たりの売値
// Ref: perl/lib/Hako/Const.pm:332
const TreeValue = 5

// CostChangeName は名前変更のコスト
// Ref: perl/lib/Hako/Const.pm:335
const CostChangeName = 500

// EatenFood は人口1単位あたりの食料消費料
// Ref: perl/lib/Hako/Const.pm:338
const EatenFood = 0.2

//----------------------------------------------------------------------
// 基地の経験値
//----------------------------------------------------------------------

// MaxExpPoint は経験値の最大値 (ただし、最大でも255まで)
// Ref: perl/lib/Hako/Const.pm:344
const MaxExpPoint = 200

// MaxBaseLevel はミサイル基地のレベルの最大値
// Ref: perl/lib/Hako/Const.pm:347
const MaxBaseLevel = 5

// MaxSBaseLevel は海底基地のレベルの最大値
// Ref: perl/lib/Hako/Const.pm:348
const MaxSBaseLevel = 3

// BaseLevelUp はミサイル基地の経験値がいくつでレベルアップか
// Ref: perl/lib/Hako/Const.pm:352
var BaseLevelUp = []int{20, 60, 120, 200}

// SBaseLevelUp は海底基地の経験値がいくつでレベルアップか
// Ref: perl/lib/Hako/Const.pm:353
var SBaseLevelUp = []int{50, 200}

//----------------------------------------------------------------------
// 防衛施設の自爆
//----------------------------------------------------------------------

// DBaseAuto は怪獣に踏まれた時自爆するかどうか (true: 自爆する, false: 自爆しない)
// Ref: perl/lib/Hako/Const.pm:359
const DBaseAuto = true

//----------------------------------------------------------------------
// 災害
//----------------------------------------------------------------------

// 通常災害発生率 (確率は0.1%単位)
// Ref: perl/lib/Hako/Const.pm:365-372

// DisEarthquake は地震の発生率
const DisEarthquake = 5

// DisTsunami は津波の発生率
const DisTsunami = 15

// DisTyphoon は台風の発生率
const DisTyphoon = 20

// DisMeteo は隕石の発生率
const DisMeteo = 15

// DisHugeMeteo は巨大隕石の発生率
const DisHugeMeteo = 5

// DisEruption は噴火の発生率
const DisEruption = 10

// DisFire は火災の発生率
const DisFire = 10

// DisMaizo は埋蔵金の発生率
const DisMaizo = 10

// 地盤沈下
// Ref: perl/lib/Hako/Const.pm:375-376

// DisFallBorder は地盤沈下の安全限界の広さ (Hex数)
const DisFallBorder = 90

// DisFalldown は安全限界を超えた場合の地盤沈下確率
const DisFalldown = 30

// 怪獣
// Ref: perl/lib/Hako/Const.pm:379-390

// DisMonsBorder1 は怪獣レベル1の人口基準
const DisMonsBorder1 = 1000

// DisMonsBorder2 は怪獣レベル2の人口基準
const DisMonsBorder2 = 2500

// DisMonsBorder3 は怪獣レベル3の人口基準
const DisMonsBorder3 = 4000

// DisMonster は怪獣の単位面積あたりの出現率 (0.01%単位)
const DisMonster = 3

// MonsterNumber は怪獣の種類数
const MonsterNumber = 8

// MonsterLevel1 はレベル1で出現する怪獣の番号の最大値 (サンジラまで)
const MonsterLevel1 = 2

// MonsterLevel2 はレベル2で出現する怪獣の番号の最大値 (いのらゴーストまで)
const MonsterLevel2 = 5

// MonsterLevel3 はレベル3で出現する怪獣の番号の最大値 (キングいのらまで=全部)
const MonsterLevel3 = 7

// MonsterName は怪獣の名前
// Ref: perl/lib/Hako/Const.pm:393-402
var MonsterName = []string{
	"メカいのら",       // 0 (人造)
	"いのら",          // 1
	"サンジラ",        // 2
	"レッドいのら",    // 3
	"ダークいのら",    // 4
	"いのらゴースト",  // 5
	"クジラ",          // 6
	"キングいのら",    // 7
}

// MonsterBHP は怪獣の最低体力
// Ref: perl/lib/Hako/Const.pm:405
var MonsterBHP = []int{2, 1, 1, 3, 2, 1, 4, 5}

// MonsterDHP は怪獣の体力の幅
// Ref: perl/lib/Hako/Const.pm:406
var MonsterDHP = []int{0, 2, 2, 2, 2, 0, 2, 2}

// MonsterSpecial は怪獣の特殊能力
// 0: 特になし, 1: 足が速い(最大2歩), 2: 足がとても速い(何歩か不明),
// 3: 奇数ターンは硬化, 4: 偶数ターンは硬化
// Ref: perl/lib/Hako/Const.pm:407
var MonsterSpecial = []int{0, 0, 3, 0, 1, 2, 4, 0}

// MonsterExp は怪獣の経験値
// Ref: perl/lib/Hako/Const.pm:408
var MonsterExp = []int{5, 5, 7, 12, 15, 10, 20, 30}

// MonsterValue は怪獣の死体の値段
// Ref: perl/lib/Hako/Const.pm:409
var MonsterValue = []int{0, 400, 500, 1000, 800, 300, 1500, 2000}

// MonsterImage は怪獣の画像ファイル
// Ref: perl/lib/Hako/Const.pm:419-422
var MonsterImage = []string{
	"monster7.gif", "monster0.gif", "monster5.gif", "monster1.gif",
	"monster2.gif", "monster8.gif", "monster6.gif", "monster3.gif",
}

// MonsterImage2 は怪獣の画像ファイルその2 (硬化中)
// Ref: perl/lib/Hako/Const.pm:425-426
var MonsterImage2 = []string{
	"", "", "monster4.gif", "", "", "", "monster4.gif", "",
}

//----------------------------------------------------------------------
// 油田
//----------------------------------------------------------------------

// OilMoney は油田の収入
// Ref: perl/lib/Hako/Const.pm:432
const OilMoney = 1000

// OilRatio は油田の枯渇確率
// Ref: perl/lib/Hako/Const.pm:435
const OilRatio = 40

//----------------------------------------------------------------------
// 記念碑
//----------------------------------------------------------------------

// MonumentNumber は記念碑の種類数
// Ref: perl/lib/Hako/Const.pm:441
const MonumentNumber = 3

// MonumentName は記念碑の名前
// Ref: perl/lib/Hako/Const.pm:444
var MonumentName = []string{"モノリス", "平和記念碑", "戦いの碑"}

// MonumentImage は記念碑の画像ファイル
// Ref: perl/lib/Hako/Const.pm:447
var MonumentImage = []string{"monument0.gif", "monument0.gif", "monument0.gif"}

//----------------------------------------------------------------------
// 賞関係
//----------------------------------------------------------------------

// TurnPrizeUnit はターン杯を何ターン毎に出すか
// Ref: perl/lib/Hako/Const.pm:453
const TurnPrizeUnit = 100

// Prize は賞の名前
// Ref: perl/lib/Hako/Const.pm:456-462
var Prize = []string{
	"ターン杯", "繁栄賞",
	"超繁栄賞", "究極繁栄賞",
	"平和賞", "超平和賞",
	"究極平和賞", "災難賞",
	"超災難賞", "究極災難賞",
}

//----------------------------------------------------------------------
// 外見関係
//----------------------------------------------------------------------

// HTMLBody は<BODY>タグのオプション
// Ref: perl/lib/Hako/Const.pm:468
const HTMLBody = `BGCOLOR="#EEFFFF"`

// Title はゲームのタイトル文字
// Ref: perl/lib/Hako/Const.pm:471
const Title = "箱庭諸島２"

// タグ
// Ref: perl/lib/Hako/Const.pm:475-516

// TagTitleBegin はタイトル文字開始タグ
const TagTitleBegin = `<FONT SIZE=7 COLOR="#8888ff">`

// TagTitleEnd はタイトル文字終了タグ
const TagTitleEnd = `</FONT>`

// TagHeaderBegin はH1タグ用開始タグ
const TagHeaderBegin = `<FONT COLOR="#4444ff">`

// TagHeaderEnd はH1タグ用終了タグ
const TagHeaderEnd = `</FONT>`

// TagBigBegin は大きい文字開始タグ
const TagBigBegin = `<FONT SIZE=6>`

// TagBigEnd は大きい文字終了タグ
const TagBigEnd = `</FONT>`

// TagNameBegin は島の名前など開始タグ
const TagNameBegin = `<FONT COLOR="#a06040"><B>`

// TagNameEnd は島の名前など終了タグ
const TagNameEnd = `</B></FONT>`

// TagName2Begin は薄くなった島の名前開始タグ
const TagName2Begin = `<FONT COLOR="#808080"><B>`

// TagName2End は薄くなった島の名前終了タグ
const TagName2End = `</B></FONT>`

// TagNumberBegin は順位の番号など開始タグ
const TagNumberBegin = `<FONT COLOR="#800000"><B>`

// TagNumberEnd は順位の番号など終了タグ
const TagNumberEnd = `</B></FONT>`

// TagTHBegin は順位表における見だし開始タグ
const TagTHBegin = `<FONT COLOR="#C00000"><B>`

// TagTHEnd は順位表における見だし終了タグ
const TagTHEnd = `</B></FONT>`

// TagComNameBegin は開発計画の名前開始タグ
const TagComNameBegin = `<FONT COLOR="#d08000"><B>`

// TagComNameEnd は開発計画の名前終了タグ
const TagComNameEnd = `</B></FONT>`

// TagDisasterBegin は災害開始タグ
const TagDisasterBegin = `<FONT COLOR="#ff0000"><B>`

// TagDisasterEnd は災害終了タグ
const TagDisasterEnd = `</B></FONT>`

// TagLbbsSSBegin はローカル掲示板、観光者の書いた文字開始タグ
const TagLbbsSSBegin = `<FONT COLOR="#0000ff"><B>`

// TagLbbsSSEnd はローカル掲示板、観光者の書いた文字終了タグ
const TagLbbsSSEnd = `</B></FONT>`

// TagLbbsOWBegin はローカル掲示板、島主の書いた文字開始タグ
const TagLbbsOWBegin = `<FONT COLOR="#ff0000"><B>`

// TagLbbsOWEnd はローカル掲示板、島主の書いた文字終了タグ
const TagLbbsOWEnd = `</B></FONT>`

// NormalColor は通常の文字色
// Ref: perl/lib/Hako/Const.pm:519
const NormalColor = "#000000"

// 順位表、セルの属性
// Ref: perl/lib/Hako/Const.pm:522-529

// BgTitleCell は順位表見出しのセル属性
const BgTitleCell = `BGCOLOR="#ccffcc"`

// BgNumberCell は順位表順位のセル属性
const BgNumberCell = `BGCOLOR="#ccffcc"`

// BgNameCell は順位表島の名前のセル属性
const BgNameCell = `BGCOLOR="#ccffff"`

// BgInfoCell は順位表島の情報のセル属性
const BgInfoCell = `BGCOLOR="#ccffff"`

// BgCommentCell は順位表コメント欄のセル属性
const BgCommentCell = `BGCOLOR="#ccffcc"`

// BgInputCell は開発計画フォームのセル属性
const BgInputCell = `BGCOLOR="#ccffcc"`

// BgMapCell は開発計画地図のセル属性
const BgMapCell = `BGCOLOR="#ccffcc"`

// BgCommandCell は開発計画入力済み計画のセル属性
const BgCommandCell = `BGCOLOR="#ccffcc"`

//----------------------------------------------------------------------
// これ以降のスクリプトは、変更されることを想定していません
//----------------------------------------------------------------------

// ThisFile はこのファイルのURL
// Ref: perl/lib/Hako/Const.pm:545
var ThisFile = BaseDir + "/hako-main.cgi"

//----------------------------------------------------------------------
// 地形番号
//----------------------------------------------------------------------
// Ref: perl/lib/Hako/Const.pm:548-562

const (
	LandSea      = 0  // 海
	LandWaste    = 1  // 荒地
	LandPlains   = 2  // 平地
	LandTown     = 3  // 町系
	LandForest   = 4  // 森
	LandFarm     = 5  // 農場
	LandFactory  = 6  // 工場
	LandBase     = 7  // ミサイル基地
	LandDefence  = 8  // 防衛施設
	LandMountain = 9  // 山
	LandMonster  = 10 // 怪獣
	LandSbase    = 11 // 海底基地
	LandOil      = 12 // 海底油田
	LandMonument = 13 // 記念碑
	LandHaribote = 14 // ハリボテ
)

//----------------------------------------------------------------------
// コマンド
//----------------------------------------------------------------------
// Ref: perl/lib/Hako/Const.pm:565-604

// CommandTotal はコマンドの種類数
const CommandTotal = 28

// コマンド番号
const (
	// 整地系
	ComPrepare  = 1  // 整地
	ComPrepare2 = 2  // 地ならし
	ComReclaim  = 3  // 埋め立て
	ComDestroy  = 4  // 掘削
	ComSellTree = 5  // 伐採

	// 作る系
	ComPlant    = 11 // 植林
	ComFarm     = 12 // 農場整備
	ComFactory  = 13 // 工場建設
	ComMountain = 14 // 採掘場整備
	ComBase     = 15 // ミサイル基地建設
	ComDbase    = 16 // 防衛施設建設
	ComSbase    = 17 // 海底基地建設
	ComMonument = 18 // 記念碑建造
	ComHaribote = 19 // ハリボテ設置

	// 発射系
	ComMissileNM   = 31 // ミサイル発射
	ComMissilePP   = 32 // PPミサイル発射
	ComMissileST   = 33 // STミサイル発射
	ComMissileLD   = 34 // 陸地破壊弾発射
	ComSendMonster = 35 // 怪獣派遣

	// 運営系
	ComDoNothing  = 41 // 資金繰り
	ComSell       = 42 // 食料輸出
	ComMoney      = 43 // 資金援助
	ComFood       = 44 // 食料援助
	ComPropaganda = 45 // 誘致活動
	ComGiveup     = 46 // 島の放棄

	// 自動入力系
	ComAutoPrepare  = 61 // フル整地
	ComAutoPrepare2 = 62 // フル地ならし
	ComAutoDelete   = 63 // 全コマンド消去
)

// ComList はコマンドの順番
// Ref: perl/lib/Hako/Const.pm:607-618
var ComList = []int{
	ComPrepare, ComSell, ComPrepare2,
	ComReclaim, ComDestroy, ComSellTree,
	ComPlant, ComFarm, ComFactory,
	ComMountain, ComBase, ComDbase,
	ComSbase, ComMonument, ComHaribote,
	ComMissileNM, ComMissilePP, ComMissileST,
	ComMissileLD, ComSendMonster, ComDoNothing,
	ComMoney, ComFood, ComPropaganda,
	ComGiveup, ComAutoPrepare, ComAutoPrepare2,
	ComAutoDelete,
}

// ComName はコマンドの名前
// Ref: perl/lib/Hako/Const.pm:623-677
var ComName = map[int]string{
	ComPrepare:      "整地",
	ComPrepare2:     "地ならし",
	ComReclaim:      "埋め立て",
	ComDestroy:      "掘削",
	ComSellTree:     "伐採",
	ComPlant:        "植林",
	ComFarm:         "農場整備",
	ComFactory:      "工場建設",
	ComMountain:     "採掘場整備",
	ComBase:         "ミサイル基地建設",
	ComDbase:        "防衛施設建設",
	ComSbase:        "海底基地建設",
	ComMonument:     "記念碑建造",
	ComHaribote:     "ハリボテ設置",
	ComMissileNM:    "ミサイル発射",
	ComMissilePP:    "PPミサイル発射",
	ComMissileST:    "STミサイル発射",
	ComMissileLD:    "陸地破壊弾発射",
	ComSendMonster:  "怪獣派遣",
	ComDoNothing:    "資金繰り",
	ComSell:         "食料輸出",
	ComMoney:        "資金援助",
	ComFood:         "食料援助",
	ComPropaganda:   "誘致活動",
	ComGiveup:       "島の放棄",
	ComAutoPrepare:  "整地自動入力",
	ComAutoPrepare2: "地ならし自動入力",
	ComAutoDelete:   "全計画を白紙撤回",
}

// ComCost はコマンドの値段
// Ref: perl/lib/Hako/Const.pm:623-678
var ComCost = map[int]int{
	ComPrepare:      5,
	ComPrepare2:     100,
	ComReclaim:      150,
	ComDestroy:      200,
	ComSellTree:     0,
	ComPlant:        50,
	ComFarm:         20,
	ComFactory:      100,
	ComMountain:     300,
	ComBase:         300,
	ComDbase:        800,
	ComSbase:        8000,
	ComMonument:     9999,
	ComHaribote:     1,
	ComMissileNM:    20,
	ComMissilePP:    50,
	ComMissileST:    50,
	ComMissileLD:    100,
	ComSendMonster:  3000,
	ComDoNothing:    0,
	ComSell:         -100,
	ComMoney:        100,
	ComFood:         -100,
	ComPropaganda:   1000,
	ComGiveup:       0,
	ComAutoPrepare:  0,
	ComAutoPrepare2: 0,
	ComAutoDelete:   0,
}

// PointNumber は島の座標数
// Ref: perl/lib/Hako/Const.pm:681
const PointNumber = IslandSize * IslandSize

// TempBack は「戻る」リンク
// Ref: perl/lib/Hako/Const.pm:684-685
var TempBack = `<A HREF="` + ThisFile + `">` + TagBigBegin + `トップへ戻る` + TagBigEnd + `</A>`

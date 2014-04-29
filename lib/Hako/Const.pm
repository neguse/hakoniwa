package Hako::Const;

use strict;
use warnings;
use utf8;

use Exporter::Easy (
  EXPORT => [qw(
    $baseDir
    $imageDir
    $masterPassword
    $HspecialPassword
    $adminName
    $email
    $bbs
    $toppage
    $HdirMode
    $HdirName
    $lockMode
    $HunitTime
    $unlockTime
    $HmaxIsland
    $HtopLogTurn
    $HlogMax
    $HbackupTurn
    $HbackupTimes
    $HhistoryMax
    $HgiveupTurn
    $HcommandMax
    $HuseLbbs
    $HlbbsMax
    $HislandSize
    $HhideMoneyMode
    $cryptOn
    $Hdebug
    $HinitialMoney
    $HinitialFood
    $HunitMoney
    $HunitFood
    $HunitPop
    $HunitArea
    $HunitTree
    $HtreeValue
    $HcostChangeName
    $HeatenFood
    $HmaxExpPoint
    $maxBaseLevel
    $maxSBaseLevel
    @baseLevelUp
    @sBaseLevelUp
    $HdBaseAuto
    $HdisEarthquake
    $HdisTsunami
    $HdisTyphoon
    $HdisMeteo
    $HdisHugeMeteo
    $HdisEruption
    $HdisFire
    $HdisMaizo
    $HdisFallBorder
    $HdisFalldown
    $HdisMonsBorder1
    $HdisMonsBorder2
    $HdisMonsBorder3
    $HdisMonster
    $HmonsterNumber
    $HmonsterLevel1
    $HmonsterLevel2
    $HmonsterLevel3
    @HmonsterName
    @HmonsterBHP
    @HmonsterDHP
    @HmonsterSpecial
    @HmonsterExp
    @HmonsterValue
    @HmonsterImage
    @HmonsterImage2
    $HoilMoney
    $HoilRatio
    $HmonumentNumber
    @HmonumentName
    @HmonumentImage
    $HturnPrizeUnit
    @Hprize
    $htmlBody
    $Htitle
    $HtagTitle_
    $H_tagTitle
    $HtagHeader_
    $H_tagHeader
    $HtagBig_
    $H_tagBig
    $HtagName_
    $H_tagName
    $HtagName2_
    $H_tagName2
    $HtagNumber_
    $H_tagNumber
    $HtagTH_
    $H_tagTH
    $HtagComName_
    $H_tagComName
    $HtagDisaster_
    $H_tagDisaster
    $HtagLbbsSS_
    $H_tagLbbsSS
    $HtagLbbsOW_
    $H_tagLbbsOW
    $HnormalColor
    $HbgTitleCell
    $HbgNumberCell
    $HbgNameCell
    $HbgInfoCell
    $HbgCommentCell
    $HbgInputCell
    $HbgMapCell
    $HbgCommandCell
    $HthisFile
    $HlandSea
    $HlandWaste
    $HlandPlains
    $HlandTown
    $HlandForest
    $HlandFarm
    $HlandFactory
    $HlandBase
    $HlandDefence
    $HlandMountain
    $HlandMonster
    $HlandSbase
    $HlandOil
    $HlandMonument
    $HlandHaribote
    $HcommandTotal
    $HcomPrepare
    $HcomPrepare2
    $HcomReclaim
    $HcomDestroy
    $HcomSellTree
    $HcomPlant
    $HcomFarm
    $HcomFactory
    $HcomMountain
    $HcomBase
    $HcomDbase
    $HcomSbase
    $HcomMonument
    $HcomHaribote
    $HcomMissileNM
    $HcomMissilePP
    $HcomMissileST
    $HcomMissileLD
    $HcomSendMonster
    $HcomDoNothing
    $HcomSell
    $HcomMoney
    $HcomFood
    $HcomPropaganda
    $HcomGiveup
    $HcomAutoPrepare
    $HcomAutoPrepare2
    $HcomAutoDelete
    @HcomList
    @HcomName
    @HcomCost
    $HpointNumber
    $HtempBack
  )],
);

#----------------------------------------------------------------------
# 各種設定値
# (これ以降の部分の各設定値を、適切な値に変更してください)
#----------------------------------------------------------------------

#----------------------------------------------------------------------
# 以下、必ず設定する部分
#----------------------------------------------------------------------

# このファイルを置くディレクトリ
# our($baseDir) = 'http://サーバー/ディレクトリ';
#
# 例)
# http://cgi2.bekkoame.ne.jp/cgi-bin/user/u5534/hakoniwa/hako-main.cgi
# として置く場合、
# our($baseDir) = 'http://cgi2.bekkoame.ne.jp/cgi-bin/user/u5534/hakoniwa';
# とする。最後にスラッシュ(/)は付けない。

# our($baseDir) = 'http://サーバー/ディレクトリ';
our ($baseDir) = 'http://localhost:5000/cgi-bin';

# 画像ファイルを置くディレクトリ
# our($imageDir) = 'http://サーバー/ディレクトリ';
our ($imageDir) = 'http://localhost:5000/images';

# マスターパスワード
# このパスワードは、すべての島のパスワードを代用できます。
# 例えば、「他の島のパスワード変更」等もできます。
our ($masterPassword) = 'yourpassword';

# 特殊パスワード
# このパスワードで「名前変更」を行うと、その島の資金、食料が最大値になります。
# (実際に名前を変える必要はありません。)
our $HspecialPassword = 'yourspecialpassword';

# 管理者名
our ($adminName) = '管理者の名前';

# 管理者のメールアドレス
our ($email) = '管理者@どこか.どこか.どこか';

# 掲示板アドレス
our ($bbs) = 'http://サーバー/掲示板.cgi';

# ホームページのアドレス
our ($toppage) = 'http://サーバー/ホームページ.html';

# ディレクトリのパーミッション
# 通常は0755でよいが、0777、0705、0704等でないとできないサーバーもあるらしい
our $HdirMode = 0755;

# データディレクトリの名前
# ここで設定した名前のディレクトリ以下にデータが格納されます。
# デフォルトでは'data'となっていますが、セキュリティのため
# なるべく違う名前に変更してください。
our $HdirName = 'data';

# データの書き込み方

# ロックの方式
# 1 ディレクトリ
# 2 システムコール(可能ならば最も望ましい)
# 3 シンボリックリンク
# 4 通常ファイル(あまりお勧めでない)
our ($lockMode) = 2;

# (注)
# 4を選択する場合には、'key-free'という、パーミション666の空のファイルを、
# このファイルと同位置に置いて下さい。

#----------------------------------------------------------------------
# 必ず設定する部分は以上
#----------------------------------------------------------------------

#----------------------------------------------------------------------
# 以下、好みによって設定する部分
#----------------------------------------------------------------------
#----------------------------------------
# ゲームの進行やファイルなど
#----------------------------------------
# 1ターンが何秒か
our $HunitTime = 21600;    # 6時間

# 異常終了基準時間
# (ロック後何秒で、強制解除するか)
our ($unlockTime) = 120;

# 島の最大数
our $HmaxIsland = 30;

# トップページに表示するログのターン数
our $HtopLogTurn = 1;

# ログファイル保持ターン数
our $HlogMax = 8;

# バックアップを何ターンおきに取るか
our $HbackupTurn = 12;

# バックアップを何回分残すか
our $HbackupTimes = 4;

# 発見ログ保持行数
our $HhistoryMax = 10;

# 放棄コマンド自動入力ターン数
our $HgiveupTurn = 28;

# コマンド入力限界数
# (ゲームが始まってから変更すると、データファイルの互換性が無くなります。)
our $HcommandMax = 20;

# ローカル掲示板行数を使用するかどうか(0:使用しない、1:使用する)
our $HuseLbbs = 0;

# ローカル掲示板行数
our $HlbbsMax = 10;

# 島の大きさ
# (変更できないかも)
our $HislandSize = 12;

# 他人から資金を見えなくするか
# 0 見えない
# 1 見える
# 2 100の位で四捨五入
our $HhideMoneyMode = 2;

# パスワードの暗号化(0だと暗号化しない、1だと暗号化する)
our ($cryptOn) = 1;

# デバッグモード(1だと、「ターンを進める」ボタンが使用できる)
our $Hdebug = 0;

#----------------------------------------
# 資金、食料などの設定値と単位
#----------------------------------------
# 初期資金
our $HinitialMoney = 100;

# 初期食料
our $HinitialFood = 100;

# お金の単位
our $HunitMoney = '億円';

# 食料の単位
our $HunitFood = '00トン';

# 人口の単位
our $HunitPop = '00人';

# 広さの単位
our $HunitArea = '00万坪';

# 木の数の単位
our $HunitTree = '00本';

# 木の単位当たりの売値
our $HtreeValue = 5;

# 名前変更のコスト
our $HcostChangeName = 500;

# 人口1単位あたりの食料消費料
our $HeatenFood = 0.2;

#----------------------------------------
# 基地の経験値
#----------------------------------------
# 経験値の最大値
our $HmaxExpPoint = 200;    # ただし、最大でも255まで

# レベルの最大値
our ($maxBaseLevel)  = 5;    # ミサイル基地
our ($maxSBaseLevel) = 3;    # 海底基地

# 経験値がいくつでレベルアップか
our ( @baseLevelUp, @sBaseLevelUp );
@baseLevelUp = ( 20, 60, 120, 200 );    # ミサイル基地
@sBaseLevelUp = ( 50, 200 );            # 海底基地

#----------------------------------------
# 防衛施設の自爆
#----------------------------------------
# 怪獣に踏まれた時自爆するなら1、しないなら0
our $HdBaseAuto = 1;

#----------------------------------------
# 災害
#----------------------------------------
# 通常災害発生率(確率は0.1%単位)
our $HdisEarthquake = 5;     # 地震
our $HdisTsunami    = 15;    # 津波
our $HdisTyphoon    = 20;    # 台風
our $HdisMeteo      = 15;    # 隕石
our $HdisHugeMeteo  = 5;     # 巨大隕石
our $HdisEruption   = 10;    # 噴火
our $HdisFire       = 10;    # 火災
our $HdisMaizo      = 10;    # 埋蔵金

# 地盤沈下
our $HdisFallBorder = 90;    # 安全限界の広さ(Hex数)
our $HdisFalldown   = 30;    # その広さを超えた場合の確率

# 怪獣
our $HdisMonsBorder1 = 1000;  # 人口基準1(怪獣レベル1)
our $HdisMonsBorder2 = 2500;  # 人口基準2(怪獣レベル2)
our $HdisMonsBorder3 = 4000;  # 人口基準3(怪獣レベル3)
our $HdisMonster     = 3;     # 単位面積あたりの出現率(0.01%単位)

# 種類
our $HmonsterNumber = 8;

# 各基準において出てくる怪獣の番号の最大値
our $HmonsterLevel1 = 2;      # サンジラまで
our $HmonsterLevel2 = 5;      # いのらゴーストまで
our $HmonsterLevel3 = 7;      # キングいのらまで(全部)

# 名前
our @HmonsterName = (
    'メカいのら',          # 0(人造)
    'いのら',                # 1
    'サンジラ',             # 2
    'レッドいのら',       # 3
    'ダークいのら',       # 4
    'いのらゴースト',    # 5
    'クジラ',                # 6
    'キングいのら'        # 7
);

# 最低体力、体力の幅、特殊能力、経験値、死体の値段
our @HmonsterBHP     = ( 2, 1,   1,   3,    2,   1,   4,    5 );
our @HmonsterDHP     = ( 0, 2,   2,   2,    2,   0,   2,    2 );
our @HmonsterSpecial = ( 0, 0,   3,   0,    1,   2,   4,    0 );
our @HmonsterExp     = ( 5, 5,   7,   12,   15,  10,  20,   30 );
our @HmonsterValue   = ( 0, 400, 500, 1000, 800, 300, 1500, 2000 );

# 特殊能力の内容は、
# 0 特になし
# 1 足が速い(最大2歩あるく)
# 2 足がとても速い(最大何歩あるくか不明)
# 3 奇数ターンは硬化
# 4 偶数ターンは硬化

# 画像ファイル
our @HmonsterImage = (
    'monster7.gif', 'monster0.gif', 'monster5.gif', 'monster1.gif',
    'monster2.gif', 'monster8.gif', 'monster6.gif', 'monster3.gif'
);

# 画像ファイルその2(硬化中)
our @HmonsterImage2
    = ( '', '', 'monster4.gif', '', '', '', 'monster4.gif', '' );

#----------------------------------------
# 油田
#----------------------------------------
# 油田の収入
our $HoilMoney = 1000;

# 油田の枯渇確率
our $HoilRatio = 40;

#----------------------------------------
# 記念碑
#----------------------------------------
# 何種類あるか
our $HmonumentNumber = 3;

# 名前
our @HmonumentName = ( 'モノリス', '平和記念碑', '戦いの碑' );

# 画像ファイル
our @HmonumentImage = ( 'monument0.gif', 'monument0.gif', 'monument0.gif' );

#----------------------------------------
# 賞関係
#----------------------------------------
# ターン杯を何ターン毎に出すか
our $HturnPrizeUnit = 100;

# 賞の名前
our @Hprize = (
    'ターン杯',    '繁栄賞',
    '超繁栄賞',    '究極繁栄賞',
    '平和賞',       '超平和賞',
    '究極平和賞', '災難賞',
    '超災難賞',    '究極災難賞',
);

#----------------------------------------
# 外見関係
#----------------------------------------
# <BODY>タグのオプション
our ($htmlBody) = 'BGCOLOR="#EEFFFF"';

# ゲームのタイトル文字
our $Htitle = '箱庭諸島２';

# タグ
# タイトル文字
our $HtagTitle_ = '<FONT SIZE=7 COLOR="#8888ff">';
our $H_tagTitle = '</FONT>';

# H1タグ用
our $HtagHeader_ = '<FONT COLOR="#4444ff">';
our $H_tagHeader = '</FONT>';

# 大きい文字
our $HtagBig_ = '<FONT SIZE=6>';
our $H_tagBig = '</FONT>';

# 島の名前など
our $HtagName_ = '<FONT COLOR="#a06040"><B>';
our $H_tagName = '</B></FONT>';

# 薄くなった島の名前
our $HtagName2_ = '<FONT COLOR="#808080"><B>';
our $H_tagName2 = '</B></FONT>';

# 順位の番号など
our $HtagNumber_ = '<FONT COLOR="#800000"><B>';
our $H_tagNumber = '</B></FONT>';

# 順位表における見だし
our $HtagTH_ = '<FONT COLOR="#C00000"><B>';
our $H_tagTH = '</B></FONT>';

# 開発計画の名前
our $HtagComName_ = '<FONT COLOR="#d08000"><B>';
our $H_tagComName = '</B></FONT>';

# 災害
our $HtagDisaster_ = '<FONT COLOR="#ff0000"><B>';
our $H_tagDisaster = '</B></FONT>';

# ローカル掲示板、観光者の書いた文字
our $HtagLbbsSS_ = '<FONT COLOR="#0000ff"><B>';
our $H_tagLbbsSS = '</B></FONT>';

# ローカル掲示板、島主の書いた文字
our $HtagLbbsOW_ = '<FONT COLOR="#ff0000"><B>';
our $H_tagLbbsOW = '</B></FONT>';

# 通常の文字色(これだけでなく、BODYタグのオプションもちゃんと変更すべし
our $HnormalColor = '#000000';

# 順位表、セルの属性
our $HbgTitleCell   = 'BGCOLOR="#ccffcc"';    # 順位表見出し
our $HbgNumberCell  = 'BGCOLOR="#ccffcc"';    # 順位表順位
our $HbgNameCell    = 'BGCOLOR="#ccffff"';    # 順位表島の名前
our $HbgInfoCell    = 'BGCOLOR="#ccffff"';    # 順位表島の情報
our $HbgCommentCell = 'BGCOLOR="#ccffcc"';    # 順位表コメント欄
our $HbgInputCell   = 'BGCOLOR="#ccffcc"';    # 開発計画フォーム
our $HbgMapCell     = 'BGCOLOR="#ccffcc"';    # 開発計画地図
our $HbgCommandCell = 'BGCOLOR="#ccffcc"';    # 開発計画入力済み計画

#----------------------------------------------------------------------
# 好みによって設定する部分は以上
#----------------------------------------------------------------------

#----------------------------------------------------------------------
# これ以降のスクリプトは、変更されることを想定していませんが、
# いじってもかまいません。
# コマンドの名前、値段などは解りやすいと思います。
#----------------------------------------------------------------------

#----------------------------------------------------------------------
# 各種定数
#----------------------------------------------------------------------
# このファイル
our $HthisFile = "$baseDir/hako-main.cgi";

# 地形番号
our $HlandSea      = 0;     # 海
our $HlandWaste    = 1;     # 荒地
our $HlandPlains   = 2;     # 平地
our $HlandTown     = 3;     # 町系
our $HlandForest   = 4;     # 森
our $HlandFarm     = 5;     # 農場
our $HlandFactory  = 6;     # 工場
our $HlandBase     = 7;     # ミサイル基地
our $HlandDefence  = 8;     # 防衛施設
our $HlandMountain = 9;     # 山
our $HlandMonster  = 10;    # 怪獣
our $HlandSbase    = 11;    # 海底基地
our $HlandOil      = 12;    # 海底油田
our $HlandMonument = 13;    # 記念碑
our $HlandHaribote = 14;    # ハリボテ

# コマンド
our $HcommandTotal = 28;    # コマンドの種類

# 計画番号の設定
# 整地系
our $HcomPrepare  = 01;    # 整地
our $HcomPrepare2 = 02;    # 地ならし
our $HcomReclaim  = 03;    # 埋め立て
our $HcomDestroy  = 04;    # 掘削
our $HcomSellTree = 05;    # 伐採

# 作る系
our $HcomPlant    = 11;    # 植林
our $HcomFarm     = 12;    # 農場整備
our $HcomFactory  = 13;    # 工場建設
our $HcomMountain = 14;    # 採掘場整備
our $HcomBase     = 15;    # ミサイル基地建設
our $HcomDbase    = 16;    # 防衛施設建設
our $HcomSbase    = 17;    # 海底基地建設
our $HcomMonument = 18;    # 記念碑建造
our $HcomHaribote = 19;    # ハリボテ設置

# 発射系
our $HcomMissileNM   = 31;    # ミサイル発射
our $HcomMissilePP   = 32;    # PPミサイル発射
our $HcomMissileST   = 33;    # STミサイル発射
our $HcomMissileLD   = 34;    # 陸地破壊弾発射
our $HcomSendMonster = 35;    # 怪獣派遣

# 運営系
our $HcomDoNothing  = 41;     # 資金繰り
our $HcomSell       = 42;     # 食料輸出
our $HcomMoney      = 43;     # 資金援助
our $HcomFood       = 44;     # 食料援助
our $HcomPropaganda = 45;     # 誘致活動
our $HcomGiveup     = 46;     # 島の放棄

# 自動入力系
our $HcomAutoPrepare  = 61;    # フル整地
our $HcomAutoPrepare2 = 62;    # フル地ならし
our $HcomAutoDelete   = 63;    # 全コマンド消去

# 順番
our @HcomList = (
    $HcomPrepare,   $HcomSell,        $HcomPrepare2,
    $HcomReclaim,   $HcomDestroy,     $HcomSellTree,
    $HcomPlant,     $HcomFarm,        $HcomFactory,
    $HcomMountain,  $HcomBase,        $HcomDbase,
    $HcomSbase,     $HcomMonument,    $HcomHaribote,
    $HcomMissileNM, $HcomMissilePP,   $HcomMissileST,
    $HcomMissileLD, $HcomSendMonster, $HcomDoNothing,
    $HcomMoney,     $HcomFood,        $HcomPropaganda,
    $HcomGiveup,    $HcomAutoPrepare, $HcomAutoPrepare2,
    $HcomAutoDelete
);

# 計画の名前と値段
our @HcomName;
our @HcomCost;
$HcomName[$HcomPrepare]      = '整地';
$HcomCost[$HcomPrepare]      = 5;
$HcomName[$HcomPrepare2]     = '地ならし';
$HcomCost[$HcomPrepare2]     = 100;
$HcomName[$HcomReclaim]      = '埋め立て';
$HcomCost[$HcomReclaim]      = 150;
$HcomName[$HcomDestroy]      = '掘削';
$HcomCost[$HcomDestroy]      = 200;
$HcomName[$HcomSellTree]     = '伐採';
$HcomCost[$HcomSellTree]     = 0;
$HcomName[$HcomPlant]        = '植林';
$HcomCost[$HcomPlant]        = 50;
$HcomName[$HcomFarm]         = '農場整備';
$HcomCost[$HcomFarm]         = 20;
$HcomName[$HcomFactory]      = '工場建設';
$HcomCost[$HcomFactory]      = 100;
$HcomName[$HcomMountain]     = '採掘場整備';
$HcomCost[$HcomMountain]     = 300;
$HcomName[$HcomBase]         = 'ミサイル基地建設';
$HcomCost[$HcomBase]         = 300;
$HcomName[$HcomDbase]        = '防衛施設建設';
$HcomCost[$HcomDbase]        = 800;
$HcomName[$HcomSbase]        = '海底基地建設';
$HcomCost[$HcomSbase]        = 8000;
$HcomName[$HcomMonument]     = '記念碑建造';
$HcomCost[$HcomMonument]     = 9999;
$HcomName[$HcomHaribote]     = 'ハリボテ設置';
$HcomCost[$HcomHaribote]     = 1;
$HcomName[$HcomMissileNM]    = 'ミサイル発射';
$HcomCost[$HcomMissileNM]    = 20;
$HcomName[$HcomMissilePP]    = 'PPミサイル発射';
$HcomCost[$HcomMissilePP]    = 50;
$HcomName[$HcomMissileST]    = 'STミサイル発射';
$HcomCost[$HcomMissileST]    = 50;
$HcomName[$HcomMissileLD]    = '陸地破壊弾発射';
$HcomCost[$HcomMissileLD]    = 100;
$HcomName[$HcomSendMonster]  = '怪獣派遣';
$HcomCost[$HcomSendMonster]  = 3000;
$HcomName[$HcomDoNothing]    = '資金繰り';
$HcomCost[$HcomDoNothing]    = 0;
$HcomName[$HcomSell]         = '食料輸出';
$HcomCost[$HcomSell]         = -100;
$HcomName[$HcomMoney]        = '資金援助';
$HcomCost[$HcomMoney]        = 100;
$HcomName[$HcomFood]         = '食料援助';
$HcomCost[$HcomFood]         = -100;
$HcomName[$HcomPropaganda]   = '誘致活動';
$HcomCost[$HcomPropaganda]   = 1000;
$HcomName[$HcomGiveup]       = '島の放棄';
$HcomCost[$HcomGiveup]       = 0;
$HcomName[$HcomAutoPrepare]  = '整地自動入力';
$HcomCost[$HcomAutoPrepare]  = 0;
$HcomName[$HcomAutoPrepare2] = '地ならし自動入力';
$HcomCost[$HcomAutoPrepare2] = 0;
$HcomName[$HcomAutoDelete]   = '全計画を白紙撤回';
$HcomCost[$HcomAutoDelete]   = 0;

# 島の座標数
our $HpointNumber = $HislandSize * $HislandSize;

# 「戻る」リンク
our $HtempBack
    = "<A HREF=\"$HthisFile\">${HtagBig_}トップへ戻る${H_tagBig}</A>";

1;


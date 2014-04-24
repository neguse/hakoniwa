#!/usr/bin/perl
# ↑はサーバーに合わせて変更して下さい。
# perl5用です。

#----------------------------------------------------------------------
# 箱庭諸島 ver2.30
# メインスクリプト(ver1.02)
# 使用条件、使用方法等は、hako-readme.txtファイルを参照
#
# 箱庭諸島のページ: http://www.bekkoame.ne.jp/~tokuoka/hakoniwa.html
#----------------------------------------------------------------------

#----------------------------------------------------------------------
# 各種設定値
# (これ以降の部分の各設定値を、適切な値に変更してください)
#----------------------------------------------------------------------

#----------------------------------------------------------------------
# 以下、必ず設定する部分
#----------------------------------------------------------------------

# このファイルを置くディレクトリ
# my($baseDir) = 'http://サーバー/ディレクトリ';
#
# 例)
# http://cgi2.bekkoame.ne.jp/cgi-bin/user/u5534/hakoniwa/hako-main.cgi
# として置く場合、
# my($baseDir) = 'http://cgi2.bekkoame.ne.jp/cgi-bin/user/u5534/hakoniwa';
# とする。最後にスラッシュ(/)は付けない。

# my($baseDir) = 'http://サーバー/ディレクトリ';
my ($baseDir) = 'http://localhost:5000/cgi-bin';

# 画像ファイルを置くディレクトリ
# my($imageDir) = 'http://サーバー/ディレクトリ';
my ($imageDir) = 'http://localhost:5000/images';

# マスターパスワード
# このパスワードは、すべての島のパスワードを代用できます。
# 例えば、「他の島のパスワード変更」等もできます。
my ($masterPassword) = 'yourpassword';

# 特殊パスワード
# このパスワードで「名前変更」を行うと、その島の資金、食料が最大値になります。
# (実際に名前を変える必要はありません。)
$HspecialPassword = 'yourspecialpassword';

# 管理者名
my ($adminName) = '管理者の名前';

# 管理者のメールアドレス
my ($email) = '管理者@どこか.どこか.どこか';

# 掲示板アドレス
my ($bbs) = 'http://サーバー/掲示板.cgi';

# ホームページのアドレス
my ($toppage) = 'http://サーバー/ホームページ.html';

# ディレクトリのパーミッション
# 通常は0755でよいが、0777、0705、0704等でないとできないサーバーもあるらしい
$HdirMode = 0755;

# データディレクトリの名前
# ここで設定した名前のディレクトリ以下にデータが格納されます。
# デフォルトでは'data'となっていますが、セキュリティのため
# なるべく違う名前に変更してください。
$HdirName = 'data';

# データの書き込み方

# ロックの方式
# 1 ディレクトリ
# 2 システムコール(可能ならば最も望ましい)
# 3 シンボリックリンク
# 4 通常ファイル(あまりお勧めでない)
my ($lockMode) = 2;

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
$HunitTime = 21600;    # 6時間

# 異常終了基準時間
# (ロック後何秒で、強制解除するか)
my ($unlockTime) = 120;

# 島の最大数
$HmaxIsland = 30;

# トップページに表示するログのターン数
$HtopLogTurn = 1;

# ログファイル保持ターン数
$HlogMax = 8;

# バックアップを何ターンおきに取るか
$HbackupTurn = 12;

# バックアップを何回分残すか
$HbackupTimes = 4;

# 発見ログ保持行数
$HhistoryMax = 10;

# 放棄コマンド自動入力ターン数
$HgiveupTurn = 28;

# コマンド入力限界数
# (ゲームが始まってから変更すると、データファイルの互換性が無くなります。)
$HcommandMax = 20;

# ローカル掲示板行数を使用するかどうか(0:使用しない、1:使用する)
$HuseLbbs = 0;

# ローカル掲示板行数
$HlbbsMax = 10;

# 島の大きさ
# (変更できないかも)
$HislandSize = 12;

# 他人から資金を見えなくするか
# 0 見えない
# 1 見える
# 2 100の位で四捨五入
$HhideMoneyMode = 2;

# パスワードの暗号化(0だと暗号化しない、1だと暗号化する)
my ($cryptOn) = 1;

# デバッグモード(1だと、「ターンを進める」ボタンが使用できる)
$Hdebug = 0;

#----------------------------------------
# 資金、食料などの設定値と単位
#----------------------------------------
# 初期資金
$HinitialMoney = 100;

# 初期食料
$HinitialFood = 100;

# お金の単位
$HunitMoney = '億円';

# 食料の単位
$HunitFood = '00トン';

# 人口の単位
$HunitPop = '00人';

# 広さの単位
$HunitArea = '00万坪';

# 木の数の単位
$HunitTree = '00本';

# 木の単位当たりの売値
$HtreeValue = 5;

# 名前変更のコスト
$HcostChangeName = 500;

# 人口1単位あたりの食料消費料
$HeatenFood = 0.2;

#----------------------------------------
# 基地の経験値
#----------------------------------------
# 経験値の最大値
$HmaxExpPoint = 200;    # ただし、最大でも255まで

# レベルの最大値
my ($maxBaseLevel)  = 5;    # ミサイル基地
my ($maxSBaseLevel) = 3;    # 海底基地

# 経験値がいくつでレベルアップか
my ( @baseLevelUp, @sBaseLevelUp );
@baseLevelUp = ( 20, 60, 120, 200 );    # ミサイル基地
@sBaseLevelUp = ( 50, 200 );            # 海底基地

#----------------------------------------
# 防衛施設の自爆
#----------------------------------------
# 怪獣に踏まれた時自爆するなら1、しないなら0
$HdBaseAuto = 1;

#----------------------------------------
# 災害
#----------------------------------------
# 通常災害発生率(確率は0.1%単位)
$HdisEarthquake = 5;     # 地震
$HdisTsunami    = 15;    # 津波
$HdisTyphoon    = 20;    # 台風
$HdisMeteo      = 15;    # 隕石
$HdisHugeMeteo  = 5;     # 巨大隕石
$HdisEruption   = 10;    # 噴火
$HdisFire       = 10;    # 火災
$HdisMaizo      = 10;    # 埋蔵金

# 地盤沈下
$HdisFallBorder = 90;    # 安全限界の広さ(Hex数)
$HdisFalldown   = 30;    # その広さを超えた場合の確率

# 怪獣
$HdisMonsBorder1 = 1000;    # 人口基準1(怪獣レベル1)
$HdisMonsBorder2 = 2500;    # 人口基準2(怪獣レベル2)
$HdisMonsBorder3 = 4000;    # 人口基準3(怪獣レベル3)
$HdisMonster     = 3;       # 単位面積あたりの出現率(0.01%単位)

# 種類
$HmonsterNumber = 8;

# 各基準において出てくる怪獣の番号の最大値
$HmonsterLevel1 = 2;        # サンジラまで
$HmonsterLevel2 = 5;        # いのらゴーストまで
$HmonsterLevel3 = 7;        # キングいのらまで(全部)

# 名前
@HmonsterName = (
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
@HmonsterBHP     = ( 2, 1,   1,   3,    2,   1,   4,    5 );
@HmonsterDHP     = ( 0, 2,   2,   2,    2,   0,   2,    2 );
@HmonsterSpecial = ( 0, 0,   3,   0,    1,   2,   4,    0 );
@HmonsterExp     = ( 5, 5,   7,   12,   15,  10,  20,   30 );
@HmonsterValue   = ( 0, 400, 500, 1000, 800, 300, 1500, 2000 );

# 特殊能力の内容は、
# 0 特になし
# 1 足が速い(最大2歩あるく)
# 2 足がとても速い(最大何歩あるくか不明)
# 3 奇数ターンは硬化
# 4 偶数ターンは硬化

# 画像ファイル
@HmonsterImage = (
    'monster7.gif', 'monster0.gif', 'monster5.gif', 'monster1.gif',
    'monster2.gif', 'monster8.gif', 'monster6.gif', 'monster3.gif'
);

# 画像ファイルその2(硬化中)
@HmonsterImage2 = ( '', '', 'monster4.gif', '', '', '', 'monster4.gif', '' );

#----------------------------------------
# 油田
#----------------------------------------
# 油田の収入
$HoilMoney = 1000;

# 油田の枯渇確率
$HoilRatio = 40;

#----------------------------------------
# 記念碑
#----------------------------------------
# 何種類あるか
$HmonumentNumber = 3;

# 名前
@HmonumentName = ( 'モノリス', '平和記念碑', '戦いの碑' );

# 画像ファイル
@HmonumentImage = ( 'monument0.gif', 'monument0.gif', 'monument0.gif' );

#----------------------------------------
# 賞関係
#----------------------------------------
# ターン杯を何ターン毎に出すか
$HturnPrizeUnit = 100;

# 賞の名前
$Hprize[0] = 'ターン杯';
$Hprize[1] = '繁栄賞';
$Hprize[2] = '超繁栄賞';
$Hprize[3] = '究極繁栄賞';
$Hprize[4] = '平和賞';
$Hprize[5] = '超平和賞';
$Hprize[6] = '究極平和賞';
$Hprize[7] = '災難賞';
$Hprize[8] = '超災難賞';
$Hprize[9] = '究極災難賞';

#----------------------------------------
# 外見関係
#----------------------------------------
# <BODY>タグのオプション
my ($htmlBody) = 'BGCOLOR="#EEFFFF"';

# ゲームのタイトル文字
$Htitle = '箱庭諸島２';

# タグ
# タイトル文字
$HtagTitle_ = '<FONT SIZE=7 COLOR="#8888ff">';
$H_tagTitle = '</FONT>';

# H1タグ用
$HtagHeader_ = '<FONT COLOR="#4444ff">';
$H_tagHeader = '</FONT>';

# 大きい文字
$HtagBig_ = '<FONT SIZE=6>';
$H_tagBig = '</FONT>';

# 島の名前など
$HtagName_ = '<FONT COLOR="#a06040"><B>';
$H_tagName = '</B></FONT>';

# 薄くなった島の名前
$HtagName2_ = '<FONT COLOR="#808080"><B>';
$H_tagName2 = '</B></FONT>';

# 順位の番号など
$HtagNumber_ = '<FONT COLOR="#800000"><B>';
$H_tagNumber = '</B></FONT>';

# 順位表における見だし
$HtagTH_ = '<FONT COLOR="#C00000"><B>';
$H_tagTH = '</B></FONT>';

# 開発計画の名前
$HtagComName_ = '<FONT COLOR="#d08000"><B>';
$H_tagComName = '</B></FONT>';

# 災害
$HtagDisaster_ = '<FONT COLOR="#ff0000"><B>';
$H_tagDisaster = '</B></FONT>';

# ローカル掲示板、観光者の書いた文字
$HtagLbbsSS_ = '<FONT COLOR="#0000ff"><B>';
$H_tagLbbsSS = '</B></FONT>';

# ローカル掲示板、島主の書いた文字
$HtagLbbsOW_ = '<FONT COLOR="#ff0000"><B>';
$H_tagLbbsOW = '</B></FONT>';

# 通常の文字色(これだけでなく、BODYタグのオプションもちゃんと変更すべし
$HnormalColor = '#000000';

# 順位表、セルの属性
$HbgTitleCell   = 'BGCOLOR="#ccffcc"';    # 順位表見出し
$HbgNumberCell  = 'BGCOLOR="#ccffcc"';    # 順位表順位
$HbgNameCell    = 'BGCOLOR="#ccffff"';    # 順位表島の名前
$HbgInfoCell    = 'BGCOLOR="#ccffff"';    # 順位表島の情報
$HbgCommentCell = 'BGCOLOR="#ccffcc"';    # 順位表コメント欄
$HbgInputCell   = 'BGCOLOR="#ccffcc"';    # 開発計画フォーム
$HbgMapCell     = 'BGCOLOR="#ccffcc"';    # 開発計画地図
$HbgCommandCell = 'BGCOLOR="#ccffcc"';    # 開発計画入力済み計画

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
$HthisFile = "$baseDir/hako-main.cgi";

# 地形番号
$HlandSea      = 0;     # 海
$HlandWaste    = 1;     # 荒地
$HlandPlains   = 2;     # 平地
$HlandTown     = 3;     # 町系
$HlandForest   = 4;     # 森
$HlandFarm     = 5;     # 農場
$HlandFactory  = 6;     # 工場
$HlandBase     = 7;     # ミサイル基地
$HlandDefence  = 8;     # 防衛施設
$HlandMountain = 9;     # 山
$HlandMonster  = 10;    # 怪獣
$HlandSbase    = 11;    # 海底基地
$HlandOil      = 12;    # 海底油田
$HlandMonument = 13;    # 記念碑
$HlandHaribote = 14;    # ハリボテ

# コマンド
$HcommandTotal = 28;    # コマンドの種類

# 計画番号の設定
# 整地系
$HcomPrepare  = 01;    # 整地
$HcomPrepare2 = 02;    # 地ならし
$HcomReclaim  = 03;    # 埋め立て
$HcomDestroy  = 04;    # 掘削
$HcomSellTree = 05;    # 伐採

# 作る系
$HcomPlant    = 11;    # 植林
$HcomFarm     = 12;    # 農場整備
$HcomFactory  = 13;    # 工場建設
$HcomMountain = 14;    # 採掘場整備
$HcomBase     = 15;    # ミサイル基地建設
$HcomDbase    = 16;    # 防衛施設建設
$HcomSbase    = 17;    # 海底基地建設
$HcomMonument = 18;    # 記念碑建造
$HcomHaribote = 19;    # ハリボテ設置

# 発射系
$HcomMissileNM   = 31;    # ミサイル発射
$HcomMissilePP   = 32;    # PPミサイル発射
$HcomMissileST   = 33;    # STミサイル発射
$HcomMissileLD   = 34;    # 陸地破壊弾発射
$HcomSendMonster = 35;    # 怪獣派遣

# 運営系
$HcomDoNothing  = 41;     # 資金繰り
$HcomSell       = 42;     # 食料輸出
$HcomMoney      = 43;     # 資金援助
$HcomFood       = 44;     # 食料援助
$HcomPropaganda = 45;     # 誘致活動
$HcomGiveup     = 46;     # 島の放棄

# 自動入力系
$HcomAutoPrepare  = 61;    # フル整地
$HcomAutoPrepare2 = 62;    # フル地ならし
$HcomAutoDelete   = 63;    # 全コマンド消去

# 順番
@HcomList = (
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

#----------------------------------------------------------------------
# 変数
#----------------------------------------------------------------------

# COOKIE
my ($defaultID);        # 島の名前
my ($defaultTarget);    # ターゲットの名前

# 島の座標数
$HpointNumber = $HislandSize * $HislandSize;

#----------------------------------------------------------------------
# メイン
#----------------------------------------------------------------------

# 「戻る」リンク
$HtempBack
    = "<A HREF=\"$HthisFile\">${HtagBig_}トップへ戻る${H_tagBig}</A>";

# ロックをかける
if ( !hakolock() ) {

    # ロック失敗
    # ヘッダ出力
    tempHeader();

    # ロック失敗メッセージ
    tempLockFail();

    # フッタ出力
    tempFooter();

    # 終了
    exit(0);
}

# 乱数の初期化
srand( time ^ $$ );

# COOKIE読みこみ
cookieInput();

# CGI読みこみ
cgiInput();

# 島データの読みこみ
if ( readIslandsFile($HcurrentID) == 0 ) {
    unlock();
    tempHeader();
    tempNoDataFile();
    tempFooter();
    exit(0);
}

# テンプレートを初期化
tempInitialize();

# COOKIE出力
cookieOutput();

# ヘッダ出力
tempHeader();

if ( $HmainMode eq 'turn' ) {

    # ターン進行
    require('hako-turn.cgi');
    require('hako-top.cgi');
    turnMain();

}
elsif ( $HmainMode eq 'new' ) {

    # 島の新規作成
    require('hako-turn.cgi');
    require('hako-map.cgi');
    newIslandMain();

}
elsif ( $HmainMode eq 'print' ) {

    # 観光モード
    require('hako-map.cgi');
    printIslandMain();

}
elsif ( $HmainMode eq 'owner' ) {

    # 開発モード
    require('hako-map.cgi');
    ownerMain();

}
elsif ( $HmainMode eq 'command' ) {

    # コマンド入力モード
    require('hako-map.cgi');
    commandMain();

}
elsif ( $HmainMode eq 'comment' ) {

    # コメント入力モード
    require('hako-map.cgi');
    commentMain();

}
elsif ( $HmainMode eq 'lbbs' ) {

    # ローカル掲示板モード
    require('hako-map.cgi');
    localBbsMain();

}
elsif ( $HmainMode eq 'change' ) {

    # 情報変更モード
    require('hako-turn.cgi');
    require('hako-top.cgi');
    changeMain();

}
else {
    # その他の場合はトップページモード
    require('hako-top.cgi');
    topPageMain();
}

# フッタ出力
tempFooter();

# 終了
exit(0);

# コマンドを前にずらす
sub slideFront {
    my ( $command, $number ) = @_;
    my ($i);

    # それぞれずらす
    splice( @$command, $number, 1 );

    # 最後に資金繰り
    $command->[ $HcommandMax - 1 ] = {
        'kind'   => $HcomDoNothing,
        'target' => 0,
        'x'      => 0,
        'y'      => 0,
        'arg'    => 0
    };
}

# コマンドを後にずらす
sub slideBack {
    my ( $command, $number ) = @_;
    my ($i);

    # それぞれずらす
    return if $number == $#$command;
    pop(@$command);
    splice( @$command, $number, 0, $command->[$number] );
}

#----------------------------------------------------------------------
# 島データ入出力
#----------------------------------------------------------------------

# 全島データ読みこみ
sub readIslandsFile {
    my ($num) = @_;    # 0だと地形読みこまず
                       # -1だと全地形を読む
                       # 番号だとその島の地形だけは読みこむ

    # データファイルを開く
    if ( !open( IN, "${HdirName}/hakojima.dat" ) ) {
        rename( "${HdirName}/hakojima.tmp", "${HdirName}/hakojima.dat" );
        if ( !open( IN, "${HdirName}/hakojima.dat" ) ) {
            return 0;
        }
    }

    # 各パラメータの読みこみ
    $HislandTurn = int(<IN>);    # ターン数
    if ( $HislandTurn == 0 ) {
        return 0;
    }
    $HislandLastTime = int(<IN>);    # 最終更新時間
    if ( $HislandLastTime == 0 ) {
        return 0;
    }
    $HislandNumber = int(<IN>);      # 島の総数
    $HislandNextID = int(<IN>);      # 次に割り当てるID

    # ターン処理判定
    my ($now) = time;
    if (   ( ( $Hdebug == 1 ) && ( $HmainMode eq 'Hdebugturn' ) )
        || ( ( $now - $HislandLastTime ) >= $HunitTime ) )
    {
        $HmainMode = 'turn';
        $num       = -1;             # 全島読みこむ
    }

    # 島の読みこみ
    my ($i);
    for ( $i = 0; $i < $HislandNumber; $i++ ) {
        $Hislands[$i] = readIsland($num);
        $HidToNumber{ $Hislands[$i]->{'id'} } = $i;
    }

    # ファイルを閉じる
    close(IN);

    return 1;
}

# 島ひとつ読みこみ
sub readIsland {
    my ($num) = @_;
    my ($name,     $id,      $prize,    $absent, $comment,
        $password, $money,   $food,     $pop,    $area,
        $farm,     $factory, $mountain, $score
    );
    $name = <IN>;    # 島の名前
    chomp($name);
    if ( $name =~ s/,(.*)$//g ) {
        $score = int($1);
    }
    else {
        $score = 0;
    }
    $id    = int(<IN>);    # ID番号
    $prize = <IN>;         # 受賞
    chomp($prize);
    $absent  = int(<IN>);    # 連続資金繰り数
    $comment = <IN>;         # コメント
    chomp($comment);
    $password = <IN>;        # 暗号化パスワード
    chomp($password);
    $money    = int(<IN>);    # 資金
    $food     = int(<IN>);    # 食料
    $pop      = int(<IN>);    # 人口
    $area     = int(<IN>);    # 広さ
    $farm     = int(<IN>);    # 農場
    $factory  = int(<IN>);    # 工場
    $mountain = int(<IN>);    # 採掘場

    # HidToNameテーブルへ保存
    $HidToName{$id} = $name;    #

    # 地形
    my ( @land, @landValue, $line, @command, @lbbs );

    if ( ( $num == -1 ) || ( $num == $id ) ) {
        if ( !open( IIN, "${HdirName}/island.$id" ) ) {
            rename( "${HdirName}/islandtmp.$id", "${HdirName}/island.$id" );
            if ( !open( IIN, "${HdirName}/island.$id" ) ) {
                exit(0);
            }
        }
        my ( $x, $y );
        for ( $y = 0; $y < $HislandSize; $y++ ) {
            $line = <IIN>;
            for ( $x = 0; $x < $HislandSize; $x++ ) {
                $line =~ s/^(.)(..)//;
                $land[$x][$y]      = hex($1);
                $landValue[$x][$y] = hex($2);
            }
        }

        # コマンド
        my ($i);
        for ( $i = 0; $i < $HcommandMax; $i++ ) {
            $line = <IIN>;
            $line =~ /^([0-9]*),([0-9]*),([0-9]*),([0-9]*),([0-9]*)$/;
            $command[$i] = {
                'kind'   => int($1),
                'target' => int($2),
                'x'      => int($3),
                'y'      => int($4),
                'arg'    => int($5)
            };
        }

        # ローカル掲示板
        for ( $i = 0; $i < $HlbbsMax; $i++ ) {
            $line = <IIN>;
            chomp($line);
            $lbbs[$i] = $line;
        }

        close(IIN);
    }

    # 島型にして返す
    return {
        'name'      => $name,
        'id'        => $id,
        'score'     => $score,
        'prize'     => $prize,
        'absent'    => $absent,
        'comment'   => $comment,
        'password'  => $password,
        'money'     => $money,
        'food'      => $food,
        'pop'       => $pop,
        'area'      => $area,
        'farm'      => $farm,
        'factory'   => $factory,
        'mountain'  => $mountain,
        'land'      => \@land,
        'landValue' => \@landValue,
        'command'   => \@command,
        'lbbs'      => \@lbbs,
    };
}

# 全島データ書き込み
sub writeIslandsFile {
    my ($num) = @_;

    # ファイルを開く
    open( OUT, ">${HdirName}/hakojima.tmp" );

    # 各パラメータ書き込み
    print OUT "$HislandTurn\n";
    print OUT "$HislandLastTime\n";
    print OUT "$HislandNumber\n";
    print OUT "$HislandNextID\n";

    # 島の書きこみ
    my ($i);
    for ( $i = 0; $i < $HislandNumber; $i++ ) {
        writeIsland( $Hislands[$i], $num );
    }

    # ファイルを閉じる
    close(OUT);

    # 本来の名前にする
    unlink("${HdirName}/hakojima.dat");
    rename( "${HdirName}/hakojima.tmp", "${HdirName}/hakojima.dat" );
}

# 島ひとつ書き込み
sub writeIsland {
    my ( $island, $num ) = @_;
    my ($score);
    $score = int( $island->{'score'} );
    print OUT $island->{'name'} . ",$score\n";
    print OUT $island->{'id'} . "\n";
    print OUT $island->{'prize'} . "\n";
    print OUT $island->{'absent'} . "\n";
    print OUT $island->{'comment'} . "\n";
    print OUT $island->{'password'} . "\n";
    print OUT $island->{'money'} . "\n";
    print OUT $island->{'food'} . "\n";
    print OUT $island->{'pop'} . "\n";
    print OUT $island->{'area'} . "\n";
    print OUT $island->{'farm'} . "\n";
    print OUT $island->{'factory'} . "\n";
    print OUT $island->{'mountain'} . "\n";

    # 地形
    if ( ( $num <= -1 ) || ( $num == $island->{'id'} ) ) {
        open( IOUT, ">${HdirName}/islandtmp.$island->{'id'}" );

        my ( $land, $landValue );
        $land      = $island->{'land'};
        $landValue = $island->{'landValue'};
        my ( $x, $y );
        for ( $y = 0; $y < $HislandSize; $y++ ) {
            for ( $x = 0; $x < $HislandSize; $x++ ) {
                printf IOUT ( "%x%02x", $land->[$x][$y],
                    $landValue->[$x][$y] );
            }
            print IOUT "\n";
        }

        # コマンド
        my ( $command, $cur, $i );
        $command = $island->{'command'};
        for ( $i = 0; $i < $HcommandMax; $i++ ) {
            printf IOUT (
                "%d,%d,%d,%d,%d\n",         $command->[$i]->{'kind'},
                $command->[$i]->{'target'}, $command->[$i]->{'x'},
                $command->[$i]->{'y'},      $command->[$i]->{'arg'}
            );
        }

        # ローカル掲示板
        my ($lbbs);
        $lbbs = $island->{'lbbs'};
        for ( $i = 0; $i < $HlbbsMax; $i++ ) {
            print IOUT $lbbs->[$i] . "\n";
        }

        close(IOUT);
        unlink("${HdirName}/island.$island->{'id'}");
        rename(
            "${HdirName}/islandtmp.$island->{'id'}",
            "${HdirName}/island.$island->{'id'}"
        );
    }
}

#----------------------------------------------------------------------
# 入出力
#----------------------------------------------------------------------

# 標準出力への出力
sub out {
    print STDOUT $_[0];
}

# デバッグログ
sub HdebugOut {
    open( DOUT, ">>debug.log" );
    print DOUT ( $_[0] );
    close(DOUT);
}

# CGIの読みこみ
sub cgiInput {
    my ( $line, $getLine );

    # 入力を受け取って日本語コードをEUCに
    $line = <>;
    $line =~ tr/+/ /;
    $line =~ s/%([a-fA-F0-9][a-fA-F0-9])/pack("C", hex($1))/eg;
    $line = $line;
    $line =~ s/[\x00-\x1f\,]//g;

    # GETのやつも受け取る
    $getLine = $ENV{'QUERY_STRING'};

    # 対象の島
    if ( $line =~ /CommandButton([0-9]+)=/ ) {

        # コマンド送信ボタンの場合
        $HcurrentID = $1;
        $defaultID  = $1;
    }

    if ( $line =~ /ISLANDNAME=([^\&]*)\&/ ) {

        # 名前指定の場合
        $HcurrentName = cutColumn( $1, 32 );
    }

    if ( $line =~ /ISLANDID=([0-9]+)\&/ ) {

        # その他の場合
        $HcurrentID = $1;
        $defaultID  = $1;
    }

    # パスワード
    if ( $line =~ /OLDPASS=([^\&]*)\&/ ) {
        $HoldPassword     = $1;
        $HdefaultPassword = $1;
    }
    if ( $line =~ /PASSWORD=([^\&]*)\&/ ) {
        $HinputPassword   = $1;
        $HdefaultPassword = $1;
    }
    if ( $line =~ /PASSWORD2=([^\&]*)\&/ ) {
        $HinputPassword2 = $1;
    }

    # メッセージ
    if ( $line =~ /MESSAGE=([^\&]*)\&/ ) {
        $Hmessage = cutColumn( $1, 80 );
    }

    # ローカル掲示板
    if ( $line =~ /LBBSNAME=([^\&]*)\&/ ) {
        $HlbbsName    = $1;
        $HdefaultName = $1;
    }
    if ( $line =~ /LBBSMESSAGE=([^\&]*)\&/ ) {
        $HlbbsMessage = cutColumn( $1, 80 );
    }

    # main modeの取得
    if ( $line =~ /TurnButton/ ) {
        if ( $Hdebug == 1 ) {
            $HmainMode = 'Hdebugturn';
        }
    }
    elsif ( $line =~ /OwnerButton/ ) {
        $HmainMode = 'owner';
    }
    elsif ( $getLine =~ /Sight=([0-9]*)/ ) {
        $HmainMode  = 'print';
        $HcurrentID = $1;
    }
    elsif ( $line =~ /NewIslandButton/ ) {
        $HmainMode = 'new';
    }
    elsif ( $line =~ /LbbsButton(..)([0-9]*)/ ) {
        $HmainMode = 'lbbs';
        if ( $1 eq 'SS' ) {

            # 観光者
            $HlbbsMode = 0;
        }
        elsif ( $1 eq 'OW' ) {

            # 島主
            $HlbbsMode = 1;
        }
        else {
            # 削除
            $HlbbsMode = 2;
        }
        $HcurrentID = $2;

        # 削除かもしれないので、番号を取得
        $line =~ /NUMBER=([^\&]*)\&/;
        $HcommandPlanNumber = $1;

    }
    elsif ( $line =~ /ChangeInfoButton/ ) {
        $HmainMode = 'change';
    }
    elsif ( $line =~ /MessageButton([0-9]*)/ ) {
        $HmainMode  = 'comment';
        $HcurrentID = $1;
    }
    elsif ( $line =~ /CommandButton/ ) {
        $HmainMode = 'command';

        # コマンドモードの場合、コマンドの取得
        $line =~ /NUMBER=([^\&]*)\&/;
        $HcommandPlanNumber = $1;
        $line =~ /COMMAND=([^\&]*)\&/;
        $HcommandKind = $1;
        $HdefaultKind = $1;
        $line =~ /AMOUNT=([^\&]*)\&/;
        $HcommandArg = $1;
        $line =~ /TARGETID=([^\&]*)\&/;
        $HcommandTarget = $1;
        $defaultTarget  = $1;
        $line =~ /POINTX=([^\&]*)\&/;
        $HcommandX = $1;
        $HdefaultX = $1;
        $line =~ /POINTY=([^\&]*)\&/;
        $HcommandY = $1;
        $HdefaultY = $1;
        $line =~ /COMMANDMODE=(write|insert|delete)/;
        $HcommandMode = $1;
    }
    else {
        $HmainMode = 'top';
    }

}

#cookie入力
sub cookieInput {
    my ($cookie);

    $cookie = $ENV{'HTTP_COOKIE'};

    if ( $cookie =~ /${HthisFile}OWNISLANDID=\(([^\)]*)\)/ ) {
        $defaultID = $1;
    }
    if ( $cookie =~ /${HthisFile}OWNISLANDPASSWORD=\(([^\)]*)\)/ ) {
        $HdefaultPassword = $1;
    }
    if ( $cookie =~ /${HthisFile}TARGETISLANDID=\(([^\)]*)\)/ ) {
        $defaultTarget = $1;
    }
    if ( $cookie =~ /${HthisFile}LBBSNAME=\(([^\)]*)\)/ ) {
        $HdefaultName = $1;
    }
    if ( $cookie =~ /${HthisFile}POINTX=\(([^\)]*)\)/ ) {
        $HdefaultX = $1;
    }
    if ( $cookie =~ /${HthisFile}POINTY=\(([^\)]*)\)/ ) {
        $HdefaultY = $1;
    }
    if ( $cookie =~ /${HthisFile}KIND=\(([^\)]*)\)/ ) {
        $HdefaultKind = $1;
    }

}

#cookie出力
sub cookieOutput {
    my ( $cookie, $info );

    # 消える期限の設定
    my ( $sec, $min, $hour, $date, $mon, $year, $day, $yday, $dummy )
        = gmtime( time + 30 * 86400 );    # 現在 + 30日

    # 2ケタ化
    $year += 1900;
    if ( $date < 10 ) { $date = "0$date"; }
    if ( $hour < 10 ) { $hour = "0$hour"; }
    if ( $min < 10 )  { $min  = "0$min"; }
    if ( $sec < 10 )  { $sec  = "0$sec"; }

    # 曜日を文字に
    $day = (
        "Sunday",   "Monday", "Tuesday", "Wednesday",
        "Thursday", "Friday", "Saturday"
    )[$day];

    # 月を文字に
    $mon = (
        "Jan", "Feb", "Mar", "Apr", "May", "Jun",
        "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"
    )[$mon];

    # パスと期限のセット
    $info   = "; expires=$day, $date\-$mon\-$year $hour:$min:$sec GMT\n";
    $cookie = '';

    if ( ($HcurrentID) && ( $HmainMode eq 'owner' ) ) {
        $cookie .= "Set-Cookie: ${HthisFile}OWNISLANDID=($HcurrentID) $info";
    }
    if ($HinputPassword) {
        $cookie
            .= "Set-Cookie: ${HthisFile}OWNISLANDPASSWORD=($HinputPassword) $info";
    }
    if ($HcommandTarget) {
        $cookie
            .= "Set-Cookie: ${HthisFile}TARGETISLANDID=($HcommandTarget) $info";
    }
    if ($HlbbsName) {
        $cookie .= "Set-Cookie: ${HthisFile}LBBSNAME=($HlbbsName) $info";
    }
    if ($HcommandX) {
        $cookie .= "Set-Cookie: ${HthisFile}POINTX=($HcommandX) $info";
    }
    if ($HcommandY) {
        $cookie .= "Set-Cookie: ${HthisFile}POINTY=($HcommandY) $info";
    }
    if ($HcommandKind) {

        # 自動系以外
        $cookie .= "Set-Cookie: ${HthisFile}KIND=($HcommandKind) $info";
    }

    out($cookie);
}

#----------------------------------------------------------------------
# ユーティリティ
#----------------------------------------------------------------------
sub hakolock {
    if ( $lockMode == 1 ) {

        # directory式ロック
        return hakolock1();

    }
    elsif ( $lockMode == 2 ) {

        # flock式ロック
        return hakolock2();
    }
    elsif ( $lockMode == 3 ) {

        # symlink式ロック
        return hakolock3();
    }
    else {
        # 通常ファイル式ロック
        return hakolock4();
    }
}

sub hakolock1 {

    # ロックを試す
    if ( mkdir( 'hakojimalock', $HdirMode ) ) {

        # 成功
        return 1;
    }
    else {
        # 失敗
        my ($b) = ( stat('hakojimalock') )[9];
        if ( ( $b > 0 ) && ( ( time() - $b ) > $unlockTime ) ) {

            # 強制解除
            unlock();

            # ヘッダ出力
            tempHeader();

            # 強制解除メッセージ
            tempUnlock();

            # フッタ出力
            tempFooter();

            # 終了
            exit(0);
        }
        return 0;
    }
}

sub hakolock2 {
    open( LOCKID, '>>hakojimalockflock' );
    if ( flock( LOCKID, 2 ) ) {

        # 成功
        return 1;
    }
    else {
        # 失敗
        return 0;
    }
}

sub hakolock3 {

    # ロックを試す
    if ( symlink( 'hakojimalockdummy', 'hakojimalock' ) ) {

        # 成功
        return 1;
    }
    else {
        # 失敗
        my ($b) = ( lstat('hakojimalock') )[9];
        if ( ( $b > 0 ) && ( ( time() - $b ) > $unlockTime ) ) {

            # 強制解除
            unlock();

            # ヘッダ出力
            tempHeader();

            # 強制解除メッセージ
            tempUnlock();

            # フッタ出力
            tempFooter();

            # 終了
            exit(0);
        }
        return 0;
    }
}

sub hakolock4 {

    # ロックを試す
    if ( unlink('key-free') ) {

        # 成功
        open( OUT, '>key-locked' );
        print OUT time;
        close(OUT);
        return 1;
    }
    else {
        # ロック時間チェック
        if ( !open( IN, 'key-locked' ) ) {
            return 0;
        }

        my ($t);
        $t = <IN>;
        close(IN);
        if ( ( $t != 0 ) && ( ( $t + $unlockTime ) < time ) ) {

            # 120秒以上経過してたら、強制的にロックを外す
            unlock();

            # ヘッダ出力
            tempHeader();

            # 強制解除メッセージ
            tempUnlock();

            # フッタ出力
            tempFooter();

            # 終了
            exit(0);
        }
        return 0;
    }
}

# ロックを外す
sub unlock {
    if ( $lockMode == 1 ) {

        # directory式ロック
        rmdir('hakojimalock');

    }
    elsif ( $lockMode == 2 ) {

        # flock式ロック
        close(LOCKID);

    }
    elsif ( $lockMode == 3 ) {

        # symlink式ロック
        unlink('hakojimalock');
    }
    else {
        # 通常ファイル式ロック
        my ($i);
        $i = rename( 'key-locked', 'key-free' );
    }
}

# 小さい方を返す
sub min {
    return ( $_[0] < $_[1] ) ? $_[0] : $_[1];
}

# パスワードエンコード
sub encode {
    if ( $cryptOn == 1 ) {
        return crypt( $_[0], 'h2' );
    }
    else {
        return $_[0];
    }
}

# パスワードチェック
sub checkPassword {
    my ( $p1, $p2 ) = @_;

    # nullチェック
    if ( $p2 eq '' ) {
        return 0;
    }

    # マスターパスワードチェック
    if ( $masterPassword eq $p2 ) {
        return 1;
    }

    # 本来のチェック
    if ( $p1 eq encode($p2) ) {
        return 1;
    }

    return 0;
}

# 1000億単位丸めルーチン
sub aboutMoney {
    my ($m) = @_;
    if ( $m < 500 ) {
        return "推定500${HunitMoney}未満";
    }
    else {
        $m = int( ( $m + 500 ) / 1000 );
        return "推定${m}000${HunitMoney}";
    }
}

# エスケープ文字の処理
sub htmlEscape {
    my ($s) = @_;
    $s =~ s/&/&amp;/g;
    $s =~ s/</&lt;/g;
    $s =~ s/>/&gt;/g;
    $s =~ s/\"/&quot;/g;    #"
    return $s;
}

# 80ケタに切り揃え
sub cutColumn {
    my ( $s, $c ) = @_;
    if ( length($s) <= $c ) {
        return $s;
    }
    else {
        # 合計80ケタになるまで切り取り
        my ($ss)    = '';
        my ($count) = 0;
        while ( $count < $c ) {
            $s =~ s/(^[\x80-\xFF][\x80-\xFF])|(^[\x00-\x7F])//;
            if ($1) {
                $ss .= $1;
                $count++;
            }
            else {
                $ss .= $2;
            }
            $count++;
        }
        return $ss;
    }
}

# 島の名前から番号を得る(IDじゃなくて番号)
sub nameToNumber {
    my ($name) = @_;

    # 全島から探す
    my ($i);
    for ( $i = 0; $i < $HislandNumber; $i++ ) {
        if ( $Hislands[$i]->{'name'} eq $name ) {
            return $i;
        }
    }

    # 見つからなかった場合
    return -1;
}

# 怪獣の情報
sub monsterSpec {
    my ($lv) = @_;

    # 種類
    my ($kind) = int( $lv / 10 );

    # 名前
    my ($name);
    $name = $HmonsterName[$kind];

    # 体力
    my ($hp) = $lv - ( $kind * 10 );

    return ( $kind, $name, $hp );
}

# 経験地からレベルを算出
sub expToLevel {
    my ( $kind, $exp ) = @_;
    my ($i);
    if ( $kind == $HlandBase ) {

        # ミサイル基地
        for ( $i = $maxBaseLevel; $i > 1; $i-- ) {
            if ( $exp >= $baseLevelUp[ $i - 2 ] ) {
                return $i;
            }
        }
        return 1;
    }
    else {
        # 海底基地
        for ( $i = $maxSBaseLevel; $i > 1; $i-- ) {
            if ( $exp >= $sBaseLevelUp[ $i - 2 ] ) {
                return $i;
            }
        }
        return 1;
    }

}

# (0,0)から(size - 1, size - 1)までの数字が一回づつ出てくるように
# (@Hrpx, @Hrpy)を設定
sub makeRandomPointArray {

    # 初期値
    my ($y);
    @Hrpx = ( 0 .. $HislandSize - 1 ) x $HislandSize;
    for ( $y = 0; $y < $HislandSize; $y++ ) {
        push( @Hrpy, ($y) x $HislandSize );
    }

    # シャッフル
    my ($i);
    for ( $i = $HpointNumber; --$i; ) {
        my ($j) = int( rand( $i + 1 ) );
        if ( $i == $j ) { next; }
        @Hrpx[ $i, $j ] = @Hrpx[ $j, $i ];
        @Hrpy[ $i, $j ] = @Hrpy[ $j, $i ];
    }
}

# 0から(n - 1)の乱数
sub random {
    return int( rand(1) * $_[0] );
}

#----------------------------------------------------------------------
# ログ表示
#----------------------------------------------------------------------
# ファイル番号指定でログ表示
sub logFilePrint {
    my ( $fileNumber, $id, $mode ) = @_;
    open( LIN, "${HdirName}/hakojima.log$_[0]" );
    my ( $line, $m, $turn, $id1, $id2, $message );
    while ( $line = <LIN> ) {
        $line =~ /^([0-9]*),([0-9]*),([0-9]*),([0-9]*),(.*)$/;
        ( $m, $turn, $id1, $id2, $message ) = ( $1, $2, $3, $4, $5 );

        # 機密関係
        if ( $m == 1 ) {
            if ( ( $mode == 0 ) || ( $id1 != $id ) ) {

                # 機密表示権利なし
                next;
            }
            $m = '<B>(機密)</B>';
        }
        else {
            $m = '';
        }

        # 表示的確か
        if ( $id != 0 ) {
            if (   ( $id != $id1 )
                && ( $id != $id2 ) )
            {
                next;
            }
        }

        # 表示
        out("<NOBR>${HtagNumber_}ターン$turn$m${H_tagNumber}：$message</NOBR><BR>\n"
        );
    }
    close(LIN);
}

#----------------------------------------------------------------------
# テンプレート
#----------------------------------------------------------------------
# 初期化
sub tempInitialize {

    # 島セレクト(デフォルト自分)
    $HislandList = getIslandList($defaultID);
    $HtargetList = getIslandList($defaultTarget);
}

# 島データのプルダウンメニュー用
sub getIslandList {
    my ($select) = @_;
    my ( $list, $name, $id, $s, $i );

    #島リストのメニュー
    $list = '';
    for ( $i = 0; $i < $HislandNumber; $i++ ) {
        $name = $Hislands[$i]->{'name'};
        $id   = $Hislands[$i]->{'id'};
        if ( $id eq $select ) {
            $s = 'SELECTED';
        }
        else {
            $s = '';
        }
        $list .= "<OPTION VALUE=\"$id\" $s>${name}島\n";
    }
    return $list;
}

# ヘッダ
sub tempHeader {
    out(<<END);
Content-type: text/html

<HTML>
<HEAD>
<META http-equiv="Content-Type" content="text/html; charset=UTF-8">
<TITLE>
$Htitle
</TITLE>
<BASE HREF="$imageDir/">
</HEAD>
<BODY $htmlBody>
<A HREF="http://www.bekkoame.ne.jp/~tokuoka/hakoniwa.html">箱庭諸島スクリプト配布元</A><HR>
END
}

# フッタ
sub tempFooter {
    out(<<END);
<HR>
<P align=center>
管理者:$adminName(<A HREF="mailto:$email">$email</A>)<BR>
掲示板(<A HREF="$bbs">$bbs</A>)<BR>
トップページ(<A HREF="$toppage">$toppage</A>)<BR>
箱庭諸島のページ(<A HREF="http://www.bekkoame.ne.jp/~tokuoka/hakoniwa.html">http://www.bekkoame.ne.jp/~tokuoka/hakoniwa.html</A>)<BR>
</P>
</BODY>
</HTML>
END
}

# ロック失敗
sub tempLockFail {

    # タイトル
    out(<<END);
${HtagBig_}同時アクセスエラーです。<BR>
ブラウザの「戻る」ボタンを押し、<BR>
しばらく待ってから再度お試し下さい。${H_tagBig}$HtempBack
END
}

# 強制解除
sub tempUnlock {

    # タイトル
    out(<<END);
${HtagBig_}前回のアクセスが異常終了だったようです。<BR>
ロックを強制解除しました。${H_tagBig}$HtempBack
END
}

# hakojima.datがない
sub tempNoDataFile {
    out(<<END);
${HtagBig_}データファイルが開けません。${H_tagBig}$HtempBack
END
}

# パスワード間違い
sub tempWrongPassword {
    out(<<END);
${HtagBig_}パスワードが違います。${H_tagBig}$HtempBack
END
}

# 何か問題発生
sub tempProblem {
    out(<<END);
${HtagBig_}問題発生、とりあえず戻ってください。${H_tagBig}$HtempBack
END
}

package Hako::Variable;

use strict;
use warnings;
use utf8;
use open ':encoding(utf8)';

use Exporter::Easy (
    EXPORT => [
        qw(
            $LOCKID
            $HcurrentNumber
            $HcommandArg
            $HcommandKind
            $HcommandMode
            $HcommandPlanNumber
            $HcommandTarget
            $HcommandX
            $HcommandY
            $HcurrentID
            $HcurrentName
            $HdefaultKind
            $HdefaultName
            $HdefaultPassword
            $HdefaultX
            $HdefaultY
            $HinputPassword
            $HinputPassword2
            $HislandLastTime
            $HislandList
            $HislandNextID
            $HislandNumber
            $HislandTurn
            $HlbbsMessage
            $HlbbsMode
            $HlbbsName
            $HmainMode
            $Hmessage
            $HoldPassword
            $HtargetList
            %HidToName
            %HidToNumber
            @Hislands
            @Hrpx
            @Hrpy
            @HdefenceHex
            @HlogPool
            @HlateLogPool
            @HsecretLogPool
            $defaultID
            $defaultTarget
            )
    ],
);

#----------------------------------------------------------------------
# ユーザ入力値
#----------------------------------------------------------------------

# 操作・表示対象の島
our $HcurrentID;

# $HcurrentIDの島番号(@Hislands内のindex)
our $HcurrentNumber;

# コマンド
our $HcommandArg;
our $HcommandKind;
our $HcommandMode;
our $HcommandPlanNumber;
our $HcommandTarget;
our $HcommandX;
our $HcommandY;

# 島新規作成・名前変更時の名前
our $HcurrentName;

# ユーザ入力パスワード
our $HinputPassword;
our $HinputPassword2;

# LBBS投稿メッセージ
our $HlbbsMessage;

# LBBS投稿者名
our $HlbbsName;

# コメントに設定するメッセージ
our $Hmessage;

# 変更前パスワード
our $HoldPassword;

# デフォルト値(Cookieに保存され、あらかじめフォームなどに設定される)
our $defaultID;
our $defaultTarget;
our $HdefaultKind;
our $HdefaultName;
our $HdefaultPassword;
our $HdefaultX;
our $HdefaultY;

#----------------------------------------------------------------------
# 島情報(ファイルから読みこんだもの)
#----------------------------------------------------------------------

# ターン数
our $HislandTurn;

# 最終ターン経過時刻
our $HislandLastTime;

# 島の総数
our $HislandNumber;

# 次に割り当てる島ID
our $HislandNextID;

# 島データのリスト
our @Hislands;

#----------------------------------------------------------------------
# テンポラリな情報
#----------------------------------------------------------------------

# モード(turn, new, print, owner, command, comment, lbbs, change, top)
our $HmainMode;

# LBBSモード(0:観光者,1:島主,2:削除)
our $HlbbsMode;

# 島セレクト用OPTIONコントロール(デフォルト自分)
our $HislandList;

# Targetの島セレクト用OPTIONコントロール
our $HtargetList;

# 島ID->名前
our %HidToName;

# 島ID->島番号
our %HidToNumber;

# 島の全座標をランダムにシャッフルしたもの
our @Hrpx;
our @Hrpy;

# 防衛施設の座標 [島ID][X][Y]
our @HdefenceHex;

# 通常ログ
our @HlogPool;

# 遅延ログ
our @HlateLogPool;

# 機密ログ
our @HsecretLogPool;

# ロック用ファイルハンドル($lockMode == 2の時に利用)
our $LOCKID;

1;


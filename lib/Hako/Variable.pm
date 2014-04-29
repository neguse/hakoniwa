package Hako::Variable;

use strict;
use warnings;
use utf8;

use Exporter::Easy (
  EXPORT => [qw(
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
  )],
);

#----------------------------------------------------------------------
# 変数
#----------------------------------------------------------------------

our $HcurrentNumber;
our $HcommandArg;
our $HcommandKind;
our $HcommandMode;
our $HcommandPlanNumber;
our $HcommandTarget;
our $HcommandX;
our $HcommandY;
our $HcurrentID;
our $HcurrentName;
our $HdefaultKind;
our $HdefaultName;
our $HdefaultPassword;
our $HdefaultX;
our $HdefaultY;
our $HinputPassword;
our $HinputPassword2;
our $HislandLastTime;
our $HislandList;
our $HislandNextID;
our $HislandNumber;
our $HislandTurn;
our $HlbbsMessage;
our $HlbbsMode;
our $HlbbsName;
our $HmainMode;
our $Hmessage;
our $HoldPassword;
our $HtargetList;
our %HidToName;
our %HidToNumber;
our @Hislands;
our @Hrpx;
our @Hrpy;
our @HdefenceHex;
our @HlogPool;
our @HlateLogPool;
our @HsecretLogPool;
# COOKIE
our ($defaultID);        # 島の名前
our ($defaultTarget);    # ターゲットの名前

1;


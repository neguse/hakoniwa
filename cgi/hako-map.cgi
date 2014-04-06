#----------------------------------------------------------------------
# 箱庭諸島 ver2.30
# 地図モードモジュール(ver1.00)
# 使用条件、使用方法等は、hako-readme.txtファイルを参照
#
# 箱庭諸島のページ: http://www.bekkoame.ne.jp/~tokuoka/hakoniwa.html
#----------------------------------------------------------------------


#----------------------------------------------------------------------
# 観光モード
#----------------------------------------------------------------------
# メイン
sub printIslandMain {
    # 開放
    unlock();

    # idから島番号を取得
    $HcurrentNumber = $HidToNumber{$HcurrentID};

    # なぜかその島がない場合
    if($HcurrentNumber eq '') {
	tempProblem();
	return;
    }

    # 名前の取得
    $HcurrentName = $Hislands[$HcurrentNumber]->{'name'};

    # 観光画面
    tempPrintIslandHead(); # ようこそ!!
    islandInfo(); # 島の情報
    islandMap(0); # 島の地図、観光モード

    # ○○島ローカル掲示板
    if($HuseLbbs) {
	tempLbbsHead();     # ローカル掲示板
	tempLbbsInput();   # 書き込みフォーム
	tempLbbsContents(); # 掲示板内容
    }

    # 近況
    tempRecent(0);
}

#----------------------------------------------------------------------
# 開発モード
#----------------------------------------------------------------------
# メイン
sub ownerMain {
    # 開放
    unlock();

    # モードを明示
    $HmainMode = 'owner';

    # idから島を取得
    $HcurrentNumber = $HidToNumber{$HcurrentID};
    my($island) = $Hislands[$HcurrentNumber];
    $HcurrentName = $island->{'name'};

    # パスワード
    if(!checkPassword($island->{'password'},$HinputPassword)) {
	# password間違い
	tempWrongPassword();
	return;
    }

    # 開発画面
    tempOwner(); # 「開発計画」

    # ○○島ローカル掲示板
    if($HuseLbbs) {
	tempLbbsHead();     # ローカル掲示板
	tempLbbsInputOW();   # 書き込みフォーム
	tempLbbsContents(); # 掲示板内容
    }

    # 近況
    tempRecent(1);
}

#----------------------------------------------------------------------
# コマンドモード
#----------------------------------------------------------------------
# メイン
sub commandMain {
    # idから島を取得
    $HcurrentNumber = $HidToNumber{$HcurrentID};
    my($island) = $Hislands[$HcurrentNumber];
    $HcurrentName = $island->{'name'};

    # パスワード
    if(!checkPassword($island->{'password'},$HinputPassword)) {
	# password間違い
	unlock();
	tempWrongPassword();
	return;
    }

    # モードで分岐
    my($command) = $island->{'command'};

    if($HcommandMode eq 'delete') {
	slideFront($command, $HcommandPlanNumber);
	tempCommandDelete();
    } elsif(($HcommandKind == $HcomAutoPrepare) ||
	    ($HcommandKind == $HcomAutoPrepare2)) {
	# フル整地、フル地ならし
	# 座標配列を作る
	makeRandomPointArray();
	my($land) = $island->{'land'};

	# コマンドの種類決定
	my($kind) = $HcomPrepare;
	if($HcommandKind == $HcomAutoPrepare2) {
	    $kind = $HcomPrepare2;
	}

	my($i) = 0;
	my($j) = 0;
	while(($j < $HpointNumber) && ($i < $HcommandMax)) {
	    my($x) = $Hrpx[$j];
	    my($y) = $Hrpy[$j];
	    if($land->[$x][$y] == $HlandWaste) {
		slideBack($command, $HcommandPlanNumber);
		$command->[$HcommandPlanNumber] = {
		    'kind' => $kind,
		    'target' => 0,
		    'x' => $x,
		    'y' => $y,
		    'arg' => 0
		    };
		$i++;
	    }
	    $j++;
	}
	tempCommandAdd();
    } elsif($HcommandKind == $HcomAutoDelete) {
	# 全消し
	my($i);
	for($i = 0; $i < $HcommandMax; $i++) {
	    slideFront($command, $HcommandPlanNumber);
	}
	tempCommandDelete();
    } else {
	if($HcommandMode eq 'insert') {
	    slideBack($command, $HcommandPlanNumber);
	}
	tempCommandAdd();
	# コマンドを登録
	$command->[$HcommandPlanNumber] = {
	    'kind' => $HcommandKind,
	    'target' => $HcommandTarget,
	    'x' => $HcommandX,
	    'y' => $HcommandY,
	    'arg' => $HcommandArg
	    };
    }

    # データの書き出し
    writeIslandsFile($HcurrentID);

    # owner modeへ
    ownerMain();

}

#----------------------------------------------------------------------
# コメント入力モード
#----------------------------------------------------------------------
# メイン
sub commentMain {
    # idから島を取得
    $HcurrentNumber = $HidToNumber{$HcurrentID};
    my($island) = $Hislands[$HcurrentNumber];
    $HcurrentName = $island->{'name'};

    # パスワード
    if(!checkPassword($island->{'password'},$HinputPassword)) {
	# password間違い
	unlock();
	tempWrongPassword();
	return;
    }

    # メッセージを更新
    $island->{'comment'} = htmlEscape($Hmessage);

    # データの書き出し
    writeIslandsFile($HcurrentID);

    # コメント更新メッセージ
    tempComment();

    # owner modeへ
    ownerMain();
}

#----------------------------------------------------------------------
# ローカル掲示板モード
#----------------------------------------------------------------------
# メイン

sub localBbsMain {
    # idから島番号を取得
    $HcurrentNumber = $HidToNumber{$HcurrentID};
    my($island) = $Hislands[$HcurrentNumber];

    # なぜかその島がない場合
    if($HcurrentNumber eq '') {
	unlock();
	tempProblem();
	return;
    }

    # 削除モードじゃなくて名前かメッセージがない場合
    if($HlbbsMode != 2) {
	if(($HlbbsName eq '') || ($HlbbsName eq '')) {
	    unlock();
	    tempLbbsNoMessage();
	    return;
	}
    }

    # 観光者モードじゃない時はパスワードチェック
    if($HlbbsMode != 0) {
	if(!checkPassword($island->{'password'},$HinputPassword)) {
	    # password間違い
	    unlock();
	    tempWrongPassword();
	    return;
	}
    }

    my($lbbs);
    $lbbs = $island->{'lbbs'};

    # モードで分岐
    if($HlbbsMode == 2) {
	# 削除モード
	# メッセージを前にずらす
	slideBackLbbsMessage($lbbs, $HcommandPlanNumber);
	tempLbbsDelete();
    } else {
	# 記帳モード
	# メッセージを後ろにずらす
	slideLbbsMessage($lbbs);

	# メッセージ書き込み
	my($message);
	if($HlbbsMode == 0) {
	    $message = '0';
	} else {
	    $message = '1';
	}
	$HlbbsName = "$HislandTurn：" . htmlEscape($HlbbsName);
	$HlbbsMessage = htmlEscape($HlbbsMessage);
	$lbbs->[0] = "$message>$HlbbsName>$HlbbsMessage";

	tempLbbsAdd();
    }

    # データ書き出し
    writeIslandsFile($HcurrentID);

    # もとのモードへ
    if($HlbbsMode == 0) {
	printIslandMain();
    } else {
	ownerMain();
    }
}

# ローカル掲示板のメッセージを一つ後ろにずらす
sub slideLbbsMessage {
    my($lbbs) = @_;
    my($i);
#    pop(@$lbbs);
#    push(@$lbbs, $lbbs->[0]);
    pop(@$lbbs);
    unshift(@$lbbs, $lbbs->[0]);
}

# ローカル掲示板のメッセージを一つ前にずらす
sub slideBackLbbsMessage {
    my($lbbs, $number) = @_;
    my($i);
    splice(@$lbbs, $number, 1);
    $lbbs->[$HlbbsMax - 1] = '0>>';
}

#----------------------------------------------------------------------
# 島の地図
#----------------------------------------------------------------------

# 情報の表示
sub islandInfo {
    my($island) = $Hislands[$HcurrentNumber];
    # 情報表示
    my($rank) = $HcurrentNumber + 1;
    my($farm) = $island->{'farm'};
    my($factory) = $island->{'factory'};
    my($mountain) = $island->{'mountain'};
    $farm = ($farm == 0) ? "保有せず" : "${farm}0$HunitPop";
    $factory = ($factory == 0) ? "保有せず" : "${factory}0$HunitPop";
    $mountain = ($mountain == 0) ? "保有せず" : "${mountain}0$HunitPop";

    my($mStr1) = '';
    my($mStr2) = '';
    if(($HhideMoneyMode == 1) || ($HmainMode eq 'owner')) {
	# 無条件またはownerモード
	$mStr1 = "<TH $HbgTitleCell nowrap=nowrap><NOBR>${HtagTH_}資金${H_tagTH}</NOBR></TH>";
	$mStr2 = "<TD $HbgInfoCell align=right nowrap=nowrap><NOBR>$island->{'money'}$HunitMoney</NOBR></TD>";
    } elsif($HhideMoneyMode == 2) {
	my($mTmp) = aboutMoney($island->{'money'});

	# 1000億単位モード
	$mStr1 = "<TH $HbgTitleCell nowrap=nowrap><NOBR>${HtagTH_}資金${H_tagTH}</NOBR></TH>";
	$mStr2 = "<TD $HbgInfoCell align=right nowrap=nowrap><NOBR>$mTmp</NOBR></TD>";
    }
    out(<<END);
<CENTER>
<TABLE BORDER>
<TR>
<TH $HbgTitleCell nowrap=nowrap><NOBR>${HtagTH_}順位${H_tagTH}</NOBR></TH>
<TH $HbgTitleCell nowrap=nowrap><NOBR>${HtagTH_}人口${H_tagTH}</NOBR></TH>
$mStr1
<TH $HbgTitleCell nowrap=nowrap><NOBR>${HtagTH_}食料${H_tagTH}</NOBR></TH>
<TH $HbgTitleCell nowrap=nowrap><NOBR>${HtagTH_}面積${H_tagTH}</NOBR></TH>
<TH $HbgTitleCell nowrap=nowrap><NOBR>${HtagTH_}農場規模${H_tagTH}</NOBR></TH>
<TH $HbgTitleCell nowrap=nowrap><NOBR>${HtagTH_}工場規模${H_tagTH}</NOBR></TH>
<TH $HbgTitleCell nowrap=nowrap><NOBR>${HtagTH_}採掘場規模${H_tagTH}</NOBR></TH>
</TR>
<TR>
<TD $HbgNumberCell align=middle nowrap=nowrap><NOBR>${HtagNumber_}$rank${H_tagNumber}</NOBR></TD>
<TD $HbgInfoCell align=right nowrap=nowrap><NOBR>$island->{'pop'}$HunitPop</NOBR></TD>
$mStr2
<TD $HbgInfoCell align=right nowrap=nowrap><NOBR>$island->{'food'}$HunitFood</NOBR></TD>
<TD $HbgInfoCell align=right nowrap=nowrap><NOBR>$island->{'area'}$HunitArea</NOBR></TD>
<TD $HbgInfoCell align=right nowrap=nowrap><NOBR>${farm}</NOBR></TD>
<TD $HbgInfoCell align=right nowrap=nowrap><NOBR>${factory}</NOBR></TD>
<TD $HbgInfoCell align=right nowrap=nowrap><NOBR>${mountain}</NOBR></TD>
</TR>
</TABLE></CENTER>
END
}

# 地図の表示
# 引数が1なら、ミサイル基地等をそのまま表示
sub islandMap {
    my($mode) = @_;
    my($island);
    $island = $Hislands[$HcurrentNumber];

    out(<<END);
<CENTER><TABLE BORDER><TR><TD>
END
    # 地形、地形値を取得
    my($land) = $island->{'land'};
    my($landValue) = $island->{'landValue'};
    my($l, $lv);

    # コマンド取得
    my($command) = $island->{'command'};
    my($com, @comStr, $i);
    if($HmainMode eq 'owner') {
	for($i = 0; $i < $HcommandMax; $i++) {
	    my($j) = $i + 1;
	    $com = $command->[$i];
	    if($com->{'kind'} < 20) {
		$comStr[$com->{'x'}][$com->{'y'}] .=
		    " [${j}]$HcomName[$com->{'kind'}]";
	    }
	}
    }

    # 座標(上)を出力
    out("<IMG SRC=\"xbar.gif\" width=400 height=16><BR>");

    # 各地形および改行を出力
    my($x, $y);
    for($y = 0; $y < $HislandSize; $y++) {
	# 偶数行目なら番号を出力
        if(($y % 2) == 0) {
	    out("<IMG SRC=\"space${y}.gif\" width=16 height=32>");
	}

	# 各地形を出力
	for($x = 0; $x < $HislandSize; $x++) {
	    $l = $land->[$x][$y];
	    $lv = $landValue->[$x][$y];
	    landString($l, $lv, $x, $y, $mode, $comStr[$x][$y]);
	}

	# 奇数行目なら番号を出力
        if(($y % 2) == 1) {
	    out("<IMG SRC=\"space${y}.gif\" width=16 height=32>");
	}

	# 改行を出力
	out("<BR>");
    }
    out("</TD></TR></TABLE></CENTER>\n");
}

sub landString {
    my($l, $lv, $x, $y, $mode, $comStr) = @_;
    my($point) = "($x,$y)";
    my($image, $alt);

    if($l == $HlandSea) {

	if($lv == 1) {
	    # 浅瀬
	    $image = 'land14.gif';
	    $alt = '海(浅瀬)';
        } else {
            # 海
	    $image = 'land0.gif';
	    $alt = '海';
        }
    } elsif($l == $HlandWaste) {
	# 荒地
	if($lv == 1) {
	    $image = 'land13.gif'; # 着弾点
	    $alt = '荒地';
	} else {
	    $image = 'land1.gif';
	    $alt = '荒地';
	}
    } elsif($l == $HlandPlains) {
	# 平地
	$image = 'land2.gif';
	$alt = '平地';
    } elsif($l == $HlandForest) {
	# 森
	if($mode == 1) {
	    $image = 'land6.gif';
	    $alt = "森(${lv}$HunitTree)";
	} else {
	    # 観光者の場合は木の本数隠す
	    $image = 'land6.gif';
	    $alt = '森';
	}
    } elsif($l == $HlandTown) {
	# 町
	my($p, $n);
	if($lv < 30) {
	    $p = 3;
	    $n = '村';
	} elsif($lv < 100) {
	    $p = 4;
	    $n = '町';
	} else {
	    $p = 5;
	    $n = '都市';
	}

	$image = "land${p}.gif";
	$alt = "$n(${lv}$HunitPop)";
    } elsif($l == $HlandFarm) {
	# 農場
	$image = 'land7.gif';
	$alt = "農場(${lv}0${HunitPop}規模)";
    } elsif($l == $HlandFactory) {
	# 工場
	$image = 'land8.gif';
	$alt = "工場(${lv}0${HunitPop}規模)";
    } elsif($l == $HlandBase) {
	if($mode == 0) {
	    # 観光者の場合は森のふり
	    $image = 'land6.gif';
	    $alt = '森';
	} else {
	    # ミサイル基地
	    my($level) = expToLevel($l, $lv);
	    $image = 'land9.gif';
	    $alt = "ミサイル基地 (レベル ${level}/経験値 $lv)";
	}
    } elsif($l == $HlandSbase) {
	# 海底基地
	if($mode == 0) {
	    # 観光者の場合は海のふり
	    $image = 'land0.gif';
	    $alt = '海';
	} else {
	    my($level) = expToLevel($l, $lv);
	    $image = 'land12.gif';
	    $alt = "海底基地 (レベル ${level}/経験値 $lv)";
	}
    } elsif($l == $HlandDefence) {
	# 防衛施設
	$image = 'land10.gif';
	$alt = '防衛施設';
    } elsif($l == $HlandHaribote) {
	# ハリボテ
	$image = 'land10.gif';
	if($mode == 0) {
	    # 観光者の場合は防衛施設のふり
	    $alt = '防衛施設';
	} else {
	    $alt = 'ハリボテ';
	}
    } elsif($l == $HlandOil) {
	# 海底油田
	$image = 'land16.gif';
	$alt = '海底油田';
    } elsif($l == $HlandMountain) {
	# 山
	my($str);
	$str = '';
	if($lv > 0) {
	    $image = 'land15.gif';
	    $alt = "山(採掘場${lv}0${HunitPop}規模)";
	} else {
	    $image = 'land11.gif';
	    $alt = '山';
	}
    } elsif($l == $HlandMonument) {
	# 記念碑
	$image = $HmonumentImage[$lv];
	$alt = $HmonumentName[$lv];
    } elsif($l == $HlandMonster) {
	# 怪獣
	my($kind, $name, $hp) = monsterSpec($lv);
	my($special) = $HmonsterSpecial[$kind];
	$image = $HmonsterImage[$kind];

	# 硬化中?
	if((($special == 3) && (($HislandTurn % 2) == 1)) ||
	   (($special == 4) && (($HislandTurn % 2) == 0))) {
	    # 硬化中
	    $image = $HmonsterImage2[$kind];
	}
	$alt = "怪獣$name(体力${hp})";
    }


    # 開発画面の場合は、座標設定
    if($mode == 1) {
	out("<A HREF=\"JavaScript:void(0);\" onclick=\"ps($x,$y)\">");
    }

    out("<IMG SRC=\"$image\" ALT=\"$point $alt $comStr\" width=32 height=32 BORDER=0>");

    # 座標設定閉じ
    if($mode == 1) {
	out("</A>");
    }
}


#----------------------------------------------------------------------
# テンプレートその他
#----------------------------------------------------------------------
# 個別ログ表示
sub logPrintLocal {
    my($mode) = @_;
    my($i);
    for($i = 0; $i < $HlogMax; $i++) {
	logFilePrint($i, $HcurrentID, $mode);
    }
}

# ○○島へようこそ！！
sub tempPrintIslandHead {
    out(<<END);
<CENTER>
${HtagBig_}${HtagName_}「${HcurrentName}島」${H_tagName}へようこそ！！${H_tagBig}<BR>
$HtempBack<BR>
</CENTER>
END
}

# ○○島開発計画
sub tempOwner {
    out(<<END);
<CENTER>
${HtagBig_}${HtagName_}${HcurrentName}島${H_tagName}開発計画${H_tagBig}<BR>
$HtempBack<BR>
</CENTER>
<SCRIPT Language="JavaScript">
<!--
function ps(x, y) {
    document.forms[0].elements[4].options[x].selected = true;
    document.forms[0].elements[5].options[y].selected = true;
    return true;
}

function ns(x) {
    document.forms[0].elements[2].options[x].selected = true;
    return true;
}

//-->
</SCRIPT>
END

    islandInfo();

    out(<<END);
<CENTER>
<TABLE BORDER>
<TR>
<TD $HbgInputCell >
<CENTER>
<FORM action="$HthisFile" method=POST>
<INPUT TYPE=submit VALUE="計画送信" NAME=CommandButton$Hislands[$HcurrentNumber]->{'id'}>
<HR>
<B>パスワード</B></BR>
<INPUT TYPE=password NAME=PASSWORD VALUE="$HdefaultPassword">
<HR>
<B>計画番号</B><SELECT NAME=NUMBER>
END
    # 計画番号
    my($j, $i);
    for($i = 0; $i < $HcommandMax; $i++) {
	$j = $i + 1;
	out("<OPTION VALUE=$i>$j\n");
    }

    out(<<END);
</SELECT><BR>
<HR>
<B>開発計画</B><BR>
<SELECT NAME=COMMAND>
END

    #コマンド
    my($kind, $cost, $s);
    for($i = 0; $i < $HcommandTotal; $i++) {
	$kind = $HcomList[$i];
	$cost = $HcomCost[$kind];
	if($cost == 0) {
	    $cost = '無料'
	} elsif($cost < 0) {
	    $cost = - $cost;
	    $cost .= $HunitFood;
	} else {
	    $cost .= $HunitMoney;
	}
	if($kind == $HdefaultKind) {
	    $s = 'SELECTED';
	} else {
	    $s = '';
	}
	out("<OPTION VALUE=$kind $s>$HcomName[$kind]($cost)\n");
    }

    out(<<END);
</SELECT>
<HR>
<B>座標(</B>
<SELECT NAME=POINTX>

END
    for($i = 0; $i < $HislandSize; $i++) {
	if($i == $HdefaultX) {
	    out("<OPTION VALUE=$i SELECTED>$i\n");
        } else {
	    out("<OPTION VALUE=$i>$i\n");
        }
    }

    out(<<END);
</SELECT>, <SELECT NAME=POINTY>
END

    for($i = 0; $i < $HislandSize; $i++) {
	if($i == $HdefaultY) {
	    out("<OPTION VALUE=$i SELECTED>$i\n");
        } else {
	    out("<OPTION VALUE=$i>$i\n");
        }
    }
    out(<<END);
</SELECT><B>)</B>
<HR>
<B>数量</B><SELECT NAME=AMOUNT>
END

    # 数量
    for($i = 0; $i < 100; $i++) {
	out("<OPTION VALUE=$i>$i\n");
    }

    out(<<END);
</SELECT>
<HR>
<B>目標の島</B><BR>
<SELECT NAME=TARGETID>
$HtargetList<BR>
</SELECT>
<HR>
<B>動作</B><BR>
<INPUT TYPE=radio NAME=COMMANDMODE VALUE=insert CHECKED>挿入
<INPUT TYPE=radio NAME=COMMANDMODE VALUE=write>上書き<BR>
<INPUT TYPE=radio NAME=COMMANDMODE VALUE=delete>削除
<HR>
<INPUT TYPE=submit VALUE="計画送信" NAME=CommandButton$Hislands[$HcurrentNumber]->{'id'}>

</CENTER>
</FORM>
</TD>
<TD $HbgMapCell>
END
    islandMap(1);    # 島の地図、所有者モード
    out(<<END);
</TD>
<TD $HbgCommandCell>
END
    for($i = 0; $i < $HcommandMax; $i++) {
	tempCommand($i, $Hislands[$HcurrentNumber]->{'command'}->[$i]);
    }

    out(<<END);

</TD>
</TR>
</TABLE>
</CENTER>
<HR>
<CENTER>
${HtagBig_}コメント更新${H_tagBig}<BR>
<FORM action="$HthisFile" method="POST">
コメント<INPUT TYPE=text NAME=MESSAGE SIZE=80><BR>
パスワード<INPUT TYPE=password NAME=PASSWORD VALUE="$HdefaultPassword">
<INPUT TYPE=submit VALUE="コメント更新" NAME=MessageButton$Hislands[$HcurrentNumber]->{'id'}>
</FORM>
</CENTER>
END

}

# 入力済みコマンド表示
sub tempCommand {
    my($number, $command) = @_;
    my($kind, $target, $x, $y, $arg) =
	(
	 $command->{'kind'},
	 $command->{'target'},
	 $command->{'x'},
	 $command->{'y'},
	 $command->{'arg'}
	 );
    my($name) = "$HtagComName_${HcomName[$kind]}$H_tagComName";
    my($point) = "$HtagName_($x,$y)$H_tagName";
    $target = $HidToName{$target};
    if($target eq '') {
	$target = "無人";
    }
    $target = "$HtagName_${target}島$H_tagName";
    my($value) = $arg * $HcomCost[$kind];
    if($value == 0) {
	$value = $HcomCost[$kind];
    }
    if($value < 0) {
	$value = -$value;
	$value = "$value$HunitFood";
    } else {
	$value = "$value$HunitMoney";
    }
    $value = "$HtagName_$value$H_tagName";

    my($j) = sprintf("%02d：", $number + 1);

    out("<A STYlE=\"text-decoration:none\" HREF=\"JavaScript:void(0);\" onClick=\"ns($number)\"><NOBR>$HtagNumber_$j$H_tagNumber<FONT COLOR=\"$HnormalColor\">");

    if(($kind == $HcomDoNothing) ||
       ($kind == $HcomGiveup)) {
	out("$name");
    } elsif(($kind == $HcomMissileNM) ||
	    ($kind == $HcomMissilePP) ||
	    ($kind == $HcomMissileST) ||
	    ($kind == $HcomMissileLD)) {
	# ミサイル系
	my($n) = ($arg == 0 ? '無制限' : "${arg}発");
	out("$target$pointへ$name($HtagName_$n$H_tagName)");
    } elsif($kind == $HcomSendMonster) {
	# 怪獣派遣
	out("$targetへ$name");
    } elsif($kind == $HcomSell) {
	# 食料輸出
	out("$name$value");
    } elsif($kind == $HcomPropaganda) {
	# 誘致活動
	out("$name");
    } elsif(($kind == $HcomMoney) ||
	    ($kind == $HcomFood)) {
	# 援助
	out("$targetへ$name$value");
    } elsif($kind == $HcomDestroy) {
	# 掘削
	if($arg != 0) {
	    out("$pointで$name(予算${value})");
	} else {
	    out("$pointで$name");
	}
    } elsif(($kind == $HcomFarm) ||
	     ($kind == $HcomFactory) ||
	     ($kind == $HcomMountain)) {	
	# 回数付き
	if($arg == 0) {
	    out("$pointで$name");
	} else {
	    out("$pointで$name($arg回)");
	}
    } else {
	# 座標付き
	out("$pointで$name");
    }

    out("</FONT></NOBR></A><BR>");
}

# ローカル掲示板
sub tempLbbsHead {
    out(<<END);
<HR>
<CENTER>
${HtagBig_}${HtagName_}${HcurrentName}島${H_tagName}観光者通信${H_tagBig}<BR>
</CENTER>
END
}

# ローカル掲示板入力フォーム
sub tempLbbsInput {
    out(<<END);
<CENTER>
<FORM action="$HthisFile" method="POST">
<TABLE BORDER>
<TR>
<TH>名前</TH>
<TH>内容</TH>
<TH>動作</TH>
</TR>
<TR>
<TD><INPUT TYPE="text" SIZE=32 MAXLENGTH=32 NAME="LBBSNAME" VALUE="$HdefaultName"></TD>
<TD><INPUT TYPE="text" SIZE=80 NAME="LBBSMESSAGE"></TD>
<TD><INPUT TYPE="submit" VALUE="記帳する" NAME="LbbsButtonSS$HcurrentID"></TD>
</TR>
</TABLE>
</FORM>
</CENTER>
END
}

# ローカル掲示板入力フォーム owner mode用
sub tempLbbsInputOW {
    out(<<END);
<CENTER>
<FORM action="$HthisFile" method="POST">
<TABLE BORDER>
<TR>
<TH>名前</TH>
<TH COLSPAN=2>内容</TH>
</TR>
<TR>
<TD><INPUT TYPE="text" SIZE=32 MAXLENGTH=32 NAME="LBBSNAME" VALUE="$HdefaultName"></TD>
<TD COLSPAN=2><INPUT TYPE="text" SIZE=80 NAME="LBBSMESSAGE"></TD>
</TR>
<TR>
<TH>パスワード</TH>
<TH COLSPAN=2>動作</TH>
</TR>
<TR>
<TD><INPUT TYPE=password SIZE=32 MAXLENGTH=32 NAME=PASSWORD VALUE="$HdefaultPassword"></TD>
<TD align=right>
<INPUT TYPE="submit" VALUE="記帳する" NAME="LbbsButtonOW$HcurrentID">
</TD>
<TD align=right>
番号
<SELECT NAME=NUMBER>
END
    # 発言番号
    my($j, $i);
    for($i = 0; $i < $HlbbsMax; $i++) {
	$j = $i + 1;
	out("<OPTION VALUE=$i>$j\n");
    }
    out(<<END);
</SELECT>
<INPUT TYPE="submit" VALUE="削除する" NAME="LbbsButtonDL$HcurrentID">
</TD>
</TR>
</TABLE>
</FORM>
</CENTER>
END
}

# ローカル掲示板内容
sub tempLbbsContents {
    my($lbbs, $line);
    $lbbs = $Hislands[$HcurrentNumber]->{'lbbs'};
    out(<<END);
<CENTER>
<TABLE BORDER>
<TR>
<TH>番号</TH>
<TH>記帳内容</TH>
</TR>
END

    my($i);
    for($i = 0; $i < $HlbbsMax; $i++) {
	$line = $lbbs->[$i];
	if($line =~ /([0-9]*)\>(.*)\>(.*)$/) {
	    my($j) = $i + 1;
	    out("<TR><TD align=center>$HtagNumber_$j$H_tagNumber</TD>");
	    if($1 == 0) {
		# 観光者
		out("<TD>$HtagLbbsSS_$2 > $3$H_tagLbbsSS</TD></TR>");
	    } else {
		# 島主
		out("<TD>$HtagLbbsOW_$2 > $3$H_tagLbbsOW</TD></TR>");
	    }
	}
    }

    out(<<END);
</TD></TR></TABLE></CENTER>
END
}

# ローカル掲示板で名前かメッセージがない場合
sub tempLbbsNoMessage {
    out(<<END);
${HtagBig_}名前または内容の欄が空欄です。${H_tagBig}$HtempBack
END
}

# 書きこみ削除
sub tempLbbsDelete {
    out(<<END);
${HtagBig_}記帳内容を削除しました${H_tagBig}<HR>
END
}

# コマンド登録
sub tempLbbsAdd {
    out(<<END);
${HtagBig_}記帳を行いました${H_tagBig}<HR>
END
}

# コマンド削除
sub tempCommandDelete {
    out(<<END);
${HtagBig_}コマンドを削除しました${H_tagBig}<HR>
END
}

# コマンド登録
sub tempCommandAdd {
    out(<<END);
${HtagBig_}コマンドを登録しました${H_tagBig}<HR>
END
}

# コメント変更成功
sub tempComment {
    out(<<END);
${HtagBig_}コメントを更新しました${H_tagBig}<HR>
END
}

# 近況
sub tempRecent {
    my($mode) = @_;
    out(<<END);
<HR>
${HtagBig_}${HtagName_}${HcurrentName}島${H_tagName}の近況${H_tagBig}<BR>
END
    logPrintLocal($mode);
}

1;

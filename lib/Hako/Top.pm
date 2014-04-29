package Hako::Top;

use strict;
use warnings;
use utf8;

use Exporter::Easy (
  EXPORT => [qw(
    topPageMain
  )],
);

use Hako::Const;
use Hako::Variable;
use Hako::Main;

#----------------------------------------------------------------------
# 箱庭諸島 ver2.30
# トップモジュール(ver1.00)
# 使用条件、使用方法等は、hako-readme.txtファイルを参照
#
# 箱庭諸島のページ: http://www.bekkoame.ne.jp/~tokuoka/hakoniwa.html
#----------------------------------------------------------------------

#----------------------------------------------------------------------
# トップページモード
#----------------------------------------------------------------------
# メイン
sub topPageMain {

    # 開放
    unlock();

    # テンプレート出力
    tempTopPage();
}

# トップページ
sub tempTopPage {

    # タイトル
    out(<<END);
${HtagTitle_}$Htitle${H_tagTitle}
END

    # デバッグモードなら「ターンを進める」ボタン
    if ( $Hdebug == 1 ) {
        out(<<END);
<FORM action="$HthisFile" method="POST">
<INPUT TYPE="submit" VALUE="ターンを進める" NAME="TurnButton">
</FORM>
END
    }

    my ($mStr1) = '';
    if ( $HhideMoneyMode != 0 ) {
        $mStr1
            = "<TH $HbgTitleCell align=center nowrap=nowrap><NOBR>${HtagTH_}資金${H_tagTH}</NOBR></TH>";
    }

    # フォーム
    out(<<END);
<H1>${HtagHeader_}ターン$HislandTurn${H_tagHeader}</H1>

<HR>
<H1>${HtagHeader_}自分の島へ${H_tagHeader}</H1>
<FORM action="$HthisFile" method="POST">
あなたの島の名前は？<BR>
<SELECT NAME="ISLANDID">
$HislandList
</SELECT><BR>

パスワードをどうぞ！！<BR>
<INPUT TYPE="password" NAME="PASSWORD" VALUE="$HdefaultPassword" SIZE=32 MAXLENGTH=32><BR>
<INPUT TYPE="submit" VALUE="開発しに行く" NAME="OwnerButton"><BR>
</FORM>

<HR>

<H1>${HtagHeader_}諸島の状況${H_tagHeader}</H1>
<P>
島の名前をクリックすると、<B>観光</B>することができます。
</P>
<TABLE BORDER>
<TR>
<TH $HbgTitleCell align=center nowrap=nowrap><NOBR>${HtagTH_}順位${H_tagTH}</NOBR></TH>
<TH $HbgTitleCell align=center nowrap=nowrap><NOBR>${HtagTH_}島${H_tagTH}</NOBR></TH>
<TH $HbgTitleCell align=center nowrap=nowrap><NOBR>${HtagTH_}人口${H_tagTH}</NOBR></TH>
<TH $HbgTitleCell align=center nowrap=nowrap><NOBR>${HtagTH_}面積${H_tagTH}</NOBR></TH>
$mStr1
<TH $HbgTitleCell align=center nowrap=nowrap><NOBR>${HtagTH_}食料${H_tagTH}</NOBR></TH>
<TH $HbgTitleCell align=center nowrap=nowrap><NOBR>${HtagTH_}農場規模${H_tagTH}</NOBR></TH>
<TH $HbgTitleCell align=center nowrap=nowrap><NOBR>${HtagTH_}工場規模${H_tagTH}</NOBR></TH>
<TH $HbgTitleCell align=center nowrap=nowrap><NOBR>${HtagTH_}採掘場規模${H_tagTH}</NOBR></TH>
</TR>
END

    my ( $island, $j, $farm, $factory, $mountain, $name, $id, $prize, $ii );
    for ( $ii = 0; $ii < $HislandNumber; $ii++ ) {
        $j      = $ii + 1;
        $island = $Hislands[$ii];

        $id       = $island->{'id'};
        $farm     = $island->{'farm'};
        $factory  = $island->{'factory'};
        $mountain = $island->{'mountain'};
        $farm     = ( $farm == 0 ) ? "保有せず" : "${farm}0$HunitPop";
        $factory
            = ( $factory == 0 ) ? "保有せず" : "${factory}0$HunitPop";
        $mountain
            = ( $mountain == 0 ) ? "保有せず" : "${mountain}0$HunitPop";

        if ( $island->{'absent'} == 0 ) {
            $name = "${HtagName_}$island->{'name'}島${H_tagName}";
        }
        else {
            $name
                = "${HtagName2_}$island->{'name'}島($island->{'absent'})${H_tagName2}";
        }

        $prize = $island->{'prize'};
        my ( $flags, $monsters, $turns );
        $prize =~ /([0-9]*),([0-9]*),(.*)/;
        $flags    = $1;
        $monsters = $2;
        $turns    = $3;
        $prize    = '';

        # ターン杯の表示
        while ( $turns =~ s/([0-9]*),// ) {
            $prize
                .= "<IMG SRC=\"prize0.gif\" ALT=\"$1${Hprize[0]}\" WIDTH=16 HEIGHT=16> ";
        }

        # 名前に賞の文字を追加
        my ($f) = 1;
        my ($i);
        for ( $i = 1; $i < 10; $i++ ) {
            if ( $flags & $f ) {
                $prize
                    .= "<IMG SRC=\"prize${i}.gif\" ALT=\"${Hprize[$i]}\" WIDTH=16 HEIGHT=16> ";
            }
            $f *= 2;
        }

        # 倒した怪獣リスト
        $f = 1;
        my ($max)       = -1;
        my ($mNameList) = '';
        for ( $i = 0; $i < $HmonsterNumber; $i++ ) {
            if ( $monsters & $f ) {
                $mNameList .= "[$HmonsterName[$i]] ";
                $max = $i;
            }
            $f *= 2;
        }
        if ( $max != -1 ) {
            $prize
                .= "<IMG SRC=\"${HmonsterImage[$max]}\" ALT=\"$mNameList\" WIDTH=16 HEIGHT=16> ";
        }

        my ($mStr1) = '';
        if ( $HhideMoneyMode == 1 ) {
            $mStr1
                = "<TD $HbgInfoCell align=right nowrap=nowrap><NOBR>$island->{'money'}$HunitMoney</NOBR></TD>";
        }
        elsif ( $HhideMoneyMode == 2 ) {
            my ($mTmp) = aboutMoney( $island->{'money'} );
            $mStr1
                = "<TD $HbgInfoCell align=right nowrap=nowrap><NOBR>$mTmp</NOBR></TD>";
        }

        out(<<END);
<TR>
<TD $HbgNumberCell ROWSPAN=2 align=center nowrap=nowrap><NOBR>${HtagNumber_}$j${H_tagNumber}</NOBR></TD>
<TD $HbgNameCell ROWSPAN=2 align=left nowrap=nowrap>
<NOBR>
<A STYlE=\"text-decoration:none\" HREF="${HthisFile}?Sight=${id}">
$name
</A>
</NOBR><BR>
$prize
</TD>
<TD $HbgInfoCell align=right nowrap=nowrap>
<NOBR>$island->{'pop'}$HunitPop</NOBR></TD>
<TD $HbgInfoCell align=right nowrap=nowrap><NOBR>$island->{'area'}$HunitArea</NOBR></TD>
$mStr1
<TD $HbgInfoCell align=right nowrap=nowrap><NOBR>$island->{'food'}$HunitFood</NOBR></TD>
<TD $HbgInfoCell align=right nowrap=nowrap><NOBR>$farm</NOBR></TD>
<TD $HbgInfoCell align=right nowrap=nowrap><NOBR>$factory</NOBR></TD>
<TD $HbgInfoCell align=right nowrap=nowrap><NOBR>$mountain</NOBR></TD>
</TR>
<TR>
<TD $HbgCommentCell COLSPAN=7 align=left nowrap=nowrap><NOBR>${HtagTH_}コメント：${H_tagTH}$island->{'comment'}</NOBR></TD>
</TR>
END
    }

    out(<<END);
</TABLE>

<HR>
<H1>${HtagHeader_}新しい島を探す${H_tagHeader}</H1>
END

    if ( $HislandNumber < $HmaxIsland ) {
        out(<<END);
<FORM action="$HthisFile" method="POST">
どんな名前をつける予定？<BR>
<INPUT TYPE="text" NAME="ISLANDNAME" SIZE=32 MAXLENGTH=32>島<BR>
パスワードは？<BR>
<INPUT TYPE="password" NAME="PASSWORD" SIZE=32 MAXLENGTH=32><BR>
念のためパスワードをもう一回<BR>
<INPUT TYPE="password" NAME="PASSWORD2" SIZE=32 MAXLENGTH=32><BR>

<INPUT TYPE="submit" VALUE="探しに行く" NAME="NewIslandButton">
</FORM>
END
    }
    else {
        out(<<END);
        島の数が最大数です・・・現在登録できません。
END
    }

    out(<<END);
<HR>
<H1>${HtagHeader_}島の名前とパスワードの変更${H_tagHeader}</H1>
<P>
(注意)名前の変更には$HcostChangeName${HunitMoney}かかります。
</P>
<FORM action="$HthisFile" method="POST">
どの島ですか？<BR>
<SELECT NAME="ISLANDID">
$HislandList
</SELECT>
<BR>
どんな名前に変えますか？(変更する場合のみ)<BR>
<INPUT TYPE="text" NAME="ISLANDNAME" SIZE=32 MAXLENGTH=32>島<BR>
パスワードは？(必須)<BR>
<INPUT TYPE="password" NAME="OLDPASS" SIZE=32 MAXLENGTH=32><BR>
新しいパスワードは？(変更する時のみ)<BR>
<INPUT TYPE="password" NAME="PASSWORD" SIZE=32 MAXLENGTH=32><BR>
念のためパスワードをもう一回(変更する時のみ)<BR>
<INPUT TYPE="password" NAME="PASSWORD2" SIZE=32 MAXLENGTH=32><BR>

<INPUT TYPE="submit" VALUE="変更する" NAME="ChangeInfoButton">
</FORM>

<HR>

<H1>${HtagHeader_}最近の出来事${H_tagHeader}</H1>
END
    logPrintTop();
    out(<<END);
<H1>${HtagHeader_}発見の記録${H_tagHeader}</H1>
END
    historyPrint();
}

# トップページ用ログ表示
sub logPrintTop {
    my ($i);
    for ( $i = 0; $i < $HtopLogTurn; $i++ ) {
        logFilePrint( $i, 0, 0 );
    }
}

# 記録ファイル表示
sub historyPrint {
    open( HIN, "${HdirName}/hakojima.his" );
    my ( @line, $l );
    while ( $l = <HIN> ) {
        chomp($l);
        push( @line, $l );
    }
    @line = reverse(@line);

    foreach $l (@line) {
        $l =~ /^([0-9]*),(.*)$/;
        out("<NOBR>${HtagNumber_}ターン${1}${H_tagNumber}：${2}</NOBR><BR>\n"
        );
    }
    close(HIN);
}

1;

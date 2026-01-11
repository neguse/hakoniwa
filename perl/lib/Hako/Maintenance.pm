package Hako::Maintenance;

use strict;
use warnings;
use utf8;
use open ':encoding(utf8)';

use Exporter::Easy (
    EXPORT => [
        qw(
            run_maintenance
            )
    ],
);

#----------------------------------------------------------------------
# 箱庭諸島 ver2.30
# メンテナンスツール(ver1.01)
# 使用条件、使用方法等は、hako-readme.txtファイルを参照
#
# 箱庭諸島のページ: http://www.bekkoame.ne.jp/~tokuoka/hakoniwa.html
#----------------------------------------------------------------------

# ――――――――――――――――――――――――――――――
# 各種設定値
# ――――――――――――――――――――――――――――――

# マスターパスワード
my ($masterpassword) = 'yourpassword';

# 1ターンが何秒か
my ($unitTime) = 21600;    # 6時間

# ディレクトリのパーミッション
my ($dirMode) = 0755;

# このファイル
# my($thisFile) = 'http://場所/cgi-bin/hako-mente.cgi';
my ($thisFile) = 'http://localhost:5000/cgi-bin/hako-mente.cgi';

# データディレクトリの名前
# hakojima.cgi中のものと合わせてください。
my ($dirName) = 'data';

# use Time::Localが使えない環境では、'use Time::Local'の行を消して下さい。
# ただし、更新時間の変更が'秒指定で変更'しかできなくなります。
use Time::Local;

# ――――――――――――――――――――――――――――――
# 設定項目は以上
# ――――――――――――――――――――――――――――――

# 各種変数
my ($mainMode);
my ($inputPass);
my ($deleteID);
my ($currentID);
my ($ctYear);
my ($ctMon);
my ($ctDate);
my ($ctHour);
my ($ctMin);
my ($ctSec);

# 標準出力への出力
sub out {
    print STDOUT Encode::encode_utf8( $_[0] );
}

sub run_maintenance {

    out <<END;
Content-type: text/html

<HTML>
<HEAD>
<META http-equiv="Content-Type" content="text/html; charset=UTF-8">
<TITLE>箱島２ メンテナンスツール</TITLE>
</HEAD>
<BODY>
END

    cgiInput();

    if ( $mainMode eq 'delete' ) {
        if ( passCheck() ) {
            deleteMode();
        }
    }
    elsif ( $mainMode eq 'current' ) {
        if ( passCheck() ) {
            currentMode();
        }
    }
    elsif ( $mainMode eq 'time' ) {
        if ( passCheck() ) {
            timeMode();
        }
    }
    elsif ( $mainMode eq 'stime' ) {
        if ( passCheck() ) {
            stimeMode();
        }
    }
    elsif ( $mainMode eq 'new' ) {
        if ( passCheck() ) {
            newMode();
        }
    }
    mainMode();

    out <<END;
</FORM>
</BODY>
</HTML>
END

}

sub myrmtree {
    my ($dn) = @_;
    opendir( my $DIN, "$dn/" );
    my ($fileName);
    while ( $fileName = readdir($DIN) ) {
        unlink("$dn/$fileName");
    }
    closedir($DIN);
    rmdir($dn);
}

sub currentMode {
    myrmtree "${dirName}";
    mkdir( "${dirName}", $dirMode );
    opendir( my $DIN, "${dirName}.bak$currentID/" );
    my ($fileName);
    while ( $fileName = readdir($DIN) ) {
        fileCopy( "${dirName}.bak$currentID/$fileName",
            "${dirName}/$fileName" );
    }
    closedir($DIN);
}

sub deleteMode {
    if ( $deleteID eq '' ) {
        myrmtree "${dirName}";
    }
    else {
        myrmtree "${dirName}.bak$deleteID";
    }
    unlink "hakojimalockflock";
}

sub newMode {
    mkdir( $dirName, $dirMode );

    # 現在の時間を取得
    my ($now) = time;
    $now = $now - ( $now % ($unitTime) );

    open( my $OUT, ">$dirName/hakojima.dat" );    # ファイルを開く
    print $OUT "1\n";                             # ターン数1
    print $OUT "$now\n";                          # 開始時間
    print $OUT "0\n";                             # 島の数
    print $OUT "1\n";                             # 次に割り当てるID

    # ファイルを閉じる
    close($OUT);
}

sub timeMode {
    $ctMon--;
    $ctYear -= 1900;
    $ctSec = timelocal( $ctSec, $ctMin, $ctHour, $ctDate, $ctMon, $ctYear );
    stimeMode();
}

sub stimeMode {
    my ($t) = $ctSec;
    open( my $IN, "${dirName}/hakojima.dat" );
    my (@lines);
    @lines = <$IN>;
    close($IN);

    $lines[1] = "$t\n";

    open( my $OUT, ">${dirName}/hakojima.dat" );
    print $OUT @lines;
    close($OUT);
}

sub mainMode {
    opendir( my $DIN, "./" );

    out <<END;
<FORM action="$thisFile" method="POST">
<H1>箱島２ メンテナンスツール</H1>
<B>パスワード:</B><INPUT TYPE=password SIZE=32 MAXLENGTH=32 NAME=PASSWORD></TD>
END

    # 現役データ
    if ( -d "${dirName}" ) {
        dataPrint("");
    }
    else {
        out <<END;
    <HR>
    <INPUT TYPE="submit" VALUE="新しいデータを作る" NAME="NEW">
END
    }

    # バックアップデータ
    my ($dn);
    while ( $dn = readdir($DIN) ) {
        if ( $dn =~ /^${dirName}.bak(.*)/ ) {
            dataPrint($1);
        }
    }
    closedir($DIN);
}

# 表示モード
sub dataPrint {
    my ($suf) = @_;

    out "<HR>";
    my $IN;
    if ( $suf eq "" ) {
        open( $IN, "${dirName}/hakojima.dat" );
        out "<H1>現役データ</H1>";
    }
    else {
        open( $IN, "${dirName}.bak$suf/hakojima.dat" );
        out "<H1>バックアップ$suf</H1>";
    }

    my ($lastTurn);
    $lastTurn = <$IN>;
    my ($lastTime);
    $lastTime = <$IN>;

    close($IN);

    my ($timeString) = timeToString($lastTime);

    out <<END;
    <B>ターン$lastTurn</B><BR>
    <B>最終更新時間</B>:$timeString<BR>
    <B>最終更新時間(秒数表示)</B>:1970年1月1日から$lastTime 秒<BR>
    <INPUT TYPE="submit" VALUE="このデータを削除" NAME="DELETE$suf">
END

    if ( $suf eq "" ) {
        my ( $sec, $min, $hour, $date, $mon, $year, $day, $yday, $dummy )
            = localtime($lastTime);
        $mon++;
        $year += 1900;

        out <<END;
    <H2>最終更新時間の変更</H2>
    <INPUT TYPE="text" SIZE=4 NAME="YEAR" VALUE="$year">年
    <INPUT TYPE="text" SIZE=2 NAME="MON" VALUE="$mon">月
    <INPUT TYPE="text" SIZE=2 NAME="DATE" VALUE="$date">日
    <INPUT TYPE="text" SIZE=2 NAME="HOUR" VALUE="$hour">時
    <INPUT TYPE="text" SIZE=2 NAME="MIN" VALUE="$min">分
    <INPUT TYPE="text" SIZE=2 NAME="NSEC" VALUE="$sec">秒
    <INPUT TYPE="submit" VALUE="変更" NAME="NTIME"><BR>
    1970年1月1日から<INPUT TYPE="text" SIZE=32 NAME="SSEC" VALUE="$lastTime">秒
    <INPUT TYPE="submit" VALUE="秒指定で変更" NAME="STIME">

END
    }
    else {
        out <<END;
	<INPUT TYPE="submit" VALUE="このデータを現役に" NAME="CURRENT$suf">
END
    }
}

sub timeToString {
    my ( $sec, $min, $hour, $date, $mon, $year, $day, $yday, $dummy )
        = localtime( $_[0] );
    $mon++;
    $year += 1900;

    return "${year}年 ${mon}月 ${date}日 ${hour}時 ${min}分 ${sec}秒";
}

# CGIの読みこみ
sub cgiInput {
    my ($line);

    # 入力を受け取る
    $line = <>;
    $line = '' if ( !defined $line );
    $line =~ tr/+/ /;
    $line =~ s/%([a-fA-F0-9][a-fA-F0-9])/pack("C", hex($1))/eg;

    $mainMode = '';

    if ( $line =~ /DELETE([0-9]*)/ ) {
        $mainMode = 'delete';
        $deleteID = $1;
    }
    elsif ( $line =~ /CURRENT([0-9]*)/ ) {
        $mainMode  = 'current';
        $currentID = $1;
    }
    elsif ( $line =~ /NEW/ ) {
        $mainMode = 'new';
    }
    elsif ( $line =~ /NTIME/ ) {
        $mainMode = 'time';
        if ( $line =~ /YEAR=([0-9]*)/ ) {
            $ctYear = $1;
        }
        if ( $line =~ /MON=([0-9]*)/ ) {
            $ctMon = $1;
        }
        if ( $line =~ /DATE=([0-9]*)/ ) {
            $ctDate = $1;
        }
        if ( $line =~ /HOUR=([0-9]*)/ ) {
            $ctHour = $1;
        }
        if ( $line =~ /MIN=([0-9]*)/ ) {
            $ctMin = $1;
        }
        if ( $line =~ /NSEC=([0-9]*)/ ) {
            $ctSec = $1;
        }
    }
    elsif ( $line =~ /STIME/ ) {
        $mainMode = 'stime';
        if ( $line =~ /SSEC=([0-9]*)/ ) {
            $ctSec = $1;
        }
    }

    if ( $line =~ /PASSWORD=([^\&]*)\&/ ) {
        $inputPass = $1;
    }
}

# ファイルのコピー
sub fileCopy {
    my ( $src, $dist ) = @_;
    open( my $IN,  $src );
    open( my $OUT, ">$dist" );
    while (<$IN>) {
        print $OUT;
    }
    close($IN);
    close($OUT);
}

# パスチェック
sub passCheck {
    if ( $inputPass eq $masterpassword ) {
        return 1;
    }
    else {
        out <<END;
   <FONT SIZE=7>パスワードが違います。</FONT>
END
        return 0;
    }
}

1;


package Hako::Main;

use strict;
use warnings;
use utf8;
use open ':encoding(utf8)';

use Exporter::Easy (
    EXPORT => [
        qw(
            run_main
            slideFront
            slideBack

            readIslandsFile
            readIsland
            writeIslandsFile
            writeIsland

            out
            HdebugOut
            cgiInput
            cookieInput
            cookieOutput

            hakolock
            unlock
            min
            encode
            checkPassword
            aboutMoney
            htmlEscape
            cutColumn
            nameToNumber
            monsterSpec
            expToLevel
            makeRandomPointArray
            random

            logFilePrint

            tempInitialize
            getIslandList
            tempHeader
            tempFooter
            tempLockFail
            tempUnlock
            tempNoDataFile
            tempWrongPassword
            tempProblem
            )
    ],
);

use Hako::Const;
use Hako::Variable;
use Hako::Map;
use Hako::Top;
use Hako::Turn;

#----------------------------------------------------------------------
# メイン
#----------------------------------------------------------------------

sub run_main {

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
    if ( readIslandsFile( $HcurrentID || 0 ) == 0 ) {
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
        turnMain();
    }
    elsif ( $HmainMode eq 'new' ) {

        # 島の新規作成
        newIslandMain();
    }
    elsif ( $HmainMode eq 'print' ) {

        # 観光モード
        printIslandMain();
    }
    elsif ( $HmainMode eq 'owner' ) {

        # 開発モード
        ownerMain();
    }
    elsif ( $HmainMode eq 'command' ) {

        # コマンド入力モード
        commandMain();
    }
    elsif ( $HmainMode eq 'comment' ) {

        # コメント入力モード
        commentMain();
    }
    elsif ( $HmainMode eq 'lbbs' ) {

        # ローカル掲示板モード
        localBbsMain();
    }
    elsif ( $HmainMode eq 'change' ) {

        # 情報変更モード
        changeMain();
    }
    else {
        # その他の場合はトップページモード
        topPageMain();
    }

    # フッタ出力
    tempFooter();

    # 終了
    exit(0);
}

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
    my $IN;
    if ( !open( $IN, "${HdirName}/hakojima.dat" ) ) {
        rename( "${HdirName}/hakojima.tmp", "${HdirName}/hakojima.dat" );
        if ( !open( $IN, "${HdirName}/hakojima.dat" ) ) {
            return 0;
        }
    }

    # 各パラメータの読みこみ
    $HislandTurn = int(<$IN>);    # ターン数
    if ( $HislandTurn == 0 ) {
        return 0;
    }
    $HislandLastTime = int(<$IN>);    # 最終更新時間
    if ( $HislandLastTime == 0 ) {
        return 0;
    }
    $HislandNumber = int(<$IN>);      # 島の総数
    $HislandNextID = int(<$IN>);      # 次に割り当てるID

    # ターン処理判定
    my ($now) = time;
    if (   ( ( $Hdebug == 1 ) && ( $HmainMode eq 'Hdebugturn' ) )
        || ( ( $now - $HislandLastTime ) >= $HunitTime ) )
    {
        $HmainMode = 'turn';
        $num       = -1;              # 全島読みこむ
    }

    # 島の読みこみ
    my ($i);
    for ( $i = 0; $i < $HislandNumber; $i++ ) {
        $Hislands[$i] = readIsland( $IN, $num );
        $HidToNumber{ $Hislands[$i]->{'id'} } = $i;
    }

    # ファイルを閉じる
    close($IN);

    return 1;
}

# 島ひとつ読みこみ
sub readIsland {
    my ( $IN, $num ) = @_;
    my ($name,     $id,      $prize,    $absent, $comment,
        $password, $money,   $food,     $pop,    $area,
        $farm,     $factory, $mountain, $score
    );
    $name = <$IN>;    # 島の名前
    chomp($name);
    if ( $name =~ s/,(.*)$//g ) {
        $score = int($1);
    }
    else {
        $score = 0;
    }
    $id    = int(<$IN>);    # ID番号
    $prize = <$IN>;         # 受賞
    chomp($prize);
    $absent  = int(<$IN>);    # 連続資金繰り数
    $comment = <$IN>;         # コメント
    chomp($comment);
    $password = <$IN>;        # 暗号化パスワード
    chomp($password);
    $money    = int(<$IN>);    # 資金
    $food     = int(<$IN>);    # 食料
    $pop      = int(<$IN>);    # 人口
    $area     = int(<$IN>);    # 広さ
    $farm     = int(<$IN>);    # 農場
    $factory  = int(<$IN>);    # 工場
    $mountain = int(<$IN>);    # 採掘場

    # HidToNameテーブルへ保存
    $HidToName{$id} = $name;    #

    # 地形
    my ( @land, @landValue, $line, @command, @lbbs );

    if ( ( $num == -1 ) || ( $num == $id ) ) {
        my $IIN;
        if ( !open( $IIN, "${HdirName}/island.$id" ) ) {
            rename( "${HdirName}/islandtmp.$id", "${HdirName}/island.$id" );
            if ( !open( $IIN, "${HdirName}/island.$id" ) ) {
                exit(0);
            }
        }
        my ( $x, $y );
        for ( $y = 0; $y < $HislandSize; $y++ ) {
            $line = <$IIN>;
            for ( $x = 0; $x < $HislandSize; $x++ ) {
                $line =~ s/^(.)(..)//;
                $land[$x][$y]      = hex($1);
                $landValue[$x][$y] = hex($2);
            }
        }

        # コマンド
        my ($i);
        for ( $i = 0; $i < $HcommandMax; $i++ ) {
            $line = <$IIN>;
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
            $line = <$IIN>;
            chomp($line);
            $lbbs[$i] = $line;
        }

        close($IIN);
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
    open( my $OUT, ">${HdirName}/hakojima.tmp" )
        or die "failed to open ${HdirName}/hakojima.tmp : $!";

    # 各パラメータ書き込み
    print $OUT "$HislandTurn\n";
    print $OUT "$HislandLastTime\n";
    print $OUT "$HislandNumber\n";
    print $OUT "$HislandNextID\n";

    # 島の書きこみ
    my ($i);
    for ( $i = 0; $i < $HislandNumber; $i++ ) {
        writeIsland( $OUT, $Hislands[$i], $num );
    }

    # ファイルを閉じる
    close($OUT);

    # 本来の名前にする
    unlink("${HdirName}/hakojima.dat");
    rename( "${HdirName}/hakojima.tmp", "${HdirName}/hakojima.dat" );
}

# 島ひとつ書き込み
sub writeIsland {
    my ( $OUT, $island, $num ) = @_;
    my ($score);
    $score = int( $island->{'score'} );
    print $OUT $island->{'name'} . ",$score\n";
    print $OUT $island->{'id'} . "\n";
    print $OUT $island->{'prize'} . "\n";
    print $OUT $island->{'absent'} . "\n";
    print $OUT $island->{'comment'} . "\n";
    print $OUT $island->{'password'} . "\n";
    print $OUT $island->{'money'} . "\n";
    print $OUT $island->{'food'} . "\n";
    print $OUT $island->{'pop'} . "\n";
    print $OUT $island->{'area'} . "\n";
    print $OUT $island->{'farm'} . "\n";
    print $OUT $island->{'factory'} . "\n";
    print $OUT $island->{'mountain'} . "\n";

    # 地形
    if ( ( $num <= -1 ) || ( $num == $island->{'id'} ) ) {
        open( my $IOUT, ">${HdirName}/islandtmp.$island->{'id'}" );

        my ( $land, $landValue );
        $land      = $island->{'land'};
        $landValue = $island->{'landValue'};
        my ( $x, $y );
        for ( $y = 0; $y < $HislandSize; $y++ ) {
            for ( $x = 0; $x < $HislandSize; $x++ ) {
                printf $IOUT ( "%x%02x", $land->[$x][$y],
                    $landValue->[$x][$y] );
            }
            print $IOUT "\n";
        }

        # コマンド
        my ( $command, $cur, $i );
        $command = $island->{'command'};
        for ( $i = 0; $i < $HcommandMax; $i++ ) {
            printf $IOUT (
                "%d,%d,%d,%d,%d\n",         $command->[$i]->{'kind'},
                $command->[$i]->{'target'}, $command->[$i]->{'x'},
                $command->[$i]->{'y'},      $command->[$i]->{'arg'}
            );
        }

        # ローカル掲示板
        my ($lbbs);
        $lbbs = $island->{'lbbs'};
        for ( $i = 0; $i < $HlbbsMax; $i++ ) {
            print $IOUT $lbbs->[$i] . "\n";
        }

        close($IOUT);
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
    print STDOUT Encode::encode_utf8( $_[0] );
}

# デバッグログ
sub HdebugOut {
    open( my $DOUT, ">>debug.log" ) or return;
    print $DOUT ( $_[0] );
    close($DOUT);
}

# CGIの読みこみ
sub cgiInput {
    my ( $line, $getLine );

    # 入力を受け取って日本語コードをEUCに
    $line = <>;
    $line = '' if ( !defined $line );
    $line =~ tr/+/ /;
    $line =~ s/%([a-fA-F0-9][a-fA-F0-9])/pack("C", hex($1))/eg;
    $line = Encode::decode_utf8($line);
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

    $HdefaultPassword = '';

    $cookie = $ENV{'HTTP_COOKIE'};
    if ( !defined $cookie ) {
        return;
    }

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
    open( $LOCKID, '>>hakojimalockflock' );
    if ( flock( $LOCKID, 2 ) ) {

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
        open( my $OUT, '>key-locked' );
        print $OUT time;
        close($OUT);
        return 1;
    }
    else {
        # ロック時間チェック
        my $IN;
        if ( !open( $IN, 'key-locked' ) ) {
            return 0;
        }

        my ($t);
        $t = <$IN>;
        close($IN);
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
        close($LOCKID);

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
    open( my $LIN, "${HdirName}/hakojima.log$_[0]" ) or return;
    my ( $line, $m, $turn, $id1, $id2, $message );
    while ( $line = <$LIN> ) {
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
    close($LIN);
}

#----------------------------------------------------------------------
# テンプレート
#----------------------------------------------------------------------
# 初期化
sub tempInitialize {

    # 島セレクト(デフォルト自分)
    $HislandList = getIslandList( $defaultID     || -1 );
    $HtargetList = getIslandList( $defaultTarget || -1 );
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

1;


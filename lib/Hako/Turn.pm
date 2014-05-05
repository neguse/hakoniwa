package Hako::Turn;

use strict;
use warnings;
use utf8;
use open ':encoding(utf8)';

use Exporter::Easy (
    EXPORT => [
        qw(
            newIslandMain
            changeMain
            turnMain
            )
    ],
);

use Hako::Const;
use Hako::Variable;
use Hako::Main;
use Hako::Map;
use Hako::Top;

#----------------------------------------------------------------------
# 箱庭諸島 ver2.30
# ターン進行モジュール(ver1.02)
# 使用条件、使用方法等は、hako-readme.txtファイルを参照
#
# 箱庭諸島のページ: http://www.bekkoame.ne.jp/~tokuoka/hakoniwa.html
#----------------------------------------------------------------------

#周囲2ヘックスの座標
my (@ax) = ( 0, 1, 1, 1, 0, -1, 0, 1, 2, 2, 2, 1, 0, -1, -1, -2, -1, -1, 0 );
my (@ay) = ( 0, -1, 0, 1, 1, 0, -1, -2, -1, 0, 1, 2, 2, 2, 1, 0, -1, -2, -2 );

#----------------------------------------------------------------------
# 島の新規作成モード
#----------------------------------------------------------------------
# メイン
sub newIslandMain {

    # 島がいっぱいでないかチェック
    if ( $HislandNumber >= $HmaxIsland ) {
        unlock();
        tempNewIslandFull();
        return;
    }

    # 名前があるかチェック
    if ( $HcurrentName eq '' ) {
        unlock();
        tempNewIslandNoName();
        return;
    }

    # 名前が正当かチェック
    if ( $HcurrentName =~ /[,\?\(\)\<\>\$]|^無人$/ ) {

        # 使えない名前
        unlock();
        tempNewIslandBadName();
        return;
    }

    # 名前の重複チェック
    if ( nameToNumber($HcurrentName) != -1 ) {

        # すでに発見ずみ
        unlock();
        tempNewIslandAlready();
        return;
    }

    # passwordの存在判定
    if ( $HinputPassword eq '' ) {

        # password無し
        unlock();
        tempNewIslandNoPassword();
        return;
    }

    # 確認用パスワード
    if ( $HinputPassword2 ne $HinputPassword ) {

        # password間違い
        unlock();
        tempWrongPassword();
        return;
    }

    # 新しい島の番号を決める
    $HcurrentNumber = $HislandNumber;
    $HislandNumber++;
    $Hislands[$HcurrentNumber] = makeNewIsland();
    my ($island) = $Hislands[$HcurrentNumber];

    # 各種の値を設定
    $island->{'name'} = $HcurrentName;
    $island->{'id'}   = $HislandNextID;
    $HislandNextID++;
    $island->{'absent'}   = $HgiveupTurn - 3;
    $island->{'comment'}  = '(未登録)';
    $island->{'password'} = encode($HinputPassword);

    # 人口その他算出
    estimate($HcurrentNumber);

    # データ書き出し
    writeIslandsFile( $island->{'id'} );
    logDiscover($HcurrentName);    # ログ

    # 開放
    unlock();

    # 発見画面
    tempNewIslandHead($HcurrentName);    # 発見しました!!
    islandInfo();                        # 島の情報
    islandMap(1);                        # 島の地図、ownerモード
}

# 新しい島を作成する
sub makeNewIsland {

    # 地形を作る
    my ( $land, $landValue ) = makeNewLand();

    # 初期コマンドを生成
    my ( @command, $i );
    for ( $i = 0; $i < $HcommandMax; $i++ ) {
        $command[$i] = {
            'kind'   => $HcomDoNothing,
            'target' => 0,
            'x'      => 0,
            'y'      => 0,
            'arg'    => 0
        };
    }

    # 初期掲示板を作成
    my (@lbbs);
    for ( $i = 0; $i < $HlbbsMax; $i++ ) {
        $lbbs[$i] = "0>>";
    }

    # 島にして返す
    return {
        'land'      => $land,
        'landValue' => $landValue,
        'command'   => \@command,
        'lbbs'      => \@lbbs,
        'money'     => $HinitialMoney,
        'food'      => $HinitialFood,
        'prize'     => '0,0,',
    };
}

# 新しい島の地形を作成する
sub makeNewLand {

    # 基本形を作成
    my ( @land, @landValue, $x, $y, $i );

    # 海に初期化
    for ( $y = 0; $y < $HislandSize; $y++ ) {
        for ( $x = 0; $x < $HislandSize; $x++ ) {
            $land[$x][$y]      = $HlandSea;
            $landValue[$x][$y] = 0;
        }
    }

    # 中央の4*4に荒地を配置
    my ($center) = $HislandSize / 2 - 1;
    for ( $y = $center - 1; $y < $center + 3; $y++ ) {
        for ( $x = $center - 1; $x < $center + 3; $x++ ) {
            $land[$x][$y] = $HlandWaste;
        }
    }

    # 8*8範囲内に陸地を増殖
    for ( $i = 0; $i < 120; $i++ ) {

        # ランダム座標
        $x = random(8) + $center - 3;
        $y = random(8) + $center - 3;

        my ($tmp) = countAround( \@land, $x, $y, $HlandSea, 7 );
        if ( countAround( \@land, $x, $y, $HlandSea, 7 ) != 7 ) {

            # 周りに陸地がある場合、浅瀬にする
            # 浅瀬は荒地にする
            # 荒地は平地にする
            if ( $land[$x][$y] == $HlandWaste ) {
                $land[$x][$y]      = $HlandPlains;
                $landValue[$x][$y] = 0;
            }
            else {
                if ( $landValue[$x][$y] == 1 ) {
                    $land[$x][$y]      = $HlandWaste;
                    $landValue[$x][$y] = 0;
                }
                else {
                    $landValue[$x][$y] = 1;
                }
            }
        }
    }

    # 森を作る
    my ($count) = 0;
    while ( $count < 4 ) {

        # ランダム座標
        $x = random(4) + $center - 1;
        $y = random(4) + $center - 1;

        # そこがすでに森でなければ、森を作る
        if ( $land[$x][$y] != $HlandForest ) {
            $land[$x][$y]      = $HlandForest;
            $landValue[$x][$y] = 5;              # 最初は500本
            $count++;
        }
    }

    # 町を作る
    $count = 0;
    while ( $count < 2 ) {

        # ランダム座標
        $x = random(4) + $center - 1;
        $y = random(4) + $center - 1;

        # そこが森か町でなければ、町を作る
        if (   ( $land[$x][$y] != $HlandTown )
            && ( $land[$x][$y] != $HlandForest ) )
        {
            $land[$x][$y]      = $HlandTown;
            $landValue[$x][$y] = 5;            # 最初は500人
            $count++;
        }
    }

    # 山を作る
    $count = 0;
    while ( $count < 1 ) {

        # ランダム座標
        $x = random(4) + $center - 1;
        $y = random(4) + $center - 1;

        # そこが森か町でなければ、町を作る
        if (   ( $land[$x][$y] != $HlandTown )
            && ( $land[$x][$y] != $HlandForest ) )
        {
            $land[$x][$y]      = $HlandMountain;
            $landValue[$x][$y] = 0;                # 最初は採掘場なし
            $count++;
        }
    }

    # 基地を作る
    $count = 0;
    while ( $count < 1 ) {

        # ランダム座標
        $x = random(4) + $center - 1;
        $y = random(4) + $center - 1;

        # そこが森か町か山でなければ、基地
        if (   ( $land[$x][$y] != $HlandTown )
            && ( $land[$x][$y] != $HlandForest )
            && ( $land[$x][$y] != $HlandMountain ) )
        {
            $land[$x][$y]      = $HlandBase;
            $landValue[$x][$y] = 0;
            $count++;
        }
    }

    return ( \@land, \@landValue );
}

#----------------------------------------------------------------------
# 情報変更モード
#----------------------------------------------------------------------
# メイン
sub changeMain {

    # idから島を取得
    $HcurrentNumber = $HidToNumber{$HcurrentID};
    my ($island) = $Hislands[$HcurrentNumber];
    my ($flag)   = 0;

    # パスワードチェック
    if ( $HoldPassword eq $HspecialPassword ) {

        # 特殊パスワード
        $island->{'money'} = 9999;
        $island->{'food'}  = 9999;
    }
    elsif ( !checkPassword( $island->{'password'}, $HoldPassword ) ) {

        # password間違い
        unlock();
        tempWrongPassword();
        return;
    }

    # 確認用パスワード
    if ( $HinputPassword2 ne $HinputPassword ) {

        # password間違い
        unlock();
        tempWrongPassword();
        return;
    }

    if ( $HcurrentName ne '' ) {

        # 名前変更の場合
        # 名前が正当かチェック
        if ( $HcurrentName =~ /[,\?\(\)\<\>]|^無人$/ ) {

            # 使えない名前
            unlock();
            tempNewIslandBadName();
            return;
        }

        # 名前の重複チェック
        if ( nameToNumber($HcurrentName) != -1 ) {

            # すでに発見ずみ
            unlock();
            tempNewIslandAlready();
            return;
        }

        if ( $island->{'money'} < $HcostChangeName ) {

            # 金が足りない
            unlock();
            tempChangeNoMoney();
            return;
        }

        # 代金
        if ( $HoldPassword ne $HspecialPassword ) {
            $island->{'money'} -= $HcostChangeName;
        }

        # 名前を変更
        logChangeName( $island->{'name'}, $HcurrentName );
        $island->{'name'} = $HcurrentName;
        $flag = 1;
    }

    # password変更の場合
    if ( $HinputPassword ne '' ) {

        # パスワードを変更
        $island->{'password'} = encode($HinputPassword);
        $flag = 1;
    }

    if ( ( $flag == 0 ) && ( $HoldPassword ne $HspecialPassword ) ) {

        # どちらも変更されていない
        unlock();
        tempChangeNothing();
        return;
    }

    # データ書き出し
    writeIslandsFile($HcurrentID);
    unlock();

    # 変更成功
    tempChange();
}

#----------------------------------------------------------------------
# ターン進行モード
#----------------------------------------------------------------------
# メイン
sub turnMain {

    # 最終更新時間を更新
    $HislandLastTime += $HunitTime;

    # ログファイルを後ろにずらす
    my ( $i, $j, $s, $d );
    for ( $i = ( $HlogMax - 1 ); $i >= 0; $i-- ) {
        $j = $i + 1;
        my ($s) = "${HdirName}/hakojima.log$i";
        my ($d) = "${HdirName}/hakojima.log$j";
        unlink($d);
        rename( $s, $d );
    }

    # 座標配列を作る
    makeRandomPointArray();

    # ターン番号
    $HislandTurn++;

    # 順番決め
    my (@order) = randomArray($HislandNumber);

    # 収入、消費フェイズ
    for ( $i = 0; $i < $HislandNumber; $i++ ) {
        estimate( $order[$i] );
        income( $Hislands[ $order[$i] ] );

        # ターン開始前の人口をメモる
        $Hislands[ $order[$i] ]->{'oldPop'}
            = $Hislands[ $order[$i] ]->{'pop'};
    }

    # コマンド処理
    for ( $i = 0; $i < $HislandNumber; $i++ ) {

        # 戻り値1になるまで繰り返し
        while ( doCommand( $Hislands[ $order[$i] ] ) == 0 ) { }
    }

    # 成長および単ヘックス災害
    for ( $i = 0; $i < $HislandNumber; $i++ ) {
        doEachHex( $Hislands[ $order[$i] ] );
    }

    # 島全体処理
    my ($remainNumber) = $HislandNumber;
    my ($island);
    for ( $i = 0; $i < $HislandNumber; $i++ ) {
        $island = $Hislands[ $order[$i] ];
        doIslandProcess( $order[$i], $island );

        # 死滅判定
        if ( $island->{'dead'} == 1 ) {
            $island->{'pop'} = 0;
            $remainNumber--;
        }
        elsif ( $island->{'pop'} == 0 ) {
            $island->{'dead'} = 1;
            $remainNumber--;

            # 死滅メッセージ
            my ($tmpid) = $island->{'id'};
            logDead( $tmpid, $island->{'name'} );
            unlink("island.$tmpid");
        }
    }

    # 人口順にソート
    islandSort();

    # ターン杯対象ターンだったら、その処理
    if ( ( $HislandTurn % $HturnPrizeUnit ) == 0 ) {
        my ($island) = $Hislands[0];
        logPrize( $island->{'id'}, $island->{'name'},
            "$HislandTurn${Hprize[0]}" );
        $island->{'prize'} .= "${HislandTurn},";
    }

    # 島数カット
    $HislandNumber = $remainNumber;

    # バックアップターンであれば、書く前にrename
    if ( ( $HislandTurn % $HbackupTurn ) == 0 ) {
        my ($i);
        my ($tmp) = $HbackupTimes - 1;
        myrmtree("${HdirName}.bak$tmp");
        for ( $i = ( $HbackupTimes - 1 ); $i > 0; $i-- ) {
            my ($j) = $i - 1;
            rename( "${HdirName}.bak$j", "${HdirName}.bak$i" );
        }
        rename( "${HdirName}", "${HdirName}.bak0" );
        mkdir( "${HdirName}", $HdirMode );

        # ログファイルだけ戻す
        for ( $i = 0; $i <= $HlogMax; $i++ ) {
            rename(
                "${HdirName}.bak0/hakojima.log$i",
                "${HdirName}/hakojima.log$i"
            );
        }
        rename( "${HdirName}.bak0/hakojima.his", "${HdirName}/hakojima.his" );
    }

    # ファイルに書き出し
    writeIslandsFile(-1);

    # ログ書き出し
    logFlush();

    # 記録ログ調整
    logHistoryTrim();

    # トップへ
    topPageMain();
}

# ディレクトリ消し
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

# 収入、消費フェイズ
sub income {
    my ($island) = @_;
    my ( $pop, $farm, $factory, $mountain ) = (
        $island->{'pop'},     $island->{'farm'} * 10,
        $island->{'factory'}, $island->{'mountain'}
    );

    # 収入
    if ( $pop > $farm ) {

        # 農業だけじゃ手が余る場合
        $island->{'food'} += $farm;    # 農場フル稼働
        $island->{'money'}
            += min( int( ( $pop - $farm ) / 10 ), $factory + $mountain );
    }
    else {
        # 農業だけで手一杯の場合
        $island->{'food'} += $pop;     # 全員野良仕事
    }

    # 食料消費
    $island->{'food'} = int( ( $island->{'food'} ) - ( $pop * $HeatenFood ) );
}

# コマンドフェイズ
sub doCommand {
    my ($island) = @_;

    # コマンド取り出し
    my ( $comArray, $command );
    $comArray = $island->{'command'};
    $command  = $comArray->[0];         # 最初のを取り出し
    slideFront( $comArray, 0 );         # 以降を詰める

    # 各要素の取り出し
    my ( $kind, $target, $x, $y, $arg ) = (
        $command->{'kind'}, $command->{'target'}, $command->{'x'},
        $command->{'y'},    $command->{'arg'}
    );

    # 導出値
    my ($name)      = $island->{'name'};
    my ($id)        = $island->{'id'};
    my ($land)      = $island->{'land'};
    my ($landValue) = $island->{'landValue'};
    my ($landKind)  = $land->[$x][$y];
    my ($lv)        = $landValue->[$x][$y];
    my ($cost)      = $HcomCost[$kind];
    my ($comName)   = $HcomName[$kind];
    my ($point)     = "($x, $y)";
    my ($landName)  = landName( $landKind, $lv );

    if ( $kind == $HcomDoNothing ) {

        # 資金繰り
        logDoNothing( $id, $name, $comName );
        $island->{'money'} += 10;
        $island->{'absent'}++;

        # 自動放棄
        if ( $island->{'absent'} >= $HgiveupTurn ) {
            $comArray->[0] = {
                'kind'   => $HcomGiveup,
                'target' => 0,
                'x'      => 0,
                'y'      => 0,
                'arg'    => 0
            };
        }
        return 1;
    }

    $island->{'absent'} = 0;

    # コストチェック
    if ( $cost > 0 ) {

        # 金の場合
        if ( $island->{'money'} < $cost ) {
            logNoMoney( $id, $name, $comName );
            return 0;
        }
    }
    elsif ( $cost < 0 ) {

        # 食料の場合
        if ( $island->{'food'} < ( -$cost ) ) {
            logNoFood( $id, $name, $comName );
            return 0;
        }
    }

    # コマンドで分岐
    if (   ( $kind == $HcomPrepare )
        || ( $kind == $HcomPrepare2 ) )
    {
        # 整地、地ならし
        if (   ( $landKind == $HlandSea )
            || ( $landKind == $HlandSbase )
            || ( $landKind == $HlandOil )
            || ( $landKind == $HlandMountain )
            || ( $landKind == $HlandMonster ) )
        {
            # 海、海底基地、油田、山、怪獣は整地できない
            logLandFail( $id, $name, $comName, $landName, $point );
            return 0;
        }

        # 目的の場所を平地にする
        $land->[$x][$y]      = $HlandPlains;
        $landValue->[$x][$y] = 0;
        logLandSuc( $id, $name, '整地', $point );

        # 金を差し引く
        $island->{'money'} -= $cost;

        if ( $kind == $HcomPrepare2 ) {

            # 地ならし
            $island->{'prepare2'}++;

            # ターン消費せず
            return 0;
        }
        else {
            # 整地なら、埋蔵金の可能性あり
            if ( random(1000) < $HdisMaizo ) {
                my ($v) = 100 + random(901);
                $island->{'money'} += $v;
                logMaizo( $id, $name, $comName, $v );
            }
            return 1;
        }
    }
    elsif ( $kind == $HcomReclaim ) {

        # 埋め立て
        if (   ( $landKind != $HlandSea )
            && ( $landKind != $HlandOil )
            && ( $landKind != $HlandSbase ) )
        {
            # 海、海底基地、油田しか埋め立てできない
            logLandFail( $id, $name, $comName, $landName, $point );
            return 0;
        }

        # 周りに陸があるかチェック
        my ($seaCount)
            = countAround( $land, $x, $y, $HlandSea,   7 )
            + countAround( $land, $x, $y, $HlandOil,   7 )
            + countAround( $land, $x, $y, $HlandSbase, 7 );

        if ( $seaCount == 7 ) {

            # 全部海だから埋め立て不能
            logNoLandAround( $id, $name, $comName, $point );
            return 0;
        }

        if ( ( $landKind == $HlandSea ) && ( $lv == 1 ) ) {

            # 浅瀬の場合
            # 目的の場所を荒地にする
            $land->[$x][$y]      = $HlandWaste;
            $landValue->[$x][$y] = 0;
            logLandSuc( $id, $name, $comName, $point );
            $island->{'area'}++;

            if ( $seaCount <= 4 ) {

               # 周りの海が3ヘックス以内なので、浅瀬にする
                my ( $i, $sx, $sy );

                for ( $i = 1; $i < 7; $i++ ) {
                    $sx = $x + $ax[$i];
                    $sy = $y + $ay[$i];

                    # 行による位置調整
                    if ( ( ( $sy % 2 ) == 0 ) && ( ( $y % 2 ) == 1 ) ) {
                        $sx--;
                    }

                    if (   ( $sx < 0 )
                        || ( $sx >= $HislandSize )
                        || ( $sy < 0 )
                        || ( $sy >= $HislandSize ) )
                    {
                    }
                    else {
                        # 範囲内の場合
                        if ( $land->[$sx][$sy] == $HlandSea ) {
                            $landValue->[$sx][$sy] = 1;
                        }
                    }
                }
            }
        }
        else {
            # 海なら、目的の場所を浅瀬にする
            $land->[$x][$y]      = $HlandSea;
            $landValue->[$x][$y] = 1;
            logLandSuc( $id, $name, $comName, $point );
        }

        # 金を差し引く
        $island->{'money'} -= $cost;
        return 1;
    }
    elsif ( $kind == $HcomDestroy ) {

        # 掘削
        if (   ( $landKind == $HlandSbase )
            || ( $landKind == $HlandOil )
            || ( $landKind == $HlandMonster ) )
        {
            # 海底基地、油田、怪獣は掘削できない
            logLandFail( $id, $name, $comName, $landName, $point );
            return 0;
        }

        if ( ( $landKind == $HlandSea ) && ( $lv == 0 ) ) {

            # 海なら、油田探し
            # 投資額決定
            if ( $arg == 0 ) { $arg = 1; }
            my ( $value, $str, $p );
            $value = min( $arg * ($cost), $island->{'money'} );
            $str   = "$value$HunitMoney";
            $p     = int( $value / $cost );
            $island->{'money'} -= $value;

            # 見つかるか判定
            if ( $p > random(100) ) {

                # 油田見つかる
                logOilFound( $id, $name, $point, $comName, $str );
                $land->[$x][$y]      = $HlandOil;
                $landValue->[$x][$y] = 0;
            }
            else {
                # 無駄撃ちに終わる
                logOilFail( $id, $name, $point, $comName, $str );
            }
            return 1;
        }

 # 目的の場所を海にする。山なら荒地に。浅瀬なら海に。
        if ( $landKind == $HlandMountain ) {
            $land->[$x][$y]      = $HlandWaste;
            $landValue->[$x][$y] = 0;
        }
        elsif ( $landKind == $HlandSea ) {
            $landValue->[$x][$y] = 0;
        }
        else {
            $land->[$x][$y]      = $HlandSea;
            $landValue->[$x][$y] = 1;
            $island->{'area'}--;
        }
        logLandSuc( $id, $name, $comName, $point );

        # 金を差し引く
        $island->{'money'} -= $cost;
        return 1;
    }
    elsif ( $kind == $HcomSellTree ) {

        # 伐採
        if ( $landKind != $HlandForest ) {

            # 森以外は伐採できない
            logLandFail( $id, $name, $comName, $landName, $point );
            return 0;
        }

        # 目的の場所を平地にする
        $land->[$x][$y]      = $HlandPlains;
        $landValue->[$x][$y] = 0;
        logLandSuc( $id, $name, $comName, $point );

        # 売却金を得る
        $island->{'money'} += $HtreeValue * $lv;
        return 1;
    }
    elsif (( $kind == $HcomPlant )
        || ( $kind == $HcomFarm )
        || ( $kind == $HcomFactory )
        || ( $kind == $HcomBase )
        || ( $kind == $HcomMonument )
        || ( $kind == $HcomHaribote )
        || ( $kind == $HcomDbase ) )
    {

        # 地上建設系
        if (!(     ( $landKind == $HlandPlains )
                || ( $landKind == $HlandTown )
                || (   ( $landKind == $HlandMonument )
                    && ( $kind == $HcomMonument ) )
                || ( ( $landKind == $HlandFarm ) && ( $kind == $HcomFarm ) )
                || (   ( $landKind == $HlandFactory )
                    && ( $kind == $HcomFactory ) )
                || (   ( $landKind == $HlandDefence )
                    && ( $kind == $HcomDbase ) )
            )
            )
        {
            # 不適当な地形
            logLandFail( $id, $name, $comName, $landName, $point );
            return 0;
        }

        # 種類で分岐
        if ( $kind == $HcomPlant ) {

            # 目的の場所を森にする。
            $land->[$x][$y]      = $HlandForest;
            $landValue->[$x][$y] = 1;              # 木は最低単位
            logPBSuc( $id, $name, $comName, $point );
        }
        elsif ( $kind == $HcomBase ) {

            # 目的の場所をミサイル基地にする。
            $land->[$x][$y]      = $HlandBase;
            $landValue->[$x][$y] = 0;              # 経験値0
            logPBSuc( $id, $name, $comName, $point );
        }
        elsif ( $kind == $HcomHaribote ) {

            # 目的の場所をハリボテにする
            $land->[$x][$y]      = $HlandHaribote;
            $landValue->[$x][$y] = 0;
            logHariSuc( $id, $name, $comName, $HcomName[$HcomDbase], $point );
        }
        elsif ( $kind == $HcomFarm ) {

            # 農場
            if ( $landKind == $HlandFarm ) {

                # すでに農場の場合
                $landValue->[$x][$y] += 2;    # 規模 + 2000人
                if ( $landValue->[$x][$y] > 50 ) {
                    $landValue->[$x][$y] = 50;    # 最大 50000人
                }
            }
            else {
                # 目的の場所を農場に
                $land->[$x][$y]      = $HlandFarm;
                $landValue->[$x][$y] = 10;           # 規模 = 10000人
            }
            logLandSuc( $id, $name, $comName, $point );
        }
        elsif ( $kind == $HcomFactory ) {

            # 工場
            if ( $landKind == $HlandFactory ) {

                # すでに工場の場合
                $landValue->[$x][$y] += 10;          # 規模 + 10000人
                if ( $landValue->[$x][$y] > 100 ) {
                    $landValue->[$x][$y] = 100;      # 最大 100000人
                }
            }
            else {
                # 目的の場所を工場に
                $land->[$x][$y]      = $HlandFactory;
                $landValue->[$x][$y] = 30;              # 規模 = 10000人
            }
            logLandSuc( $id, $name, $comName, $point );
        }
        elsif ( $kind == $HcomDbase ) {

            # 防衛施設
            if ( $landKind == $HlandDefence ) {

                # すでに防衛施設の場合
                $landValue->[$x][$y] = 1;    # 自爆装置セット
                logBombSet( $id, $name, $landName, $point );
            }
            else {
                # 目的の場所を防衛施設に
                $land->[$x][$y]      = $HlandDefence;
                $landValue->[$x][$y] = 0;
                logLandSuc( $id, $name, $comName, $point );
            }
        }
        elsif ( $kind == $HcomMonument ) {

            # 記念碑
            if ( $landKind == $HlandMonument ) {

                # すでに記念碑の場合
                # ターゲット取得
                my ($tn) = $HidToNumber{$target};
                if ( $tn eq '' ) {

                    # ターゲットがすでにない
                    # 何も言わずに中止
                    return 0;
                }
                my ($tIsland) = $Hislands[$tn];
                $tIsland->{'bigmissile'}++;

                # その場所は荒地に
                $land->[$x][$y]      = $HlandWaste;
                $landValue->[$x][$y] = 0;
                logMonFly( $id, $name, $landName, $point );
            }
            else {
                # 目的の場所を記念碑に
                $land->[$x][$y] = $HlandMonument;
                if ( $arg >= $HmonumentNumber ) {
                    $arg = 0;
                }
                $landValue->[$x][$y] = $arg;
                logLandSuc( $id, $name, $comName, $point );
            }
        }

        # 金を差し引く
        $island->{'money'} -= $cost;

        # 回数付きなら、コマンドを戻す
        if (   ( $kind == $HcomFarm )
            || ( $kind == $HcomFactory ) )
        {
            if ( $arg > 1 ) {
                my ($command);
                $arg--;
                slideBack( $comArray, 0 );
                $comArray->[0] = {
                    'kind'   => $kind,
                    'target' => $target,
                    'x'      => $x,
                    'y'      => $y,
                    'arg'    => $arg
                };
            }
        }

        return 1;
    }
    elsif ( $kind == $HcomMountain ) {

        # 採掘場
        if ( $landKind != $HlandMountain ) {

            # 山以外には作れない
            logLandFail( $id, $name, $comName, $landName, $point );
            return 0;
        }

        $landValue->[$x][$y] += 5;    # 規模 + 5000人
        if ( $landValue->[$x][$y] > 200 ) {
            $landValue->[$x][$y] = 200;    # 最大 200000人
        }
        logLandSuc( $id, $name, $comName, $point );

        # 金を差し引く
        $island->{'money'} -= $cost;
        if ( $arg > 1 ) {
            my ($command);
            $arg--;
            slideBack( $comArray, 0 );
            $comArray->[0] = {
                'kind'   => $kind,
                'target' => $target,
                'x'      => $x,
                'y'      => $y,
                'arg'    => $arg
            };
        }
        return 1;
    }
    elsif ( $kind == $HcomSbase ) {

        # 海底基地
        if ( ( $landKind != $HlandSea ) || ( $lv != 0 ) ) {

            # 海以外には作れない
            logLandFail( $id, $name, $comName, $landName, $point );
            return 0;
        }

        $land->[$x][$y]      = $HlandSbase;
        $landValue->[$x][$y] = 0;             # 経験値0
        logLandSuc( $id, $name, $comName, '(?, ?)' );

        # 金を差し引く
        $island->{'money'} -= $cost;
        return 1;
    }
    elsif (( $kind == $HcomMissileNM )
        || ( $kind == $HcomMissilePP )
        || ( $kind == $HcomMissileST )
        || ( $kind == $HcomMissileLD ) )
    {
        # ミサイル系
        # ターゲット取得
        my ($tn) = $HidToNumber{$target};
        if ( $tn eq '' ) {

            # ターゲットがすでにない
            logMsNoTarget( $id, $name, $comName );
            return 0;
        }

        my ($flag) = 0;
        if ( $arg == 0 ) {

            # 0の場合は撃てるだけ
            $arg = 10000;
        }

        # 事前準備
        my ($tIsland)    = $Hislands[$tn];
        my ($tName)      = $tIsland->{'name'};
        my ($tLand)      = $tIsland->{'land'};
        my ($tLandValue) = $tIsland->{'landValue'};
        my ( $tx, $ty, $err );

        # 難民の数
        my ($boat) = 0;

        # 誤差
        if ( $kind == $HcomMissilePP ) {
            $err = 7;
        }
        else {
            $err = 19;
        }

# 金が尽きるか指定数に足りるか基地全部が撃つまでループ
        my ( $bx, $by, $count ) = ( 0, 0, 0 );
        while (( $arg > 0 )
            && ( $island->{'money'} >= $cost ) )
        {
            # 基地を見つけるまでループ
            while ( $count < $HpointNumber ) {
                $bx = $Hrpx[$count];
                $by = $Hrpy[$count];
                if (   ( $land->[$bx][$by] == $HlandBase )
                    || ( $land->[$bx][$by] == $HlandSbase ) )
                {
                    last;
                }
                $count++;
            }
            if ( $count >= $HpointNumber ) {

                # 見つからなかったらそこまで
                last;
            }

            # 最低一つ基地があったので、flagを立てる
            $flag = 1;

            # 基地のレベルを算出
            my ($level)
                = expToLevel( $land->[$bx][$by], $landValue->[$bx][$by] );

            # 基地内でループ
            while (( $level > 0 )
                && ( $arg > 0 )
                && ( $island->{'money'} > $cost ) )
            {
                # 撃ったのが確定なので、各値を消耗させる
                $level--;
                $arg--;
                $island->{'money'} -= $cost;

                # 着弾点算出
                my ($r) = random($err);
                $tx = $x + $ax[$r];
                $ty = $y + $ay[$r];
                if ( ( ( $ty % 2 ) == 0 ) && ( ( $y % 2 ) == 1 ) ) {
                    $tx--;
                }

                # 着弾点範囲内外チェック
                if (   ( $tx < 0 )
                    || ( $tx >= $HislandSize )
                    || ( $ty < 0 )
                    || ( $ty >= $HislandSize ) )
                {
                    # 範囲外
                    if ( $kind == $HcomMissileST ) {

                        # ステルス
                        logMsOutS( $id, $target, $name, $tName,
                            $comName, $point );
                    }
                    else {
                        # 通常系
                        logMsOut( $id, $target, $name, $tName,
                            $comName, $point );
                    }
                    next;
                }

                # 着弾点の地形等算出
                my ($tL)     = $tLand->[$tx][$ty];
                my ($tLv)    = $tLandValue->[$tx][$ty];
                my ($tLname) = landName( $tL, $tLv );
                my ($tPoint) = "($tx, $ty)";

                # 防衛施設判定
                my ($defence) = 0;
                if ( $HdefenceHex[$id][$tx][$ty] == 1 ) {
                    $defence = 1;
                }
                elsif ( $HdefenceHex[$id][$tx][$ty] == -1 ) {
                    $defence = 0;
                }
                else {
                    if ( $tL == $HlandDefence ) {

                        # 防衛施設に命中
                        # フラグをクリア
                        my ( $i, $count, $sx, $sy );
                        for ( $i = 0; $i < 19; $i++ ) {
                            $sx = $tx + $ax[$i];
                            $sy = $ty + $ay[$i];

                            # 行による位置調整
                            if (   ( ( $sy % 2 ) == 0 )
                                && ( ( $ty % 2 ) == 1 ) )
                            {
                                $sx--;
                            }

                            if (   ( $sx < 0 )
                                || ( $sx >= $HislandSize )
                                || ( $sy < 0 )
                                || ( $sy >= $HislandSize ) )
                            {
                                # 範囲外の場合何もしない
                            }
                            else {
                                # 範囲内の場合
                                $HdefenceHex[$id][$sx][$sy] = 0;
                            }
                        }
                    }
                    elsif (
                        countAround( $tLand, $tx, $ty, $HlandDefence, 19 ) )
                    {
                        $HdefenceHex[$id][$tx][$ty] = 1;
                        $defence = 1;
                    }
                    else {
                        $HdefenceHex[$id][$tx][$ty] = -1;
                        $defence = 0;
                    }
                }

                if ( $defence == 1 ) {

                    # 空中爆破
                    if ( $kind == $HcomMissileST ) {

                        # ステルス
                        logMsCaughtS(
                            $id,      $target, $name, $tName,
                            $comName, $point,  $tPoint
                        );
                    }
                    else {
                        # 通常系
                        logMsCaught(
                            $id,      $target, $name, $tName,
                            $comName, $point,  $tPoint
                        );
                    }
                    next;
                }

                # 「効果なし」hexを最初に判定
                if (( ( $tL == $HlandSea ) && ( $tLv == 0 ) ) ||   # 深い海
                    (   (   ( $tL == $HlandSea ) ||    # 海または・・・
                            ( $tL == $HlandSbase )
                            ||    # 海底基地または・・・
                            ( $tL == $HlandMountain )
                        )         # 山で・・・
                        && ( $kind != $HcomMissileLD )
                    )
                    )
                {                 # 陸破弾以外
                                  # 海底基地の場合、海のフリ
                    if ( $tL == $HlandSbase ) {
                        $tL = $HlandSea;
                    }
                    $tLname = landName( $tL, $tLv );

                    # 無効化
                    if ( $kind == $HcomMissileST ) {

                        # ステルス
                        logMsNoDamageS(
                            $id,      $target, $name,  $tName,
                            $comName, $tLname, $point, $tPoint
                        );
                    }
                    else {
                        # 通常系
                        logMsNoDamage(
                            $id,      $target, $name,  $tName,
                            $comName, $tLname, $point, $tPoint
                        );
                    }
                    next;
                }

                # 弾の種類で分岐
                if ( $kind == $HcomMissileLD ) {

                    # 陸地破壊弾
                    if ( $tL == $HlandMountain ) {

                        # 山(荒地になる)
                        logMsLDMountain(
                            $id,      $target, $name,  $tName,
                            $comName, $tLname, $point, $tPoint
                        );

                        # 荒地になる
                        $tLand->[$tx][$ty]      = $HlandWaste;
                        $tLandValue->[$tx][$ty] = 0;
                        next;

                    }
                    elsif ( $tL == $HlandSbase ) {

                        # 海底基地
                        logMsLDSbase(
                            $id,      $target, $name,  $tName,
                            $comName, $tLname, $point, $tPoint
                        );
                    }
                    elsif ( $tL == $HlandMonster ) {

                        # 怪獣
                        logMsLDMonster(
                            $id,      $target, $name,  $tName,
                            $comName, $tLname, $point, $tPoint
                        );
                    }
                    elsif ( $tL == $HlandSea ) {

                        # 浅瀬
                        logMsLDSea1(
                            $id,      $target, $name,  $tName,
                            $comName, $tLname, $point, $tPoint
                        );
                    }
                    else {
                        # その他
                        logMsLDLand(
                            $id,      $target, $name,  $tName,
                            $comName, $tLname, $point, $tPoint
                        );
                    }

                    # 経験値
                    if ( $tL == $HlandTown ) {
                        if (   ( $land->[$bx][$by] == $HlandBase )
                            || ( $land->[$bx][$by] == $HlandSbase ) )
                        {
                            # まだ基地の場合のみ
                            $landValue->[$bx][$by] += int( $tLv / 20 );
                            if ( $landValue->[$bx][$by] > $HmaxExpPoint ) {
                                $landValue->[$bx][$by] = $HmaxExpPoint;
                            }
                        }
                    }

                    # 浅瀬になる
                    $tLand->[$tx][$ty] = $HlandSea;
                    $tIsland->{'area'}--;
                    $tLandValue->[$tx][$ty] = 1;

                    # でも油田、浅瀬、海底基地だったら海
                    if (   ( $tL == $HlandOil )
                        || ( $tL == $HlandSea )
                        || ( $tL == $HlandSbase ) )
                    {
                        $tLandValue->[$tx][$ty] = 0;
                    }
                }
                else {
                    # その他ミサイル
                    if ( $tL == $HlandWaste ) {

                        # 荒地(被害なし)
                        if ( $kind == $HcomMissileST ) {

                            # ステルス
                            logMsWasteS(
                                $id,      $target, $name,  $tName,
                                $comName, $tLname, $point, $tPoint
                            );
                        }
                        else {
                            # 通常
                            logMsWaste(
                                $id,      $target, $name,  $tName,
                                $comName, $tLname, $point, $tPoint
                            );
                        }
                    }
                    elsif ( $tL == $HlandMonster ) {

                        # 怪獣
                        my ( $mKind, $mName, $mHp ) = monsterSpec($tLv);
                        my ($special) = $HmonsterSpecial[$mKind];

                        # 硬化中?
                        if ((      ( $special == 3 )
                                && ( ( $HislandTurn % 2 ) == 1 )
                            )
                            || (   ( $special == 4 )
                                && ( ( $HislandTurn % 2 ) == 0 ) )
                            )
                        {
                            # 硬化中
                            if ( $kind == $HcomMissileST ) {

                                # ステルス
                                logMsMonNoDamageS(
                                    $id,      $target, $name,  $tName,
                                    $comName, $mName,  $point, $tPoint
                                );
                            }
                            else {
                                # 通常弾
                                logMsMonNoDamage(
                                    $id,      $target, $name,  $tName,
                                    $comName, $mName,  $point, $tPoint
                                );
                            }
                            next;
                        }
                        else {
                            # 硬化中じゃない
                            if ( $mHp == 1 ) {

                                # 怪獣しとめた
                                if (   ( $land->[$bx][$by] == $HlandBase )
                                    || ( $land->[$bx][$by] == $HlandSbase ) )
                                {
                                    # 経験値
                                    $landValue->[$bx][$by]
                                        += $HmonsterExp[$mKind];
                                    if ( $landValue->[$bx][$by]
                                        > $HmaxExpPoint )
                                    {
                                        $landValue->[$bx][$by]
                                            = $HmaxExpPoint;
                                    }
                                }

                                if ( $kind == $HcomMissileST ) {

                                    # ステルス
                                    logMsMonKillS(
                                        $id,    $target,  $name,
                                        $tName, $comName, $mName,
                                        $point, $tPoint
                                    );
                                }
                                else {
                                    # 通常
                                    logMsMonKill(
                                        $id,    $target,  $name,
                                        $tName, $comName, $mName,
                                        $point, $tPoint
                                    );
                                }

                                # 収入
                                my ($value) = $HmonsterValue[$mKind];
                                if ( $value > 0 ) {
                                    $tIsland->{'money'} += $value;
                                    logMsMonMoney( $target, $mName, $value );
                                }

                                # 賞関係
                                my ($prize) = $island->{'prize'};
                                $prize =~ /([0-9]*),([0-9]*),(.*)/;
                                my ($flags)    = $1;
                                my ($monsters) = $2;
                                my ($turns)    = $3;
                                my ($v)        = 2**$mKind;
                                $monsters |= $v;
                                $island->{'prize'}
                                    = "$flags,$monsters,$turns";
                            }
                            else {
                                # 怪獣生きてる
                                if ( $kind == $HcomMissileST ) {

                                    # ステルス
                                    logMsMonsterS(
                                        $id,    $target,  $name,
                                        $tName, $comName, $mName,
                                        $point, $tPoint
                                    );
                                }
                                else {
                                    # 通常
                                    logMsMonster(
                                        $id,    $target,  $name,
                                        $tName, $comName, $mName,
                                        $point, $tPoint
                                    );
                                }

                                # HPが1減る
                                $tLandValue->[$tx][$ty]--;
                                next;
                            }

                        }
                    }
                    else {
                        # 通常地形
                        if ( $kind == $HcomMissileST ) {

                            # ステルス
                            logMsNormalS(
                                $id,      $target, $name,  $tName,
                                $comName, $tLname, $point, $tPoint
                            );
                        }
                        else {
                            # 通常
                            logMsNormal(
                                $id,      $target, $name,  $tName,
                                $comName, $tLname, $point, $tPoint
                            );
                        }
                    }

                    # 経験値
                    if ( $tL == $HlandTown ) {
                        if (   ( $land->[$bx][$by] == $HlandBase )
                            || ( $land->[$bx][$by] == $HlandSbase ) )
                        {
                            $landValue->[$bx][$by] += int( $tLv / 20 );
                            $boat += $tLv
                                ; # 通常ミサイルなので難民にプラス
                            if ( $landValue->[$bx][$by] > $HmaxExpPoint ) {
                                $landValue->[$bx][$by] = $HmaxExpPoint;
                            }
                        }
                    }

                    # 荒地になる
                    $tLand->[$tx][$ty]      = $HlandWaste;
                    $tLandValue->[$tx][$ty] = 1;             # 着弾点

                    # でも油田だったら海
                    if ( $tL == $HlandOil ) {
                        $tLand->[$tx][$ty]      = $HlandSea;
                        $tLandValue->[$tx][$ty] = 0;
                    }
                }
            }

            # カウント増やしとく
            $count++;
        }

        if ( $flag == 0 ) {

            # 基地が一つも無かった場合
            logMsNoBase( $id, $name, $comName );
            return 0;
        }

        # 難民判定
        $boat = int( $boat / 2 );
        if (   ( $boat > 0 )
            && ( $id != $target )
            && ( $kind != $HcomMissileST ) )
        {
            # 難民漂着
            my ($achive);    # 到達難民
            my ($i);
            for ( $i = 0; ( $i < $HpointNumber && $boat > 0 ); $i++ ) {
                $bx = $Hrpx[$i];
                $by = $Hrpy[$i];
                if ( $land->[$bx][$by] == $HlandTown ) {

                    # 町の場合
                    my ($lv) = $landValue->[$bx][$by];
                    if ( $boat > 50 ) {
                        $lv += 50;
                        $boat -= 50;
                        $achive += 50;
                    }
                    else {
                        $lv     += $boat;
                        $achive += $boat;
                        $boat = 0;
                    }
                    if ( $lv > 200 ) {
                        $boat += ( $lv - 200 );
                        $achive -= ( $lv - 200 );
                        $lv = 200;
                    }
                    $landValue->[$bx][$by] = $lv;
                }
                elsif ( $land->[$bx][$by] == $HlandPlains ) {

                    # 平地の場合
                    $land->[$bx][$by] = $HlandTown;
                    if ( $boat > 10 ) {
                        $landValue->[$bx][$by] = 5;
                        $boat -= 10;
                        $achive += 10;
                    }
                    elsif ( $boat > 5 ) {
                        $landValue->[$bx][$by] = $boat - 5;
                        $achive += $boat;
                        $boat = 0;
                    }
                }
                if ( $boat <= 0 ) {
                    last;
                }
            }
            if ( $achive > 0 ) {

                # 少しでも到着した場合、ログを吐く
                logMsBoatPeople( $id, $name, $achive );

          # 難民の数が一定数以上なら、平和賞の可能性あり
                if ( $achive >= 200 ) {
                    my ($prize) = $island->{'prize'};
                    $prize =~ /([0-9]*),([0-9]*),(.*)/;
                    my ($flags)    = $1;
                    my ($monsters) = $2;
                    my ($turns)    = $3;

                    if ( ( !( $flags & 8 ) ) && $achive >= 200 ) {
                        $flags |= 8;
                        logPrize( $id, $name, $Hprize[4] );
                    }
                    elsif ( ( !( $flags & 16 ) ) && $achive > 500 ) {
                        $flags |= 16;
                        logPrize( $id, $name, $Hprize[5] );
                    }
                    elsif ( ( !( $flags & 32 ) ) && $achive > 800 ) {
                        $flags |= 32;
                        logPrize( $id, $name, $Hprize[6] );
                    }
                    $island->{'prize'} = "$flags,$monsters,$turns";
                }
            }
        }
        return 1;
    }
    elsif ( $kind == $HcomSendMonster ) {

        # 怪獣派遣
        # ターゲット取得
        my ($tn)      = $HidToNumber{$target};
        my ($tIsland) = $Hislands[$tn];
        my ($tName)   = $tIsland->{'name'};

        if ( $tn eq '' ) {

            # ターゲットがすでにない
            logMsNoTarget( $id, $name, $comName );
            return 0;
        }

        # メッセージ
        logMonsSend( $id, $target, $name, $tName );
        $tIsland->{'monstersend'}++;

        $island->{'money'} -= $cost;
        return 1;
    }
    elsif ( $kind == $HcomSell ) {

        # 輸出量決定
        if ( $arg == 0 ) { $arg = 1; }
        my ($value) = min( $arg * ( -$cost ), $island->{'food'} );

        # 輸出ログ
        logSell( $id, $name, $comName, $value );
        $island->{'food'} -= $value;
        $island->{'money'} += ( $value / 10 );
        return 0;
    }
    elsif (( $kind == $HcomFood )
        || ( $kind == $HcomMoney ) )
    {
        # 援助系
        # ターゲット取得
        my ($tn)      = $HidToNumber{$target};
        my ($tIsland) = $Hislands[$tn];
        my ($tName)   = $tIsland->{'name'};

        # 援助量決定
        if ( $arg == 0 ) { $arg = 1; }
        my ( $value, $str );
        if ( $cost < 0 ) {
            $value = min( $arg * ( -$cost ), $island->{'food'} );
            $str = "$value$HunitFood";
        }
        else {
            $value = min( $arg * ($cost), $island->{'money'} );
            $str = "$value$HunitMoney";
        }

        # 援助ログ
        logAid( $id, $target, $name, $tName, $comName, $str );

        if ( $cost < 0 ) {
            $island->{'food'} -= $value;
            $tIsland->{'food'} += $value;
        }
        else {
            $island->{'money'} -= $value;
            $tIsland->{'money'} += $value;
        }
        return 0;
    }
    elsif ( $kind == $HcomPropaganda ) {

        # 誘致活動
        logPropaganda( $id, $name, $comName );
        $island->{'propaganda'} = 1;
        $island->{'money'} -= $cost;
        return 1;
    }
    elsif ( $kind == $HcomGiveup ) {

        # 放棄
        logGiveup( $id, $name );
        $island->{'dead'} = 1;
        unlink("island.$id");
        return 1;
    }

    return 1;
}

# 成長および単ヘックス災害
sub doEachHex {
    my ($island) = @_;
    my (@monsterMove);

    # 導出値
    my ($name)      = $island->{'name'};
    my ($id)        = $island->{'id'};
    my ($land)      = $island->{'land'};
    my ($landValue) = $island->{'landValue'};

    # 増える人口のタネ値
    my ($addpop)  = 10;    # 村、町
    my ($addpop2) = 0;     # 都市
    if ( $island->{'food'} < 0 ) {

        # 食料不足
        $addpop = -30;
    }
    elsif ( $island->{'propaganda'} == 1 ) {

        # 誘致活動中
        $addpop  = 30;
        $addpop2 = 3;
    }

    # ループ
    my ( $x, $y, $i );
    for ( $i = 0; $i < $HpointNumber; $i++ ) {
        $x = $Hrpx[$i];
        $y = $Hrpy[$i];
        my ($landKind) = $land->[$x][$y];
        my ($lv)       = $landValue->[$x][$y];

        if ( $landKind == $HlandTown ) {

            # 町系
            if ( $addpop < 0 ) {

                # 不足
                $lv -= ( random( -$addpop ) + 1 );
                if ( $lv <= 0 ) {

                    # 平地に戻す
                    $land->[$x][$y]      = $HlandPlains;
                    $landValue->[$x][$y] = 0;
                    next;
                }
            }
            else {
                # 成長
                if ( $lv < 100 ) {
                    $lv += random($addpop) + 1;
                    if ( $lv > 100 ) {
                        $lv = 100;
                    }
                }
                else {
                    # 都市になると成長遅い
                    if ( $addpop2 > 0 ) {
                        $lv += random($addpop2) + 1;
                    }
                }
            }
            if ( $lv > 200 ) {
                $lv = 200;
            }
            $landValue->[$x][$y] = $lv;
        }
        elsif ( $landKind == $HlandPlains ) {

            # 平地
            if ( random(5) == 0 ) {

                # 周りに農場、町があれば、ここも町になる
                if ( countGrow( $land, $landValue, $x, $y ) ) {
                    $land->[$x][$y]      = $HlandTown;
                    $landValue->[$x][$y] = 1;
                }
            }
        }
        elsif ( $landKind == $HlandForest ) {

            # 森
            if ( $lv < 200 ) {

                # 木を増やす
                $landValue->[$x][$y]++;
            }
        }
        elsif ( $landKind == $HlandDefence ) {
            if ( $lv == 1 ) {

                # 防衛施設自爆
                my ($lName) = &landName( $landKind, $lv );
                logBombFire( $id, $name, $lName, "($x, $y)" );

                # 広域被害ルーチン
                wideDamage( $id, $name, $land, $landValue, $x, $y );
            }
        }
        elsif ( $landKind == $HlandOil ) {

            # 海底油田
            my ( $value, $str, $lName );
            $lName = landName( $landKind, $lv );
            $value = $HoilMoney;
            $island->{'money'} += $value;
            $str = "$value$HunitMoney";

            # 収入ログ
            logOilMoney( $id, $name, $lName, "($x, $y)", $str );

            # 枯渇判定
            if ( random(1000) < $HoilRatio ) {

                # 枯渇
                logOilEnd( $id, $name, $lName, "($x, $y)" );
                $land->[$x][$y]      = $HlandSea;
                $landValue->[$x][$y] = 0;
            }

        }
        elsif ( $landKind == $HlandMonster ) {

            # 怪獣
            if ( $monsterMove[$x][$y] == 2 ) {

                # すでに動いた後
                next;
            }

            # 各要素の取り出し
            my ( $mKind, $mName, $mHp ) = monsterSpec( $landValue->[$x][$y] );
            my ($special) = $HmonsterSpecial[$mKind];

            # 硬化中?
            if (   ( ( $special == 3 ) && ( ( $HislandTurn % 2 ) == 1 ) )
                || ( ( $special == 4 ) && ( ( $HislandTurn % 2 ) == 0 ) ) )
            {
                # 硬化中
                next;
            }

            # 動く方向を決定
            my ( $d, $sx, $sy );
            my ($i);
            for ( $i = 0; $i < 3; $i++ ) {
                $d  = random(6) + 1;
                $sx = $x + $ax[$d];
                $sy = $y + $ay[$d];

                # 行による位置調整
                if ( ( ( $sy % 2 ) == 0 ) && ( ( $y % 2 ) == 1 ) ) {
                    $sx--;
                }

                # 範囲外判定
                if (   ( $sx < 0 )
                    || ( $sx >= $HislandSize )
                    || ( $sy < 0 )
                    || ( $sy >= $HislandSize ) )
                {
                    next;
                }

                # 海、海基、油田、怪獣、山、記念碑以外
                if (   ( $land->[$sx][$sy] != $HlandSea )
                    && ( $land->[$sx][$sy] != $HlandSbase )
                    && ( $land->[$sx][$sy] != $HlandOil )
                    && ( $land->[$sx][$sy] != $HlandMountain )
                    && ( $land->[$sx][$sy] != $HlandMonument )
                    && ( $land->[$sx][$sy] != $HlandMonster ) )
                {
                    last;
                }
            }

            if ( $i == 3 ) {

                # 動かなかった
                next;
            }

            # 動いた先の地形によりメッセージ
            my ($l)     = $land->[$sx][$sy];
            my ($lv)    = $landValue->[$sx][$sy];
            my ($lName) = landName( $l, $lv );
            my ($point) = "($sx, $sy)";

            # 移動
            $land->[$sx][$sy]      = $land->[$x][$y];
            $landValue->[$sx][$sy] = $landValue->[$x][$y];

            # もと居た位置を荒地に
            $land->[$x][$y]      = $HlandWaste;
            $landValue->[$x][$y] = 0;

            # 移動済みフラグ
            if ( $HmonsterSpecial[$mKind] == 2 ) {

                # 移動済みフラグは立てない
            }
            elsif ( $HmonsterSpecial[$mKind] == 1 ) {

                # 速い怪獣
                $monsterMove[$sx][$sy] = $monsterMove[$x][$y] + 1;
            }
            else {
                # 普通の怪獣
                $monsterMove[$sx][$sy] = 2;
            }

            if ( ( $l == $HlandDefence ) && ( $HdBaseAuto == 1 ) ) {

                # 防衛施設を踏んだ
                logMonsMoveDefence( $id, $name, $lName, $point, $mName );

                # 広域被害ルーチン
                wideDamage( $id, $name, $land, $landValue, $sx, $sy );
            }
            else {
                # 行き先が荒地になる
                logMonsMove( $id, $name, $lName, $point, $mName );
            }
        }

        # 火災判定
        if (   ( ( $landKind == $HlandTown ) && ( $lv > 30 ) )
            || ( $landKind == $HlandHaribote )
            || ( $landKind == $HlandFactory ) )
        {
            if ( random(1000) < $HdisFire ) {

                # 周囲の森と記念碑を数える
                if ((     countAround( $land, $x, $y, $HlandForest, 7 )
                        + countAround( $land, $x, $y, $HlandMonument, 7 )
                    ) == 0
                    )
                {
                    # 無かった場合、火災で壊滅
                    my ($l)     = $land->[$x][$y];
                    my ($lv)    = $landValue->[$x][$y];
                    my ($point) = "($x, $y)";
                    my ($lName) = landName( $l, $lv );
                    logFire( $id, $name, $lName, $point );
                    $land->[$x][$y]      = $HlandWaste;
                    $landValue->[$x][$y] = 0;
                }
            }
        }
    }
}

# 周囲の町、農場があるか判定
sub countGrow {
    my ( $land, $landValue, $x, $y ) = @_;
    my ( $i, $sx, $sy );
    for ( $i = 1; $i < 7; $i++ ) {
        $sx = $x + $ax[$i];
        $sy = $y + $ay[$i];

        # 行による位置調整
        if ( ( ( $sy % 2 ) == 0 ) && ( ( $y % 2 ) == 1 ) ) {
            $sx--;
        }

        if (   ( $sx < 0 )
            || ( $sx >= $HislandSize )
            || ( $sy < 0 )
            || ( $sy >= $HislandSize ) )
        {
        }
        else {
            # 範囲内の場合
            if (   ( $land->[$sx][$sy] == $HlandTown )
                || ( $land->[$sx][$sy] == $HlandFarm ) )
            {
                if ( $landValue->[$sx][$sy] != 1 ) {
                    return 1;
                }
            }
        }
    }
    return 0;
}

# 島全体
sub doIslandProcess {
    my ( $number, $island ) = @_;

    # 導出値
    my ($name)      = $island->{'name'};
    my ($id)        = $island->{'id'};
    my ($land)      = $island->{'land'};
    my ($landValue) = $island->{'landValue'};

    # 地震判定
    if ( random(1000) < ( ( $island->{'prepare2'} + 1 ) * $HdisEarthquake ) )
    {
        # 地震発生
        logEarthquake( $id, $name );

        my ( $x, $y, $landKind, $lv, $i );
        for ( $i = 0; $i < $HpointNumber; $i++ ) {
            $x        = $Hrpx[$i];
            $y        = $Hrpy[$i];
            $landKind = $land->[$x][$y];
            $lv       = $landValue->[$x][$y];

            if (   ( ( $landKind == $HlandTown ) && ( $lv >= 100 ) )
                || ( $landKind == $HlandHaribote )
                || ( $landKind == $HlandFactory ) )
            {
                # 1/4で壊滅
                if ( random(4) == 0 ) {
                    logEQDamage( $id, $name, landName( $landKind, $lv ),
                        "($x, $y)" );
                    $land->[$x][$y]      = $HlandWaste;
                    $landValue->[$x][$y] = 0;
                }
            }

        }
    }

    # 食料不足
    if ( $island->{'food'} <= 0 ) {

        # 不足メッセージ
        logStarve( $id, $name );
        $island->{'food'} = 0;

        my ( $x, $y, $landKind, $lv, $i );
        for ( $i = 0; $i < $HpointNumber; $i++ ) {
            $x        = $Hrpx[$i];
            $y        = $Hrpy[$i];
            $landKind = $land->[$x][$y];
            $lv       = $landValue->[$x][$y];

            if (   ( $landKind == $HlandFarm )
                || ( $landKind == $HlandFactory )
                || ( $landKind == $HlandBase )
                || ( $landKind == $HlandDefence ) )
            {
                # 1/4で壊滅
                if ( random(4) == 0 ) {
                    logSvDamage( $id, $name, landName( $landKind, $lv ),
                        "($x, $y)" );
                    $land->[$x][$y]      = $HlandWaste;
                    $landValue->[$x][$y] = 0;
                }
            }
        }
    }

    # 津波判定
    if ( random(1000) < $HdisTsunami ) {

        # 津波発生
        logTsunami( $id, $name );

        my ( $x, $y, $landKind, $lv, $i );
        for ( $i = 0; $i < $HpointNumber; $i++ ) {
            $x        = $Hrpx[$i];
            $y        = $Hrpy[$i];
            $landKind = $land->[$x][$y];
            $lv       = $landValue->[$x][$y];

            if (   ( $landKind == $HlandTown )
                || ( $landKind == $HlandFarm )
                || ( $landKind == $HlandFactory )
                || ( $landKind == $HlandBase )
                || ( $landKind == $HlandDefence )
                || ( $landKind == $HlandHaribote ) )
            {
                # 1d12 <= (周囲の海 - 1) で崩壊
                if (random(12) < (
                              countAround( $land, $x, $y, $HlandOil, 7 )
                            + countAround( $land, $x, $y, $HlandSbase, 7 )
                            + countAround( $land, $x, $y, $HlandSea,   7 )
                            - 1
                    )
                    )
                {
                    logTsunamiDamage( $id, $name, landName( $landKind, $lv ),
                        "($x, $y)" );
                    $land->[$x][$y]      = $HlandWaste;
                    $landValue->[$x][$y] = 0;
                }
            }

        }
    }

    # 怪獣判定
    my ($r)   = random(10000);
    my ($pop) = $island->{'pop'};
    do {
        if ((      ( $r < ( $HdisMonster * $island->{'area'} ) )
                && ( $pop >= $HdisMonsBorder1 )
            )
            || ( $island->{'monstersend'} > 0 )
            )
        {
            # 怪獣出現
            # 種類を決める
            my ( $lv, $kind );
            if ( $island->{'monstersend'} > 0 ) {

                # 人造
                $kind = 0;
                $island->{'monstersend'}--;
            }
            elsif ( $pop >= $HdisMonsBorder3 ) {

                # level3まで
                $kind = random($HmonsterLevel3) + 1;
            }
            elsif ( $pop >= $HdisMonsBorder2 ) {

                # level2まで
                $kind = random($HmonsterLevel2) + 1;
            }
            else {
                # level1のみ
                $kind = random($HmonsterLevel1) + 1;
            }

            # lvの値を決める
            $lv
                = $kind * 10
                + $HmonsterBHP[$kind]
                + random( $HmonsterDHP[$kind] );

            # どこに現れるか決める
            my ( $bx, $by, $i );
            for ( $i = 0; $i < $HpointNumber; $i++ ) {
                $bx = $Hrpx[$i];
                $by = $Hrpy[$i];
                if ( $land->[$bx][$by] == $HlandTown ) {

                    # 地形名
                    my ($lName)
                        = landName( $HlandTown, $landValue->[$bx][$by] );

                    # そのヘックスを怪獣に
                    $land->[$bx][$by]      = $HlandMonster;
                    $landValue->[$bx][$by] = $lv;

                    # 怪獣情報
                    my ( $mKind, $mName, $mHp ) = monsterSpec($lv);

                    # メッセージ
                    logMonsCome( $id, $name, $mName, "($bx, $by)", $lName );
                    last;
                }
            }
        }
    } while ( $island->{'monstersend'} > 0 );

    # 地盤沈下判定
    if (   ( $island->{'area'} > $HdisFallBorder )
        && ( random(1000) < $HdisFalldown ) )
    {
        # 地盤沈下発生
        logFalldown( $id, $name );

        my ( $x, $y, $landKind, $lv, $i );
        for ( $i = 0; $i < $HpointNumber; $i++ ) {
            $x        = $Hrpx[$i];
            $y        = $Hrpy[$i];
            $landKind = $land->[$x][$y];
            $lv       = $landValue->[$x][$y];

            if (   ( $landKind != $HlandSea )
                && ( $landKind != $HlandSbase )
                && ( $landKind != $HlandOil )
                && ( $landKind != $HlandMountain ) )
            {

                # 周囲に海があれば、値を-1に
                if (  countAround( $land, $x, $y, $HlandSea, 7 )
                    + countAround( $land, $x, $y, $HlandSbase, 7 ) )
                {
                    logFalldownLand( $id, $name, landName( $landKind, $lv ),
                        "($x, $y)" );
                    $land->[$x][$y]      = -1;
                    $landValue->[$x][$y] = 0;
                }
            }
        }

        for ( $i = 0; $i < $HpointNumber; $i++ ) {
            $x        = $Hrpx[$i];
            $y        = $Hrpy[$i];
            $landKind = $land->[$x][$y];

            if ( $landKind == -1 ) {

                # -1になっている所を浅瀬に
                $land->[$x][$y]      = $HlandSea;
                $landValue->[$x][$y] = 1;
            }
            elsif ( $landKind == $HlandSea ) {

                # 浅瀬は海に
                $landValue->[$x][$y] = 0;
            }

        }
    }

    # 台風判定
    if ( random(1000) < $HdisTyphoon ) {

        # 台風発生
        logTyphoon( $id, $name );

        my ( $x, $y, $landKind, $lv, $i );
        for ( $i = 0; $i < $HpointNumber; $i++ ) {
            $x        = $Hrpx[$i];
            $y        = $Hrpy[$i];
            $landKind = $land->[$x][$y];
            $lv       = $landValue->[$x][$y];

            if (   ( $landKind == $HlandFarm )
                || ( $landKind == $HlandHaribote ) )
            {

                # 1d12 <= (6 - 周囲の森) で崩壊
                if (random(12) < (
                        6 - countAround( $land, $x, $y, $HlandForest, 7 )
                            - countAround( $land, $x, $y, $HlandMonument, 7 )
                    )
                    )
                {
                    logTyphoonDamage( $id, $name, landName( $landKind, $lv ),
                        "($x, $y)" );
                    $land->[$x][$y]      = $HlandPlains;
                    $landValue->[$x][$y] = 0;
                }
            }

        }
    }

    # 巨大隕石判定
    if ( random(1000) < $HdisHugeMeteo ) {
        my ( $x, $y, $landKind, $lv, $point );

        # 落下
        $x        = random($HislandSize);
        $y        = random($HislandSize);
        $landKind = $land->[$x][$y];
        $lv       = $landValue->[$x][$y];
        $point    = "($x, $y)";

        # メッセージ
        logHugeMeteo( $id, $name, $point );

        # 広域被害ルーチン
        wideDamage( $id, $name, $land, $landValue, $x, $y );
    }

    # 巨大ミサイル判定
    while ( $island->{'bigmissile'} > 0 ) {
        $island->{'bigmissile'}--;

        my ( $x, $y, $landKind, $lv, $point );

        # 落下
        $x        = random($HislandSize);
        $y        = random($HislandSize);
        $landKind = $land->[$x][$y];
        $lv       = $landValue->[$x][$y];
        $point    = "($x, $y)";

        # メッセージ
        logMonDamage( $id, $name, $point );

        # 広域被害ルーチン
        wideDamage( $id, $name, $land, $landValue, $x, $y );
    }

    # 隕石判定
    if ( random(1000) < $HdisMeteo ) {
        my ( $x, $y, $landKind, $lv, $point, $first );
        $first = 1;
        while ( ( random(2) == 0 ) || ( $first == 1 ) ) {
            $first = 0;

            # 落下
            $x        = random($HislandSize);
            $y        = random($HislandSize);
            $landKind = $land->[$x][$y];
            $lv       = $landValue->[$x][$y];
            $point    = "($x, $y)";

            if ( ( $landKind == $HlandSea ) && ( $lv == 0 ) ) {

                # 海ポチャ
                logMeteoSea( $id, $name, landName( $landKind, $lv ), $point );
            }
            elsif ( $landKind == $HlandMountain ) {

                # 山破壊
                logMeteoMountain( $id, $name, landName( $landKind, $lv ),
                    $point );
                $land->[$x][$y]      = $HlandWaste;
                $landValue->[$x][$y] = 0;
                next;
            }
            elsif ( $landKind == $HlandSbase ) {
                logMeteoSbase( $id, $name, landName( $landKind, $lv ),
                    $point );
            }
            elsif ( $landKind == $HlandMonster ) {
                logMeteoMonster( $id, $name, landName( $landKind, $lv ),
                    $point );
            }
            elsif ( $landKind == $HlandSea ) {

                # 浅瀬
                logMeteoSea1( $id, $name, landName( $landKind, $lv ),
                    $point );
            }
            else {
                logMeteoNormal( $id, $name, landName( $landKind, $lv ),
                    $point );
            }
            $land->[$x][$y]      = $HlandSea;
            $landValue->[$x][$y] = 0;
        }
    }

    # 噴火判定
    if ( random(1000) < $HdisEruption ) {
        my ( $x, $y, $sx, $sy, $i, $landKind, $lv, $point );
        $x        = random($HislandSize);
        $y        = random($HislandSize);
        $landKind = $land->[$x][$y];
        $lv       = $landValue->[$x][$y];
        $point    = "($x, $y)";
        logEruption( $id, $name, landName( $landKind, $lv ), $point );
        $land->[$x][$y]      = $HlandMountain;
        $landValue->[$x][$y] = 0;

        for ( $i = 1; $i < 7; $i++ ) {
            $sx = $x + $ax[$i];
            $sy = $y + $ay[$i];

            # 行による位置調整
            if ( ( ( $sy % 2 ) == 0 ) && ( ( $y % 2 ) == 1 ) ) {
                $sx--;
            }

            $landKind = $land->[$sx][$sy];
            $lv       = $landValue->[$sx][$sy];
            $point    = "($sx, $sy)";

            if (   ( $sx < 0 )
                || ( $sx >= $HislandSize )
                || ( $sy < 0 )
                || ( $sy >= $HislandSize ) )
            {
            }
            else {
                # 範囲内の場合
                $landKind = $land->[$sx][$sy];
                $lv       = $landValue->[$sx][$sy];
                $point    = "($sx, $sy)";
                if (   ( $landKind == $HlandSea )
                    || ( $landKind == $HlandOil )
                    || ( $landKind == $HlandSbase ) )
                {
                    # 海の場合
                    if ( $lv == 1 ) {

                        # 浅瀬
                        logEruptionSea1( $id, $name,
                            landName( $landKind, $lv ), $point );
                    }
                    else {
                        logEruptionSea( $id, $name,
                            landName( $landKind, $lv ), $point );
                        $land->[$sx][$sy]      = $HlandSea;
                        $landValue->[$sx][$sy] = 1;
                        next;
                    }
                }
                elsif (( $landKind == $HlandMountain )
                    || ( $landKind == $HlandMonster )
                    || ( $landKind == $HlandWaste ) )
                {
                    next;
                }
                else {
                    # それ以外の場合
                    logEruptionNormal( $id, $name, landName( $landKind, $lv ),
                        $point );
                }
                $land->[$sx][$sy]      = $HlandWaste;
                $landValue->[$sx][$sy] = 0;
            }
        }
    }

    # 食料があふれてたら換金
    if ( $island->{'food'} > 9999 ) {
        $island->{'money'} += int( ( $island->{'food'} - 9999 ) / 10 );
        $island->{'food'} = 9999;
    }

    # 金があふれてたら切り捨て
    if ( $island->{'money'} > 9999 ) {
        $island->{'money'} = 9999;
    }

    # 各種の値を計算
    estimate($number);

    # 繁栄、災難賞
    $pop = $island->{'pop'};
    my ($damage) = $island->{'oldPop'} - $pop;
    my ($prize)  = $island->{'prize'};
    $prize =~ /([0-9]*),([0-9]*),(.*)/;
    my ($flags)    = $1;
    my ($monsters) = $2;
    my ($turns)    = $3;

    # 繁栄賞
    if ( ( !( $flags & 1 ) ) && $pop >= 3000 ) {
        $flags |= 1;
        logPrize( $id, $name, $Hprize[1] );
    }
    elsif ( ( !( $flags & 2 ) ) && $pop >= 5000 ) {
        $flags |= 2;
        logPrize( $id, $name, $Hprize[2] );
    }
    elsif ( ( !( $flags & 4 ) ) && $pop >= 10000 ) {
        $flags |= 4;
        logPrize( $id, $name, $Hprize[3] );
    }

    # 災難賞
    if ( ( !( $flags & 64 ) ) && $damage >= 500 ) {
        $flags |= 64;
        logPrize( $id, $name, $Hprize[7] );
    }
    elsif ( ( !( $flags & 128 ) ) && $damage >= 1000 ) {
        $flags |= 128;
        logPrize( $id, $name, $Hprize[8] );
    }
    elsif ( ( !( $flags & 256 ) ) && $damage >= 2000 ) {
        $flags |= 256;
        logPrize( $id, $name, $Hprize[9] );
    }

    $island->{'prize'} = "$flags,$monsters,$turns";
}

# 人口順にソート
sub islandSort {
    my ( $flag, $i, $tmp );

    # 人口が同じときは直前のターンの順番のまま
    my @idx = ( 0 .. $#Hislands );
    @idx = sort {
               $Hislands[$b]->{'pop'} <=> $Hislands[$a]->{'pop'}
            || $a <=> $b
    } @idx;
    @Hislands = @Hislands[@idx];
}

# 広域被害ルーチン
sub wideDamage {
    my ( $id, $name, $land, $landValue, $x, $y ) = @_;
    my ( $sx, $sy, $i, $landKind, $landName, $lv, $point );

    for ( $i = 0; $i < 19; $i++ ) {
        $sx = $x + $ax[$i];
        $sy = $y + $ay[$i];

        # 行による位置調整
        if ( ( ( $sy % 2 ) == 0 ) && ( ( $y % 2 ) == 1 ) ) {
            $sx--;
        }

        $landKind = $land->[$sx][$sy];
        $lv       = $landValue->[$sx][$sy];
        $landName = landName( $landKind, $lv );
        $point    = "($sx, $sy)";

        # 範囲外判定
        if (   ( $sx < 0 )
            || ( $sx >= $HislandSize )
            || ( $sy < 0 )
            || ( $sy >= $HislandSize ) )
        {
            next;
        }

        # 範囲による分岐
        if ( $i < 7 ) {

            # 中心、および1ヘックス
            if ( $landKind == $HlandSea ) {
                $landValue->[$sx][$sy] = 0;
                next;
            }
            elsif (( $landKind == $HlandSbase )
                || ( $landKind == $HlandOil ) )
            {
                logWideDamageSea2( $id, $name, $landName, $point );
                $land->[$sx][$sy]      = $HlandSea;
                $landValue->[$sx][$sy] = 0;
            }
            else {
                if ( $landKind == $HlandMonster ) {
                    logWideDamageMonsterSea( $id, $name, $landName, $point );
                }
                else {
                    logWideDamageSea( $id, $name, $landName, $point );
                }
                $land->[$sx][$sy] = $HlandSea;
                if ( $i == 0 ) {

                    # 海
                    $landValue->[$sx][$sy] = 0;
                }
                else {
                    # 浅瀬
                    $landValue->[$sx][$sy] = 1;
                }
            }
        }
        else {
            # 2ヘックス
            if (   ( $landKind == $HlandSea )
                || ( $landKind == $HlandOil )
                || ( $landKind == $HlandWaste )
                || ( $landKind == $HlandMountain )
                || ( $landKind == $HlandSbase ) )
            {
                next;
            }
            elsif ( $landKind == $HlandMonster ) {
                logWideDamageMonster( $id, $name, $landName, $point );
                $land->[$sx][$sy]      = $HlandWaste;
                $landValue->[$sx][$sy] = 0;
            }
            else {
                logWideDamageWaste( $id, $name, $landName, $point );
                $land->[$sx][$sy]      = $HlandWaste;
                $landValue->[$sx][$sy] = 0;
            }
        }
    }
}

# ログへの出力
# 第1引数:メッセージ
# 第2引数:当事者
# 第3引数:相手
# 通常ログ
sub logOut {
    push( @HlogPool, "0,$HislandTurn,$_[1],$_[2],$_[0]" );
}

# 遅延ログ
sub logLate {
    push( @HlateLogPool, "0,$HislandTurn,$_[1],$_[2],$_[0]" );
}

# 機密ログ
sub logSecret {
    push( @HsecretLogPool, "1,$HislandTurn,$_[1],$_[2],$_[0]" );
}

# 記録ログ
sub logHistory {
    open( my $HOUT, ">>${HdirName}/hakojima.his" );
    print $HOUT "$HislandTurn,$_[0]\n";
    close($HOUT);
}

# 記録ログ調整
sub logHistoryTrim {
    open( my $HIN, "${HdirName}/hakojima.his" );
    my ( @line, $l, $count );
    $count = 0;
    while ( $l = <$HIN> ) {
        chomp($l);
        push( @line, $l );
        $count++;
    }
    close($HIN);

    if ( $count > $HhistoryMax ) {
        open( my $HOUT, ">${HdirName}/hakojima.his" );
        my ($i);
        for ( $i = ( $count - $HhistoryMax ); $i < $count; $i++ ) {
            print $HOUT "$line[$i]\n";
        }
        close($HOUT);
    }
}

# ログ書き出し
sub logFlush {
    open( my $LOUT, ">${HdirName}/hakojima.log0" );

    # 全部逆順にして書き出す
    my ($i);
    for ( $i = $#HsecretLogPool; $i >= 0; $i-- ) {
        print $LOUT $HsecretLogPool[$i];
        print $LOUT "\n";
    }
    for ( $i = $#HlateLogPool; $i >= 0; $i-- ) {
        print $LOUT $HlateLogPool[$i];
        print $LOUT "\n";
    }
    for ( $i = $#HlogPool; $i >= 0; $i-- ) {
        print $LOUT $HlogPool[$i];
        print $LOUT "\n";
    }
    close($LOUT);
}

#----------------------------------------------------------------------
# ログテンプレート
#----------------------------------------------------------------------
# 資金足りない
sub logNoMoney {
    my ( $id, $name, $comName ) = @_;
    logOut(
        "${HtagName_}${name}島${H_tagName}で予定されていた${HtagComName_}$comName${H_tagComName}は、資金不足のため中止されました。",
        $id
    );
}

# 食料足りない
sub logNoFood {
    my ( $id, $name, $comName ) = @_;
    logOut(
        "${HtagName_}${name}島${H_tagName}で予定されていた${HtagComName_}$comName${H_tagComName}は、備蓄食料不足のため中止されました。",
        $id
    );
}

# 対象地形の種類による失敗
sub logLandFail {
    my ( $id, $name, $comName, $kind, $point ) = @_;
    logOut(
        "${HtagName_}${name}島${H_tagName}で予定されていた${HtagComName_}$comName${H_tagComName}は、予定地の${HtagName_}$point${H_tagName}が<B>$kind</B>だったため中止されました。",
        $id
    );

    END;
}

# 周りに陸がなくて埋め立て失敗
sub logNoLandAround {
    my ( $id, $name, $comName, $point ) = @_;
    logOut(
        "${HtagName_}${name}島${H_tagName}で予定されていた${HtagComName_}$comName${H_tagComName}は、予定地の${HtagName_}$point${H_tagName}の周辺に陸地がなかったため中止されました。",
        $id
    );

    END;
}

# 整地系成功
sub logLandSuc {
    my ( $id, $name, $comName, $point ) = @_;
    logOut(
        "${HtagName_}${name}島$point${H_tagName}で${HtagComName_}${comName}${H_tagComName}が行われました。",
        $id
    );

    END;
}

# 油田発見
sub logOilFound {
    my ( $id, $name, $point, $comName, $str ) = @_;
    logOut(
        "${HtagName_}${name}島$point${H_tagName}で<B>$str</B>の予算をつぎ込んだ${HtagComName_}${comName}${H_tagComName}が行われ、<B>油田が掘り当てられました</B>。",
        $id
    );

    END;
}

# 油田発見ならず
sub logOilFail {
    my ( $id, $name, $point, $comName, $str ) = @_;
    logOut(
        "${HtagName_}${name}島$point${H_tagName}で<B>$str</B>の予算をつぎ込んだ${HtagComName_}${comName}${H_tagComName}が行われましたが、油田は見つかりませんでした。",
        $id
    );

    END;
}

# 油田からの収入
sub logOilMoney {
    my ( $id, $name, $lName, $point, $str ) = @_;
    logOut(
        "${HtagName_}${name}島$point${H_tagName}の<B>$lName</B>から、<B>$str</B>の収益が上がりました。",
        $id
    );

    END;
}

# 油田枯渇
sub logOilEnd {
    my ( $id, $name, $lName, $point ) = @_;
    logOut(
        "${HtagName_}${name}島$point${H_tagName}の<B>$lName</B>は枯渇したようです。",
        $id
    );

    END;
}

# 防衛施設、自爆セット
sub logBombSet {
    my ( $id, $name, $lName, $point ) = @_;
    logOut(
        "${HtagName_}${name}島$point${H_tagName}の<B>$lName</B>の<B>自爆装置がセット</B>されました。",
        $id
    );

    END;
}

# 防衛施設、自爆作動
sub logBombFire {
    my ( $id, $name, $lName, $point ) = @_;
    logOut(
        "${HtagName_}${name}島$point${H_tagName}の<B>$lName</B>、${HtagDisaster_}自爆装置作動！！${H_tagDisaster}",
        $id
    );

    END;
}

# 記念碑、発射
sub logMonFly {
    my ( $id, $name, $lName, $point ) = @_;
    logOut(
        "${HtagName_}${name}島$point${H_tagName}の<B>$lName</B>が<B>轟音とともに飛び立ちました</B>。",
        $id
    );

    END;
}

# 記念碑、落下
sub logMonDamage {
    my ( $id, $name, $point ) = @_;
    logOut(
        "<B>何かとてつもないもの</B>が${HtagName_}${name}島$point${H_tagName}地点に落下しました！！",
        $id
    );
}

# 植林orミサイル基地
sub logPBSuc {
    my ( $id, $name, $comName, $point ) = @_;
    logSecret(
        "${HtagName_}${name}島$point${H_tagName}で${HtagComName_}${comName}${H_tagComName}が行われました。",
        $id
    );
    logOut(
        "こころなしか、${HtagName_}${name}島${H_tagName}の<B>森</B>が増えたようです。",
        $id
    );

    END;
}

# ハリボテ
sub logHariSuc {
    my ( $id, $name, $comName, $comName2, $point ) = @_;
    logSecret(
        "${HtagName_}${name}島$point${H_tagName}で${HtagComName_}${comName}${H_tagComName}が行われました。",
        $id
    );
    logLandSuc( $id, $name, $comName2, $point );

    END;
}

# ミサイル撃とうとした(or 怪獣派遣しようとした)がターゲットがいない
sub logMsNoTarget {
    my ( $id, $name, $comName ) = @_;
    logOut(
        "${HtagName_}${name}島${H_tagName}で予定されていた${HtagComName_}${comName}${H_tagComName}は、目標の島に人が見当たらないため中止されました。",
        $id
    );

    END;
}

# ミサイル撃とうとしたが基地がない
sub logMsNoBase {
    my ( $id, $name, $comName ) = @_;
    logOut(
        "${HtagName_}${name}島${H_tagName}で予定されていた${HtagComName_}${comName}${H_tagComName}は、<B>ミサイル設備を保有していない</B>ために実行できませんでした。",
        $id
    );

    END;
}

# ミサイル撃ったが範囲外
sub logMsOut {
    my ( $id, $tId, $name, $tName, $comName, $point ) = @_;
    logOut(
        "${HtagName_}${name}島${H_tagName}が${HtagName_}${tName}島$point${H_tagName}地点に向けて${HtagComName_}${comName}${H_tagComName}を行いましたが、<B>領域外の海</B>に落ちた模様です。",
        $id, $tId
    );
}

# ステルスミサイル撃ったが範囲外
sub logMsOutS {
    my ( $id, $tId, $name, $tName, $comName, $point ) = @_;
    logSecret(
        "${HtagName_}${name}島${H_tagName}が${HtagName_}${tName}島$point${H_tagName}地点に向けて${HtagComName_}${comName}${H_tagComName}を行いましたが、<B>領域外の海</B>に落ちた模様です。",
        $id, $tId
    );
    logLate(
        "<B>何者か</B>が${HtagName_}${tName}島$point${H_tagName}へ向けて${HtagComName_}${comName}${H_tagComName}を行いましたが、<B>領域外の海</B>に落ちた模様です。",
        $tId
    );
}

# ミサイル撃ったが防衛施設でキャッチ
sub logMsCaught {
    my ( $id, $tId, $name, $tName, $comName, $point, $tPoint ) = @_;
    logOut(
        "${HtagName_}${name}島${H_tagName}が${HtagName_}${tName}島$point${H_tagName}地点に向けて${HtagComName_}${comName}${H_tagComName}を行いましたが、${HtagName_}$tPoint${H_tagName}地点上空にて力場に捉えられ、<B>空中爆発</B>しました。",
        $id, $tId
    );
}

# ステルスミサイル撃ったが防衛施設でキャッチ
sub logMsCaughtS {
    my ( $id, $tId, $name, $tName, $comName, $point, $tPoint ) = @_;
    logSecret(
        "${HtagName_}${name}島${H_tagName}が${HtagName_}${tName}島$point${H_tagName}地点に向けて${HtagComName_}${comName}${H_tagComName}を行いましたが、${HtagName_}$tPoint${H_tagName}地点上空にて力場に捉えられ、<B>空中爆発</B>しました。",
        $id, $tId
    );
    logLate(
        "<B>何者か</B>が${HtagName_}${tName}島$point${H_tagName}へ向けて${HtagComName_}${comName}${H_tagComName}を行いましたが、${HtagName_}$tPoint${H_tagName}地点上空にて力場に捉えられ、<B>空中爆発</B>しました。",
        $tId
    );
}

# ミサイル撃ったが効果なし
sub logMsNoDamage {
    my ( $id, $tId, $name, $tName, $comName, $tLname, $point, $tPoint ) = @_;
    logOut(
        "${HtagName_}${name}島${H_tagName}が${HtagName_}${tName}島$point${H_tagName}地点に向けて${HtagComName_}${comName}${H_tagComName}を行いましたが、${HtagName_}$tPoint${H_tagName}の<B>$tLname</B>に落ちたので被害がありませんでした。",
        $id, $tId
    );
}

# ステルスミサイル撃ったが効果なし
sub logMsNoDamageS {
    my ( $id, $tId, $name, $tName, $comName, $tLname, $point, $tPoint ) = @_;
    logSecret(
        "${HtagName_}${name}島${H_tagName}が${HtagName_}${tName}島$point${H_tagName}地点に向けて${HtagComName_}${comName}${H_tagComName}を行いましたが、${HtagName_}$tPoint${H_tagName}の<B>$tLname</B>に落ちたので被害がありませんでした。",
        $id, $tId
    );

    logLate(
        "<B>何者か</B>が${HtagName_}${tName}島$point${H_tagName}地点に向けて${HtagComName_}${comName}${H_tagComName}を行いましたが、${HtagName_}$tPoint${H_tagName}の<B>$tLname</B>に落ちたので被害がありませんでした。",
        $tId
    );
}

# 陸地破壊弾、山に命中
sub logMsLDMountain {
    my ( $id, $tId, $name, $tName, $comName, $tLname, $point, $tPoint ) = @_;
    logOut(
        "${HtagName_}${name}島${H_tagName}が${HtagName_}${tName}島$point${H_tagName}地点に向けて${HtagComName_}${comName}${H_tagComName}を行い、${HtagName_}$tPoint${H_tagName}の<B>$tLname</B>に命中。<B>$tLname</B>は消し飛び、荒地と化しました。",
        $id, $tId
    );
}

# 陸地破壊弾、海底基地に命中
sub logMsLDSbase {
    my ( $id, $tId, $name, $tName, $comName, $tLname, $point, $tPoint ) = @_;
    logOut(
        "${HtagName_}${name}島${H_tagName}が${HtagName_}${tName}島$point${H_tagName}地点に向けて${HtagComName_}${comName}${H_tagComName}を行い、${HtagName_}$tPoint${H_tagName}に着水後爆発、同地点にあった<B>$tLname</B>は跡形もなく吹き飛びました。",
        $id, $tId
    );
}

# 陸地破壊弾、怪獣に命中
sub logMsLDMonster {
    my ( $id, $tId, $name, $tName, $comName, $tLname, $point, $tPoint ) = @_;
    logOut(
        "${HtagName_}${name}島${H_tagName}が${HtagName_}${tName}島$point${H_tagName}地点に向けて${HtagComName_}${comName}${H_tagComName}を行い、${HtagName_}$tPoint${H_tagName}に着弾し爆発。陸地は<B>怪獣$tLname</B>もろとも水没しました。",
        $id, $tId
    );
}

# 陸地破壊弾、浅瀬に命中
sub logMsLDSea1 {
    my ( $id, $tId, $name, $tName, $comName, $tLname, $point, $tPoint ) = @_;
    logOut(
        "${HtagName_}${name}島${H_tagName}が${HtagName_}${tName}島$point${H_tagName}地点に向けて${HtagComName_}${comName}${H_tagComName}を行い、${HtagName_}$tPoint${H_tagName}の<B>$tLname</B>に着弾。海底がえぐられました。",
        $id, $tId
    );
}

# 陸地破壊弾、その他の地形に命中
sub logMsLDLand {
    my ( $id, $tId, $name, $tName, $comName, $tLname, $point, $tPoint ) = @_;
    logOut(
        "${HtagName_}${name}島${H_tagName}が${HtagName_}${tName}島$point${H_tagName}地点に向けて${HtagComName_}${comName}${H_tagComName}を行い、${HtagName_}$tPoint${H_tagName}の<B>$tLname</B>に着弾。陸地は水没しました。",
        $id, $tId
    );
}

# 通常ミサイル、荒地に着弾
sub logMsWaste {
    my ( $id, $tId, $name, $tName, $comName, $tLname, $point, $tPoint ) = @_;
    logOut(
        "${HtagName_}${name}島${H_tagName}が${HtagName_}${tName}島$point${H_tagName}地点に向けて${HtagComName_}${comName}${H_tagComName}を行いましたが、${HtagName_}$tPoint${H_tagName}の<B>$tLname</B>に落ちました。",
        $id, $tId
    );
}

# ステルスミサイル、荒地に着弾
sub logMsWasteS {
    my ( $id, $tId, $name, $tName, $comName, $tLname, $point, $tPoint ) = @_;
    logSecret(
        "${HtagName_}${name}島${H_tagName}が${HtagName_}${tName}島$point${H_tagName}地点に向けて${HtagComName_}${comName}${H_tagComName}を行いましたが、${HtagName_}$tPoint${H_tagName}の<B>$tLname</B>に落ちました。",
        $id, $tId
    );
    logLate(
        "<B>何者か</B>が${HtagName_}${tName}島$point${H_tagName}地点に向けて${HtagComName_}${comName}${H_tagComName}を行いましたが、${HtagName_}$tPoint${H_tagName}の<B>$tLname</B>に落ちました。",
        $tId
    );
}

# 通常ミサイル、怪獣に命中、硬化中にて無傷
sub logMsMonNoDamage {
    my ( $id, $tId, $name, $tName, $comName, $tLname, $point, $tPoint ) = @_;
    logOut(
        "${HtagName_}${name}島${H_tagName}が${HtagName_}${tName}島$point${H_tagName}地点に向けて${HtagComName_}${comName}${H_tagComName}を行い、${HtagName_}$tPoint${H_tagName}の<B>怪獣$tLname</B>に命中、しかし硬化状態だったため効果がありませんでした。",
        $id, $tId
    );
}

# ステルスミサイル、怪獣に命中、硬化中にて無傷
sub logMsMonNoDamageS {
    my ( $id, $tId, $name, $tName, $comName, $tLname, $point, $tPoint ) = @_;
    logSecret(
        "${HtagName_}${name}島${H_tagName}が${HtagName_}${tName}島$point${H_tagName}地点に向けて${HtagComName_}${comName}${H_tagComName}を行い、${HtagName_}$tPoint${H_tagName}の<B>怪獣$tLname</B>に命中、しかし硬化状態だったため効果がありませんでした。",
        $id, $tId
    );
    logOut(
        "<B>何者か</B>が${HtagName_}${tName}島$point${H_tagName}地点に向けて${HtagComName_}${comName}${H_tagComName}を行い、${HtagName_}$tPoint${H_tagName}の<B>怪獣$tLname</B>に命中、しかし硬化状態だったため効果がありませんでした。",
        $tId
    );
}

# 通常ミサイル、怪獣に命中、殺傷
sub logMsMonKill {
    my ( $id, $tId, $name, $tName, $comName, $tLname, $point, $tPoint ) = @_;
    logOut(
        "${HtagName_}${name}島${H_tagName}が${HtagName_}${tName}島$point${H_tagName}地点に向けて${HtagComName_}${comName}${H_tagComName}を行い、${HtagName_}$tPoint${H_tagName}の<B>怪獣$tLname</B>に命中。<B>怪獣$tLname</B>は力尽き、倒れました。",
        $id, $tId
    );
}

# ステルスミサイル、怪獣に命中、殺傷
sub logMsMonKillS {
    my ( $id, $tId, $name, $tName, $comName, $tLname, $point, $tPoint ) = @_;
    logSecret(
        "${HtagName_}${name}島${H_tagName}が${HtagName_}${tName}島$point${H_tagName}地点に向けて${HtagComName_}${comName}${H_tagComName}を行い、${HtagName_}$tPoint${H_tagName}の<B>怪獣$tLname</B>に命中。<B>怪獣$tLname</B>は力尽き、倒れました。",
        $id, $tId
    );
    logLate(
        "<B>何者か</B>が${HtagName_}${tName}島$point${H_tagName}地点に向けて${HtagComName_}${comName}${H_tagComName}を行い、${HtagName_}$tPoint${H_tagName}の<B>怪獣$tLname</B>に命中。<B>怪獣$tLname</B>は力尽き、倒れました。",
        $tId
    );
}

# 通常ミサイル、怪獣に命中、ダメージ
sub logMsMonster {
    my ( $id, $tId, $name, $tName, $comName, $tLname, $point, $tPoint ) = @_;
    logOut(
        "${HtagName_}${name}島${H_tagName}が${HtagName_}${tName}島$point${H_tagName}地点に向けて${HtagComName_}${comName}${H_tagComName}を行い、${HtagName_}$tPoint${H_tagName}の<B>怪獣$tLname</B>に命中。<B>怪獣$tLname</B>は苦しそうに咆哮しました。",
        $id, $tId
    );
}

# ステルスミサイル、怪獣に命中、ダメージ
sub logMsMonsterS {
    my ( $id, $tId, $name, $tName, $comName, $tLname, $point, $tPoint ) = @_;
    logSecret(
        "${HtagName_}${name}島${H_tagName}が${HtagName_}${tName}島$point${H_tagName}地点に向けて${HtagComName_}${comName}${H_tagComName}を行い、${HtagName_}$tPoint${H_tagName}の<B>怪獣$tLname</B>に命中。<B>怪獣$tLname</B>は苦しそうに咆哮しました。",
        $id, $tId
    );
    logLate(
        "<B>何者か</B>が${HtagName_}${tName}島$point${H_tagName}地点に向けて${HtagComName_}${comName}${H_tagComName}を行い、${HtagName_}$tPoint${H_tagName}の<B>怪獣$tLname</B>に命中。<B>怪獣$tLname</B>は苦しそうに咆哮しました。",
        $tId
    );
}

# 怪獣の死体
sub logMsMonMoney {
    my ( $tId, $mName, $value ) = @_;
    logOut(
        "<B>怪獣$mName</B>の残骸には、<B>$value$HunitMoney</B>の値が付きました。",
        $tId
    );
}

# 通常ミサイル通常地形に命中
sub logMsNormal {
    my ( $id, $tId, $name, $tName, $comName, $tLname, $point, $tPoint ) = @_;
    logOut(
        "${HtagName_}${name}島${H_tagName}が${HtagName_}${tName}島$point${H_tagName}地点に向けて${HtagComName_}${comName}${H_tagComName}を行い、${HtagName_}$tPoint${H_tagName}の<B>$tLname</B>に命中、一帯が壊滅しました。",
        $id, $tId
    );
}

# ステルスミサイル通常地形に命中
sub logMsNormalS {
    my ( $id, $tId, $name, $tName, $comName, $tLname, $point, $tPoint ) = @_;
    logSecret(
        "${HtagName_}${name}島${H_tagName}が${HtagName_}${tName}島$point${H_tagName}地点に向けて${HtagComName_}${comName}${H_tagComName}を行い、${HtagName_}$tPoint${H_tagName}の<B>$tLname</B>に命中、一帯が壊滅しました。",
        $id, $tId
    );
    logLate(
        "<B>何者か</B>が${HtagName_}${tName}島$point${H_tagName}地点に向けて${HtagComName_}${comName}${H_tagComName}を行い、${HtagName_}$tPoint${H_tagName}の<B>$tLname</B>に命中、一帯が壊滅しました。",
        $tId
    );
}

# ミサイル難民到着
sub logMsBoatPeople {
    my ( $id, $name, $achive ) = @_;
    logOut(
        "${HtagName_}${name}島${H_tagName}にどこからともなく<B>$achive${HunitPop}もの難民</B>が漂着しました。${HtagName_}${name}島${H_tagName}は快く受け入れたようです。",
        $id
    );
}

# 怪獣派遣
sub logMonsSend {
    my ( $id, $tId, $name, $tName ) = @_;
    logOut(
        "${HtagName_}${name}島${H_tagName}が<B>人造怪獣</B>を建造。${HtagName_}${tName}島${H_tagName}へ送りこみました。",
        $id, $tId
    );
}

# 資金繰り
sub logDoNothing {
    my ( $id, $name, $comName ) = @_;

#    logOut("${HtagName_}${name}島${H_tagName}で${HtagComName_}${comName}${H_tagComName}が行われました。",$id);
}

# 輸出
sub logSell {
    my ( $id, $name, $comName, $value ) = @_;
    logOut(
        "${HtagName_}${name}島${H_tagName}が<B>$value$HunitFood</B>の${HtagComName_}${comName}${H_tagComName}を行いました。",
        $id
    );
}

# 援助
sub logAid {
    my ( $id, $tId, $name, $tName, $comName, $str ) = @_;
    logOut(
        "${HtagName_}${name}島${H_tagName}が${HtagName_}${tName}島${H_tagName}へ<B>$str</B>の${HtagComName_}${comName}${H_tagComName}を行いました。",
        $id, $tId
    );
}

# 誘致活動
sub logPropaganda {
    my ( $id, $name, $comName ) = @_;
    logOut(
        "${HtagName_}${name}島${H_tagName}で${HtagComName_}${comName}${H_tagComName}が行われました。",
        $id
    );
}

# 放棄
sub logGiveup {
    my ( $id, $name ) = @_;
    logOut(
        "${HtagName_}${name}島${H_tagName}は放棄され、<B>無人島</B>になりました。",
        $id
    );
    logHistory(
        "${HtagName_}${name}島${H_tagName}、放棄され<B>無人島</B>となる。"
    );
}

# 死滅
sub logDead {
    my ( $id, $name ) = @_;
    logOut(
        "${HtagName_}${name}島${H_tagName}から人がいなくなり、<B>無人島</B>になりました。",
        $id
    );
    logHistory(
        "${HtagName_}${name}島${H_tagName}、人がいなくなり<B>無人島</B>となる。"
    );
}

# 発見
sub logDiscover {
    my ($name) = @_;
    logHistory("${HtagName_}${name}島${H_tagName}が発見される。");
}

# 名前の変更
sub logChangeName {
    my ( $name1, $name2 ) = @_;
    logHistory(
        "${HtagName_}${name1}島${H_tagName}、名称を${HtagName_}${name2}島${H_tagName}に変更する。"
    );
}

# 飢餓
sub logStarve {
    my ( $id, $name ) = @_;
    logOut(
        "${HtagName_}${name}島${H_tagName}の${HtagDisaster_}食料が不足${H_tagDisaster}しています！！",
        $id
    );
}

# 怪獣現る
sub logMonsCome {
    my ( $id, $name, $mName, $point, $lName ) = @_;
    logOut(
        "${HtagName_}${name}島${H_tagName}に<B>怪獣$mName</B>出現！！${HtagName_}$point${H_tagName}の<B>$lName</B>が踏み荒らされました。",
        $id
    );
}

# 怪獣動く
sub logMonsMove {
    my ( $id, $name, $lName, $point, $mName ) = @_;
    logOut(
        "${HtagName_}${name}島$point${H_tagName}の<B>$lName</B>が<B>怪獣$mName</B>に踏み荒らされました。",
        $id
    );
}

# 怪獣、防衛施設を踏む
sub logMonsMoveDefence {
    my ( $id, $name, $lName, $point, $mName ) = @_;
    logOut(
        "<B>怪獣$mName</B>が${HtagName_}${name}島$point${H_tagName}の<B>$lName</B>へ到達、<B>${lName}の自爆装置が作動！！</B>",
        $id
    );
}

# 火災
sub logFire {
    my ( $id, $name, $lName, $point ) = @_;
    logOut(
        "${HtagName_}${name}島$point${H_tagName}の<B>$lName</B>が${HtagDisaster_}火災${H_tagDisaster}により壊滅しました。",
        $id
    );
}

# 埋蔵金
sub logMaizo {
    my ( $id, $name, $comName, $value ) = @_;
    logOut(
        "${HtagName_}${name}島${H_tagName}での${HtagComName_}$comName${H_tagComName}中に、<B>$value$HunitMoney"
            . "もの埋蔵金</B>が発見されました。",
        $id
    );
}

# 地震発生
sub logEarthquake {
    my ( $id, $name ) = @_;
    logOut(
        "${HtagName_}${name}島${H_tagName}で大規模な${HtagDisaster_}地震${H_tagDisaster}が発生！！",
        $id
    );
}

# 地震被害
sub logEQDamage {
    my ( $id, $name, $lName, $point ) = @_;
    logOut(
        "${HtagName_}${name}島$point${H_tagName}の<B>$lName</B>は${HtagDisaster_}地震${H_tagDisaster}により壊滅しました。",
        $id
    );
}

# 食料不足被害
sub logSvDamage {
    my ( $id, $name, $lName, $point ) = @_;
    logOut(
        "${HtagName_}${name}島$point${H_tagName}の<B>$lName</B>に<B>食料を求めて住民が殺到</B>。<B>$lName</B>は壊滅しました。",
        $id
    );
}

# 津波発生
sub logTsunami {
    my ( $id, $name ) = @_;
    logOut(
        "${HtagName_}${name}島${H_tagName}付近で${HtagDisaster_}津波${H_tagDisaster}発生！！",
        $id
    );
}

# 津波被害
sub logTsunamiDamage {
    my ( $id, $name, $lName, $point ) = @_;
    logOut(
        "${HtagName_}${name}島$point${H_tagName}の<B>$lName</B>は${HtagDisaster_}津波${H_tagDisaster}により崩壊しました。",
        $id
    );
}

# 台風発生
sub logTyphoon {
    my ( $id, $name ) = @_;
    logOut(
        "${HtagName_}${name}島${H_tagName}に${HtagDisaster_}台風${H_tagDisaster}上陸！！",
        $id
    );
}

# 台風被害
sub logTyphoonDamage {
    my ( $id, $name, $lName, $point ) = @_;
    logOut(
        "${HtagName_}${name}島$point${H_tagName}の<B>$lName</B>は${HtagDisaster_}台風${H_tagDisaster}で飛ばされました。",
        $id
    );
}

# 隕石、海
sub logMeteoSea {
    my ( $id, $name, $lName, $point ) = @_;
    logOut(
        "${HtagName_}${name}島$point${H_tagName}の<B>$lName</B>に${HtagDisaster_}隕石${H_tagDisaster}が落下しました。",
        $id
    );
}

# 隕石、山
sub logMeteoMountain {
    my ( $id, $name, $lName, $point ) = @_;
    logOut(
        "${HtagName_}${name}島$point${H_tagName}の<B>$lName</B>に${HtagDisaster_}隕石${H_tagDisaster}が落下、<B>$lName</B>は消し飛びました。",
        $id
    );
}

# 隕石、海底基地
sub logMeteoSbase {
    my ( $id, $name, $lName, $point ) = @_;
    logOut(
        "${HtagName_}${name}島$point${H_tagName}の<B>$lName</B>に${HtagDisaster_}隕石${H_tagDisaster}が落下、<B>$lName</B>は崩壊しました。",
        $id
    );
}

# 隕石、怪獣
sub logMeteoMonster {
    my ( $id, $name, $lName, $point ) = @_;
    logOut(
        "<B>怪獣$lName</B>がいた${HtagName_}${name}島$point${H_tagName}地点に${HtagDisaster_}隕石${H_tagDisaster}が落下、陸地は<B>怪獣$lName</B>もろとも水没しました。",
        $id
    );
}

# 隕石、浅瀬
sub logMeteoSea1 {
    my ( $id, $name, $lName, $point ) = @_;
    logOut(
        "${HtagName_}${name}島$point${H_tagName}地点に${HtagDisaster_}隕石${H_tagDisaster}が落下、海底がえぐられました。",
        $id
    );
}

# 隕石、その他
sub logMeteoNormal {
    my ( $id, $name, $lName, $point ) = @_;
    logOut(
        "${HtagName_}${name}島$point${H_tagName}地点の<B>$lName</B>に${HtagDisaster_}隕石${H_tagDisaster}が落下、一帯が水没しました。",
        $id
    );
}

# 隕石、その他
sub logHugeMeteo {
    my ( $id, $name, $point ) = @_;
    logOut(
        "${HtagName_}${name}島$point${H_tagName}地点に${HtagDisaster_}巨大隕石${H_tagDisaster}が落下！！",
        $id
    );
}

# 噴火
sub logEruption {
    my ( $id, $name, $lName, $point ) = @_;
    logOut(
        "${HtagName_}${name}島$point${H_tagName}地点で${HtagDisaster_}火山が噴火${H_tagDisaster}、<B>山</B>が出来ました。",
        $id
    );
}

# 噴火、浅瀬
sub logEruptionSea1 {
    my ( $id, $name, $lName, $point ) = @_;
    logOut(
        "${HtagName_}${name}島$point${H_tagName}地点の<B>$lName</B>は、${HtagDisaster_}噴火${H_tagDisaster}の影響で陸地になりました。",
        $id
    );
}

# 噴火、海or海基
sub logEruptionSea {
    my ( $id, $name, $lName, $point ) = @_;
    logOut(
        "${HtagName_}${name}島$point${H_tagName}地点の<B>$lName</B>は、${HtagDisaster_}噴火${H_tagDisaster}の影響で海底が隆起、浅瀬になりました。",
        $id
    );
}

# 噴火、その他
sub logEruptionNormal {
    my ( $id, $name, $lName, $point ) = @_;
    logOut(
        "${HtagName_}${name}島$point${H_tagName}地点の<B>$lName</B>は、${HtagDisaster_}噴火${H_tagDisaster}の影響で壊滅しました。",
        $id
    );
}

# 地盤沈下発生
sub logFalldown {
    my ( $id, $name ) = @_;
    logOut(
        "${HtagName_}${name}島${H_tagName}で${HtagDisaster_}地盤沈下${H_tagDisaster}が発生しました！！",
        $id
    );
}

# 地盤沈下被害
sub logFalldownLand {
    my ( $id, $name, $lName, $point ) = @_;
    logOut(
        "${HtagName_}${name}島$point${H_tagName}の<B>$lName</B>は海の中へ沈みました。",
        $id
    );
}

# 広域被害、水没
sub logWideDamageSea {
    my ( $id, $name, $lName, $point ) = @_;
    logOut(
        "${HtagName_}${name}島$point${H_tagName}の<B>$lName</B>は<B>水没</B>しました。",
        $id
    );
}

# 広域被害、海の建設
sub logWideDamageSea2 {
    my ( $id, $name, $lName, $point ) = @_;
    logOut(
        "${HtagName_}${name}島$point${H_tagName}の<B>$lName</B>は跡形もなくなりました。",
        $id
    );
}

# 広域被害、怪獣水没
sub logWideDamageMonsterSea {
    my ( $id, $name, $lName, $point ) = @_;
    logOut(
        "${HtagName_}${name}島$point${H_tagName}の陸地は<B>怪獣$lName</B>もろとも水没しました。",
        $id
    );
}

# 広域被害、怪獣
sub logWideDamageMonster {
    my ( $id, $name, $lName, $point ) = @_;
    logOut(
        "${HtagName_}${name}島$point${H_tagName}の<B>怪獣$lName</B>は消し飛びました。",
        $id
    );
}

# 広域被害、荒地
sub logWideDamageWaste {
    my ( $id, $name, $lName, $point ) = @_;
    logOut(
        "${HtagName_}${name}島$point${H_tagName}の<B>$lName</B>は一瞬にして<B>荒地</B>と化しました。",
        $id
    );
}

# 受賞
sub logPrize {
    my ( $id, $name, $pName ) = @_;
    logOut(
        "${HtagName_}${name}島${H_tagName}が<B>$pName</B>を受賞しました。",
        $id
    );
    logHistory("${HtagName_}${name}島${H_tagName}、<B>$pName</B>を受賞");
}

# 島がいっぱいな場合
sub tempNewIslandFull {
    out(<<END);
${HtagBig_}申し訳ありません、島が一杯で登録できません！！${H_tagBig}$HtempBack
END
}

# 新規で名前がない場合
sub tempNewIslandNoName {
    out(<<END);
${HtagBig_}島につける名前が必要です。${H_tagBig}$HtempBack
END
}

# 新規で名前が不正な場合
sub tempNewIslandBadName {
    out(<<END);
${HtagBig_}',?()<>\$'とか入ってたり、「無人島」とかいった変な名前はやめましょうよ〜${H_tagBig}$HtempBack
END
}

# すでにその名前の島がある場合
sub tempNewIslandAlready {
    out(<<END);
${HtagBig_}その島ならすでに発見されています。${H_tagBig}$HtempBack
END
}

# パスワードがない場合
sub tempNewIslandNoPassword {
    out(<<END);
${HtagBig_}パスワードが必要です。${H_tagBig}$HtempBack
END
}

# 島を発見しました!!
sub tempNewIslandHead {
    out(<<END);
<CENTER>
${HtagBig_}島を発見しました！！${H_tagBig}<BR>
${HtagBig_}${HtagName_}「${HcurrentName}島」${H_tagName}と命名します。${H_tagBig}<BR>
$HtempBack<BR>
</CENTER>
END
}

# 地形の呼び方
sub landName {
    my ( $land, $lv ) = @_;
    if ( $land == $HlandSea ) {
        if ( $lv == 1 ) {
            return '浅瀬';
        }
        else {
            return '海';
        }
    }
    elsif ( $land == $HlandWaste ) {
        return '荒地';
    }
    elsif ( $land == $HlandPlains ) {
        return '平地';
    }
    elsif ( $land == $HlandTown ) {
        if ( $lv < 30 ) {
            return '村';
        }
        elsif ( $lv < 100 ) {
            return '町';
        }
        else {
            return '都市';
        }
    }
    elsif ( $land == $HlandForest ) {
        return '森';
    }
    elsif ( $land == $HlandFarm ) {
        return '農場';
    }
    elsif ( $land == $HlandFactory ) {
        return '工場';
    }
    elsif ( $land == $HlandBase ) {
        return 'ミサイル基地';
    }
    elsif ( $land == $HlandDefence ) {
        return '防衛施設';
    }
    elsif ( $land == $HlandMountain ) {
        return '山';
    }
    elsif ( $land == $HlandMonster ) {
        my ( $kind, $name, $hp ) = monsterSpec($lv);
        return $name;
    }
    elsif ( $land == $HlandSbase ) {
        return '海底基地';
    }
    elsif ( $land == $HlandOil ) {
        return '海底油田';
    }
    elsif ( $land == $HlandMonument ) {
        return $HmonumentName[$lv];
    }
    elsif ( $land == $HlandHaribote ) {
        return 'ハリボテ';
    }
}

# 人口その他の値を算出
sub estimate {
    my ($number) = $_[0];
    my ($island);
    my ( $pop, $area, $farm, $factory, $mountain ) = ( 0, 0, 0, 0, 0, 0 );

    # 地形を取得
    $island = $Hislands[$number];
    my ($land)      = $island->{'land'};
    my ($landValue) = $island->{'landValue'};

    # 数える
    my ( $x, $y, $kind, $value );
    for ( $y = 0; $y < $HislandSize; $y++ ) {
        for ( $x = 0; $x < $HislandSize; $x++ ) {
            $kind  = $land->[$x][$y];
            $value = $landValue->[$x][$y];
            if (   ( $kind != $HlandSea )
                && ( $kind != $HlandSbase )
                && ( $kind != $HlandOil ) )
            {
                $area++;
                if ( $kind == $HlandTown ) {

                    # 町
                    $pop += $value;
                }
                elsif ( $kind == $HlandFarm ) {

                    # 農場
                    $farm += $value;
                }
                elsif ( $kind == $HlandFactory ) {

                    # 工場
                    $factory += $value;
                }
                elsif ( $kind == $HlandMountain ) {

                    # 山
                    $mountain += $value;
                }
            }
        }
    }

    # 代入
    $island->{'pop'}      = $pop;
    $island->{'area'}     = $area;
    $island->{'farm'}     = $farm;
    $island->{'factory'}  = $factory;
    $island->{'mountain'} = $mountain;
}

# 範囲内の地形を数える
sub countAround {
    my ( $land, $x, $y, $kind, $range ) = @_;
    my ( $i, $count, $sx, $sy );
    $count = 0;
    for ( $i = 0; $i < $range; $i++ ) {
        $sx = $x + $ax[$i];
        $sy = $y + $ay[$i];

        # 行による位置調整
        if ( ( ( $sy % 2 ) == 0 ) && ( ( $y % 2 ) == 1 ) ) {
            $sx--;
        }

        if (   ( $sx < 0 )
            || ( $sx >= $HislandSize )
            || ( $sy < 0 )
            || ( $sy >= $HislandSize ) )
        {
            # 範囲外の場合
            if ( $kind == $HlandSea ) {

                # 海なら加算
                $count++;
            }
        }
        else {
            # 範囲内の場合
            if ( $land->[$sx][$sy] == $kind ) {
                $count++;
            }
        }
    }
    return $count;
}

# 0から(n - 1)までの数字が一回づつ出てくる数列を作る
sub randomArray {
    my ($n) = @_;
    my ( @list, $i );

    # 初期値
    if ( $n == 0 ) {
        $n = 1;
    }
    @list = ( 0 .. $n - 1 );

    # シャッフル
    for ( $i = $n; --$i; ) {
        my ($j) = int( rand( $i + 1 ) );
        if ( $i == $j ) { next; }
        @list[ $i, $j ] = @list[ $j, $i ];
    }

    return @list;
}

# 名前変更失敗
sub tempChangeNothing {
    out(<<END);
${HtagBig_}名前、パスワードともに空欄です${H_tagBig}$HtempBack
END
}

# 名前変更資金足りず
sub tempChangeNoMoney {
    out(<<END);
${HtagBig_}資金不足のため変更できません${H_tagBig}$HtempBack
END
}

# 名前変更成功
sub tempChange {
    out(<<END);
${HtagBig_}変更完了しました${H_tagBig}$HtempBack
END
}

1;

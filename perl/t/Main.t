use strict;
use warnings;
use utf8;

use Test::More;
use Test::WWW::Mechanize::PSGI;
use Plack::Test;
use Plack::Util;

my $path       = 'http://localhost:5000/cgi-bin/hako-main.cgi';
my $mente_path = 'http://localhost:5000/cgi-bin/hako-mente.cgi';

my $app = Plack::Util::load_psgi('./app.psgi');
my $mech = Test::WWW::Mechanize::PSGI->new( app => $app );

my $island_name      = 'てすとじま';
my $island_name_full = 'てすとじま島';
my $pass             = 'testpass';
my $mente_pass       = 'yourpassword';

sub initialize_data {
    system "rm -fr cgi/data*";
    $mech->get_ok($mente_path);
    $mech->submit_form_ok(
        +{  fields => +{ PASSWORD => $mente_pass, },
            button => 'NEW',
        }
    );
}

sub new_data {
    initialize_data;
    $mech->get_ok($path);
    $mech->content_like(qr|箱庭諸島２|);
    $mech->submit_form_ok(
        +{  form_number => 2,
            fields      => +{
                ISLANDNAME => $island_name,
                PASSWORD   => $pass,
                PASSWORD2  => $pass,
            },
            button => 'NewIslandButton',
        }
    );
}

subtest 'トップページが表示できる', sub {
    initialize_data;
    $mech->get_ok($path);
    $mech->content_like(qr|箱庭諸島２|);
    $mech->content_like(qr|ターン1|);
    $mech->content_like(qr|自分の島へ|);
    $mech->content_like(qr|諸島の状況|);
    $mech->content_like(qr|新しい島を探す|);
    $mech->content_like(qr|島の名前とパスワードの変更|);
    $mech->content_like(qr|最近の出来事|);
    $mech->content_like(qr|発見の記録|);
    is @{ $mech->forms }, 3;
};

subtest '新しい島を探す', sub {
    initialize_data;
    $mech->get_ok($path);
    $mech->content_like(qr|箱庭諸島２|);
    $mech->submit_form_ok(
        +{  form_number => 2,
            fields      => +{
                ISLANDNAME => $island_name,
                PASSWORD   => $pass,
                PASSWORD2  => $pass,
            },
            button => 'NewIslandButton',
        }
    );
    $mech->content_like(qr|島を発見しました！！|);
    $mech->content_like(qr|「${island_name_full}」|);
};

subtest '自分の島へ', sub {
    new_data;
    $mech->get_ok($path);
    $mech->submit_form_ok(
        +{  form_number => 1,
            fields      => +{
                ISLANDID => 1,
                PASSWORD => $pass,
            },
            button => 'OwnerButton',
        }
    );
    $mech->content_like(qr|${island_name_full}.*開発計画|);
    $mech->content_like(qr|${island_name_full}.*の近況|);
};

subtest '観光', sub {
    new_data;
    $mech->get_ok($path);
    $mech->follow_link_ok( +{ url_regex => qr|Sight| } );
    $mech->content_like(qr|「${island_name_full}」.*へようこそ！！|);
    $mech->content_like(qr|${island_name_full}.*の近況|);
};

subtest '島の名前とパスワードの変更', sub {
    new_data;
    $mech->get_ok($path);
    $mech->submit_form_ok(
        +{  form_number => 3,
            fields      => +{
                ISLANDID   => 1,
                ISLANDNAME => 'かわったなまえ',
                OLDPASS    => $pass,
            },
            button => 'ChangeInfoButton',
        }
    );
    $mech->content_like(qr|資金不足のため変更できません|);
};

subtest 'コマンドの登録', sub {
    new_data;
    $mech->get_ok($path);
    $mech->submit_form_ok(
        +{  form_number => 1,
            fields      => +{
                ISLANDID => 1,
                PASSWORD => $pass,
            },
            button => 'OwnerButton',
        }
    );
    $mech->content_like(qr|${island_name_full}.*開発計画|);
    $mech->content_like(qr|${island_name_full}.*の近況|);
    $mech->submit_form_ok(
        +{  form_number => 1,
            fields      => +{

                NUMBER  => 0,
                COMMAND => 1,
                POINTX  => 0,
                POINTY  => 0,
                AMOUNT  => 0,
            },
            button => 'CommandButton1',
        }
    );
    $mech->content_like(qr|コマンドを登録しました|);
    $mech->content_like(qr|(0,0).*で.*整地|);
};

subtest 'コメント更新', sub {
    new_data;
    $mech->get_ok($path);
    $mech->submit_form_ok(
        +{  form_number => 1,
            fields      => +{
                ISLANDID => 1,
                PASSWORD => $pass,
            },
            button => 'OwnerButton',
        }
    );
    $mech->content_like(qr|${island_name_full}.*開発計画|);
    $mech->content_like(qr|${island_name_full}.*の近況|);
    $mech->submit_form_ok(
        +{  form_number => 2,
            fields      => +{
                MESSAGE  => 'めっせーじ',
                PASSWORD => $pass,
            },
            button => 'MessageButton1',
        }
    );
    $mech->content_like(qr|コメントを更新しました|);
    $mech->get_ok($path);
    $mech->content_like(qr|コメント：.*めっせーじ|);
};

subtest '島を10個作る', sub {
    initialize_data;

    for ( my $i = 0; $i < 10; $i++ ) {
        $mech->get_ok($path);
        $mech->content_like(qr|箱庭諸島２|);
        $mech->submit_form_ok(
            +{  form_number => 2,
                fields      => +{
                    ISLANDNAME => $island_name . $i,
                    PASSWORD   => $pass . $i,
                    PASSWORD2  => $pass . $i,
                },
                button => 'NewIslandButton',
            }
        );
        $mech->content_like(qr|島を発見しました！！|);
        $mech->content_like(qr|「${island_name}${i}島」|);
    }
};

done_testing;


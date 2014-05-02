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

my $island_name = 'てすとじま';
my $island_name_full = 'てすとじま島';
my $pass        = 'testpass';
my $mente_pass  = 'yourpassword';

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
        +{
            form_number => 1,
            fields => +{
                ISLANDID => 1,
                PASSWORD => $pass,
            },
            button => 'OwnerButton',
        }
    );
    $mech->content_like(qr|${island_name_full}.*開発計画|);
    $mech->content_like(qr|${island_name_full}.*の近況|);
};

done_testing;


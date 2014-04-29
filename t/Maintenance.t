use strict;
use warnings;
use utf8;

use Test::More;
use Test::WWW::Mechanize::PSGI;
use Plack::Test;
use Plack::Util;

my $path       = 'http://localhost:5000/cgi-bin/hako-mente.cgi';
my $wrong_pass = 'hogehoge';
my $pass       = 'yourpassword';

my $app = Plack::Util::load_psgi('./app.psgi');
my $mech = Test::WWW::Mechanize::PSGI->new( app => $app );

sub clean_data {
    system "rm -fr cgi/data*";
}

sub make_backup_data {
    system "cp -r cgi/data cgi/data.bak1";
}

subtest 'アクセスする', sub {
    clean_data;
    $mech->get_ok($path);
    $mech->content_like(qr|箱島２ メンテナンスツール|);
};

subtest 'データの新規作成に成功する', sub {
    clean_data;
    $mech->get_ok($path);
    $mech->submit_form_ok(
        +{  fields => +{ PASSWORD => $pass, },
            button => 'NEW',
        }
    );
    $mech->content_like(qr|現役データ|);
};

subtest
    'パスワードを間違えた場合、データの新規作成に失敗する',
    sub {
    clean_data;
    $mech->get_ok($path);
    $mech->submit_form_ok(
        +{  fields => +{ PASSWORD => $wrong_pass, },
            button => 'NEW',
        }
    );
    $mech->content_like(qr|パスワードが違います。|);
    };

subtest 'データの削除に成功する', sub {
    clean_data;
    $mech->get_ok($path);
    $mech->submit_form_ok(
        +{  fields => +{ PASSWORD => $pass, },
            button => 'NEW',
        }
    );
    $mech->content_like(qr|現役データ|);
    $mech->submit_form_ok(
        +{  fields => +{ PASSWORD => $pass, },
            button => 'DELETE',
        }
    );
    $mech->content_like(qr|新しいデータを作る|);
};

done_testing;


#!/usr/bin/env plackup
use Plack::Builder;
use Plack::App::WrapCGI;
use Plack::App::File;

builder {
	mount '/images' => Plack::App::File->new(root => './public/images');
	mount '/cgi-bin/hako-main.cgi'
			=> Plack::App::WrapCGI->new(script => './cgi/hako-main.cgi', execute => 1)->to_app;
	mount '/cgi-bin/hako-mente.cgi'
			=> Plack::App::WrapCGI->new(script => './cgi/hako-mente.cgi', execute => 1)->to_app;
};


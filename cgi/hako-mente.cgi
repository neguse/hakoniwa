#!/usr/bin/perl

use strict;
use warnings;
use utf8;

BEGIN { push @INC, '../lib'; }

use Hako::Maintenance;

run_maintenance;


#!/bin/sh
carton exec perltidy cgi/*.cgi
carton exec perltidy t/*.t
carton exec perltidy app.psgi

#!/bin/sh

webcronpack="webcron-`date +%Y%m%d%H%M%S`"

mv conf/app.conf .app.conf
mv database/cron.sqlite3 .cron.sqlite3
rm -rf webcron-*.tgz

echo "开始打包$webcronpack.tgz ......"

export GOARCH=amd64
export GOOS=linux

bee pack

mv webcron.tar.gz $webcronpack.tgz

mv .app.conf conf/app.conf
mv .cron.sqlite3 database/cron.sqlite3

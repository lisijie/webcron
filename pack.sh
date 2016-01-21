#!/bin/sh

tarfile="webcron-$1.tar.gz"

echo "开始打包$tarfile..."

export GOARCH=amd64
export GOOS=linux

bee pack

mv webcron.tar.gz $tarfile

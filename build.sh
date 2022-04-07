#!/bin/bash

PROG="srun4-webcron"
UPX=/usr/bin/upx

build() {
  echo "dos2unix ./*.sh"
  dos2unix ./*.sh

  echo "Start build $PROG"
  go build -ldflags="-s -w" -o $PROG main.go

  echo "Start upx $PROG"
  if [ -x $UPX ]; then
    $UPX -9 $PROG
  else
    upx -9 $PROG
  fi

  echo "Build Success"

  echo "Start Translate to /srun3/bin"
  systemctl stop srun-webcron
  rm -f /srun3/bin/$PROG
  mv ./$PROG /srun3/bin

  systemctl start srun-webcron
}

build

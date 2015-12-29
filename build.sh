#!/bin/sh

go build -o ./bin/webcron ./

cp -r ./views ./conf ./static ./run.sh ./bin/

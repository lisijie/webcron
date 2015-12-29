#!/bin/sh

go build -o ./webcron ../

cp -r ../views ../conf ../static .

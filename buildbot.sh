#!/bin/bash

export GOPATH=$PWD/gopath

export GOOS=$1
export GOARCH=$2
TAGS=$3

go get -u -d github.com/Nightgunner5/gogame/main || exit $?
go test -v github.com/Nightgunner5/gogame/... || exit $?
go build -x -o gogame -tags "$TAGS" github.com/Nightgunner5/gogame/main || exit $?

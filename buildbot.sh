#!/bin/bash

export GOPATH=$PWD/gopath

export GOOS=$1
export GOARCH=$2
TAGS=$3

`go env`

go get -d github.com/Nightgunner5/gogame/main || exit $?
go test -v github.com/Nightgunner5/gogame/... || exit $?
go build -x -o gogame$GOEXT -tags "$TAGS" github.com/Nightgunner5/gogame/main || exit $?

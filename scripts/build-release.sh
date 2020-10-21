#!/bin/sh

cd $(dirname "$0")
set -e
cd ..

export GO111MODULE=on
export GOPROXY=https://goproxy.cn

mkdir -p release/osx
mkdir -p release/linux
GOARCH=amd64 GOOS=darwin ./scripts/go-mini-build.sh github.com/leizongmin/jssh release/osx
GOARCH=amd64 GOOS=linux ./scripts/go-mini-build.sh github.com/leizongmin/jssh release/linux
ls -alh release

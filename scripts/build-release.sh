#!/bin/sh

cd $(dirname "$0")
set -e
cd ..

export GO111MODULE=on
export GOPROXY=https://goproxy.cn

mkdir -p release/osx
mkdir -p release/linux
export mini_build=./scripts/go-mini-build.sh
export main_package=github.com/leizongmin/jssh
CGO_ENABLED=1 GOARCH=amd64 GOOS=darwin $mini_build $main_package release/osx
CGO_ENABLED=1 GOARCH=amd64 GOOS=linux $mini_build $main_package release/linux
ls -alh release

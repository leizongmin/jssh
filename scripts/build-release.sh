#!/bin/sh

cd "$(dirname "$0")"
set -e
cd ..

export GO111MODULE=on
export GOPROXY=https://goproxy.cn

mkdir -p release/osx
mkdir -p release/linux
export mini_build=./scripts/go-mini-build.sh
export main_package=github.com/leizongmin/jssh

./scripts/update-build-info.sh

export sys_info=$(uname -a)
if [[ $sys_info =~ "Darwin" ]]; then
  $mini_build $main_package release/osx
  # 通过Docker构建Linux版本
  export cache_dir="$(pwd)/release/cross_compile_cache"
  mkdir -p "${cache_dir}"
  docker run --rm -it -v "${cache_dir}:/go" -v $(pwd):$(pwd) -w $(pwd) -e GO111MODULE=on -e GOPROXY=https://goproxy.cn golang:1.15 $mini_build $main_package release/linux
elif [[ $sys_info =~ "Linux" ]]; then
  # 仅构建Linux版本
  $mini_build $main_package release/linux
else
  echo "not supported OS type: ${sys_info}"
fi

ls -alh release

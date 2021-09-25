#!/usr/bin/env bash

cargo build --release
cp target/release/jssh target/release/jssh-min
upx --best target/release/jssh-min
ls -lh target/release/jssh*

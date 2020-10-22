#!/bin/sh

file=internal/pkginfo/pkginfo.go

cat ${file} | sed s/build-[1,2][0-9]*/build-$(date +%Y%m%d%H%M)/g >${file}.tmp
mv ${file}.tmp ${file}

cat ${file} | sed s/commit-[0-9a-f]*/commit-$(git rev-parse --short HEAD)/g >${file}.tmp
mv ${file}.tmp ${file}

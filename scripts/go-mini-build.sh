#!/bin/sh

# from https://gist.github.com/leizongmin/f8795cb9e514fbe695c07867e4a6b46d

usage() {
  echo "usage: go-mini-build <package-name> [output-path]"
}

file_size() {
  echo $(ls -lh "$1" | awk '{print $5}')
}

if [ -z "$1" ]; then
  usage
  exit 1
else
  export PACKAGE="$1"
fi

if [ -z "$2" ]; then
  export OUTPUT="out"
else
  export OUTPUT="$2"
fi

export PREFIX="go-mini-build"

echo "$PREFIX> package:\t$PACKAGE"
echo "$PREFIX> output: \t$OUTPUT"

go build -v -ldflags "-s -w" -o $OUTPUT $PACKAGE && \
echo "$PREFIX> original size:\t$(file_size $OUTPUT)" && \
gzexe "$OUTPUT" && \
rm -f "$OUTPUT~" && \
echo "$PREFIX> final size:\t$(file_size $OUTPUT)" && \
echo "$PREFIX> done."

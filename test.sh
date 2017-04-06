#!/usr/bin/env bash
#set -eou -pipefail

WORKDIR="/usr/local/go/src/github.com/misakwa/gozyre"
BUILD_IMAGE="gozyre:test"

docker build --rm --force-rm -t $BUILD_IMAGE .

docker run -it --rm -v `pwd`:$WORKDIR -w $WORKDIR $BUILD_IMAGE go test

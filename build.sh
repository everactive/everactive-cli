#!/bin/bash
set -e

CODEPATH="/go/src/gitlab.com/everactive/everactive-cli"
BINPATH="/go/bin"
GOOS="linux"
GOARCH="amd64"
GO111MODULE=on
GOPRIVATE=gitlab.com/everactive/*
RELEASE_PATH="release"
#GOPATH="$(go env GOPATH)"
#PATH="$PATH:$GOPATH/bin"
BINARY_VERSION=$(cat version.txt)

#mkdir -p $CODEPATH && mkdir -p $BINPATH
mkdir -p ${RELEASE_PATH}

TARGET_GOOS=(darwin linux windows arm)
_GOARCH=amd64

for _GOOS in "${TARGET_GOOS[@]}"
do
  EXT=""
  if [[ "${_GOOS}" == "arm" ]]; then
    _GOOS="linux"
    _GOARCH="arm"
  else
    _GOARCH="amd64"
  fi
  if [[ "${_GOOS}" == "windows" ]]; then
    EXT=".exe"
  fi
set -x
env GOOS=${_GOOS} GOARCH=${_GOARCH} \
go build \
-o ${RELEASE_PATH}/everactive-cli-${_GOOS}-${_GOARCH}${EXT} \
-ldflags \
"-X gitlab.com/everactive/everactive-cli/lib.Version=${BINARY_VERSION}" \
main.go

tar -czf "${RELEASE_PATH}/everactive-cli-${_GOOS}-${_GOARCH}-${BINARY_VERSION}.tar.gz" "${RELEASE_PATH}/everactive-cli-${_GOOS}-${_GOARCH}${EXT}"
set +x
done
realpath ${RELEASE_PATH}
ls -l ${RELEASE_PATH}


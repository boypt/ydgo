#!/bin/bash
if [[ -f .ydkey ]]; then
  source .ydkey
fi

if [[ -z $YDAPPID ]] || [[ -z $YDAPPSEC ]]; then
  echo "APPID/SEC unset"
  exit 1
fi

VER=$(git describe --tags || git rev-parse --short HEAD)
LDFLAGS="-s -w -X main.VERSION=$VER -X main.YDAppId=$YDAPPID -X main.YDAppSec=$YDAPPSEC"

BIN=ydgo
rm -fv ydgo_*

# for normal unix env
for OS in darwin linux; do
  CGO_ENABLED=0 GOARCH=amd64 GOOS=$OS go build -o ${BIN}_${OS} -ldflags "${LDFLAGS}"
done
CGO_ENABLED=0 GOARCH=amd64 GOOS=windows go build -o ${BIN}_windows.exe -ldflags "${LDFLAGS}"


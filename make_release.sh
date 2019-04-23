#!/bin/bash

GITVER=$(git rev-parse --short HEAD)
VER=$1

if [[ -z $VER ]]; then
  VER=$GITVER
fi


BIN=ydgo
rm -fv ydgo_*

# for normal unix env
for OS in darwin linux; do
  CGO_ENABLED=0 GOARCH=amd64 GOOS=$OS go build -o ${BIN}_${OS} -ldflags "-s -w -X main.VERSION=$VER"
done
CGO_ENABLED=0 GOARCH=amd64 GOOS=windows go build -o ${BIN}_windows.exe -ldflags "-s -w -X main.VERSION=$VER"


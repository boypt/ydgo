#!/bin/bash


BIN=ydgo
GITVER=$(git rev-parse --short HEAD)

rm -fv ydgo_*

if [[ -e /proc/sys/kernel/osrelease ]] &&  grep -q -i microsoft /proc/sys/kernel/osrelease && command -v go.exe; then
  # for WSL 
  for OS in darwin linux; do
    env WSLENV=CGO_ENABLE:GOARCH:GOOS CGO_ENABLED=0 GOARCH=amd64 GOOS=$OS cmd.exe /C go.exe build -o ${BIN}_${OS} -ldflags "-s -w -X main.VERSION=$GITVER"
  done
  env WSLENV=CGO_ENABLE:GOARCH:GOOS CGO_ENABLED=0 GOARCH=amd64 GOOS=windows cmd.exe /C go.exe build -o ${BIN}_windows.exe -ldflags "-s -w -X main.VERSION=$GITVER"
else

  # for normal unix env
  for OS in darwin linux; do
    CGO_ENABLED=0 GOARCH=amd64 GOOS=$OS go build -o ${BIN}_${OS} -ldflags "-s -w -X main.VERSION=$GITVER"
  done
  CGO_ENABLED=0 GOARCH=amd64 GOOS=windows go build -o ${BIN}_windows.exe -ldflags "-s -w -X main.VERSION=$GITVER"
fi	


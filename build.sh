#!/bin/bash
set -x
rm drone-cowpoke
export CGO_ENABLED=0
export GOOS=linux
export GOARCH=amd64
go build -v -a -tags netgo
docker build --rm -t leankit/drone-cowpoke .

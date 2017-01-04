#!/bin/bash
export CGO_ENABLED=0
export GOOS=linux
export GOARCH=amd64

gb build all -f -tags netgo
docker build --rm=true -t leankit/drone-cowpoke .
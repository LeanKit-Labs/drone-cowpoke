#!/bin/bash
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o drone-cowpoke .
docker rmi -f jaymedavis/drone-cowpoke
docker build --no-cache --rm --force-rm=true -t jaymedavis/drone-cowpoke .
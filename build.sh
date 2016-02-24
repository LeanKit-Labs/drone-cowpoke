#!/bin/bash
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o drone-cowpoke .
docker rmi -f leankit/drone-cowpoke
docker build --no-cache --rm --force-rm=true -t leankit/drone-cowpoke .
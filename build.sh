#!/bin/bash
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o drone-cowpoke .
docker build --no-cache -t jaymedavis/drone-cowpoke .
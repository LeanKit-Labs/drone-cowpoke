FROM alpine

RUN apk update && apk add bash ca-certificates git

ADD ./bin/drone-cowpoke-linux-amd64 /usr/local/bin

ENTRYPOINT [ "/usr/local/bin/drone-cowpoke-linux-amd64" ]

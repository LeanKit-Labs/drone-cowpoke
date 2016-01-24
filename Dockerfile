FROM gliderlabs/alpine:3.2
RUN apk add --update \
  ca-certificates

ADD drone-cowpoke /bin/
ENTRYPOINT ["/bin/drone-cowpoke"]

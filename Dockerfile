FROM alpine:3.3
ADD drone-cowpoke /usr/local/bin/
ENTRYPOINT [ "/usr/local/bin/drone-cowpoke" ]
FROM alpine
RUN apk update && apk add bash ca-certificates git
#ADD . /app
#WORKDIR /app
ADD drone-cowpoke /usr/local/bin
ENTRYPOINT [ "/usr/local/bin/drone-cowpoke" ]

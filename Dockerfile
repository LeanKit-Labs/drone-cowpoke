FROM rancher/docker:1.9.1

ADD drone-cowpoke /go/bin/
VOLUME /var/lib/docker
ENTRYPOINT ["/usr/bin/dockerlaunch", "/go/bin/drone-cowpoke"]

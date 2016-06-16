#!/bin/bash
#drone-cowpoke <<EOF
docker run --dns 10.0.1.200 -i --rm -v ${HOME}/gopath/src/github.com/LeanKit-Labs/drone-cowpoke/docker:/app/docker leankit/drone-cowpoke <<EOF
{
	"workspace": {
		"path": "/app/docker/src"
	},
	"repo": {
		"name": "core-leankit-api",
		"owner": "BanditSoftware"
	},
	"vargs": {
		"cowpoke_url": "https://cowpoke.leankit.io",
		"cowpoke_port": 8000,
		"docker_owner": "leankit"
	}
}
EOF

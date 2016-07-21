#!/bin/bash
#drone-cowpoke <<EOF
docker run --dns 10.0.1.200 -i --rm -v /Users/evans/Desktop/drone-cowpoke/docker:/app/docker leankit/drone-cowpoke <<EOF
{
	"workspace": {
		"path": "/app/docker/src"
	},
	"repo": {
		"name": "core-leankit-api",
		"owner": "BanditSoftware"
	},
	"vargs": {
		"cowpoke_url": "http://cowpoke.leankit.io",
		"cowpoke_port": 9000,
		"docker_owner": "leankit", 
		"cowpoke_catalog_upgrade": true,
		"cowpoke_catalog" : "rancher-core-leankit-api"
	}
}
EOF

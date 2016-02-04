#!/bin/bash
./drone-cowpoke <<EOF
{
	"workspace": {
		"path": "./docker/src"
	},
	"vargs": {
		"cowpoke_url": "http://cowpoke",
		"cowpoke_port": 8080
	}
}
EOF

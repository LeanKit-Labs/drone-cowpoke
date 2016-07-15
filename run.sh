#!/bin/bash
	docker run -i --rm \
	-v ${GOPATH}/src/github.com/LeanKit-Labs/drone-cowpoke/docker:/app/docker \
	-v ${GOPATH}/src/github.com/LeanKit-Labs/drone-rancher-catalog/.test-data/rancher-catalog:/rancher-catalog \
	leankit/drone-cowpoke<<EOF
	{
		"build": {
			"Number": 56
		},
		"workspace": {
			"path": "/app"
		},
		"repo": {
			"name": "cowpoke-integration-test",
			"owner": "leankit-labs"
		},
		"vargs": {
			"docker_username": "$DOCKER_USER",
			"docker_password": "$DOCKER_PASS",
			"docker_repo": "leankit/cowpoke-integration-test",
			"catalog_repo": "BanditSoftware/rancher-cowpoke-integration-test",
			"github_token": "$GITHUB_TOKEN",
			"github_user": "$GITHUB_NAME",
			"github_email": "$GITHUB_EMAIL",
			"cowpoke_url": "http://cowpoke.leankit.io:9000",
			"rancher_catalog_name": "cowpoke-integration-test"
		}
	}
	EOF


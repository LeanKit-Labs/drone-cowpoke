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
			"docker_username": "mydockerhubname",
			"docker_password": "mysecretpassword",
			"docker_repo": "leankit/cowpoke-integration-test",
			"catalog_repo": "BanditSoftware/rancher-cowpoke-integration-test",
			"github_token": "mySecretToken",
			"github_user": "myGithubName",
			"github_email": "me@eample.com",
			"cowpoke_url": "https://cowpoke.yourdomain.com",
			"rancher_catalog_name": "cowpoke-integration-test"
			"bearer_token": "token"
		}
	}
	EOF


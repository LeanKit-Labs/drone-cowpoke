drone-cowpoke
=============

This drone plugin provides a service to post a docker image name to Rancher via  [cowpoke](https://github.com/leankit-labs/cowpoke).

## Prerequisites
 * Use [buildgoggles](git://github.com/arobson/buildgoggles) or create a .droneTags.yml file.
 * Use leankit/drone-docker plugin to publish the tags in .droneTags.yml.

Example `.droneTags.yml`
```yaml
tags:
  - latest
  - v1.0.1
  - my_really_descriptive_very_long_tag
```

## Example:
This shows a sample `.drone.yml` that builds a Node.JS project:
 * cowpoke_url:
 *

```yaml
debug: true
build:
  image: node:0.12.10
  commands:
   - npm install
   - npm install git://github.com/arobson/buildgoggles -g
   - buildgoggles

cache:
  mount:
    - node_modules
    - .git

publish:
  docker:
    image: leankit/drone-docker:latest
    environment:
      - DOCKER_LAUNCH_DEBUG=true
    username: $$DOCKER_USER
    password: $$DOCKER_PASS
    email: $$DOCKER_EMAIL
    repo: yourDockerRepo/yourProjectName

deploy:
  cowpoke:
      image: leankit/drone-cowpoke:feature-catalog-update
      environment:
        - DOCKER_LAUNCH_DEBUG=true
      docker_username: mydockerhubname,
      docker_password: mysecretpassword,
      docker_repo: leankit/cowpoke-integration-test,
      catalog_repo: BanditSoftware/rancher-cowpoke-integration-test,
      github_token: mySecretToken,
      github_user: myGithubName,
      github_email: me@eample.com,
      cowpoke_url: https://cowpoke.yourdomain.com,
      rancher_catalog_name": cowpoke-integration-test
```

## Testing via script

```sh
./drone-cowpoke <<EOF
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
		}
	}
	EOF
EOF
```

## Compile

```sh
export GO15VENDOREXPERIMENT=1
go get
go build -a -tags netgo
```

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
      image: leankit/drone-cowpoke
      environment:
        - DOCKER_LAUNCH_DEBUG=true
      cowpoke_url: https://cowpoke.yourdomain.com
      cowpoke_port: 443
      drone_owner: yourDockerRepo # optional: use if your docker owner doesn't match github
      drone_repo: yourProjectName # optional: use if your docker repo doesn't match github
```

## Testing via script

```sh
./drone-cowpoke <<EOF
{
	"workspace": {
		"path": "/app/docker/src"
	},
	"repo": {
		"name": "my-repo",
		"owner": "LeanKit-Labs"
	},
	"vargs": {
		"cowpoke_url": "https://cowpoke",
		"cowpoke_port": 8000,
		"docker_owner": "leankit"
	}
}
EOF
```

## Compile

```sh
export GO15VENDOREXPERIMENT=1
go get
go build -a -tags netgo
```

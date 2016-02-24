drone-cowpoke
=============

This drone plugin provides a service to post a docker image name to [cowpoke](https://github.com/leankit-labs/cowpoke).

## Prerequisites
 * use of [buildgoggles](git://github.com/arobson/buildgoggles) or an equivalent step as part of the build process
 * use of our fork of [drone-docker](https://github.com/LeanKit-Labs/drone-docker)

## Example:
This shows a sample `.drone.yml` that builds a Node.JS project:

```yaml
debug: true
build:
  image: nodesource/jessie:5.2
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
  cowpoke:
      image: leankit/drone-cowpoke
      environment:
        - DOCKER_LAUNCH_DEBUG=true
      cowpoke_url: https://cowpoke.yourdomain.com
      cowpoke_port: 443
```

## Testing via script

```sh
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
```

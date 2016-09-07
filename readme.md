drone-cowpoke
=============

This drone plugin provides a service to post a docker image name to Rancher via  [cowpoke](https://github.com/leankit-labs/cowpoke). It also builds a rancher catalog for the service in order to provide cowpoke with an up to date catalog.

This plugin will read all of the tags in the specified Docker Hub Repo and create the catalog entries entries based on 2 styles of tags.

* **Release Tags:** `v1.0.1`  
* **Feature Branch Tags:** `jgreat_core-api_feature-branch_1.0.3_3_aaaa1234`

#### Feature Branch Tags Format
`<GithubOwner>_<GithubProject>_<BranchName>_<Version>_<Build>_<SHA>`  
We use [buildgoggles](https://www.npmjs.com/package/buildgoggles) to generate these tags for our nodejs projects and a customized version of [drone-docker](https://hub.docker.com/r/leankit/drone-docker/) to publish programmatic tags to https://hub.docker.com.

#### Catalog Format
This creates a catalog for a single github project. Creates a Entry for each branch with a version of each entry for every build ordered by version/build in each Docker Hub tag.

```bash
base/                        # Go Templates (Create/Edit These)
  |_ catalogIcon.(png|svg)   # Copy for branch/catalogIcon.(png|svg)
  |_ config.tmpl             # Template for branch/config.yml
  |_ rancher-compose.tmpl    # Template for branch/0/rancher-compose.yml
  |_ docker-compose.tmpl     # Template for branch/0/docker-compose.yml

templates/                   # Generated Catalog (don't manually edit)
  |_ master/                 # based on branch name (master for Release Tag)
    |_ config.yml            # Catalog Entry config
    |_ catalogIcon.png       # icon image
    |_ 0/                    # builds for a branch
      |_ docker-compose.yml  # Entry/Build docker-compose.yml
      |_ rancher-compose.yml # Entry/Build rancher-compose.yml
    |_ 1/
      |_ docker-compose.yml
      |_ rancher-compose.yml
  |_ feature-branch/
    |_ config.yml
    |_ catalogIcon.png
    |_ 0/
      |_ docker-compose.yml
      |_ rancher-compose.yml
  ...
```


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
    image: leankit/drone-cowpoke:logging
    docker_username: $$DOCKER_USER
    docker_password: $$DOCKER_PASS
    docker_repo: your-docker-repository
    catalog_repo: your-catalog-repository
    github_token: $$GITHUB_TOKEN
    github_user: $$GITHUB_USER
    github_email: $$GITHUB_EMAIL
    cowpoke_url: cowpoke-url
    rancher_catalog_name: catalog-name-in-rancher
    bearer_token: the-cowpoke-api-key
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

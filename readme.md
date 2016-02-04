drone-cowpoke
=============

This drone plugin provides a service to post a docker image name to cowpoke.

## Build
A build script is included which will produce a static Go binary that will work with Alpine. It then creates the docker image for you.

## Example

```sh
./drone-cowpoke <<EOF
{
	"workspace": {
		"path": "./docker/src"
	},
	"vargs": {
		"cowpoke_url": "http://cowpoke",
		"cowpoke_port": "8080"
	}
}
EOF
```

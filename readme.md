drone-cowpoke
=============

This drone plugin provides a service to post a docker image name to cowpoke.

## Example

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

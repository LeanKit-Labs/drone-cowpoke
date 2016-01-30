drone-cowpoke
=============

This drone plugin provides a service to post a docker image name to cowpoke.

## Notes
*AS IS:*  
The cowpoke_url and the image value loaded from the .docker.json will be appended. There is no magic parsing. As the example reads below, the PUT request will be posted to `http://cowpoke.com/your/image:tag`. This will obviously need some tweaking. ;)

## Example

```sh
./drone-cowpoke <<EOF
{
  "vargs": {
    "cowpoke_url": "http://cowpoke.com/",
    "docker_json": "./testdata/.docker.json"
  }
}
EOF

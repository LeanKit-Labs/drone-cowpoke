package main

import (
  "encoding/json"
  "fmt"
  "io/ioutil"
  "net/http"
  "net/url"
  "os"
  "path/filepath"

  "github.com/drone/drone-plugin-go/plugin"
)

type DroneCowpoke struct {
  Url string `json:"cowpoke_url"`
  Port string `json:"cowpoke_port"`
}

type ImageJson struct {
	Image string `json:"image"`
}

func main() {
  workspace := plugin.Workspace{}
  vargs := DroneCowpoke{}

  plugin.Param("workspace", &workspace)
  plugin.Param("vargs", &vargs)
  plugin.MustParse()

  if len(vargs.Url) == 0 {
    fmt.Println("no cowpoke url was specified")
    os.Exit(1)
  }

  if len(vargs.Port) == 0 {
    fmt.Println("no cowpoke url was specified")
    os.Exit(1)
  }

  fmt.Println("loading image data from", filepath.Join(workspace.Path, ".docker.json"))
  image := GetImageName(filepath.Join(workspace.Path, ".docker.json"))

  var cowpokeUrl = vargs.Url + ":" + vargs.Port + "/"
  ExecutePut(cowpokeUrl + url.QueryEscape(image));
}

func ExecutePut(putUrl string) {
  fmt.Println("executing a PUT request for:", putUrl)

  client := &http.Client{}
  request, err := http.NewRequest("PUT", putUrl, nil)

  response, err := client.Do(request)
  if err != nil {
    fmt.Println("error executing request:", err)
    os.Exit(1)
  }

  defer response.Body.Close()
    contents, err := ioutil.ReadAll(response.Body)
    if err != nil {
      fmt.Println("error reading response:", err)
    }

    fmt.Println("response status code:", response.StatusCode)
    fmt.Println("content:", string(contents))
}

func GetImageName(path string) string {
  file, err := ioutil.ReadFile(path)

  if err != nil {
		fmt.Println("error opening json file", err)
	}

	var jsonobject ImageJson
	json.Unmarshal(file, &jsonobject)

	return jsonobject.Image
}

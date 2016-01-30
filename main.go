package main

import (
  "fmt"
  "io/ioutil"
  "net/http"
  "os"

  "github.com/drone/drone-plugin-go/plugin"
)

type Cowpoke struct {
	Url string `json:"cowpoke_url"`
  DockerJson string `json:"docker_json"`
}

func main() {
  vargs := Cowpoke{}
	plugin.Param("vargs", &vargs)
	plugin.MustParse()

  if len(vargs.Url) == 0 {
    fmt.Println("no cowpoke url was specified")
    os.Exit(1)
	}

  if len(vargs.DockerJson) == 0 {
    fmt.Println("no docker json path was specified")
    os.Exit(1)
  }

  image := GetImageName(vargs.DockerJson)

  // this will need encoding or some other means of parsing to send it correctly
  ExecutePut(vargs.Url + image);
}

func ExecutePut(url string) {
  fmt.Println("executing a PUT request for:", url)

  client := &http.Client{}
  request, err := http.NewRequest("PUT", url, nil)

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

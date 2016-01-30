package main

import (
  "fmt"
  "io/ioutil"
  "net/http"
  "os"

  "github.com/jaymedavis/drone-cowpoke"
)

func executePut(url string) {
  fmt.Println("executing a PUT request for:", url)

  client := &http.Client{}
  request, err := http.NewRequest("PUT", url, nil)
  // do we need auth?
  // request.SetBasicAuth("admin", "admin")
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

func main() {
  var image string = "arob/demo-2:arobson_demo-2_master_0.1.0_2_abcdef";
  executePut("http://url.com/" + image);
}

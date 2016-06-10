package main

import (
  "encoding/json"
  "fmt"
  "io/ioutil"
  "strings"
  "net/http"
  "net/url"
  "os"
  "path/filepath"
  "github.com/drone/drone-plugin-go/plugin"
)

type Cowpoke struct {
  Url string `json:"cowpoke_url"`
  Port int `json:"cowpoke_port"`
}

type ImageJson struct {
  Image string `json:"image"`
}

func main() {
  fmt.Println("starting drone-cowpoke...")

  workspace := plugin.Workspace{}
  vargs := Cowpoke{}

  plugin.Param("workspace", &workspace)
  plugin.Param("vargs", &vargs)
  plugin.MustParse()

  if len(vargs.Url) == 0 {
    fmt.Println("no cowpoke url was specified")
    os.Exit(1)
  }

  if vargs.Port == 0 {
    fmt.Println("no cowpoke port was specified")
    os.Exit(1)
  }

  fmt.Println("loading image data from", filepath.Join(workspace.Path, ".docker.json"))
  image := GetImageName(filepath.Join(workspace.Path, ".docker.json"))

  if(len(image) <= 0) {
    fmt.Println("image load failed from .docker.json")
    os.Exit(1)
  }
  SendCowpoke(image, vargs.Url, vargs.Port)
}


  
//SendCowpoke ... Send the request to copwoke
func SendCowpoke(image string, baseURL string, port int ) int {
  //leankit%2Fcore-leankit-api%3ABanditSoftware_core-leankit-api_feature-exit-when-no-redis_4.7.1_11_6563de12
  imageParts := strings.Split(image, ":");

  //if there is a tag
  if (len(imageParts) == 2) {
    tag := imageParts[1]; //get it
    //get the parts expecting OWNER_REPO_BRANCH_VERSION_BUILD_COMMIT
    //as the branch could have underscores it needs to be parsed last.
    index := strings.Index(tag, "_");
    if index == -1 {
      return 1;
    }
    owner := tag[:index];
    remainingTag := tag[index + 1:]

    index = strings.Index(remainingTag, "_")   
    if index == -1 {
      return 1;
    }
    repo := remainingTag[:index]
    remainingTag = tag[index + 1:]

    //get the end part of the tag
    index = strings.LastIndex(remainingTag, "_" )
    if index == -1 {
      return 1;
    }
    commit := remainingTag[index + 1:]
    remainingTag = tag[:index]
    
    index = strings.LastIndex(remainingTag, "_" )
    if index == -1 {
      return 1;
    }
    build := remainingTag[index + 1:]
    remainingTag = tag[:index]
    
    index = strings.LastIndex(remainingTag, "_" )
    if index == -1 {
      return 1;
    }
    version := remainingTag[index:]
    remainingTag = tag[:index]


    //now we can get the branch
    branch := remainingTag
    if (( owner != "" ) && ( repo != "" ) && ( branch != "" ) && ( version != "" ) && ( build != "" ) && ( commit != "" )) {
        cowpokeURL := fmt.Sprintf("%s:%d/api/environment/", baseURL, port)
        fmt.Println("cowpoke url set to:", cowpokeURL)
        fmt.Println(".docker.json value being posted:", image)
        ExecutePut(cowpokeURL + url.QueryEscape(image));
        fmt.Println("finished drone-cowpoke.")
        return 0;
    } else {
      fmt.Println("Production tag detected no cowpoke request made.")
      return 1
    }
  } else {
    fmt.Println("No tag specified cowpoke request not made.")
    return 2
  }
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

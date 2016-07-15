package main

import (
	"fmt"
	"github.com/drone/drone-plugin-go/plugin"
	"gopkg.in/yaml.v2"
	"strings"
	"io/ioutil"
	"net/http"
	"net/url"
	"bytes"
	"os"
	"time"
	"path/filepath"
)

type Cowpoke struct {
	Url         string `json:"cowpoke_url"`
	Port        int    `json:"cowpoke_port"`
	DockerOwner string `json:"docker_owner"`
	DockerRepo  string `json:"docker_repo"`
	CatalogUpgrade  bool `json:"cowpoke_catalog_upgrade"`
	Catalog  string `json:"cowpoke_catalog"`
}

type TagsYaml struct {
	Tags []string `yaml:"tags"`
}

func CheckImage( image string) bool {
	imageParts := strings.Split(image, ":");

	//if there is a tag
	if (len(imageParts) != 2) {
		return false;
	}

	tag := imageParts[1]; //get it
	//get the parts expecting OWNER_REPO_BRANCH_VERSION_BUILD_COMMIT
	//as the branch could have underscores it needs to be parsed last.
	index := strings.Index(tag, "_");
	if index == -1 {
		return false;
	}

	owner := tag[:index];
	remainingTag := tag[index + 1:]

	index = strings.Index(remainingTag, "_")

	if index == -1 {
	  return false;
	}
    repo := remainingTag[:index]
    remainingTag = remainingTag[index + 1:]

    //get the end part of the tag
    index = strings.LastIndex(remainingTag, "_" )
    if index == -1 {
      return false;
    }
    commit := remainingTag[index + 1:]
    remainingTag = remainingTag[:index]
    
    index = strings.LastIndex(remainingTag, "_" )
    if index == -1 {
      return false;
    }
    build := remainingTag[index + 1:]
    remainingTag = remainingTag[:index]
    
    index = strings.LastIndex(remainingTag, "_" )
    if index == -1 {
      return false;
    }
    version := remainingTag[index:]
    branch := remainingTag[:index]

	if (owner == "") || (repo == "") || (branch == "") || (version == "") || (build == "") || (commit =="")  {
		return false;
	}

	return true
}

func main() {
	fmt.Println("starting drone-cowpoke...")

	var workspace = plugin.Workspace{}
	var droneRepo = plugin.Repo{}
	var cowpoke = Cowpoke{}
	var image string
	var owner string
	var repo string

	plugin.Param("workspace", &workspace)
	plugin.Param("repo", &droneRepo)
	plugin.Param("vargs", &cowpoke)
	plugin.MustParse()

	if len(cowpoke.Url) == 0 {
		fmt.Println("no cowpoke url was specified")
		os.Exit(1)
	}

	if cowpoke.Port == 0 {
		fmt.Println("no cowpoke port was specified")
		os.Exit(1)
	}

	if len(cowpoke.DockerOwner) != 0 {
		owner = cowpoke.DockerOwner
	} else {
		owner = droneRepo.Owner
	}

	if len(cowpoke.DockerRepo) != 0 {
		repo = cowpoke.DockerRepo
	} else {
		repo = droneRepo.Name
	}

	var cowpokeUrl = fmt.Sprintf("%s:%d/api/environment/", cowpoke.Url, cowpoke.Port)
	fmt.Println("Cowpoke url set to:", cowpokeUrl)
	fmt.Println("Loading tags from: ", filepath.Join(workspace.Path, ".droneTags.yml"))
	tags := GetTags(filepath.Join(workspace.Path, ".droneTags.yml"))
	if len(tags) == 0 {
		fmt.Println("No tags found. Nothing to poke.")
	}
	for _ , tag := range tags {
		image = fmt.Sprintf("%s/%s:%s", owner, repo, tag)
		if CheckImage(image) {
			fmt.Println("Poking environments with image:", image)
			if (cowpoke.CatalogUpgrade == true) {
				jsonStr := fmt.Sprintf("{\"rancher_catalog\": \"%s\", \"docker_image\" : \"%s\"}", cowpoke.Catalog, image);
				ExecutePut(cowpokeUrl + "catalog", jsonStr)
			} else {
				ExecutePut(cowpokeUrl + url.QueryEscape(image), "{}")
			}
		} else {
			fmt.Println("Tag not formated like dev and no services will be upgraded with image: ", image)
		}
		
	}
	fmt.Println("finished drone-cowpoke.")
}

func ExecutePut(putUrl string, jsonStr string) {
	fmt.Println("executing a PUT request for:", putUrl)

	client := &http.Client{
	    Timeout: time.Second * 60,
	}
	request, err := http.NewRequest("PUT", putUrl, bytes.NewBuffer([]byte(jsonStr)))
	request.Close = true
	request.Header.Set("Content-Type", "application/json")

	response, err := client.Do(request)
	if err != nil {
		fmt.Println("error executing request:", response, err)
		os.Exit(0)
	}

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("error reading response:", err)
	}

	fmt.Println("response status code:", response.StatusCode)
	fmt.Println("content:", string(contents))
}

>>>>>>> Added code to do catalog update
func GetTags(path string) []string {
	file, err := ioutil.ReadFile(path)

	if err != nil {
		fmt.Println("error opening .droneTags.yml file", err)
		os.Exit(1)
	}

	var droneTags TagsYaml
	yaml.Unmarshal(file, &droneTags)
	return droneTags.Tags
}

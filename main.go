package main

import (
	"fmt"
	"github.com/drone/drone-plugin-go/plugin"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

type Cowpoke struct {
	Url         string `json:"cowpoke_url"`
	Port        int    `json:"cowpoke_port"`
	DockerOwner string `json:"docker_owner"`
	DockerRepo  string `json:"docker_repo"`
}

type TagsYaml struct {
	Tags []string `yaml:"tags"`
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
	for _, tag := range tags {
		image = fmt.Sprintf("%s/%s:%s", owner, repo, tag)

		fmt.Println("Poking environments with image:", image)
		ExecutePut(cowpokeUrl + url.QueryEscape(image))
	}
	fmt.Println("finished drone-cowpoke.")
}

func ExecutePut(putUrl string) {
	fmt.Println("executing a PUT request for:", putUrl)

	client := &http.Client{}
	request, err := http.NewRequest("PUT", putUrl, nil)
	request.Close = true

	response, err := client.Do(request)
	if err != nil {
		fmt.Println("error executing request:", response, err)
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

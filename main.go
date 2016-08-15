package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"text/template"
	"time"

	yaml "gopkg.in/yaml.v2"

	"github.com/blang/semver"
	"github.com/drone/drone-plugin-go/plugin"
	"github.com/heroku/docker-registry-client/registry"
)

const (
	baseDir                    string = "/rancher-catalog"
	repoDir                    string = "/rancher-catalog/repo"
	templateDir                string = "/rancher-catalog/repo/base"
	dockerComposeTemplateFile  string = "/rancher-catalog/repo/base/docker-compose.tmpl"
	rancherComposeTemplateFile string = "/rancher-catalog/repo/base/rancher-compose.tmpl"
	configTemplateFile         string = "/rancher-catalog/repo/base/config.tmpl"
	iconFileBase               string = "/rancher-catalog/repo/base/catalogIcon"
)

// catalog struct
type catalog struct {
	vargs     vargs
	workspace plugin.Workspace
	repo      plugin.Repo
	build     plugin.Build
}

// vargs strct
type vargs struct {
	DockerRepo         string `json:"docker_repo"`
	DockerUsername     string `json:"docker_username"`
	DockerPassword     string `json:"docker_password"`
	DockerURL          string `json:"docker_url"`
	CatalogRepo        string `json:"catalog_repo"`
	GitHubToken        string `json:"github_token"`
	GitHubUser         string `json:"github_user"`
	GitHubEmail        string `json:"github_email"`
	CowpokeURL         string `json:"cowpoke_url"`
	RancherCatalogName string `json:"rancher_catalog_name"`
}

// tagsByBranch struct
type tagsByBranch struct {
	branches map[string]branch
}

type TagsYaml struct {
	Tags []string `yaml:"tags"`
}

// branch struct
type branch struct {
	versions map[string]version
}

// version struct
type version struct {
	builds map[int]*Tag
}

// Tag struct
type Tag struct {
	Tag     string
	Count   int
	Owner   string
	Project string
	Branch  string
	Version string
	Build   int
	SHA     string
}

func getTagsFromYaml(workspace plugin.Workspace) []string {
	path := filepath.Join(workspace.Path, "docker/src/.droneTags.yml")
	file, err := ioutil.ReadFile(path)

	if err != nil {
		fmt.Println("error opening .droneTags.yml file", err)
		os.Exit(1)
	}

	var droneTags TagsYaml
	yaml.Unmarshal(file, &droneTags)
	return droneTags.Tags
}

func stringInSlice(str string, list []string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}

func check(val string, err string) {
	if len(val) == 0 {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	fmt.Println("starting drone-rancher-catalog...")

	var catalog = catalog{}
	plugin.Param("workspace", &catalog.workspace)
	plugin.Param("repo", &catalog.repo)
	plugin.Param("build", &catalog.build)
	plugin.Param("vargs", &catalog.vargs)
	plugin.MustParse()

	check(catalog.vargs.DockerRepo, "ERROR: docker_repo: Docker Registry Repo to read tags from, not specified")
	check(catalog.vargs.DockerUsername, "ERROR: docker_username: Docker Registry Username not specified")
	check(catalog.vargs.DockerPassword, "ERROR: docker_password: Docker Registry Password not specified")
	check(catalog.vargs.CatalogRepo, "ERROR: catalog_repo: GitHub Catalog Repo not specified")
	check(catalog.vargs.GitHubToken, "ERROR: github_token: GitHub User Token not specified")
	check(catalog.vargs.CowpokeURL, "ERROR: cowpoke_url: cowpoke url not specified")
	check(catalog.vargs.RancherCatalogName, "ERROR: rancher_catalog_name: catalog name in rancher is not specified")

	if len(catalog.vargs.DockerURL) == 0 {
		catalog.vargs.DockerURL = "https://registry.hub.docker.com/"
	}
	if len(catalog.vargs.GitHubUser) == 0 {
		catalog.vargs.GitHubUser = catalog.build.Author
	}
	if len(catalog.vargs.GitHubEmail) == 0 {
		catalog.vargs.GitHubEmail = catalog.build.Email
	}

	// create a dir outside the workspace
	if !exists(baseDir) {
		os.Mkdir(baseDir, 0755)
	}

	catalog.cloneCatalogRepo()
	os.Chdir(repoDir)
	catalog.gitConfigureEmail()
	catalog.gitConfigureUser()

	if !exists("./templates") {
		os.Mkdir("./templates", 0755)
	}

	dockerComposeTmpl := catalog.parseTemplateFile(dockerComposeTemplateFile)
	rancherComposeTmpl := catalog.parseTemplateFile(rancherComposeTemplateFile)
	configTmpl := catalog.parseTemplateFile(configTemplateFile)
	upgradeTags := getTagsFromYaml(catalog.workspace)
	tags := catalog.getTags()
	tbb := catalog.TagsByBranch(tags)

	var cowpokeRequests []*http.Request

	fmt.Println("Creating Catalog Templates for:")
	for branch := range tbb.branches {
		var count int
		var last *Tag

		// create branch dir
		branchDir := fmt.Sprintf("./templates/%s", branch)
		if !exists(branchDir) {
			os.Mkdir(branchDir, 0755)
		}

		// sort semver so we can count builds in a feature branch
		var vKeys []semver.Version
		for k := range tbb.branches[branch].versions {
			version, err := semver.Parse(k)
			if err != nil {
				fmt.Printf("Error parsing version %v \n", err)
				continue
			}
			vKeys = append(vKeys, version)
		}
		semver.Sort(vKeys)
		for _, version := range vKeys {
			// sort builds to count in order
			var bKeys []int
			ver := version.String()
			for k := range tbb.branches[branch].versions[ver].builds {
				bKeys = append(bKeys, k)
			}
			sort.Ints(bKeys)

			for _, build := range bKeys {
				tbb.branches[branch].versions[ver].builds[build].Count = count

				// create dir structure
				buildDir := fmt.Sprintf("%s/%d", branchDir, count)
				if !exists(buildDir) {
					fmt.Printf("  %d:%s %s-%d\n", count, branch, ver, build)
					os.Mkdir(buildDir, 0755)
				}

				// create docker-compose.yml and rancher-compose.yml from template
				// don't generate files if they already exist
				dockerComposeTarget := fmt.Sprintf("%s/docker-compose.yml", buildDir)
				if !exists(dockerComposeTarget) {
					catalog.executeTemplate(dockerComposeTarget, dockerComposeTmpl, tbb.branches[branch].versions[ver].builds[build])
				}
				rancherComposeTarget := fmt.Sprintf("%s/rancher-compose.yml", buildDir)
				if !exists(rancherComposeTarget) {
					catalog.executeTemplate(rancherComposeTarget, rancherComposeTmpl, tbb.branches[branch].versions[ver].builds[build])
				}

				last = tbb.branches[branch].versions[ver].builds[build]
				count++
			}
			if stringInSlice(last.Tag, upgradeTags) {
				cowpokeRequests = append(cowpokeRequests, cowpokeRequest(count, branch, catalog.vargs.CatalogRepo, catalog.vargs.RancherCatalogName, catalog.vargs.GitHubToken, catalog.vargs.CowpokeURL))
			}

		}

		// create config.yml from temlplate
		configTarget := fmt.Sprintf("%s/config.yml", branchDir)
		catalog.executeTemplate(configTarget, configTmpl, last)

		// Icon file
		copyIcon(iconFileBase, branchDir)
	}
	// TODO: Delete dir/files if tags don't exist anymore. Need to maintian build dir numbering

	if catalog.gitChanged() {
		catalog.addCatalogRepo()
		catalog.commitCatalogRepo()
		catalog.pushCatalogRepo()
	}
	client := &http.Client{
		Timeout: time.Second * 60,
	}
	for _, request := range cowpokeRequests {
		response, err := client.Do(request)
		if err != nil {
			fmt.Println("error executing request:", response, err)
			os.Exit(0)
		}
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Println("error reading response:", err)
		}
		fmt.Println("response status code:", response.StatusCode)
		fmt.Println("content:", string(contents))
		response.Body.Close()
	}
	fmt.Println("... Finished drone-rancher-catalog")
}

//calls cowpoke after catalog is built
func cowpokeRequest(catalogNo int, branchName string, CatalogRepo string, rancherCatalogName string, token string, CowpokeURL string) *http.Request {
	var jsonStr = []byte(fmt.Sprintf(`{"catalog":"%s","rancherCatalogName":"%s","githubToken":"%s","catalogVersion":"%s","branch":"%s"}`, CatalogRepo, rancherCatalogName, token, strconv.Itoa(catalogNo), branchName))
	request, err := http.NewRequest("PATCH", CowpokeURL+"/api/stack", bytes.NewBuffer(jsonStr))
	request.Header.Set("Content-Type", "application/json")
	if err != nil {
		fmt.Println("Error making request object to cowpoke")
		panic(err)
	}
	request.Close = true
	return request
}

func (c *catalog) getTags() []string {
	hub, err := registry.New(c.vargs.DockerURL, c.vargs.DockerUsername, c.vargs.DockerPassword)
	if err != nil {
		fmt.Println("ERROR: Could not Contact Docker Registry", err)
		os.Exit(1)
	}
	tags, err := hub.Tags(c.vargs.DockerRepo)
	if err != nil {
		fmt.Println("ERROR: Getting tags", err)
		os.Exit(1)
	}
	return tags
}

// parseTag Returns a Tag object from a buildgoogles style tag
func (c *catalog) parseTag(t string) *Tag {
	var tag = &Tag{}
	featureRe := regexp.MustCompile(fmt.Sprintf(`^%s_%s_`, c.repo.Owner, c.repo.Name))
	releaseRe := regexp.MustCompile(`^v\d+\.\d+\.\d+$`)
	// Skip forks and other nonsense tags
	switch {
	case featureRe.MatchString(t):
		var build string
		// fmt.Println("Found Feature Branch Tag", t)
		tagParts := strings.Split(t, "_")
		// shift the owner and project from the front
		// pop the sha, build, and version from the back
		// join whats left into the branch
		tag.Tag = t
		tag.Owner, tagParts = tagParts[0], tagParts[1:]
		tag.Project, tagParts = tagParts[0], tagParts[1:]
		tag.SHA, tagParts = tagParts[len(tagParts)-1], tagParts[:len(tagParts)-1]
		build, tagParts = tagParts[len(tagParts)-1], tagParts[:len(tagParts)-1]
		tag.Build, _ = strconv.Atoi(build)
		tag.Version, tagParts = tagParts[len(tagParts)-1], tagParts[:len(tagParts)-1]
		tag.Branch = strings.Join(tagParts, "_")
	case releaseRe.MatchString(t):
		// fmt.Println("Found Release Tag", t)
		tag.Tag = t
		tag.Owner = c.repo.Owner
		tag.Project = c.repo.Name
		tag.Branch = "master"
		tag.Build = 1
		tag.SHA = ""
		versionRe := regexp.MustCompile(`^v`)
		tag.Version = versionRe.ReplaceAllString(t, "")
	default:
		return nil
	}
	return tag
}

// tagsByBranch break down tag list and return a tagsByBranch object
func (c *catalog) TagsByBranch(tags []string) *tagsByBranch {
	tbb := &tagsByBranch{}
	tbb.branches = make(map[string]branch)
	for _, tg := range tags {
		t := c.parseTag(tg)
		if t == nil {
			continue
		}
		if _, present := tbb.branches[t.Branch]; !present {
			tbb.branches[t.Branch] = branch{
				versions: make(map[string]version),
			}
		}
		if _, present := tbb.branches[t.Branch].versions[t.Version]; !present {
			tbb.branches[t.Branch].versions[t.Version] = version{
				builds: make(map[int]*Tag),
			}
		}
		if _, present := tbb.branches[t.Branch].versions[t.Version].builds[t.Build]; !present {
			tbb.branches[t.Branch].versions[t.Version].builds[t.Build] = t
		}
	}
	return tbb
}

func exists(f string) bool {
	if _, err := os.Stat(f); os.IsNotExist(err) {
		return false
	}
	return true
}

func (c *catalog) cloneCatalogRepo() {
	gitHubURL := fmt.Sprintf("https://%s:x-oauth-basic@github.com/%s.git", c.vargs.GitHubToken, c.vargs.CatalogRepo)

	fmt.Println("Cloning Rancher-Catalog repo:", c.vargs.CatalogRepo)
	// clear if existing and git clone target repo
	os.RemoveAll(repoDir)

	out, err := exec.Command("git", "clone", gitHubURL, repoDir).CombinedOutput()
	fmt.Println(string(out))
	if err != nil {
		fmt.Printf("ERROR: Failed to Clone Repo %v\n", err)
		os.Exit(1)
	}
}

func (c *catalog) addCatalogRepo() {
	cmd := exec.Command("git", "add", "-A")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Printf("ERROR: Failed to git add %v\n", err)
		os.Exit(1)
	}
}

func (c *catalog) commitCatalogRepo() {
	message := fmt.Sprintf("'Update from Drone Build: %d'", c.build.Number)
	cmd := exec.Command("git", "commit", "-m", message)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Printf("ERROR: Failed to git commit %v\n", err)
		os.Exit(1)
	}
}

func (c *catalog) pushCatalogRepo() {
	cmd := exec.Command("git", "push")
	err := cmd.Run()
	// Not showing output, bleeds the API key
	if err != nil {
		fmt.Printf("ERROR: Failed to git push %v\n", err)
		os.Exit(1)
	}
}

func (c *catalog) parseTemplateFile(file string) *template.Template {
	name := filepath.Base(file)
	tmpl, err := template.New(name).ParseFiles(file)
	if err != nil {
		fmt.Printf("ERROR: Failed parse template %v\n", err)
		os.Exit(1)
	}
	return tmpl
}

func (c *catalog) executeTemplate(target string, tmpl *template.Template, tag *Tag) {
	targetFile, err := os.Create(target)
	if err != nil {
		fmt.Printf("ERROR: Failed to open file %v\n", err)
		os.Exit(1)
	}
	err = tmpl.Execute(targetFile, tag)
	if err != nil {
		fmt.Printf("ERROR: Failed execute template %v\n", err)
		os.Exit(1)
	}
	targetFile.Close()
}

// copy src.* (repo/base/catalogIcon.*) to dest directory
func copy(src string, dest string) {
	cmd := exec.Command("cp", src, dest)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Printf("ERROR: Failed to cp %v\n", err)
		os.Exit(1)
	}
}

func copyIcon(src string, dest string) {
	dir := filepath.Dir(src)
	base := filepath.Base(src)
	// find files in dir that match base
	iconRe := regexp.MustCompile(fmt.Sprintf(`^%s`, base))
	files, _ := ioutil.ReadDir(dir)
	for _, f := range files {
		if iconRe.MatchString(f.Name()) {
			name := fmt.Sprintf("%s/%s", dir, f.Name())
			copy(name, dest)
		}
	}
}

func (c *catalog) gitConfigureEmail() {
	cmd := exec.Command("git", "config", "user.email", c.vargs.GitHubEmail)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Printf("ERROR: Failed to git config %v\n", err)
		os.Exit(1)
	}
}

func (c *catalog) gitConfigureUser() {
	cmd := exec.Command("git", "config", "user.name", c.vargs.GitHubUser)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Printf("ERROR: Failed to git config %v\n", err)
		os.Exit(1)
	}
}

// returns true if there are files that need to be commited.
func (c *catalog) gitChanged() bool {
	cmd := exec.Command("git", "status", "--porcelain")
	out, err := cmd.Output()
	if err != nil {
		fmt.Printf("ERROR: Failed to git status %v\n", err)
		os.Exit(1)
	}
	// no output means no changes.
	if len(out) == 0 {
		fmt.Println("No files changed.")
		return false
	}
	fmt.Println("Files changed, add/commit/push changes.")
	return true
}

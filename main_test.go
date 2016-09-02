//tests will be reworked after this proof of concept plugin
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"testing"

	"github.com/drone/drone-plugin-go/plugin"
	"github.com/franela/goblin"
)

func TestHookImage(t *testing.T) {

	g := goblin.Goblin(t)

	g.Describe("file exists", func() {
		g.It("should find file", func() {
			g.Assert(exists("./main.go")).Equal(true)
		})
		g.It("should not find file", func() {
			g.Assert(exists("./DNE.go")).Equal(false)
		})
	})

	g.Describe("string in slice", func() {
		g.It("should find string", func() {
			g.Assert(stringInSlice("findme", []string{"findme"})).Equal(true)
		})
		g.It("should not find string", func() {
			g.Assert(stringInSlice("findme", []string{"nope"})).Equal(false)
		})
	})

	g.Describe("Make a request object to github to check for catalog existance", func() {
		g.It("should return the correct request", func() {
			catalogNo := 1
			branchName := "test"
			CatalogRepo := "owner/repo"
			token := "secret"
			req := checkForRepCreationRequestBuilder(CatalogRepo, branchName, catalogNo, token)
			g.Assert(req.URL.String()).Equal(fmt.Sprintf("https://api.github.com/repos/%s/contents/templates/%s/%d", CatalogRepo, branchName, catalogNo))
			username, password, good := req.BasicAuth()
			g.Assert(good).Equal(true)
			g.Assert(username).Equal(token)
			g.Assert(password).Equal("x-oauth-basic")
		})
	})

	g.Describe("Make a request object for a request to cowpoke", func() {
		g.It("should return the correct request", func() {
			catalogNo := 1
			branchName := "test"
			CatalogRepo := "repo"
			rancherCatalogName := "catalog"
			token := "secret"
			CowpokeURL := "cowpoke.mydomain.io"
			BearerToken := "token"
			var args map[string]interface{}
			req := cowpokeRequest(catalogNo, branchName, CatalogRepo, rancherCatalogName, token, CowpokeURL, BearerToken)
			body, _ := ioutil.ReadAll(req.Body)
			json.Unmarshal(body, &args)
			g.Assert(req.Header.Get("bearer")).Equal(BearerToken)
			g.Assert(req.Header.Get("Content-Type")).Equal("application/json")
			g.Assert(args["catalog"].(string)).Equal(CatalogRepo)
			g.Assert(args["rancherCatalogName"].(string)).Equal(rancherCatalogName)
			g.Assert(args["githubToken"].(string)).Equal(token)
			g.Assert(args["catalogVersion"].(string)).Equal("1")
			g.Assert(args["branch"].(string)).Equal(branchName)
		})
	})

	g.Describe("Parse tag correctly", func() {
		g.It("should get the parts on a dev formated tag", func() {
			var catalog = catalog{}
			catalog.repo = plugin.Repo{}
			catalog.repo.Owner = "Leankit-Labs"
			catalog.repo.Name = "cowpoke-integration-test"
			tag := "Leankit-Labs_cowpoke-integration-test_master_0.0.1_28_452c8fe7"
			tagInfo := catalog.parseTag(tag)
			g.Assert(tagInfo.Tag).Equal(tag)
			g.Assert(tagInfo.Owner).Equal("Leankit-Labs")
			g.Assert(tagInfo.Project).Equal("cowpoke-integration-test")
			g.Assert(tagInfo.SHA).Equal("452c8fe7")
			g.Assert(tagInfo.Build).Equal(28)
			g.Assert(tagInfo.Version).Equal("0.0.1")
			g.Assert(tagInfo.Branch).Equal("master")
		})

		g.It("should get the parts on a prod tag", func() {
			var catalog = catalog{}
			catalog.repo = plugin.Repo{}
			catalog.repo.Owner = "Leankit-Labs"
			catalog.repo.Name = "cowpoke-integration-test"
			tag := "v5.13.0"
			tagInfo := catalog.parseTag(tag)
			g.Assert(tagInfo.Tag).Equal(tag)
			g.Assert(tagInfo.Owner).Equal("Leankit-Labs")
			g.Assert(tagInfo.Project).Equal("cowpoke-integration-test")
			g.Assert(tagInfo.SHA).Equal("")
			g.Assert(tagInfo.Build).Equal(1)
			g.Assert(tagInfo.Version).Equal("5.13.0")
			g.Assert(tagInfo.Branch).Equal("master")
		})
	})

	g.Describe("Tags by branch", func() {
		var tbb *tagsByBranch
		g.Before(func() {
			var catalog = catalog{}
			catalog.repo = plugin.Repo{}
			catalog.repo.Owner = "leankit-labs"
			catalog.repo.Name = "cowpoke-integration-test"
			taglist := []string{
				"v0.0.1",
				"latest",
				"testy_test",
				"leankit-labs_cowpoke-integration-test_lots_of_under_scores_0.1.1_33_cd9f3615",
				"leankit-labs_cowpoke-integration-test_lots_of_under_scores_0.1.2_33_cd9f3615",
				"leankit-labs_cowpoke-integration-test_under_score_0.1.0_33_cd9f3615",
				"leankit-labs_cowpoke-integration-test_under_score_0.1.0_34_cd9f3615"}
			tbb = catalog.TagsByBranch(taglist)
		})
		g.It("should have all branchs", func() {

			branchKeys := []string{}
			for k := range tbb.branches {
				branchKeys = append(branchKeys, k)
			}
			g.Assert(stringInSlice("master", branchKeys)).Equal(true)
			g.Assert(stringInSlice("under_score", branchKeys)).Equal(true)
			g.Assert(stringInSlice("lots_of_under_scores", branchKeys)).Equal(true)
			g.Assert(len(branchKeys)).Equal(3)

		})
		g.It("should have all versions", func() {
			versionKeys := []string{}
			for k := range tbb.branches["lots_of_under_scores"].versions {
				versionKeys = append(versionKeys, k)
			}
			g.Assert(stringInSlice("0.1.1", versionKeys)).Equal(true)
			g.Assert(stringInSlice("0.1.2", versionKeys)).Equal(true)
			g.Assert(len(versionKeys)).Equal(2)
		})
		g.It("should have all builds", func() {
			buildKeys := []string{}
			for k := range tbb.branches["under_score"].versions["0.1.0"].builds {
				buildKeys = append(buildKeys, strconv.Itoa(k))
			}
			g.Assert(stringInSlice("33", buildKeys)).Equal(true)
			g.Assert(stringInSlice("34", buildKeys)).Equal(true)
			g.Assert(len(buildKeys)).Equal(2)
		})
	})
}

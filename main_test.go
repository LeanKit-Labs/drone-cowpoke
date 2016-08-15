//tests will be reworked after this proof of concept plugin
package main

import (
	"encoding/json"
	"io/ioutil"
	"strconv"
	"testing"

	"github.com/drone/drone-plugin-go/plugin"
	"github.com/franela/goblin"
)

func TestHookImage(t *testing.T) {

	g := goblin.Goblin(t)

	g.Describe("Make a request object for a request to cowpoke", func() {
		g.It("should return the correct request", func() {
			catalogNo := 1
			branchName := "test"
			CatalogRepo := "repo"
			rancherCatalogName := "catalog"
			token := "secret"
			CowpokeURL := "cowpoke.mydomain.io"
			var args map[string]interface{}
			req := cowpokeRequest(catalogNo, branchName, CatalogRepo, rancherCatalogName, token, CowpokeURL)
			body, _ := ioutil.ReadAll(req.Body)
			json.Unmarshal(body, &args)
			g.Assert(args["catalog"].(string) == CatalogRepo)
			g.Assert(args["rancherCatalogName"].(string) == rancherCatalogName)
			g.Assert(args["githubToken"].(string) == token)
			g.Assert(args["catalogVersion"].(string) == string(catalogNo))
			g.Assert(args["branch"].(string) == branchName)

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
			g.Assert(tagInfo.Tag == tag)
			g.Assert(tagInfo.Owner == "Leankit-Labs")
			g.Assert(tagInfo.Project == "cowpoke-integration-test")
			g.Assert(tagInfo.SHA == "452c8fe7")
			g.Assert(tagInfo.Build == 28)
			g.Assert(tagInfo.Version == "0.0.1")
			g.Assert(tagInfo.Branch == "master")
		})

		g.It("should get the parts on a prod tag", func() {
			var catalog = catalog{}
			catalog.repo = plugin.Repo{}
			catalog.repo.Owner = "Leankit-Labs"
			catalog.repo.Name = "cowpoke-integration-test"
			tag := "v5.13.0"
			tagInfo := catalog.parseTag(tag)
			g.Assert(tagInfo.Tag == tag)
			g.Assert(tagInfo.Owner == "Leankit-Labs")
			g.Assert(tagInfo.Project == "cowpoke-integration-test")
			g.Assert(tagInfo.SHA == "")
			g.Assert(tagInfo.Build == 1)
			g.Assert(tagInfo.Version == "5.13.0")
			g.Assert(tagInfo.Branch == "master")
		})
	})

	g.Describe("file exists", func() {
		g.It("should find file", func() {
			g.Assert(exists("./main.go"))
		})
		g.It("should not find file", func() {
			g.Assert(exists("./DNE.go"))
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
			g.Assert(stringInSlice("master", branchKeys))
			g.Assert(stringInSlice("under_score", branchKeys))
			g.Assert(stringInSlice("lots_of_under_scores", branchKeys))
			g.Assert(len(branchKeys) == 3)

		})
		g.It("should have all versions", func() {
			versionKeys := []string{}
			for k := range tbb.branches["lots_of_under_scores"].versions {
				versionKeys = append(versionKeys, k)
			}
			g.Assert(stringInSlice("0.1.1", versionKeys))
			g.Assert(stringInSlice("0.1.2", versionKeys))
			g.Assert(len(versionKeys) == 2)
		})
		g.It("should have all builds", func() {
			buildKeys := []string{}
			for k := range tbb.branches["under_score"].versions["0.1.0"].builds {
				buildKeys = append(buildKeys, strconv.Itoa(k))
			}
			g.Assert(stringInSlice("33", buildKeys))
			g.Assert(stringInSlice("34", buildKeys))
			g.Assert(len(buildKeys) == 2)
		})
	})
}

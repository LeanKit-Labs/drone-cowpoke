//tests will be reworked after this proof of concept plugin
package main

import (
	"encoding/json"
	"io/ioutil"
	"testing"

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
}

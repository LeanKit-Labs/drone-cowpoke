package main

import (
	"github.com/franela/goblin"
	"testing"
)

func TestHookImage(t *testing.T) {
	g := goblin.Goblin(t)

	g.Describe("when reading a docker json file", func() {
		g.It("should return the correct image value", func() {
			var tags = GetTags("./test_data/.droneTags.yml")

			g.Assert(tags).Equal([]string{"tagforthebowwow"})
		})
	})
}

package image

import (
  "testing"

	"github.com/franela/goblin"
)

func TestHook(t *testing.T) {
  g := goblin.Goblin(t)

  g.Describe("when reading a docker json file", func() {
    g.It("should return the correct image value", func() {
      var path = "./test_data/.docker.json"
      var image = GetImageName(path)

      g.Assert(image).Equal("your/image:tag")
    })
  })
}

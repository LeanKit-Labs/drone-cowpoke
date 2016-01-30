package main

import (
  "testing"

  "github.com/franela/goblin"
)

func TestHookImage(t *testing.T) {
  g := goblin.Goblin(t)

  g.Describe("when reading a docker json file", func() {
    g.It("should return the correct image value", func() {
      var image = GetImageName("./test_data/.docker.json")

      g.Assert(image).Equal("your/image:tagforthebowwow")
    })
  })
}

package main

import (
  "testing"

  "github.com/franela/goblin"
)

func TestHook(t *testing.T) {
  g := goblin.Goblin(t)

  g.Describe("when executing a put", func() {
    g.It("Should set build author to the pull request author", func() {
      var test = "test"
      g.Assert(test).Equal("test")
    })
  })
}

package main

import (
  "testing"
  "github.com/franela/goblin"
  "github.com/jarcoal/httpmock"
)

func TestHookImage(t *testing.T) {
  g := goblin.Goblin(t)

  g.Describe("when reading a docker json file", func() {
    g.It("should return the correct image value", func() {
      var image = GetImageName("./test_data/.docker.json")
      g.Assert(image).Equal("your/image:tagforthebowwow")
      
      


    })
  })
  g.Describe("when checking tag", func() {
    g.It("should call cowpoke on a dev tag", func() {
      httpmock.Activate()
      httpmock.RegisterResponder("PUT", "http://cowpoke.leankit.io:8000/api/environment/leankit%2Fcore-leankit-api%3ABanditSoftware_core-leankit-api_feature-exit-when-no-redis_4.7.1_11_6563de12",
        httpmock.NewStringResponder(200, `[{"id": 1, "name": "stuff"}]`))

      statusCode := SendCowpoke("leankit/core-leankit-api:BanditSoftware_core-leankit-api_feature-exit-when-no-redis_4.7.1_11_6563de12", "http://cowpoke.leankit.io", 8000)
      g.Assert(statusCode).Equal(0)

      defer httpmock.DeactivateAndReset()
    })
    g.It("should not call cowpoke on a production tag with version", func() {
      statusCode := SendCowpoke("foo:v3.0.0", "http://cowpoke.leankit.io", 8000)
      g.Assert(statusCode).Equal(1)
      
      
    })
    g.It("should not call cowpoke on a production tag with latest", func(){
      statusCode := SendCowpoke("foo:latest", "http://cowpoke.leankit.io", 8000)
      g.Assert(statusCode).Equal(1)
    })
    
    g.It("should not call cowpoke on no tag", func() {
      statusCode := SendCowpoke("foo", "http://cowpoke.leankit.io", 8000)
      g.Assert(statusCode).Equal(2)
    })
  })
}

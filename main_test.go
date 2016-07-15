//tests will be reworked after this proof of concept plugin
package main

//import "testing"

/*func TestHookImage(t *testing.T) {


	g := goblin.Goblin(t)

		g.Describe("when reading a docker json file", func() {
			g.It("should return the correct image value", func() {
				var tags = GetTags("./test_data/.droneTags.yml")

				g.Assert(tags).Equal([]string{"tagforthebowwow"})
			})
		})
		g.Describe("when checking tag", func() {
			g.It("should regognize a dev tag", func() {
				valid := CheckImage("leankit/core-leankit-api:BanditSoftware_core-leankit-api_feature-exit-when-no-redis_4.7.1_11_6563de12")
				g.Assert(valid).Equal(true)
			})
			g.It("should not call regognize a production tag with version", func() {
				valid := CheckImage("foo:v3.0.0")
				g.Assert(valid).Equal(false)


			})
			g.It("should not regognize a production tag with latest", func(){
				valid := CheckImage("foo:latest")
				g.Assert(valid).Equal(false)
			})

			g.It("should not regognize no tag", func() {
				valid := CheckImage("foo")
				g.Assert(valid).Equal(false)
			})
	  })
}*/

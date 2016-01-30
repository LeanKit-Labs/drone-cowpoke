package image

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

var ImageJson struct {
    Image string `json:"image"`
}

func GetImageName(path string) {
  file, err := ioutil.ReadFile(path)

  if err != nil {
		fmt.Println("error opening json file", err)
	}

	var jsonobject ImageJson
	json.Unmarshal(file, &jsonobject)
	fmt.Println("got image value", jsonboject.Image)

	return jsonobject.Image
}

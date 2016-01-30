package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type ImageJson struct {
	Image string `json:"image"`
}

func GetImageName(path string) string {
  file, err := ioutil.ReadFile(path)

  if err != nil {
		fmt.Println("error opening json file", err)
	}

	var jsonobject ImageJson
	json.Unmarshal(file, &jsonobject)

	return jsonobject.Image
}

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

type DockerHubImageTags struct {
	Count    int         `json:"count"`
	Next     string      `json:"next"`
	Previous interface{} `json:"previous"`
	Results  []struct {
		Name     string `json:"name"`
		FullSize int    `json:"full_size"`
		Images   []struct {
			Size         int64       `json:"size"`
			Architecture string      `json:"architecture"`
			Variant      interface{} `json:"variant"`
			Features     interface{} `json:"features"`
			Os           string      `json:"os"`
			OsVersion    string      `json:"os_version"`
			OsFeatures   interface{} `json:"os_features"`
		} `json:"images"`
		ID          int         `json:"id"`
		Repository  int         `json:"repository"`
		Creator     int         `json:"creator"`
		LastUpdater int         `json:"last_updater"`
		LastUpdated time.Time   `json:"last_updated"`
		ImageID     interface{} `json:"image_id"`
		V2          bool        `json:"v2"`
	} `json:"results"`
}

func main() {
	imageName := os.Args[1]

	// /がない場合library/をつける
	fullImageName := addPrefix(imageName)

	// hub.docker.comのAPIを実行
	data := getResponse(fullImageName)

	// コンソールに出力
	for n := range data.Results {
		fmt.Println(data.Results[n].Name)
	}

}

func addPrefix(name string) string {
	if (strings.Count(name, "/") == 0) {
		return "library/" + name
	}
	return name
}

func getResponse(name string) *DockerHubImageTags {
	url := "https://hub.docker.com/v2/repositories/" + name + "/tags/?page_size=200"
	resp, _ := http.Get(url)
	defer resp.Body.Close()

	byteArray, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("ioutil.ReadAll error:", err)
		return nil
	}
	data := new(DockerHubImageTags)

	if err := json.Unmarshal(byteArray, data); err != nil {
		fmt.Println("JSON Unmarshal error:", err)
		return nil
	}

	return data
}

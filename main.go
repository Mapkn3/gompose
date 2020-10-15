package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/mapkn3/gompose/model"
	"github.com/mapkn3/gompose/net"
	"github.com/mapkn3/gompose/util"
)

func getInfo(baseURL string, username string, password string, c chan model.BuildDescription) {
	APIURLPostfix := "api/json"
	projectInfoResponse := net.DoRequestWithBasicAuth(baseURL+APIURLPostfix, username, password)
	var projectInfo model.ProjectInfo
	err := json.Unmarshal(projectInfoResponse, &projectInfo)
	util.Check(err)
	buildDescriptionResponse := net.DoRequestWithBasicAuth(projectInfo.LastSuccessfulBuild.URL+APIURLPostfix, username, password)
	var buildDescription model.BuildDescription
	err = json.Unmarshal(buildDescriptionResponse, &buildDescription)
	util.Check(err)
	c <- buildDescription
}

func main() {
	wd, err := os.Getwd()
	util.Check(err)
	configPath := filepath.Join(wd, "config.json")

	var config *model.Config
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Println("We're having problem opening file: ", configPath, ". Please, create a file using the 'example.config.json'")
		return
	}
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Println("The config file does not conform to the correct format, error:", err)
		return
	}

	rawCompose, err := ioutil.ReadFile(config.ComposePath)
	if err != nil {
		log.Println("Invalid path to docker-compose file in 'config.json'", config.ComposePath)
		return
	}
	composeStr := string(rawCompose)

	c := make(chan model.BuildDescription, 100)
	for _, job := range config.TrackedJobs {
		go getInfo(job.URL, config.Credential.Username, config.Credential.Token, c)
	}
	for i := 0; i < len(config.TrackedJobs); i++ {
		buildDescription := <-c
		for _, image := range buildDescription.GetImages() {
			re := regexp.MustCompile(strings.Split(image, ":")[0] + `:\S+`)
			composeStr = re.ReplaceAllString(composeStr, image)
		}
	}
	err = ioutil.WriteFile(config.ComposePath, []byte(composeStr), 0666)
	util.Check(err)
}

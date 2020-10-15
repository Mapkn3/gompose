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
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Panicf("Config file '%s' does not exist. Please, create a file using the 'example.config.json'", configPath)
	}
	config := model.GetConfigFromFile(configPath)

	if _, err := os.Stat(config.ComposePath); os.IsNotExist(err) {
		log.Panicf("Invalid path to docker-compose file in 'config.json': %s", config.ComposePath)
	}
	rawCompose, err := ioutil.ReadFile(config.ComposePath)
	util.Check(err)
	composeStr := string(rawCompose)

	login, token := config.Credential.Username, config.Credential.Token
	c := make(chan model.BuildDescription)
	for _, job := range config.TrackedJobs {
		go getInfo(job.URL, login, token, c)
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

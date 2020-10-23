package main

import (
	"encoding/json"
	"flag"
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
	log.Printf("Begin processing: %s", baseURL)
	APIURLPostfix := "api/json"
	projectInfoResponse := net.DoRequestWithBasicAuth(baseURL+APIURLPostfix, username, password)
	var projectInfo model.ProjectInfo
	err := json.Unmarshal(projectInfoResponse, &projectInfo)
	util.Check(err)
	buildDescriptionResponse := net.DoRequestWithBasicAuth(projectInfo.LastSuccessfulBuild.URL+APIURLPostfix, username, password)
	var buildDescription model.BuildDescription
	err = json.Unmarshal(buildDescriptionResponse, &buildDescription)
	util.Check(err)
	log.Printf("End processing: %s", baseURL)
	c <- buildDescription
}

func doProjectProcessing(project model.Project, done chan string) {
	rawCompose, err := ioutil.ReadFile(project.ComposePath)
	if err != nil {
		log.Println("Invalid path to docker-compose file [", project.ComposePath, "] in 'config.json' in project", project.Name)
		return
	}
	composeStr := string(rawCompose)

	c := make(chan model.BuildDescription, 100)
	for _, job := range project.TrackedJobs {
		go getInfo(job.URL, config.Credential.Username, config.Credential.Token, c)
	}
	for i := 0; i < len(project.TrackedJobs); i++ {
		buildDescription := <-c
		for _, image := range buildDescription.GetImages() {
			re := regexp.MustCompile(strings.Split(image, ":")[0] + `:\S+`)
			composeStr = re.ReplaceAllString(composeStr, image)
		}
	}
	err = ioutil.WriteFile(project.ComposePath, []byte(composeStr), 0666)
	util.Check(err)
	done <- project.Name
}

var config *model.Config

func main() {
	wd, err := os.Getwd()
	util.Check(err)
	defaultConfigPath := filepath.Join(wd, "config.json")

	configPath := flag.String("config", defaultConfigPath, "the path to the docker-compose.yaml")
	flag.Parse()

	log.Printf("The path to the config file: %s", *configPath)
	data, err := ioutil.ReadFile(*configPath)
	if err != nil {
		log.Println("We're having problem opening file: ", configPath, ". Please, create a file using the 'example.config.json'")
		return
	}
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Println("The config file does not conform to the correct format, error:", err)
		return
	}

	done := make(chan string, 100)
	for _, project := range config.Projects {
		log.Printf("Begin project processing for %s", project.Name)
		go doProjectProcessing(project, done)
	}
	for i := 0; i < len(config.Projects); i++ {
		projectName := <-done
		log.Printf("%s - done!", projectName)
	}
}

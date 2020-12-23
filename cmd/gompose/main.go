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
	"sync"

	"github.com/mapkn3/gompose/internal/model"
	"github.com/mapkn3/gompose/pkg/net"
	"github.com/mapkn3/gompose/pkg/util"
)

func getBuildDescription(jobURL string, username string, password string, buildDescriptionChan chan model.BuildDescription) {
	log.Printf("Begin getting build description: %s", jobURL)
	APIURLPostfix := "api/json"
	projectInfoResponse := net.DoRequestWithBasicAuth(jobURL+APIURLPostfix, username, password)
	var projectInfo model.ProjectInfo
	err := json.Unmarshal(projectInfoResponse, &projectInfo)
	util.Check(err)
	buildDescriptionResponse := net.DoRequestWithBasicAuth(projectInfo.LastSuccessfulBuild.URL+APIURLPostfix, username, password)
	var buildDescription model.BuildDescription
	err = json.Unmarshal(buildDescriptionResponse, &buildDescription)
	util.Check(err)
	log.Printf("End getting build description: %s", jobURL)
	buildDescriptionChan <- buildDescription
}

func projectProcessing(project model.Project, projectWG *sync.WaitGroup) {
	defer projectWG.Done()

	rawCompose, err := ioutil.ReadFile(project.ComposePath)
	if err != nil {
		log.Println("Invalid path to docker-compose file [", project.ComposePath, "] in 'config.json' in project", project.Name)
		return
	}
	composeStr := string(rawCompose)

	buildDescriptionChan := make(chan model.BuildDescription, 30)
	var buildDescriptionWG sync.WaitGroup

	for _, job := range project.TrackedJobs {
		buildDescriptionWG.Add(1)
		url := job.URL
		log.Printf("Get job URL: %s", url)
		go getBuildDescription(url, config.Credential.Username, config.Credential.Token, buildDescriptionChan)

	}
	go func() {
		buildDescriptionWG.Wait()
		close(buildDescriptionChan)
	}()
	for buildDescription := range buildDescriptionChan {
		n := len(buildDescription.GetImages())
		log.Printf("Get %d docker images: %v", n, buildDescription.GetImages())
		for _, image := range buildDescription.GetImages() {
			re := regexp.MustCompile(strings.Split(image, ":")[0] + `:\S+`)
			composeStr = re.ReplaceAllString(composeStr, image)
		}
		buildDescriptionWG.Done()
	}

	err = ioutil.WriteFile(project.ComposePath, []byte(composeStr), 0666)
	util.Check(err)
	log.Printf("%s - done!", project.Name)
}

var config *model.Config

func main() {
	wd, err := os.Getwd()
	util.Check(err)
	defaultConfigPath := filepath.Join(wd, "config.json")

	configPath := flag.String("config", defaultConfigPath, "the path to the config.json")
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

	var projectWG sync.WaitGroup
	for _, project := range config.Projects {
		projectWG.Add(1)
		log.Printf("Begin project processing for %s", project.Name)
		go projectProcessing(project, &projectWG)
	}
	projectWG.Wait()
}

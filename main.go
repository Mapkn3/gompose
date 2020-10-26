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

func jobProcessing(username string, password string, jobChan chan string, imageChan chan string, imageCountChan chan int, quit chan int) {
	log.Printf("Begin job processing")
	buildDescriptionChan := make(chan model.BuildDescription, 100)
	for {
		select {
		case url := <-jobChan:
			log.Printf("Get job URL: %s", url)
			go getBuildDescription(url, username, password, buildDescriptionChan)
		case buildDescription := <-buildDescriptionChan:
			n := len(buildDescription.GetImages())
			log.Printf("Get %d docker images: %v", n, buildDescription.GetImages())
			imageCountChan <- n
			for _, image := range buildDescription.GetImages() {
				imageChan <- image
			}
		case <-quit:
			log.Printf("End job processing")
			return
		}
	}
}

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

func projectProcessing(project model.Project, done chan string) {
	rawCompose, err := ioutil.ReadFile(project.ComposePath)
	if err != nil {
		log.Println("Invalid path to docker-compose file [", project.ComposePath, "] in 'config.json' in project", project.Name)
		return
	}
	composeStr := string(rawCompose)

	jobChan := make(chan string, 100)
	imageChan := make(chan string, 100)
	imageCountChan := make(chan int, 100)
	quit := make(chan int, 100)
	go jobProcessing(config.Credential.Username, config.Credential.Token, jobChan, imageChan, imageCountChan, quit)

	for _, job := range project.TrackedJobs {
		jobChan <- job.URL
	}
	for i := 0; i < len(project.TrackedJobs); i++ {
		imageCount := <-imageCountChan
		for k := 0; k < imageCount; k++ {
			image := <-imageChan
			re := regexp.MustCompile(strings.Split(image, ":")[0] + `:\S+`)
			composeStr = re.ReplaceAllString(composeStr, image)
		}
	}
	quit <- 0
	err = ioutil.WriteFile(project.ComposePath, []byte(composeStr), 0666)
	util.Check(err)
	done <- project.Name
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

	done := make(chan string, 100)
	for _, project := range config.Projects {
		log.Printf("Begin project processing for %s", project.Name)
		go projectProcessing(project, done)
	}
	for i := 0; i < len(config.Projects); i++ {
		projectName := <-done
		log.Printf("%s - done!", projectName)
	}
}

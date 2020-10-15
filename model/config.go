package model

import (
	"encoding/json"
	"io/ioutil"

	"github.com/mapkn3/gompose/util"
)

// JenkinsCredential class for provide credential for Jenkins
type JenkinsCredential struct {
	Username string
	Token    string
}

// JenkinsJob contains info about tracking job
type JenkinsJob struct {
	URL string
}

// Config class for represent configuration file
type Config struct {
	Credential  JenkinsCredential
	ComposePath string
	TrackedJobs []JenkinsJob
}

// GetConfigFromFile unmarshal JSON-config from file
func GetConfigFromFile(path string) (config *Config) {
	data, err := ioutil.ReadFile(path)
	util.Check(err)
	err = json.Unmarshal(data, &config)
	util.Check(err)
	return
}

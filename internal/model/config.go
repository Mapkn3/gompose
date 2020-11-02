package model

// JenkinsCredential class for provide credential for Jenkins
type JenkinsCredential struct {
	Username string `json:"username,required"`
	Token    string `json:"token,required"`
}

// JenkinsJob contains info about tracking job
type JenkinsJob struct {
	Name string `json:"name,required"`
	URL  string `json:"url,required"`
}

// Project contains info about certain project
type Project struct {
	Name        string       `json:"name,required"`
	ComposePath string       `json:"composePath,required"`
	TrackedJobs []JenkinsJob `json:"trackedJobs,required"`
}

// Config class for represent configuration file
type Config struct {
	Credential JenkinsCredential `json:"credential,required"`
	Projects   []Project         `json:"projects,required"`
}

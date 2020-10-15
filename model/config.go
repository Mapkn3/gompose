package model

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

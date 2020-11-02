package model

import "strings"

// BuildInfo class for collect build info from jenkins
type BuildInfo struct {
	Number int    `json:"number"`
	URL    string `json:"url"`
}

// ProjectInfo class for collect project info from jenkins
type ProjectInfo struct {
	LastSuccessfulBuild BuildInfo `json:"lastSuccessfulBuild"`
}

// BuildDescription class for persist build comment
type BuildDescription struct {
	Description string `json:"description"`
}

// GetImages returns array of images names
func (d *BuildDescription) GetImages() []string {
	return strings.Fields(d.Description)
}

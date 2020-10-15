package model

import "strings"

// BuildInfo class for collect build info from jenkins
type BuildInfo struct {
	Number int
	URL    string
}

// ProjectInfo class for collect project info from jenkins
type ProjectInfo struct {
	LastSuccessfulBuild BuildInfo
}

// BuildDescription class for persist build comment
type BuildDescription struct {
	Description string
}

// GetImages returns array of images names
func (d *BuildDescription) GetImages() []string {
	return strings.Fields(d.Description)
}

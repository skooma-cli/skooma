// Package types contains shared type definitions for Skooma
package types

// Template represents a project template
type Template struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Author      string `json:"author"`
	RepoURL     string `json:"repo"`
}

// Config represents the Skooma configuration
type Config struct {
	Templates map[string]Template `json:"templates"`
}

// ProjectData holds the data collected from the user to populate the project templates.
type ProjectData struct {
	Name         string
	RootDir      string
	TemplateName string
	Template     Template
	Database     string
	RepoURL      string
	Author       string
}

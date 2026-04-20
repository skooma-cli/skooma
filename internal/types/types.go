// Package types contains shared type definitions for Skooma
package types

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Config represents the Skooma configuration
type Config struct {
	Templates map[string]Template `json:"templates"`
}

// Template represents a project template
type Template struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	RepoURL     Repository `json:"repo_url"`
	Author      string     `json:"author"`
}

// Repository represents a git repository (e.g., "github.com/owner/repo")
type Repository struct {
	Host  string `json:"host"` // e.g., "github.com"
	Owner string `json:"owner"`
	Name  string `json:"name"`
	Ref   string `json:"ref"` // branch, tag, or "latest"
}

// ProjectData holds the data collected from the user to populate the project templates
type ProjectData struct {
	Name      string     `json:"name"`
	Template  Template   `json:"template"`
	RepoURL   Repository `json:"repo_url"`
	Author    string     `json:"author"`
	Directory string     `json:"directory"`
}

// String returns the string representation of the repository URL
func (r Repository) String() string {
	if r.IsEmpty() {
		return ""
	}

	ref := r.Ref
	if ref == "" {
		ref = "latest"
	}

	// Default to https if no host is set but others are (or just construct it)
	if r.Host == "" {
		return fmt.Sprintf("%s/%s@%s", r.Owner, r.Name, ref)
	}
	return fmt.Sprintf("https://%s/%s/%s@%s", r.Host, r.Owner, r.Name, ref)
}

// IsEmpty returns true if the repository is empty
func (r Repository) IsEmpty() bool {
	return r.Host == "" && r.Owner == "" && r.Name == ""
}

// MarshalJSON implements json.Marshaler
func (r Repository) MarshalJSON() ([]byte, error) {
	if r.IsEmpty() {
		return json.Marshal("")
	}
	return json.Marshal(r.String())
}

// UnmarshalJSON implements json.Unmarshaler
func (r *Repository) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	parsed := ParseRepository(s)
	*r = parsed
	return nil
}

// ParseRepository parses a repository string into a Repository struct
func ParseRepository(s string) Repository {
	if s == "" {
		return Repository{}
	}

	// Ensure it has a scheme for url.Parse to work correctly if it's just a domain
	// e.g., github.com/owner/repo
	parseStr := s
	if !strings.HasPrefix(parseStr, "http://") && !strings.HasPrefix(parseStr, "https://") {
		parseStr = "https://" + parseStr
	}

	u, err := url.Parse(parseStr)
	if err != nil {
		return Repository{}
	}

	r := Repository{
		Host: u.Host,
	}

	parts := strings.Split(strings.Trim(u.Path, "/"), "/")
	if len(parts) >= 2 {
		r.Owner = parts[0]
		namePart := strings.TrimSuffix(parts[1], ".git")
		if idx := strings.LastIndex(namePart, "@"); idx != -1 {
			r.Name = namePart[:idx]
			r.Ref = namePart[idx+1:]
		} else {
			r.Name = namePart
			r.Ref = "latest"
		}
	} else if len(parts) == 1 && parts[0] != "" {
		namePart := strings.TrimSuffix(parts[0], ".git")
		if idx := strings.LastIndex(namePart, "@"); idx != -1 {
			r.Name = namePart[:idx]
			r.Ref = namePart[idx+1:]
		} else {
			r.Name = namePart
			r.Ref = "latest"
		}
	} else {
		r.Ref = "latest"
	}

	return r
}

// Download clones the repository to the specified destination
func (r Repository) Download(dest string) error {
	fmt.Println("---------------------------------------------------------")
	fmt.Printf("Downloading template from %s\n", r.String())
	fmt.Printf("Destination: %s\n", dest)
	fmt.Println("---------------------------------------------------------")

	if _, err := os.Stat(dest); !os.IsNotExist(err) {
		fmt.Printf("Directory %s already exists, skipping download.\n", dest)
		return nil
	}

	cloneURL := r.String()
	// Strip the @ref from the clone URL if it exists
	if idx := strings.LastIndex(cloneURL, "@"); idx != -1 {
		cloneURL = cloneURL[:idx]
	}

	args := []string{"clone", "--depth=1"}
	if r.Ref != "" && r.Ref != "latest" {
		args = append(args, "--branch", r.Ref)
	}
	args = append(args, cloneURL, dest)

	cmd := exec.Command("git", args...)
	cmd.Stdout = nil
	cmd.Stderr = nil
	err := cmd.Run()
	if err != nil {
		return err
	}

	return os.RemoveAll(filepath.Join(dest, ".git"))
}

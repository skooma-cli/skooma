// Package templates provides functions for managing project templates.
package templates

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/skooma-cli/skooma/internal/config"
	"github.com/skooma-cli/skooma/internal/logger"
	"github.com/skooma-cli/skooma/internal/types"
)

// GetTemplates returns all templates from the configuration.
func GetTemplates() (map[string]types.Template, error) {
	cfg, err := config.GetConfig()
	if err != nil {
		return nil, err
	}

	return cfg.Templates, nil
}

// GetTemplateByName returns a template by name, or nil if it doesn't exist.
func GetTemplateByName(name string) (*types.Template, error) {
	cfg, err := config.GetConfig()
	if err != nil {
		return nil, err
	}

	tpl, exists := cfg.Templates[name]
	if !exists {
		return nil, nil
	}

	return &tpl, nil
}

// AddTemplate adds a template to the configuration and saves it.
func AddTemplate(template types.Template) error {
	if template.Name == "" {
		return errors.New("template name is required")
	}
	if template.Description == "" {
		return errors.New("template description is required")
	}
	if template.RepoURL.IsEmpty() {
		return errors.New("template repository URL is required")
	}

	cfg, err := config.GetConfig()
	if err != nil {
		return err
	}

	cfg.Templates[template.Name] = template
	return config.SaveConfig(cfg)
}

// SaveTemplate updates a template in the configuration and saves it.
func SaveTemplate(t types.Template) error {
	cfg, err := config.GetConfig()
	if err != nil {
		return err
	}

	cfg.Templates[t.Name] = t
	return config.SaveConfig(cfg)
}

// RemoveTemplate removes a template from the configuration and saves it.
func RemoveTemplate(name string) error {
	cfg, err := config.GetConfig()
	if err != nil {
		return err
	}

	// Get template
	t, err := GetTemplateByName(name)
	if err != nil {
		return err
	}
	if t == nil {
		return errors.New("template not found")
	}

	// Delete template from config
	delete(cfg.Templates, name)
	err = config.SaveConfig(cfg)
	if err != nil {
		return err
	}

	// TODO: check if any other templates use the same repo/ref, if so skip deleting the directory

	// Get template directory
	templateDir, err := GetTemplateDirectory(*t)
	if err != nil {
		return err
	}

	// Delete template directory
	err = os.RemoveAll(templateDir)
	if err != nil {
		return err
	}

	return nil
}

// GetTemplateDirectory returns the path to the templates directory for a given template.
func GetTemplateDirectory(t types.Template) (string, error) {
	repo := t.RepoURL

	if repo.IsEmpty() {
		return "", errors.New("template repository URL is required")
	}

	// Get the directory where templates are stored
	templatesDir, err := config.GetTemplatesDirectory()
	if err != nil {
		return "", err
	}

	// Return the path to the templates directory for a given template
	return filepath.Join(templatesDir, repo.Owner, repo.Name, repo.Ref), nil
}

// Download clones the template's repository to the specified destination
func RepositoryDownload(t *types.Template) error {
	// Get template directory destination
	dest, err := GetTemplateDirectory(*t)
	if err != nil {
		return fmt.Errorf("error getting template directory: %w", err)
	}

	logger.Debug("Downloading template", "template", t.Name, "repo", t.RepoURL.String(), "dest", dest)

	cloneURL := t.RepoURL.String()
	// Strip the @ref from the clone URL if it exists
	if idx := strings.LastIndex(cloneURL, "@"); idx != -1 {
		cloneURL = cloneURL[:idx]
	}

	args := []string{"clone", "--depth=1"}
	if t.RepoURL.Ref != "" && t.RepoURL.Ref != "latest" {
		args = append(args, "--branch", t.RepoURL.Ref)
	}
	args = append(args, cloneURL, dest)

	if _, err := os.Stat(dest); !os.IsNotExist(err) {
		logger.Debug("Directory already exists, skipping download", "path", dest)
	} else {
		cmd := exec.Command("git", args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("error running git clone: %w", err)
		}
	}

	templateConfigPath := filepath.Join(dest, "skooma.config.json")
	if _, err := os.Stat(templateConfigPath); err != nil {
		return fmt.Errorf("template config file not found: %s", templateConfigPath)
	}

	configBytes, err := os.ReadFile(templateConfigPath)
	if err != nil {
		return err
	}

	var templateConfig *types.TemplateConfig
	if err := json.Unmarshal(configBytes, &templateConfig); err != nil {
		return err
	}

	if t.Config == nil {
		t.Config = templateConfig
	}

	// Save the template with the config data
	SaveTemplate(*t)

	return os.RemoveAll(filepath.Join(dest, ".git"))
}

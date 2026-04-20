// Package templates provides functions for managing project templates.
package templates

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/skooma-cli/skooma/internal/config"
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

	tmpl, exists := cfg.Templates[name]
	if !exists {
		return nil, nil
	}

	return &tmpl, nil
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
func GetTemplateDirectory(template types.Template) (string, error) {
	repo := template.RepoURL

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

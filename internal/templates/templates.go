// Package templates provides functions for managing project templates.
package templates

import (
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
func AddTemplate(name string, template types.Template) error {
	cfg, err := config.GetConfig()
	if err != nil {
		return err
	}

	cfg.Templates[name] = template
	return config.SaveConfig(cfg)
}

// RemoveTemplate removes a template from the configuration and saves it.
func RemoveTemplate(name string) error {
	cfg, err := config.GetConfig()
	if err != nil {
		return err
	}

	delete(cfg.Templates, name)
	return config.SaveConfig(cfg)
}

// Package config provides functions for managing the Skooma configuration file.
package config

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/skooma-cli/skooma/internal/logger"
	"github.com/skooma-cli/skooma/internal/types"
)

// Init creates the default config file and templates directory if they don't exist.
func Init() error {
	// Get Skooma directory, create if it doesn't exist
	skoomaDir, err := GetSkoomaDirectory()
	if err != nil {
		return err
	}
	if _, err := os.Stat(skoomaDir); os.IsNotExist(err) {
		err = os.MkdirAll(skoomaDir, 0755)
		if err != nil {
			return err
		}
	}

	// Get Skooma config, create if it doesn't exist
	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}
	if _, err := os.Stat(configPath); err != nil {
		// File doesn't exist, create it with default config
		// TODO: storing the name in the key is redundant, but it makes the TUI easier to build for now. We can refactor later if needed.
		defaultConfig := &types.Config{
			Templates: map[string]types.Template{
				"default": {
					Name:        "default",
					Description: "A default template with Go, React, Tailwind, and Vite",
					RepoURL:     types.ParseRepository("github.com/skooma-cli/skooma-template-default@latest"),
					Author:      "Mark Rodgers <mark@marknrodgers.com>",
				},
			},
		}

		// Write default config to file
		file, err := os.Create(configPath)
		if err != nil {
			return err
		}
		defer file.Close()

		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "\t")
		encoder.SetEscapeHTML(false)

		err = encoder.Encode(defaultConfig)
		if err != nil {
			return err
		}
	}

	// Get templates directory, create if it doesn't exist
	templatesDir, err := GetTemplatesDirectory()
	if err != nil {
		return err
	}
	if _, err := os.Stat(templatesDir); os.IsNotExist(err) {
		err = os.MkdirAll(templatesDir, 0755)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetConfig retrieves the config object from the config file
func GetConfig() (*types.Config, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config types.Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// SaveConfig writes the configuration to disk.
func SaveConfig(config *types.Config) error {
	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}

	file, err := os.Create(configPath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "\t")
	encoder.SetEscapeHTML(false)

	err = encoder.Encode(config)
	if err != nil {
		return err
	}

	return nil
}

// ViewConfig opens the configuration file in the user's default pager.
func ViewConfig() error {
	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}

	pager := os.Getenv("PAGER")
	if pager == "" {
		switch runtime.GOOS {
		case "windows":
			pager = "more"
		default:
			pager = "less"
		}
	}

	file, err := os.Open(configPath)
	if err != nil {
		return err
	}
	defer file.Close()

	cmd := exec.Command(pager)
	cmd.Stdin = file
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	logger.Debug("Opening config file", "cmd", cmd.String())

	return cmd.Run()
}

// GetSkoomaDirectory returns the path to the Skooma directory
func GetSkoomaDirectory() (string, error) {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(userConfigDir, "skooma"), nil
}

// GetConfigPath returns the path to the Skooma config file
func GetConfigPath() (string, error) {
	skoomaDir, err := GetSkoomaDirectory()
	if err != nil {
		return "", err
	}

	return filepath.Join(skoomaDir, "config.json"), nil
}

// GetTemplatesDirectory returns the path to the Skooma templates directory
func GetTemplatesDirectory() (string, error) {
	skoomaDir, err := GetSkoomaDirectory()
	if err != nil {
		return "", err
	}

	return filepath.Join(skoomaDir, "templates"), nil
}

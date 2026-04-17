// Package config provides functions for managing the Skooma configuration file.
package config

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/skooma-cli/skooma/internal/types"
)

// GetConfigPath returns the path to the Skooma config file, creating it with default config if it doesn't exist
func GetConfigPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	// Ensure the skooma config directory exists
	skoomaDir := filepath.Join(configDir, "skooma")
	if _, err := os.Stat(skoomaDir); os.IsNotExist(err) {
		err = os.MkdirAll(skoomaDir, 0755)
		if err != nil {
			return "", err
		}
	}

	// If config file exists, return its filepath
	configPath := filepath.Join(skoomaDir, "config.json")
	if _, err := os.Stat(configPath); err != nil {
		// File doesn't exist, create it with default config
		defaultConfig := &types.Config{
			Templates: map[string]types.Template{
				"default": {
					Description: "A default template with Go, React, Tailwind, and Vite",
					Author:      "Mark Rodgers <mark@marknrodgers.com>",
					RepoURL:     "github.com/skooma-cli/skooma-default-template",
				},
			},
		}

		// Write default config to file
		file, err := os.Create(configPath)
		if err != nil {
			return "", err
		}
		defer file.Close()

		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "\t")
		encoder.SetEscapeHTML(false)

		err = encoder.Encode(defaultConfig)
		if err != nil {
			return "", err
		}
	}

	return configPath, nil
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

// Open opens the configuration file in the user's default editor.
func Open() error {
	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}

	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", "", configPath)
	case "darwin":
		cmd = exec.Command("open", configPath)
	default:
		// Try editor environment variables, fall back to xdg-open
		if editor := getEditor(); editor != "" {
			cmd = exec.Command(editor, configPath)
		} else {
			cmd = exec.Command("xdg-open", configPath)
		}
	}

	return cmd.Start()
}

// getEditor returns the first available editor from environment variables
func getEditor() string {
	for _, env := range []string{"EDITOR", "VISUAL"} {
		if editor := os.Getenv(env); editor != "" {
			return editor
		}
	}
	return ""
}

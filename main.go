// Skooma is a platform-agnostic CLI tool that scaffolds projects in seconds from user-defined template repositories. Any language, any framework, any stack.
package main

import (
	"os"

	"github.com/skooma-cli/skooma/cmd"
	"github.com/skooma-cli/skooma/internal/config"
	"github.com/skooma-cli/skooma/internal/logger"
)

var version = "0.3.0-dev"

func main() {
	os.Setenv("SKOOMA_VERSION", version)

	// Initialize config
	err := config.Init()
	if err != nil {
		panic(err)
	}

	// Initialize logger
	err = logger.Init()
	if err != nil {
		panic(err)
	}

	// Execute commands
	cmd.Execute()
}

// Skooma is a platform-agnostic CLI tool that scaffolds projects in seconds from user-defined template repositories. Any language, any framework, any stack.
package main

import (
	"os"

	"github.com/skooma-cli/skooma/cmd"
	"github.com/skooma-cli/skooma/internal/config"
)

var version = "0.2.0-dev"

func main() {
	os.Setenv("SKOOMA_VERSION", version)

	// Load config to ensure it exists and is valid before executing any commands
	_, err := config.GetConfig()
	if err != nil {
		panic(err)
	}

	cmd.Execute()
}

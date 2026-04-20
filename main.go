// Skooma is a platform-agnostic CLI tool that scaffolds projects in seconds from user-defined template repositories. Any language, any framework, any stack.
package main

import (
	"os"

	"github.com/skooma-cli/skooma/cmd"
	"github.com/skooma-cli/skooma/internal/config"
)

// TODO: implement proper logger

var version = "0.2.0"

func main() {
	os.Setenv("SKOOMA_VERSION", version)

	err := config.Init()
	if err != nil {
		panic(err)
	}

	cmd.Execute()
}

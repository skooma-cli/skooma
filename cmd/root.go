package cmd

import (
	"os"
	"strings"

	"github.com/skooma-cli/skooma/internal/logger"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "skooma",
	Short: "Skooma is a CLI tool to scaffold new projects with a touch of Khajiit flair.",
	Long:  `Skooma is a CLI tool to scaffold new projects with a touch of Khajiit flair.`,
}

// Execute runs the root command, which is the entry point for the CLI application.
func Execute() {
	logger.Info("Command executed", "command", strings.Join(os.Args, " "))

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

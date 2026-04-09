package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "skooma",
	Short: "Skooma is a CLI tool to scaffold new projects with a touch of Khajiit flair.",
	Long:  `Skooma is a CLI tool to scaffold new projects with a touch of Khajiit flair.`,
}

// Execute runs the root command, which is the entry point for the CLI application.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// func init() {}

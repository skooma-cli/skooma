package cmd

import (
	"fmt"

	"github.com/skooma-cli/skooma/internal/logger"
	"github.com/spf13/cobra"
)

// TODO: add flag to tail file

// configCmd represents the config command
var logCmd = &cobra.Command{
	Use:   "log",
	Short: "View Skooma log",
	Long:  `View Skooma log file.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := logger.ViewLog(); err != nil {
			fmt.Printf("Error viewing log: %v\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(logCmd)
}

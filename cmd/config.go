package cmd

import (
	"fmt"

	"github.com/skooma-cli/skooma/internal/config"
	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Open Skooma configuration",
	Long: `Open Skooma configuration file.
This allows you to manage your templates and other settings.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := config.ViewConfig()
		if err != nil {
			fmt.Printf("Error opening config file: %v\n", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}

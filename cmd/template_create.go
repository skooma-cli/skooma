package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var templateCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new template",
	Long:  `Create a new template to be used with the brew command.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("TODO: implement create subcommand")
	},
}

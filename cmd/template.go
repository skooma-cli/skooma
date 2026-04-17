package cmd

import (
	"github.com/spf13/cobra"
)

// templateCmd represents the template command
var templateCmd = &cobra.Command{
	Use:     "template",
	Short:   "Manage Skooma templates",
	Long:    `Manage Skooma templates, which are used to scaffold projects with the brew command. You can list, create, add, and remove templates.`,
	Aliases: []string{"tpl"},
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(templateCmd)
	templateCmd.AddCommand(templateLsCmd)
	templateCmd.AddCommand(templateCreateCmd)
	templateCmd.AddCommand(templateAddCmd)
	templateAddCmd.Flags().StringVarP(&templateAddTemplate.Name, "template", "t", "", "Template name")
	templateAddCmd.Flags().StringVarP(&templateAddTemplate.RepoURL, "repo", "r", "", "Repository URL (e.g., github.com/user/repo)")

	templateCmd.AddCommand(templateRmCmd)
}

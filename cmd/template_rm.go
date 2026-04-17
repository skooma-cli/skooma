package cmd

import (
	"fmt"

	"github.com/skooma-cli/skooma/internal/templates"
	"github.com/spf13/cobra"
)

var templateRmCmd = &cobra.Command{
	Use:     "rm TEMPLATE_NAME",
	Short:   "Remove a template",
	Long:    `Remove an existing template from the available templates.`,
	Args:    cobra.ExactArgs(1),
	Aliases: []string{"remove"},
	Run: func(cmd *cobra.Command, args []string) {
		templateName := args[0]

		tmpls, err := templates.GetTemplates()
		if err != nil {
			fmt.Printf("Error loading templates: %v\n", err)
			return
		}

		if _, exists := tmpls[templateName]; !exists {
			fmt.Printf("Template '%s' does not exist.\n", templateName)
			return
		}

		err = templates.RemoveTemplate(templateName)
		if err != nil {
			fmt.Printf("Error removing template: %v\n", err)
			return
		}

		fmt.Printf("Template '%s' has been removed.\n", templateName)
	},
}

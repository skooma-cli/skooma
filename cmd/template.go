package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"charm.land/huh/v2"
	"github.com/mark-rodgers/skooma/internal/templates"
	"github.com/mark-rodgers/skooma/internal/types"
	"github.com/spf13/cobra"
)

// templateCmd represents the template command
var templateCmd = &cobra.Command{
	Use:   "template",
	Short: "Manage Skooma templates",
	Long:  `Manage Skooma templates, which are used to scaffold projects with the brew command. You can list, create, add, and remove templates.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List available templates",
	Long:  `List all available templates that can be used with the brew command.`,
	Run: func(cmd *cobra.Command, args []string) {
		templates, err := templates.GetTemplates()
		if err != nil {
			fmt.Printf("Error loading templates: %v\n", err)
			return
		}

		if len(templates) == 0 {
			fmt.Println("No templates available. Use 'skooma template add' to add a template.")
			return
		}

		// Create tabwriter with tab-separated columns
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

		// Print header
		fmt.Fprintln(w, "NAME\tDESCRIPTION\tREPO\tAUTHOR")

		// Print templates
		for name, tmpl := range templates {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", name, tmpl.Description, tmpl.Repo, tmpl.Author)
		}

		// Flush to ensure output is written
		w.Flush()
	},
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new template",
	Long:  `Create a new template to be used with the brew command.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("TODO: implement create subcommand")
	},
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add an existing template",
	Long:  `Add an existing template to the available templates.`,
	Run: func(cmd *cobra.Command, args []string) {
		tmplName := ""
		tmpl := types.Template{}

		form := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Template name:").
					Value(&tmplName).
					Validate(func(str string) error {
						// Check if empty
						if strings.TrimSpace(str) == "" {
							return errors.New("Template name can't be empty")
						}
						// Check for spaces
						if strings.Contains(str, " ") {
							return errors.New("Template name can't contain spaces")
						}
						return nil
					}),
				huh.NewInput().
					Title("Description:").
					Value(&tmpl.Description).
					Validate(func(str string) error {
						// TODO: add validation
						return nil
					}),
				huh.NewInput().
					Title("Repo URL:").
					Value(&tmpl.Repo).
					Validate(func(str string) error {
						// TODO: add validation
						return nil
					}),
				huh.NewInput().
					Title("Author:").
					Value(&tmpl.Author).
					Validate(func(str string) error {
						// TODO: add validation
						return nil
					}),
			),
		)

		if err := form.Run(); err != nil {
			fmt.Printf("Error running form: %v\n", err)
			return
		}

		err := templates.AddTemplate(tmplName, tmpl)
		if err != nil {
			fmt.Printf("Error adding template: %v\n", err)
			return
		}

		fmt.Printf("Name: %s\nDescription: %s\nRepo: %s\nAuthor: %s\n", tmplName, tmpl.Description, tmpl.Repo, tmpl.Author)
	},
}

var rmCmd = &cobra.Command{
	Use:   "rm TEMPLATE_NAME",
	Short: "Remove a template",
	Long:  `Remove an existing template from the available templates.`,
	Args:  cobra.ExactArgs(1),

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

func init() {
	rootCmd.AddCommand(templateCmd)
	templateCmd.AddCommand(lsCmd)
	templateCmd.AddCommand(createCmd)
	templateCmd.AddCommand(addCmd)
	templateCmd.AddCommand(rmCmd)
}

package cmd

import (
	"fmt"

	"charm.land/huh/v2"
	"github.com/skooma-cli/skooma/internal/sanitize"
	"github.com/skooma-cli/skooma/internal/types"
	"github.com/skooma-cli/skooma/internal/validators"
	"github.com/spf13/cobra"
)

var templateAddTemplate = types.Template{
	Name:        "",
	Description: "",
	Author:      "",
	RepoURL:     "",
}

var templateAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add an existing template",
	Long:  `Add an existing template to the available templates.`,
	Run: func(cmd *cobra.Command, args []string) {

		form := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Template name:").
					Value(&templateAddTemplate.Name).
					Validate(validators.All(
						validators.NotEmpty("Template name"),
						validators.NoSpaces("Template name"),
						validators.NoUnderscores("Template name"),
					)),
				huh.NewInput().
					Title("Repository URL:").
					Value(&templateAddTemplate.RepoURL).
					Validate(validators.All(
						validators.NotEmpty("Repository URL"),
						validators.ValidURL("Repository URL"),
					)),
			),
		)

		if err := form.Run(); err != nil {
			fmt.Printf("Error running form: %v\n", err)
			return
		}

		templateAddTemplate.RepoURL = sanitize.StripHTTPPrefix(templateAddTemplate.RepoURL)

		fmt.Printf("%+v\n", templateAddTemplate)

		// err := templates.AddTemplate(tmplName, template)
		// if err != nil {
		// 	fmt.Printf("Error adding template: %v\n", err)
		// 	return
		// }

		// fmt.Printf("Name: %s\nDescription: %s\nRepo: %s\nAuthor: %s\n", tmplName, template.Description, template.Repo, template.Author)
	},
}

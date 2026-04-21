package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"charm.land/huh/v2"
	"github.com/skooma-cli/skooma/internal/logger"
	"github.com/skooma-cli/skooma/internal/sanitize"
	"github.com/skooma-cli/skooma/internal/templates"
	"github.com/skooma-cli/skooma/internal/types"
	"github.com/skooma-cli/skooma/internal/validators"
	"github.com/spf13/cobra"
)

var templateAddTemplateNameArg string
var templateAddDescriptionFlag string
var templateAddRepoUrlFlag string
var templateAddAuthorFlag string

var templateAddCmd = &cobra.Command{
	Use:   "add TEMPLATE_NAME",
	Short: "Add an existing template",
	Long:  `Add an existing template to the available templates.`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			templateAddTemplateNameArg = args[0]
		}

		groups := []*huh.Group{}

		// Validators for the template name input
		templateNameValidators := []types.ValidatorFunc{
			validators.NotEmpty("Template name"), // only meaningful in the TUI, redundant if a flag is provided 🤷
			validators.NoSpaces("Template name"),
			validators.NoUnderscores("Template name"),
		}
		// If no template name was provided, prompt the user; otherwise validate the provided value
		if templateAddTemplateNameArg == "" {
			groups = append(groups, huh.NewGroup(
				huh.NewInput().
					Title("Template name:").
					Value(&templateAddTemplateNameArg).
					Validate(validators.All(templateNameValidators...)),
			))
		} else {
			if err := validators.All(templateNameValidators...)(templateAddTemplateNameArg); err != nil {
				logger.Fatal("Invalid template name", "error", err)
			}
		}

		// Validators for the description input
		descriptionValidators := []types.ValidatorFunc{
			validators.NotEmpty("Description"),
		}
		// If no description was provided, prompt the user; otherwise validate the provided value
		if templateAddDescriptionFlag == "" {
			groups = append(groups, huh.NewGroup(
				huh.NewInput().
					Title("Description:").
					Value(&templateAddDescriptionFlag).
					Validate(validators.All(descriptionValidators...)),
			))
		} else {
			if err := validators.All(descriptionValidators...)(templateAddDescriptionFlag); err != nil {
				logger.Fatal("Invalid description", "error", err)
			}
		}

		// Validators for the repository URL input
		repoUrlValidators := []types.ValidatorFunc{
			validators.NotEmpty("Repository URL"),
			validators.NoSpaces("Repository URL"),
			validators.ValidURL("Repository URL"),
		}
		// If no repository URL was provided, prompt the user; otherwise validate the provided value
		if templateAddRepoUrlFlag == "" {
			groups = append(groups, huh.NewGroup(
				huh.NewInput().
					Title("Repository URL (e.g., github.com/user/repo):").
					Value(&templateAddRepoUrlFlag).
					Validate(validators.All(repoUrlValidators...)),
			))
		} else {
			if err := validators.All(repoUrlValidators...)(templateAddRepoUrlFlag); err != nil {
				logger.Fatal("Invalid repository URL", "error", err)
			}
		}

		// Validators for the author name input
		authorValidators := []types.ValidatorFunc{
			validators.RFC5322Address("Author"),
		}
		// If no author was provided, prompt the user; otherwise validate the provided value
		if templateAddAuthorFlag == "" {
			groups = append(groups, huh.NewGroup(
				huh.NewInput().
					Title("Author name (e.g., Name <email@example.com>):").
					Value(&templateAddAuthorFlag).
					Validate(validators.AllowEmpty(authorValidators...)),
			))
		} else {
			if err := validators.All(authorValidators...)(templateAddAuthorFlag); err != nil {
				logger.Fatal("Invalid author name", "error", err)
			}
		}

		form := huh.NewForm(groups...)

		// Run the form to collect user input
		err := form.Run()
		if err != nil {
			logger.Fatal("Failed to run form", "error", err)
		}

		// Build project data struct to pass to the brewing process
		template := types.Template{
			Name:        templateAddTemplateNameArg,
			Description: sanitize.TrimWhitespace(templateAddDescriptionFlag),
			RepoURL:     types.ParseRepository(sanitize.StripHTTPPrefix(templateAddRepoUrlFlag)),
			Author:      templateAddAuthorFlag,
		}

		// Check if template already exists before adding
		t, err := templates.GetTemplateByName(template.Name)
		if err == nil && t != nil {
			logger.Fatal(fmt.Sprintf("Template '%s' already exists", template.Name))
		}

		// Download template from git repository
		err = templates.RepositoryDownload(&template)
		if err != nil {
			logger.Fatal("Error downloading template", "error", err)
		}

		// Add the template
		err = templates.AddTemplate(template)
		if err != nil {
			logger.Fatal("Error adding template", "error", err)
		}

		fmt.Printf("\n✅ '%s' has been added successfully!\n\n", template.Name)

		// Print template details
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintf(w, "Name\t%s\n", template.Name)
		fmt.Fprintf(w, "Description\t%s\n", template.Description)
		fmt.Fprintf(w, "Repository\t%s\n", template.RepoURL)
		fmt.Fprintf(w, "Author\t%s\n", template.Author)
		w.Flush()
	},
}

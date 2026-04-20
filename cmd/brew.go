package cmd

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"charm.land/huh/v2"
	// "github.com/briandowns/spinner"
	"github.com/skooma-cli/skooma/internal/brew"
	"github.com/skooma-cli/skooma/internal/sanitize"
	"github.com/skooma-cli/skooma/internal/templates"
	"github.com/skooma-cli/skooma/internal/types"
	"github.com/skooma-cli/skooma/internal/utils"
	"github.com/skooma-cli/skooma/internal/validators"
	"github.com/spf13/cobra"
)

var brewProjectNameArg string
var brewTemplateFlag string
var brewRepoUrlFlag string
var brewAuthorFlag string

var brewCmd = &cobra.Command{
	Use:   "brew PROJECT_NAME",
	Short: "Brew a new project",
	Long: `Brew a new project with the given name.
This command will create a new directory with the project name and generate
the necessary files for a basic project structure.`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s\n\n", utils.GetRandomKhajiitPhrase())

		if len(args) > 0 {
			brewProjectNameArg = args[0]
		}

		groups := []*huh.Group{}

		// Validators for the project name input
		projectNameValidators := []func(string) error{
			validators.NotEmpty("Project name"), // only meaningful in the TUI, redundant if a flag is provided 🤷
			validators.NoSpaces("Project name"),
			validators.NoUnderscores("Project name"),
		}
		// If no project name was provided, prompt the user; otherwise validate the provided value
		if brewProjectNameArg == "" {
			groups = append(groups, huh.NewGroup(
				huh.NewInput().
					Title("Project name:").
					Value(&brewProjectNameArg).
					Validate(validators.All(projectNameValidators...)),
			))
		} else {
			if err := validators.All(projectNameValidators...)(brewProjectNameArg); err != nil {
				log.Fatalf("❌ Invalid project name: %v\n", err)
			}
		}

		// Load templates to build options for the template selection prompt
		templates, err := templates.GetTemplates()
		if err != nil {
			log.Fatalf("❌ Error loading templates: %v\n", err)
		}

		// If no template was provided, prompt the user; otherwise validate the provided template name exists
		if brewTemplateFlag == "" {
			templateOptions := make([]huh.Option[string], 0, len(templates))
			for name, tmpl := range templates {
				templateOptions = append(templateOptions, huh.NewOption(name+" - "+tmpl.Description, name))
			}
			groups = append(groups, huh.NewGroup(
				huh.NewSelect[string]().
					Title("Template").
					Options(templateOptions...).
					Value(&brewTemplateFlag),
			))
		} else if _, ok := templates[brewTemplateFlag]; !ok {
			log.Fatalf("❌ Invalid template name: '%s'. Use 'skooma template ls' to see available templates.\n", brewTemplateFlag)
		}

		// Validators for the repository URL input
		repoUrlValidators := []func(string) error{
			validators.NoSpaces("Repository URL"),
			validators.ValidURL("Repository URL"),
		}
		// If no repository URL was provided, prompt the user; otherwise validate the provided value
		if brewRepoUrlFlag == "" {
			groups = append(groups, huh.NewGroup(
				huh.NewInput().
					Title("Repository URL (e.g., github.com/user/repo):").
					Value(&brewRepoUrlFlag).
					Validate(validators.AllowEmpty(repoUrlValidators...)),
			))
		} else {
			if err := validators.All(repoUrlValidators...)(brewRepoUrlFlag); err != nil {
				log.Fatalf("❌ Invalid repository URL: %v\n", err)
			}
		}

		// Validators for the author name input
		authorValidators := []func(string) error{
			validators.RFC5322Address("Author"),
		}
		// If no author was provided, prompt the user; otherwise validate the provided value
		if brewAuthorFlag == "" {
			groups = append(groups, huh.NewGroup(
				huh.NewInput().
					Title("Author name (e.g., Name <email@example.com>):").
					Value(&brewAuthorFlag).
					Validate(validators.AllowEmpty(authorValidators...)),
			))
		} else {
			if err := validators.All(authorValidators...)(brewAuthorFlag); err != nil {
				log.Fatalf("❌ Invalid author name: %v\n", err)
			}
		}

		form := huh.NewForm(groups...)

		// Run the form to collect user input
		err = form.Run()
		if err != nil {
			log.Fatalf("❌ Failed to run form: %v\n", err)
		}

		// Build project data struct to pass to the brewing process
		project := types.ProjectData{
			Name:     brewProjectNameArg,
			Template: templates[brewTemplateFlag],
			RepoURL:  types.ParseRepository(sanitize.StripHTTPPrefix(brewRepoUrlFlag)),
			Author:   brewAuthorFlag,
		}

		// // Start brewing spinner
		// s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
		// s.Suffix = " Brewing..."
		// s.Start()

		// Brew project
		err = brew.ScaffoldProject(&project)
		if err != nil {
			log.Fatalf("❌ Failed to brew project\n\n%v\n", err)
		}

		// Stop spinner and print success message
		// s.Stop()
		fmt.Printf("\n✅ '%s' has finished brewing!\n\n", project.Name)

		// Print project details
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintf(w, "Template\t%s - %s\n", project.Template.Name, project.Template.Description)
		fmt.Fprintf(w, "Repository\t%s\n", project.RepoURL)
		fmt.Fprintf(w, "Author\t%s\n", project.Author)
		fmt.Fprintf(w, "Directory\t%s\n", project.Directory)
		w.Flush()
	},
}

// init registers the brew command and its flags with the root command.
func init() {
	rootCmd.AddCommand(brewCmd)
	brewCmd.Flags().StringVarP(&brewTemplateFlag, "template", "t", "", "Template name")
	brewCmd.Flags().StringVarP(&brewRepoUrlFlag, "repo", "r", "", "Repository URL (e.g., github.com/user/repo)")
	brewCmd.Flags().StringVarP(&brewAuthorFlag, "author", "a", "", "Author name")
}

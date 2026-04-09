// brew.go defines the "brew" command for the Skooma CLI, which scaffolds a new
// fullstack project with Go, TypeScript, React, Tailwind, and Vite.
// It collects user input for project details, creates the necessary directory
// structure, and generates files based on embedded templates. The command also
// includes a fun brewing message and a spinner to enhance the user experience
// while the project is being set up.
package cmd

import (
	"embed"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
	"time"

	"charm.land/huh/v2"
	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
)

// ProjectData holds the data collected from the user to populate the project templates.
type ProjectData struct {
	Name     string
	RootDir  string
	RepoURL  string
	Author   string
	Database string
}

var project = ProjectData{
	Name:     "",
	RootDir:  "",
	RepoURL:  "",
	Author:   "",
	Database: "file",
}

//go:embed templates/*
var templateFS embed.FS

// getRandomBrewMessage returns a random message to display while brewing the project.
func getRandomBrewMessage() string {
	messages := []string{
		"🧪 This one is brewing a fresh batch of Skooma...",
		"🦁 Khajiit has wares, if you have coin...",
		"🌙 By Azura! This one crafts magical elixir...",
		"🏝️ May your roads lead you to warm sands...",
		"🧙 This one mixes moon sugar and nightshade...",
		"🏺 Psst! Khajiit knows you come for the good stuff...",
	}
	return messages[rand.Intn(len(messages))]
}

var brewCmd = &cobra.Command{
	Use:   "brew <project_name>",
	Short: "Brew a new project",
	Long:  `Brew a new project with the given name. This command will create a new directory with the project name and generate the necessary files for a basic project structure.`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s\n\n", getRandomBrewMessage())

		if len(args) > 0 {
			project.Name = args[0]
		}

		form := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Project name:").
					Value(&project.Name).
					Validate(func(str string) error {
						// Check if empty
						if strings.TrimSpace(str) == "" {
							return errors.New("Project name can't be empty")
						}
						// Check for spaces
						if strings.Contains(str, " ") {
							return errors.New("Project name can't contain spaces")
						}
						// Check for underscores
						if strings.Contains(str, "_") {
							return errors.New("Project name can't contain underscores")
						}
						return nil
					}),
				huh.NewInput().
					Title("Repository URL (e.g., github.com/username/repo):").
					Value(&project.RepoURL).
					Validate(func(str string) error {
						// Check for spaces
						if strings.Contains(str, " ") {
							return errors.New("Repository URL can't contain spaces")
						}
						// Check for http:// or https://
						if strings.HasPrefix(str, "http://") || strings.HasPrefix(str, "https://") {
							return errors.New("Repository URL can't contain http:// or https://")
						}
						return nil
					}),
				huh.NewInput().
					Title("Author name (e.g., John Doe <john.doe@example.com>):").
					Value(&project.Author).
					Validate(func(str string) error {
						// If author name is provided, enforce the format "Name <email>" via regex
						if strings.TrimSpace(str) != "" {
							pattern := `^[^<>]+ <[^@\s]+@[^@\s]+\.[^@\s]+>$`
							matched, err := regexp.MatchString(pattern, str)
							if err != nil || !matched {
								return errors.New("Author must be in format: Name <email@domain.com>")
							}
						}
						return nil
					}),
				huh.NewSelect[string]().
					Title("Database").
					Options(
						huh.NewOption("Flat File", "file"),
						huh.NewOption("Microsoft SQL", "mssql"),
						huh.NewOption("PostgreSQL", "postgres"),
					).
					Value(&project.Database),
			),
		)

		err := form.Run()
		if err != nil {
			log.Fatal(err)
		}

		// Get current working directory
		cwd, err := os.Getwd()
		if err != nil {
			log.Fatalf("❌ Failed to get current working directory: %v\n", err)
		}
		project.RootDir = filepath.Join(cwd, project.Name)

		// Early return if project directory already exists
		if _, err := os.Stat(project.RootDir); !os.IsNotExist(err) {
			log.Fatalf("❌ Directory '%s' already exists\n", project.RootDir)
		}

		// Start brewing spinner
		s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
		s.Suffix = " Brewing..."
		s.Start()

		// Prepare project data for templating
		err = scaffoldProject()
		if err != nil {
			log.Fatalf("❌ Failed to brew project\n\n%v\n", err)
		}

		// Simulate scaffolding work
		time.Sleep(2 * time.Second)

		s.Stop()
		fmt.Printf("\n✅ '%s' has finished brewing!\n\n", project.Name)

		// Print project details
		fmt.Printf("Directory: %s\n", project.RootDir)
		fmt.Printf("Repository: https://%s\n", project.RepoURL)
		if project.Author != "" {
			fmt.Printf("Author: %s\n", project.Author)
		}
	},
}

// init registers the brew command and its flags with the root command.
func init() {
	rootCmd.AddCommand(brewCmd)
	brewCmd.Flags().StringVarP(&project.RepoURL, "repo", "r", "", "Repository URL (e.g., github.com/username/repo)")
	brewCmd.Flags().StringVarP(&project.Author, "author", "a", "", "Author name")
	brewCmd.Flags().StringVarP(&project.Database, "database", "d", "", "Database type (\"file\", \"mssql\", \"postgres\")")
}

// scaffoldProject creates the project directory structure and generates files based on templates.
func scaffoldProject() error {
	err := createProjectRoot()
	if err != nil {
		return fmt.Errorf("failed to brew project: %w", err)
	}

	err = createBackend()
	if err != nil {
		return fmt.Errorf("failed to brew project: %w", err)
	}

	err = createFrontend()
	if err != nil {
		return fmt.Errorf("failed to brew project: %w", err)
	}
	return nil
}

// createProjectRoot creates the root project directory and processes root-level templates.
func createProjectRoot() error {
	projectRoot := project.RootDir

	// Create project root directory
	err := os.Mkdir(projectRoot, 0755)
	if err != nil {
		return fmt.Errorf("failed to create project root directory: %w", err)
	}

	// Process root-level templates
	err = processTemplate("templates/docker-compose.yml.tmpl", filepath.Join(projectRoot, "docker-compose.yml"))
	if err != nil {
		return fmt.Errorf("failed to process root-level templates: %w", err)
	}
	return nil
}

// createBackend creates the backend directory and generates files based on templates.
func createBackend() error {
	backendPath := filepath.Join(project.RootDir, "backend")
	err := os.Mkdir(backendPath, 0755)
	if err != nil {
		return fmt.Errorf("failed to create backend directory: %w", err)
	}

	// Process backend templates
	templates := []struct {
		src, dst string
	}{
		{"templates/backend/go.mod.tmpl", filepath.Join(backendPath, "go.mod")},
		{"templates/backend/main.go.tmpl", filepath.Join(backendPath, "main.go")},
		{"templates/backend/Makefile.tmpl", filepath.Join(backendPath, "Makefile")},
	}

	for _, tmpl := range templates {
		if err := processTemplate(tmpl.src, tmpl.dst); err != nil {
			return fmt.Errorf("failed to process template %s: %w", tmpl.src, err)
		}
	}
	return nil
}

// createFrontend creates the frontend directory, subdirectories, and generates files based on templates.
func createFrontend() error {
	frontendPath := filepath.Join(project.RootDir, "frontend")
	err := os.Mkdir(frontendPath, 0755)
	if err != nil {
		return fmt.Errorf("failed to create frontend directory: %w", err)
	}

	subdirs := []string{"src", "src/assets", "public"}
	for _, subdir := range subdirs {
		err := os.Mkdir(filepath.Join(frontendPath, subdir), 0755)
		if err != nil {
			return fmt.Errorf("failed to create frontend subdirectory %s: %w", subdir, err)
		}
	}

	// Copy static asset files that don't require templating
	staticFiles := []struct {
		src, dst string
	}{
		// Public directory static files
		{"templates/frontend/public/favicon.svg", filepath.Join(frontendPath, "public", "favicon.svg")},
		{"templates/frontend/public/khajiit.webp", filepath.Join(frontendPath, "public", "khajiit.webp")},
	}
	for _, file := range staticFiles {
		if err := copyFile(file.src, file.dst); err != nil {
			return fmt.Errorf("failed to copy static file %s: %w", file.src, err)
		}
	}

	// Process frontend templates
	templates := []struct {
		src, dst string
	}{
		{"templates/frontend/gitignore.tmpl", filepath.Join(frontendPath, ".gitignore")},
		{"templates/frontend/eslint.config.js.tmpl", filepath.Join(frontendPath, "eslint.config.js")},
		{"templates/frontend/index.html.tmpl", filepath.Join(frontendPath, "index.html")},
		{"templates/frontend/package.json.tmpl", filepath.Join(frontendPath, "package.json")},
		{"templates/frontend/README.md.tmpl", filepath.Join(frontendPath, "README.md")},
		{"templates/frontend/tsconfig.json.tmpl", filepath.Join(frontendPath, "tsconfig.json")},
		{"templates/frontend/tsconfig.app.json.tmpl", filepath.Join(frontendPath, "tsconfig.app.json")},
		{"templates/frontend/tsconfig.node.json.tmpl", filepath.Join(frontendPath, "tsconfig.node.json")},
		{"templates/frontend/vite.config.ts.tmpl", filepath.Join(frontendPath, "vite.config.ts")},
		{"templates/frontend/src/App.css.tmpl", filepath.Join(frontendPath, "src", "App.css")},
		{"templates/frontend/src/App.tsx.tmpl", filepath.Join(frontendPath, "src", "App.tsx")},
		{"templates/frontend/src/index.css.tmpl", filepath.Join(frontendPath, "src", "index.css")},
		{"templates/frontend/src/main.tsx.tmpl", filepath.Join(frontendPath, "src", "main.tsx")},
	}

	for _, tmpl := range templates {
		if err := processTemplate(tmpl.src, tmpl.dst); err != nil {
			return fmt.Errorf("failed to process template %s: %w", tmpl.src, err)
		}
	}
	return nil
}

// copyFile reads a file from the embedded filesystem and writes it to the specified destination path.
func copyFile(src, dst string) error {
	// Read file content from embedded filesystem
	content, err := templateFS.ReadFile(src)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", src, err)
	}

	// Write content to destination path
	err = os.WriteFile(dst, content, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file %s: %w", dst, err)
	}
	return nil
}

// processTemplate reads a template from the embedded filesystem, executes it with the project data, and writes the output to the specified path.
func processTemplate(templatePath, outputPath string) error {
	// Read template from embedded filesystem
	content, err := templateFS.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("failed to read template %s: %w", templatePath, err)
	}

	// Parse and execute template
	tmpl, err := template.New(filepath.Base(templatePath)).Parse(string(content))
	if err != nil {
		return fmt.Errorf("failed to parse template %s: %w", templatePath, err)
	}

	// Create output file
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file %s: %w", outputPath, err)
	}
	defer outputFile.Close()

	// Execute template with data
	if err := tmpl.Execute(outputFile, project); err != nil {
		return fmt.Errorf("failed to execute template %s: %w", templatePath, err)
	}
	return nil
}

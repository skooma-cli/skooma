// Package brew handles the brewing (creation) of new projects.
package brew

import (
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"charm.land/huh/v2"
	"github.com/skooma-cli/skooma/internal/logger"
	"github.com/skooma-cli/skooma/internal/templates"
	"github.com/skooma-cli/skooma/internal/types"
	"github.com/skooma-cli/skooma/internal/validators"
)

// ScaffoldProject scaffolds a new project with the given project data.
func ScaffoldProject(project *types.ProjectData) error {
	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working directory: %w", err)
	}

	// Check if project root directory already exists before starting the brewing process
	project.Directory = filepath.Join(cwd, project.Name)
	if err = createProjectRoot(project.Directory); err != nil {
		return fmt.Errorf("error creating project directory: %w", err)
	}

	// Get template directory
	src, err := templates.GetTemplateDirectory(project.Template)
	if err != nil {
		return fmt.Errorf("failed to get template directory: %w", err)
	}

	// Download template repository
	if err = templates.RepositoryDownload(&project.Template); err != nil {
		return fmt.Errorf("error downloading template: %v\n", err)
	}

	// Copy all static files from the template directory to the project directory
	if err = copyStaticFiles(src, project.Directory); err != nil {
		return fmt.Errorf("error copying static files: %w", err)
	}

	// Process template files, replacing variables with project data
	if err = processTemplateFiles(src, project.Directory, *project); err != nil {
		return fmt.Errorf("error processing template files: %w", err)
	}

	// Scaffolding is too fast, add a delay for the 'brew' effect
	time.Sleep(1 * time.Second)

	return nil
}

func BuildTemplateVariableInputGroups(variables *[]types.TemplateConfigVariable) ([]*huh.Group, error) {
	groups := []*huh.Group{}

	for i := range *variables {
		variable := &(*variables)[i]

		// Copy the default value to the value field
		variable.Value = variable.Default

		// Resolve validators
		validatorFns, err := validators.ResolveValidators(*variable)
		if err != nil {
			return nil, fmt.Errorf("error resolving validators: %w", err)
		}

		// Add appropriate input field based on variable type
		switch variable.Type {
		case "text":
			// Add text input to groups
			groups = append(groups, huh.NewGroup(
				huh.NewInput().
					Title(variable.Prompt).
					Description(variable.Description).
					Value(&variable.Value).
					Validate(validatorFns),
			))
		case "select":
			// Build options for the select input
			options := make([]huh.Option[string], 0, len(variable.Options))
			for _, v := range variable.Options {
				options = append(options, huh.NewOption(v.Label, v.Value))
			}
			// Add select input to groups
			groups = append(groups, huh.NewGroup(
				huh.NewSelect[string]().
					Title(variable.Prompt).
					Description(variable.Description).
					Options(options...).
					Value(&variable.Value).
					Validate(validatorFns),
			))
		default:
			return nil, fmt.Errorf("unsupported variable type: %s", variable.Type)
		}
	}

	return groups, nil
}

func createProjectRoot(dst string) error {
	if _, err := os.Stat(dst); os.IsNotExist(err) {
		logger.Debug("Creating directory", "path", dst)
		err = os.MkdirAll(dst, 0755)
		if err != nil {
			return err
		}
	} else {
		logger.Fatal("Project directory already exists", "path", dst)
	}
	return nil
}

func copyStaticFiles(src, dst string) error {
	// logger.Debug("Copying static files", "src", src, "dst", dst)

	filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Get paths
		relPath := strings.TrimPrefix(path, src)
		dstPath := filepath.Join(dst, relPath)

		// Skip the root directory itself
		if relPath == "" {
			return nil
		}

		// Create directory in destination directory
		if d.IsDir() {
			logger.Debug("Creating directory", "path", dstPath)
			if err := os.MkdirAll(dstPath, 0755); err != nil {
				return err
			}
			return nil
		}

		// Skip skooma.config.json file
		if d.Name() == "skooma.config.json" {
			return nil
		}

		// Skip files with .tmpl extension; those will be processed later
		ext := filepath.Ext(d.Name())
		if ext == ".tmpl" {
			return nil
		}

		logger.Debug("Copying file", "path", dstPath)

		// Read source file
		in, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("opening source file: %w", err)
		}
		defer in.Close()

		// Create destination file
		out, err := os.Create(dstPath)
		if err != nil {
			return fmt.Errorf("creating destination file: %w", err)
		}
		defer out.Close()

		// Copy contents of source file to destination file
		if _, err = io.Copy(out, in); err != nil {
			return fmt.Errorf("copying file: %w", err)
		}

		return nil
	})

	return nil
}

func processTemplateFiles(src, dst string, project types.ProjectData) error {
	// Build base variables map with your hardcoded fields
	variables := map[string]any{
		"Name":         project.Name,
		"RepoURL":      project.RepoURL.String(),
		"Author":       project.Author,
		"Directory":    project.Directory,
		"GoModulePath": fmt.Sprintf("%s/%s/%s", project.RepoURL.Host, project.RepoURL.Owner, project.RepoURL.Name),
	}

	// Append dynamic template config variables
	if project.Template.Config != nil {
		for _, v := range project.Template.Config.Variables {
			variables[v.Name] = v.Value
		}
	}

	filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Get paths
		relPath := strings.TrimPrefix(path, src)
		dstPath := filepath.Join(dst, relPath)

		// Skip directories
		if d.IsDir() {
			return nil
		}

		// Skip non-template files
		ext := filepath.Ext(d.Name())
		if ext != ".tmpl" {
			return nil
		}

		// Parse template file
		tmpl, err := template.ParseFiles(path)
		if err != nil {
			return fmt.Errorf("parsing template: %w", err)
		}

		// Create destination file, stripping the .tmpl extension from the filename
		dstPath = strings.TrimSuffix(dstPath, ".tmpl")
		f, err := os.Create(dstPath)
		if err != nil {
			return fmt.Errorf("creating destination file: %w", err)
		}
		defer f.Close()

		// Process template file, replacing variables with project data
		logger.Debug("Processing template file", "path", dstPath)
		if err := tmpl.Execute(f, variables); err != nil {
			return fmt.Errorf("executing template: %w", err)
		}

		return nil
	})
	return nil
}

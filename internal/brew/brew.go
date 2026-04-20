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

	"github.com/skooma-cli/skooma/internal/templates"
	"github.com/skooma-cli/skooma/internal/types"
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
	if err = project.Template.RepoURL.Download(src); err != nil {
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

	return nil
}

func createProjectRoot(dst string) error {
	// TODO: this should error out if directory already exists
	if _, err := os.Stat(dst); os.IsNotExist(err) {
		err = os.MkdirAll(dst, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

func copyStaticFiles(src, dst string) error {
	fmt.Println(src)
	fmt.Println(dst)
	fmt.Println("---------------------------------------")

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
			fmt.Printf("Creating directory: %s\n", dstPath)
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

		fmt.Printf("Copying file: %s\n", dstPath)

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
		fmt.Printf("Processing template file: %s\n", dstPath)
		if err := tmpl.Execute(f, project); err != nil {
			return fmt.Errorf("executing template: %w", err)
		}

		return nil
	})
	return nil
}

package services

import (
	"fmt"
	"io"
	"os"
	"strings"
)

type FilesystemReaderWriter interface {
	HasDirectory(directory string) (bool, error)
	CreateDirectory(directory string) error
	CreateFile(filename string) (*os.File, error)
}

type TemplatesWriter interface {
	ExecuteTemplate(wr io.Writer, name string, data any) error
}

type InitialiserServiceConfig struct {
	ProjectName    string
	GithubUsername string
}

// Service for hanlding initialising new projects
type InitialiserService struct {
	filesystem FilesystemReaderWriter
	templates  TemplatesWriter
}

func NewInitialiserService(filesystem FilesystemReaderWriter, templates TemplatesWriter) InitialiserService {
	return InitialiserService{filesystem, templates}
}

func (s InitialiserService) InitProject(username string, projectName string) error {
	if hasDir, err := s.filesystem.HasDirectory(projectName); err != nil || hasDir {
		return fmt.Errorf("%v already exists in this directory", projectName)
	}

	if err := s.filesystem.CreateDirectory(projectName); err != nil {
		return fmt.Errorf("failed to create '%v' directory", projectName)
	}

	// Creating required directories
	directories := []string{
		"controllers",
		"frontend",
		"frontend/src",
		"frontend/src/models",
		"frontend/src/views",
		"frontend/src/views/pages",
	}

	for _, dir := range directories {
		dirWithProject := fmt.Sprintf("%v/%v", projectName, dir)

		if err := s.filesystem.CreateDirectory(dirWithProject); err != nil {
			return fmt.Errorf("failed to create '%v' directory", dirWithProject)
		}
	}

	// Creating files
	files := []string{
		".gitignore",
		// Backend stuff
		"go.mod",
		"main.go",
		"controllers/hello.go",
		// Frontend stuff
		"frontend/favicon.ico",
		"frontend/global.css",
		"frontend/index.html",
		"frontend/package.json",
		"frontend/tailwind.config.js",
		"frontend/tsconfig.json",
		"frontend/src/index.ts",
		"frontend/src/models/hello.ts",
		"frontend/src/views/Button.ts",
		"frontend/src/views/pages/HomePage.ts",
	}

	config := InitialiserServiceConfig{
		ProjectName:    projectName,
		GithubUsername: username,
	}

	for _, filepath := range files {
		filepathWithProject := fmt.Sprintf("%v/%v", projectName, filepath)
		file, err := s.filesystem.CreateFile(filepathWithProject)
		if err != nil {
			return fmt.Errorf("failed to create '%v' file", filepath)
		}
		defer file.Close()

		parts := strings.Split(filepath, "/")
		filename := parts[len(parts)-1] // Last part will be the file
		template := fmt.Sprintf("%v.tmpl", filename)

		s.templates.ExecuteTemplate(file, template, config)
	}

	return nil
}

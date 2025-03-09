package services

import (
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
)

// Mock implementation of the filesystem, does not do any error checking
type mockFilesystem struct {
	hasDirectoryOrFileReturn bool
}

func (m mockFilesystem) ReadDirRecursive(directory string) ([]string, error) {
	return make([]string, 0), nil
}
func (m mockFilesystem) ReadFile(filename string) (string, error) {
	return "Mock contents", nil
}

func (m mockFilesystem) HasDirectoryOrFile(directory string) (bool, error) {
	return m.hasDirectoryOrFileReturn, nil
}

func (m mockFilesystem) CreateDirectory(directory string) error {
	return os.Mkdir(directory, 0755)
}

func (m mockFilesystem) CreateFile(filename string) (*os.File, error) {
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	return file, nil
}

func (m mockFilesystem) DeleteFileRecursive(filename string) error {
	return os.RemoveAll(filename)
}

func (m mockFilesystem) Root() string {
	return "."
}

func (m mockFilesystem) FromRoot(path string) string {
	return fmt.Sprintf("./%v", strings.TrimPrefix(path, "/"))
}

// Mock implementation of templates repository, just writes "mock-data" to the given file
type mockTemplates struct{}

func (m mockTemplates) WriteTemplate(wr io.Writer, name string, data any) error {
	_, err := wr.Write([]byte("mock-data"))
	return err
}

func TestInitialiseProjectCreatesExpectedFiles(t *testing.T) {
	// Arrange
	filesystem := mockFilesystem{hasDirectoryOrFileReturn: false}
	templates := mockTemplates{}
	initService := NewInitialiserService(filesystem, templates)

	username := "mock-user"
	projectName := "testproject"
	projectDir := "testproject"

	// Act
	initService.InitProject(username, projectName, projectDir)

	// Assert
	expectedFiles := []string{
		".gitignore",
		"go.mod",
		"main.go",
		"controllers/hello.go",
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

	for _, file := range expectedFiles {
		_, err := os.Stat(fmt.Sprintf("%v/%v", projectDir, file))
		if err != nil {
			if removalErr := os.RemoveAll(projectDir); removalErr != nil {
				panic(removalErr)
			}
			panic(err)
		}
	}

	if err := os.RemoveAll(projectDir); err != nil {
		panic(err)
	}
}

func TestInitialiseProjectCreatesExpectedFilesIfProjectNameDiffersFromDirName(t *testing.T) {
	// Arrange
	filesystem := mockFilesystem{hasDirectoryOrFileReturn: false}
	templates := mockTemplates{}
	initService := NewInitialiserService(filesystem, templates)

	username := "mock-user"
	projectName := "testproject"
	projectDir := "differentdirectory"

	// Act
	initService.InitProject(username, projectName, projectDir)

	// Assert
	expectedFiles := []string{
		".gitignore",
		"go.mod",
		"main.go",
		"controllers/hello.go",
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

	for _, file := range expectedFiles {
		_, err := os.Stat(fmt.Sprintf("%v/%v", projectDir, file))
		if err != nil {
			if removalErr := os.RemoveAll(projectDir); removalErr != nil {
				panic(removalErr)
			}
			panic(err)
		}
	}

	if err := os.RemoveAll(projectDir); err != nil {
		panic(err)
	}
}

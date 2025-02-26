package services

import (
	"io"
	"os"
	"testing"
)

// Mock implementation of the filesystem, does not do any error checking
type mockFilesystem struct {
	hasDirectoryOrFileReturn bool
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

// Mock implementation of templates repository, just writes "mock-data" to the given file
type mockTemplates struct{}

func (m mockTemplates) ExecuteTemplate(wr io.Writer, name string, data any) error {
	_, err := wr.Write([]byte("mock-data"))
	return err
}

func TestInitialiseProjectCreatesExpectedFiles(t *testing.T) {
	// Arrange
	filesystem := mockFilesystem{hasDirectoryOrFileReturn: false}
	templates := mockTemplates{}
	initService := NewInitialiserService(filesystem, templates)

	// Act
	initService.InitProject("mock-user", "testproject")

	// Assert
	expectedFiles := []string{
		"testproject/.gitignore",
		"testproject/go.mod",
		"testproject/main.go",
		"testproject/controllers/hello.go",
		"testproject/frontend/favicon.ico",
		"testproject/frontend/global.css",
		"testproject/frontend/index.html",
		"testproject/frontend/package.json",
		"testproject/frontend/tailwind.config.js",
		"testproject/frontend/tsconfig.json",
		"testproject/frontend/src/index.ts",
		"testproject/frontend/src/models/hello.ts",
		"testproject/frontend/src/views/Button.ts",
		"testproject/frontend/src/views/pages/HomePage.ts",
	}

	for _, file := range expectedFiles {
		_, err := os.Stat(file)
		if err != nil {
			if removalErr := os.RemoveAll("testproject"); removalErr != nil {
				panic(removalErr)
			}
			panic(err)
		}
	}

	if err := os.RemoveAll("testproject"); err != nil {
		panic(err)
	}
}

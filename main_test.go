package main

import (
	"os"
	"testing"
)

// Deletes a file, useful for cleanup
func delete(filename string) {
	if err := os.RemoveAll(filename); err != nil {
		panic(err)
	}
}

// End-to-end test of the "new" command
func TestRunWithNewCmd(t *testing.T) {
	// Arrange
	args := []string{"new", "testproject"}

	// Act
	run(args)

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
			delete("testproject")
			panic(err)
		}
	}
	delete("testproject")
}

type addCmdConfig struct {
	args         []string
	expectedFile string
	cleanup      string
}

// End-to-end test for the "add" command
func TestRunWithAddCmd(t *testing.T) {
	// Arrange
	inputs := []addCmdConfig{
		// args, expected generated file
		{args: []string{"add", "controller", "vehicle"}, expectedFile: "controllers/vehicle.go", cleanup: "controllers/vehicle.go"},
		{args: []string{"add", "service", "vehicle"}, expectedFile: "services/vehicle.go", cleanup: "services/vehicle.go"},
		{args: []string{"add", "repository", "vehicle"}, expectedFile: "repositories/vehicle.go", cleanup: "repositories/vehicle.go"},
		{args: []string{"add", "model", "vehicle"}, expectedFile: "frontend/src/models/vehicle.ts", cleanup: "frontend"},
		{args: []string{"add", "view", "vehicle"}, expectedFile: "frontend/src/views/vehicle.ts", cleanup: "frontend"},
	}

	for _, input := range inputs {

		// Act
		run(input.args)

		// Assert
		if _, err := os.Stat(input.expectedFile); err != nil {
			panic(err)
		}

		bytes, err := os.ReadFile(input.expectedFile)
		if err != nil {
			delete(input.cleanup)
			panic(err)
		}

		delete(input.cleanup)

		if string(bytes) == "" {
			t.Fatalf("%v did not have any contents", input.expectedFile)
		}

	}
}

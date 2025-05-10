package controllers

import (
	"errors"
	"fmt"
	"strings"
)

type projectInitialiser interface {
	InitProject(username, projectName, projectDir string) error
}

type NewController struct {
	initialiser projectInitialiser
	filesystem  FilesystemRoot
}

func NewNewController(initialiser projectInitialiser, filesystem FilesystemRoot) NewController {
	return NewController{initialiser, filesystem}
}

func (c NewController) Handle(args []string) error {
	if len(args) == 0 || args[0] != "new" {
		return errors.New("passed to incorrect controller! Passed to `new` controller")
	}

	if len(args) < 2 || args[1] == "" {
		return errors.New("expected argument [project-name]")
	}

	projectName := strings.TrimSuffix(args[1], "/") // Ensuring no path is accidentally included

	if err := c.initialiser.InitProject("danielronalds", projectName, projectName); err != nil { // TODO: make user configurable username
		return fmt.Errorf("unable to create project \"%v\": %v", projectName, err)
	}

	fmt.Printf("Created \"%v\" project", projectName)

	return nil
}

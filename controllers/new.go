package controllers

import (
	"errors"
	"fmt"
	"strings"
)

type ProjectInitialiser interface {
	InitProject(username string, directoryName string) error
}

type NewController struct {
	initialiser ProjectInitialiser
}

func NewNewController(initialiser ProjectInitialiser) NewController {
	return NewController{initialiser}
}

func (c NewController) Handle(args []string) error {
	if len(args) == 0 || args[0] != "new" {
		return errors.New("passed to incorrect controller! Passed to `new` controller")
	}

	if len(args) < 2 || args[1] == "" {
		return errors.New("expected argument [project-name]")
	}

	projectName := strings.TrimSuffix(args[1], "/") // Ensuring no path is accidentally included

	c.initialiser.InitProject("danielronalds", projectName) // TODO: make user configurable username

	fmt.Printf("Created \"%v\" project", projectName)

	return nil
}

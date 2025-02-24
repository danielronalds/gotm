package controllers

import (
	"errors"
	"fmt"
)

type ProjectInitialiser interface {
	InitProject(directoryName string) error
}

type NewController struct {
	initialiser ProjectInitialiser
}

func NewNewController(initialiser ProjectInitialiser) NewController {
	return NewController{initialiser}
}

func (c NewController) HandleCmd(args []string) error {
	if len(args) == 0 || args[0] != "new" {
		return errors.New("passed to incorrect controller! Passed to `new` controller")
	}

	if len(args) < 2 || args[1] == "" {
		return errors.New("expected argument [project-name]")
	}

	projectName := args[1]

	c.initialiser.InitProject(projectName)

	fmt.Printf("Created \"%v\" project", projectName)

	return nil
}

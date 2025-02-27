package controllers

import (
	"errors"
	"fmt"
	"strings"
)

type InitController struct {
	initialiser ProjectInitialiser
}

func NewInitController(initialiser ProjectInitialiser) InitController {
	return InitController{initialiser}
}

func (c InitController) Handle(args []string) error {
	if len(args) == 0 || args[0] != "init" {
		return errors.New("passed to incorrect controller! Passed to `init` controller")
	}

	if len(args) < 2 || args[1] == "" {
		return errors.New("expected argument [project-name]")
	}

	projectName := strings.TrimSuffix(args[1], "/") // Ensuring no path is accidentally included

	if err := c.initialiser.InitProject("danielronalds", projectName, "."); err != nil { // TODO: make user configurable username
		return fmt.Errorf("unable to create project \"%v\": %v", projectName, err)
	}

	fmt.Printf("Initialised \"%v\" project", projectName)

	return nil
}

package controllers

import (
	"errors"
	"fmt"
)

type ProjectCreater interface {
	CreateProject(directoryName string) error
}

type FilesystemReaderWriter interface {
	HasDirectory(directory string) (bool, error)
	CreateDirectory(directory string) error
}

type NewController struct {
	creater    ProjectCreater
	filesystem FilesystemReaderWriter
}

func NewNewController(filesystem FilesystemReaderWriter, creater ProjectCreater) NewController {
	return NewController{creater, filesystem}
}

func (c NewController) HandleCmd(args []string) error {
	if len(args) == 0 || args[0] != "new" {
		return errors.New("passed to incorrect controller! Passed to `new` controller")
	}

	if len(args) < 2 || args[1] == "" {
		return errors.New("expected argument [project-name]")
	}

	projectName := args[1]

	if hasDir, err := c.filesystem.HasDirectory(projectName); err != nil || hasDir {
		return fmt.Errorf("%v already exists in this directory", projectName)
	}

	c.filesystem.CreateDirectory(projectName)

	c.creater.CreateProject(projectName)

	fmt.Printf("Created \"%v\" project", projectName)

	return nil
}

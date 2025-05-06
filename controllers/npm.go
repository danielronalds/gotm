package controllers

import "errors"

type npmRunner interface {
	RunNpm(args []string) error
}

type NpmController struct {
	npm npmRunner
}

func NewNpmController(npm npmRunner) NpmController {
	return NpmController{npm}
}

func (c NpmController) Handle(args []string) error {
	if len(args) == 0 || args[0] != "npm" {
		return errors.New("passed to incorrect controller! Passed to `npm` controller")
	}

	c.npm.RunNpm(args[1:])

	return nil
}

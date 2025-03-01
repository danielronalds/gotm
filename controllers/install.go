package controllers

import (
	"errors"
	"fmt"
)

type ShellService interface {
	ExecuteCmdWithPipedOutput(dir, program string, args ...string) error
}

type InstallController struct {
	shell ShellService
}

func NewInstallController(shell ShellService) InstallController {
	return InstallController{shell}
}

func (c InstallController) Handle(args []string) error {
	if len(args) == 0 || args[0] != "install" {
		return errors.New("passed to incorrect controller! Passed to `install` controller")
	}

	fmt.Println("Installing npm deps")
	if err := c.shell.ExecuteCmdWithPipedOutput("frontend", "npm", "install"); err != nil {
		return fmt.Errorf("unable to install project deps: %v", err.Error())
	}

	fmt.Println("\nInstalling go deps")
	if err := c.shell.ExecuteCmdWithPipedOutput(".", "go", "mod", "tidy"); err != nil {
		return fmt.Errorf("unable to install project deps: %v", err.Error())
	}

	return nil
}

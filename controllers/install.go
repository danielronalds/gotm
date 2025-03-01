package controllers

import (
	"errors"
	"fmt"
)

type DependencyInstaller interface {
	InstallNpmDeps(projectRoot string) error
	InstallGoDeps(projectRoot string) error
}

type InstallController struct {
	installer DependencyInstaller
}

func NewInstallController(shell DependencyInstaller) InstallController {
	return InstallController{shell}
}

func (c InstallController) Handle(args []string) error {
	if len(args) == 0 || args[0] != "install" {
		return errors.New("passed to incorrect controller! Passed to `install` controller")
	}

	fmt.Println("Installing npm deps")
	if err := c.installer.InstallNpmDeps("."); err != nil {
		return fmt.Errorf("unable to install project deps: %v", err.Error())
	}

	fmt.Println("\nInstalling go deps")
	if err := c.installer.InstallGoDeps("."); err != nil {
		return fmt.Errorf("unable to install project deps: %v", err.Error())
	}

	return nil
}

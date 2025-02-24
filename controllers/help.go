package controllers

import (
	"fmt"
	"os"
)

type HelpController struct{}

func NewHelpController() HelpController {
	return HelpController{}
}

func (c HelpController) HandleCmd(args []string) error {
	if len(args) > 0 && args[0] != "help" {
		fmt.Fprintf(os.Stderr, "\"%v\" is an unknown command\n\n", args[0])
	}

	help := `gotm v0.0.1

A cli tool building opinionated full stack web applications with the GOTM stack

Commands
  help        Show this menu
`

	fmt.Println(help)

	return nil
}

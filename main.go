package main

import (
	"fmt"
	"os"
	"strings"

	c "github.com/danielronalds/gotm/controllers"
)

type Controller interface {
	Handle(args []string) error
}

func main() {
	args := os.Args[1:] // Removing program name

	var controller Controller

	command := "help"
	if len(args) != 0 {
		command = strings.ToLower(args[0])
	}

	switch command {
	default:
		controller = c.HelpController{}
	}

	err := controller.Handle(args)

	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err.Error())
		os.Exit(1)
	}

}

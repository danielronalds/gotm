package main

import (
	"fmt"
	"os"
	"strings"

	c "github.com/danielronalds/gotm/controllers"
	r "github.com/danielronalds/gotm/repositories"
	s "github.com/danielronalds/gotm/services"
)

type Controller interface {
	Handle(args []string) error
}

func run(args []string) {
	filesystem := r.NewFilesystemRepository()
	templates := r.NewTemplatesRepository()

	initService := s.NewInitialiserService(filesystem, templates)
	componentService := s.NewComponentService(filesystem, templates)

	cmd := "help" // Default command is the help command
	if len(args) != 0 {
		cmd = strings.ToLower(args[0])
	}

	controllerMap := map[string]Controller{
		"new": c.NewNewController(initService),
		"add": c.NewAddController(componentService),
	}
	controller, ok := controllerMap[cmd]
	if !ok {
		controller = c.NewHelpController()
	}

	if err := controller.Handle(args); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err.Error())
		os.Exit(1)
	}
}

func main() {
	args := os.Args[1:] // Removing program name

	run(args)
}

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
	HandleCmd(args []string) error
}

func main() {
	filesystem := r.NewFilesystemRepository()
	templates := r.NewTemplatesRepository()

	initService := s.NewInitialiserService(filesystem, templates)

	args := os.Args[1:] // Removing program name

	var controller Controller

	command := "help"
	if len(args) != 0 {
		command = strings.ToLower(args[0])
	}

	switch command {
	case "new":
		controller = c.NewNewController(initService)
	default:
		controller = c.NewHelpController()
	}

	err := controller.HandleCmd(args)

	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err.Error())
		os.Exit(1)
	}

}

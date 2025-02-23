package main

import (
	"fmt"
	"os"
	"strings"

	c "github.com/danielronalds/gotm/controllers"
	s "github.com/danielronalds/gotm/services"
)

type Controller interface {
	HandleCmd(args []string) error
}

func main() {
	filesystemService := s.NewFilesystemService()
	initService := s.NewInitialiserService()

	args := os.Args[1:] // Removing program name

	var controller Controller

	command := "help"
	if len(args) != 0 {
		command = strings.ToLower(args[0])
	}

	switch command {
	case "new":
		controller = c.NewNewController(filesystemService, initService)
	default:
		controller = c.NewHelpController()
	}

	err := controller.HandleCmd(args)

	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err.Error())
		os.Exit(1)
	}

}

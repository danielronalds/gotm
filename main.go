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
	shell := r.NewShellRepository()

	projectRoot := "."

	initService := s.NewInitialiserService(filesystem, templates)
	componentService := s.NewComponentService(filesystem, templates)
	buildService := s.NewBuildService(filesystem, shell, projectRoot)
	filewatcherService := s.NewFilewatcherService(filesystem)
	runnerService := s.NewRunnerService(projectRoot)

	cmd := "help" // Default command is the help command
	if len(args) != 0 {
		cmd = strings.ToLower(args[0])
	}

	controllerMap := map[string]Controller{
		"new":     c.NewNewController(initService),
		"init":    c.NewInitController(initService),
		"install": c.NewInstallController(buildService),
		"add":     c.NewAddController(componentService),
		"watch":   c.NewWatchController(filewatcherService, buildService, &runnerService),
	}
	controller, ok := controllerMap[cmd]
	if !ok {
		controller = c.NewHelpController()
	}

	if err := controller.Handle(args); err != nil {
		fmt.Fprintf(os.Stderr, "\n%v\n", err.Error())
		os.Exit(1)
	}
}

func main() {
	args := os.Args[1:] // Removing program name

	run(args)
}

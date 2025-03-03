package controllers

import (
	"fmt"
	"os"
	"time"
)

type FileWatcher interface {
	UpdateCache(directory string) error
	HaveFilesChanged(directory string) (bool, error)
}

type ProjectBuilder interface {
	Build(dev bool, outputDir string) error
}

type ProjectRunner interface {
	Run() error
	Stop() error
}

type WatchController struct {
	filewatcher FileWatcher
	builder     ProjectBuilder
	runner      ProjectRunner
}

func NewWatchController(watcher FileWatcher, builder ProjectBuilder, runner ProjectRunner) WatchController {
	return WatchController{watcher, builder, runner}
}

func (c WatchController) Handle(args []string) error {
	fmt.Println("Watching project")

	// FIXME: Crashes if file is deleted?

	for {
		time.Sleep(16 * time.Microsecond)

		projectChanged, err := c.filewatcher.HaveFilesChanged(".")
		if err != nil {
			return fmt.Errorf("failed to detect project changes: %v", err.Error())
		}

		if projectChanged {
			fmt.Println("\nDetected changes, rebuilding project")

			if err := c.builder.Build(true, "build"); err != nil {
				fmt.Fprintf(os.Stderr, "failed to build project:\n %v", err.Error())
				continue
			}
			if err := c.filewatcher.UpdateCache("."); err != nil {
				fmt.Fprintf(os.Stderr, "failed to update file cache: %v\n", err.Error())
			}

			c.runner.Stop()
			c.runner.Run()
		}

	}
}

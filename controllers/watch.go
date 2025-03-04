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
	DevBuild(projectRoot string) error
}

type ProjectRunner interface {
	Run() error
	Stop() error
}

type WatchController struct {
	filewatcher FileWatcher
	builder     ProjectBuilder
	runner      ProjectRunner
	filesystem  FilesystemRoot
}

func NewWatchController(watcher FileWatcher, builder ProjectBuilder, runner ProjectRunner, filesystem FilesystemRoot) WatchController {
	return WatchController{watcher, builder, runner, filesystem}
}

func (c WatchController) Handle(args []string) error {
	fmt.Println("Watching project")

	// FIXME: Crashes if file is deleted?

	for {
		time.Sleep(50 * time.Microsecond)

		projectChanged, err := c.filewatcher.HaveFilesChanged(c.filesystem.Root())
		if err != nil {
			return fmt.Errorf("failed to detect project changes: %v", err.Error())
		}

		if projectChanged {
			fmt.Println("\nDetected changes, rebuilding project")

			if err := c.builder.DevBuild(c.filesystem.Root()); err != nil {
				fmt.Fprintf(os.Stderr, "failed to build project:\n %v", err.Error())
				continue
			}
			if err := c.filewatcher.UpdateCache(c.filesystem.Root()); err != nil {
				fmt.Fprintf(os.Stderr, "failed to update file cache: %v\n", err.Error())
			}

			fmt.Println()

			c.runner.Stop()
			c.runner.Run()
		}

	}
}

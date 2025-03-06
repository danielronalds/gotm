package controllers

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type FileWatcher interface {
	UpdateCache(directory string) error
	HaveFilesChanged(directory string) (bool, error)
}

type ProjectBuilder interface {
	BuildDev() error
	CleanupDev() error
}

type ProjectRunner interface {
	Run() error
	Stop() error
}

type FileSystemWriter interface {
	DeleteFileRecursive(filename string) error
}

type FilesystemRootWriter interface {
	FilesystemRoot
	FileSystemWriter
}

type WatchController struct {
	filewatcher FileWatcher
	builder     ProjectBuilder
	runner      ProjectRunner
	filesystem  FilesystemRootWriter
}

func NewWatchController(watcher FileWatcher, builder ProjectBuilder, runner ProjectRunner, filesystem FilesystemRootWriter) WatchController {
	return WatchController{watcher, builder, runner, filesystem}
}

// Function for actions needing be run before exiting the application
func (c WatchController) cleanup() {
	c.runner.Stop()
	c.builder.CleanupDev()
}

func (c WatchController) Handle(args []string) error {
	fmt.Println("Watching project")

	ch := make(chan os.Signal, 2)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-ch
		c.cleanup()
		os.Exit(1)
	}()

	for {
		time.Sleep(50 * time.Microsecond)

		projectChanged, err := c.filewatcher.HaveFilesChanged(c.filesystem.Root())
		if err != nil {
			return fmt.Errorf("failed to detect project changes: %v", err.Error())
		}

		if projectChanged {
			fmt.Println("\nDetected changes, rebuilding project")

			if err := c.builder.BuildDev(); err != nil {
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

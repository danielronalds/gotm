package controllers

import (
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"syscall"
	"time"
)

type FileWatcher interface {
	UpdateCache(directory string) error
	HaveFilesChanged(directory string) ([]string, error)
}

type ProjectBuilder interface {
	BuildDev(frontend, backend bool) error
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

		filesChanged, err := c.filewatcher.HaveFilesChanged(c.filesystem.Root())
		if err != nil {
			return fmt.Errorf("failed to detect project changes: %v", err.Error())
		}

		if len(filesChanged) != 0 {
			buildFrontend := anyMatchRegex(filesChanged, `.*\.(js|ts)$`)
			buildBackend := anyMatchRegex(filesChanged, `.*\.go$`)

			if buildFrontend || buildBackend {
				fmt.Println("\nDetected changes, rebuilding project")
			}

			if err := c.builder.BuildDev(buildFrontend, buildBackend); err != nil {
				fmt.Fprintf(os.Stderr, "failed to build project:\n %v", err.Error())
				continue
			}
			if err := c.filewatcher.UpdateCache(c.filesystem.Root()); err != nil {
				fmt.Fprintf(os.Stderr, "failed to update file cache: %v\n", err.Error())
			}

			fmt.Println()

			// Only restart if changes have happend to the backend
			if buildBackend {
				c.runner.Stop()
				c.runner.Run()
			}
		}

	}
}

func anyMatchRegex(slice []string, regex string) bool {
	rgx, err := regexp.Compile(regex)
	if err != nil {
		// Panicing if theres an error in the regex as this should be defined at comp time
		panic(err)
	}

	for _, s := range slice {
		if rgx.MatchString(s) {
			return true
		}
	}

	return false
}

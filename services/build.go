package services

import (
	"fmt"
)

type BuildService struct {
	filesystem FilesystemReaderWriter
	shell      Shell
}

func NewBuildService(filesystem FilesystemReaderWriter, shell Shell) BuildService {
	return BuildService{filesystem, shell}
}

func (s BuildService) InstallNpmDeps() error {
	return s.shell.ExecuteCmdWithPipedOutput(s.filesystem.FromRoot("frontend"), "npm", "install")
}

func (s BuildService) InstallGoDeps() error {
	return s.shell.ExecuteCmdWithPipedOutput(s.filesystem.Root(), "go", "mod", "tidy")
}

func (s BuildService) buildGoBin(binName string) error {
	return s.shell.ExecuteCmdWithPipedOutput(s.filesystem.Root(), "go", "build", "-o", binName, ".")
}

func (s BuildService) buildFrontend() error {
	return s.shell.ExecuteCmdWithPipedOutput(s.filesystem.FromRoot("frontend"), "npm", "run", "build")
}

func (s BuildService) DevBuild(outputDir string) error {
	// Building go project
	binName := ".main.tmp"
	if err := s.buildGoBin(binName); err != nil {
		return fmt.Errorf("unable to build go binary: %v", err)
	}

	// Building frontend
	if err := s.buildFrontend(); err != nil {
		return fmt.Errorf("unable to build frontend: %v", err)
	}

	return nil
}

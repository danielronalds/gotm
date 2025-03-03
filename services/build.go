package services

import (
	"fmt"
	"strings"
)

type BuildService struct {
	filesystem  FilesystemReaderWriter
	shell       Shell
	projectRoot string
}

func NewBuildService(filesystem FilesystemReaderWriter, shell Shell, projectRoot string) BuildService {
	return BuildService{filesystem, shell, projectRoot}
}

// Util function for appending the project root to the passed in path
func (s BuildService) fromRoot(path string) string {
	return fmt.Sprintf("%v/%v", s.projectRoot, strings.TrimPrefix(path, "/"))
}

func (s BuildService) InstallNpmDeps() error {
	return s.shell.ExecuteCmdWithPipedOutput(s.fromRoot("frontend"), "npm", "install")
}

func (s BuildService) InstallGoDeps() error {
	return s.shell.ExecuteCmdWithPipedOutput(s.projectRoot, "go", "mod", "tidy")
}

func (s BuildService) buildGoBin(binName string) error {
	return s.shell.ExecuteCmdWithPipedOutput(s.projectRoot, "go", "build", "-o", binName, ".")
}

func (s BuildService) buildFrontend() error {
	return s.shell.ExecuteCmdWithPipedOutput(s.fromRoot("frontend"), "npm", "run", "build")
}

func (s BuildService) Build(dev bool, outputDir string) error {
	// Creating Output dir
	hasDir, err := s.filesystem.HasDirectoryOrFile(s.fromRoot(outputDir))
	if err != nil {
		return fmt.Errorf("unable to check if \"%v\" dir exists: %v", outputDir, err.Error())
	}
	if hasDir {
		if err := s.filesystem.DeleteFileRecursive(outputDir); err != nil {
			return fmt.Errorf("unable to delete \"%v\" dir: %v", outputDir, err)
		}
	}
	if err := s.filesystem.CreateDirectory(s.fromRoot(outputDir)); err != nil {
		return fmt.Errorf("unable to create \"%v\" dir: %v", outputDir, err.Error())
	}

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

package services

import (
	"fmt"
)

const DEV_BULID_GO_BIN = ".main.tmp"

type BuildServiceFilesystem interface {
	ProjectRoot
	DirReader
	FileDeleter
}

type BuildService struct {
	filesystem BuildServiceFilesystem
	shell      CmdRunner
}

func NewBuildService(filesystem BuildServiceFilesystem, shell CmdRunner) BuildService {
	return BuildService{filesystem, shell}
}

func (s BuildService) InstallNpmDeps() error {
	return s.shell.RunCmdWithPipedOutput(s.filesystem.FromRoot("frontend"), "npm", "install")
}

func (s BuildService) InstallGoDeps() error {
	return s.shell.RunCmdWithPipedOutput(s.filesystem.Root(), "go", "mod", "tidy")
}

func (s BuildService) buildGoBin(binName string) error {
	return s.shell.RunCmdWithPipedOutput(s.filesystem.Root(), "go", "build", "-o", binName, ".")
}

func (s BuildService) buildFrontend() error {
	return s.shell.RunCmdWithPipedOutput(s.filesystem.FromRoot("frontend"), "npm", "run", "build")
}

func (s BuildService) BuildDev() error {
	// Building go project
	if err := s.buildGoBin(DEV_BULID_GO_BIN); err != nil {
		return fmt.Errorf("unable to build go binary: %v", err)
	}

	// Building frontend
	if err := s.buildFrontend(); err != nil {
		return fmt.Errorf("unable to build frontend: %v", err)
	}

	return nil
}

func (s BuildService) CleanupDev() error {
	hasFile, err := s.filesystem.HasDirectoryOrFile(s.filesystem.FromRoot(DEV_BULID_GO_BIN))
	if err != nil {
		return fmt.Errorf("failed to detect if dev go binary exists: %v", err)
	}
	if !hasFile {
		return nil
	}

	if err := s.filesystem.DeleteFileRecursive(DEV_BULID_GO_BIN); err != nil {
		return fmt.Errorf("failed to delete dev go binary: %v", err)
	}

	return nil
}

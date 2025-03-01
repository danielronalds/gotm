package services

import (
	"fmt"
)

type BuildService struct {
	shell Shell
}

func NewBuildService(shell Shell) BuildService {
	return BuildService{shell}
}

func (s BuildService) InstallNpmDeps(projectRoot string) error {
	dir := fmt.Sprintf("%v/frontend", projectRoot)
	return s.shell.ExecuteCmdWithPipedOutput(dir, "npm", "install")
}

func (s BuildService) InstallGoDeps(projectRoot string) error {
	return s.shell.ExecuteCmdWithPipedOutput(projectRoot, "go", "mody", "tidy")
}

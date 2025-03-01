package services

import (
	"os"
	"os/exec"
)

type ShellService struct{}

func NewShellService() ShellService {
	return ShellService{}
}

func (s ShellService) ExecuteCmdWithPipedOutput(dir, program string, args ...string) error {
	workdir, err := os.Getwd()
	if err != nil {
		return err
	}
	if err := os.Chdir(dir); err != nil {
		return err
	}

	cmd := exec.Command(program, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	return os.Chdir(workdir)
}

package repositories

import (
	"os"
	"os/exec"
)

type ShellRepository struct{}

func NewShellRepository() ShellRepository {
	return ShellRepository{}
}

func (r ShellRepository) ExecuteCmdWithWithOutput(dir, program string, args ...string) (string, error) {
	workdir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	if err := os.Chdir(dir); err != nil {
		return "", err
	}

	output := ""

	cmd := exec.Command(program, args...)
	cmd.CombinedOutput()

	return output, os.Chdir(workdir)
}

func (r ShellRepository) ExecuteCmdWithPipedOutput(dir, program string, args ...string) error {
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

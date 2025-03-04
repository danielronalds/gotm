package services

import (
	"fmt"
	"os"
	"os/exec"
)

type RunnerService struct {
	process     *exec.Cmd
	projectRoot string
}

func NewRunnerService(projectRoot string) RunnerService {
	var process *exec.Cmd = nil
	return RunnerService{process, projectRoot}
}

func (s *RunnerService) Run() error {
	if s.process != nil {
		if err := s.Stop(); err != nil {
			return err
		}
	}

	cmd := exec.Command("bash", "-c", fmt.Sprintf("\"%v/.main.tmp\"", s.projectRoot))
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("unable to run project: %v", err)
	}

	s.process = cmd

	return nil
}

func (s *RunnerService) Stop() error {
	if s.process == nil {
		return nil
	}

	if err := s.process.Process.Kill(); err != nil {
		return fmt.Errorf("unable to stop project process: %v", err)
	}

	s.process = nil

	return nil
}

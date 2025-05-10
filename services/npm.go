package services

type NpmService struct {
	filesystem ProjectRoot
	shell      cmdRunner
}

func NewNpmService(filesystem ProjectRoot, shell cmdRunner) NpmService {
	return NpmService{filesystem, shell}
}

func (s NpmService) RunNpm(args []string) error {
	return s.shell.RunCmdWithPipedOutput(s.filesystem.FromRoot("frontend"), "npm", args...)
}

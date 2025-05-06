package services

type NpmService struct {
	filesystem ProjectRoot
	shell      CmdRunner
}

func NewNpmService(filesystem ProjectRoot, shell CmdRunner) NpmService {
	return NpmService{filesystem, shell}
}

func (s NpmService) RunNpm(args []string) error {
	return s.shell.RunCmdWithPipedOutput(s.filesystem.FromRoot("frontend"), "npm", args...)
}

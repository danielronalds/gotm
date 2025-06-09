package services

type SqlcService struct {
	filesystem ProjectRoot
	shell      CmdRunner
}

func NewSqlcService(filesystem ProjectRoot, shell CmdRunner) SqlcService {
	return SqlcService{filesystem, shell}
}

func (s SqlcService) RunSqlc() error {
	return s.shell.RunCmdWithPipedOutput(s.filesystem.Root(), "sqlc", "generate")
}

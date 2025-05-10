package services

import (
	"errors"
	"fmt"
)

type sqliteServiceFilesystem interface {
	ProjectRoot
	DirCreater
	FileCreater
	DirReader
}

type SqliteService struct {
	filesystem sqliteServiceFilesystem
	shell      cmdRunner
	templates  TemplatesWriter
}

func NewSqliteService(filesystem sqliteServiceFilesystem, shell cmdRunner, templates TemplatesWriter) SqliteService {
	return SqliteService{filesystem, shell, templates}
}

func (s SqliteService) Init() error {
	dbRepositoryPath := s.filesystem.FromRoot(fmt.Sprintf("%v/sqlite.go", s.filesystem.FromRoot(REPOSITORIES_DIR)))

	// Checking if sqlite repository already exists
	isDbPresent, err := s.filesystem.HasDirectoryOrFile(dbRepositoryPath)
	if err != nil {
		return fmt.Errorf("unable to check if database already exists: %v", err.Error())
	}
	if isDbPresent {
		return errors.New("database already initialised")
	}

	// First installing the database driver, before creating anything (Fail early approach)
	err = s.shell.RunCmdWithPipedOutput(s.filesystem.Root(), "go", "get", "github.com/mattn/go-sqlite3")
	if err != nil {
		return fmt.Errorf("unable to install sqlite driver: %v", err.Error())
	}

	// Creating database template file
	hasDir, err := s.filesystem.HasDirectoryOrFile(REPOSITORIES_DIR)
	if err != nil {
		return fmt.Errorf("unable to check if repositories directory exists: %v", err.Error())
	}
	if !hasDir {
		if err := s.filesystem.CreateDirectory(REPOSITORIES_DIR); err != nil {
			return fmt.Errorf("unable to create %v directory: %v", REPOSITORIES_DIR, err.Error())
		}
	}

	file, err := s.filesystem.CreateFile(dbRepositoryPath)
	if err != nil {
		return fmt.Errorf("unable to create sqlite repository file: %v", err.Error())
	}
	defer file.Close()

	if err = s.templates.WriteTemplate(file, "sqlite.go.tmpl", nil); err != nil {
		return fmt.Errorf("unable to write to sqlite repository file: %v", err.Error())
	}

	return nil
}

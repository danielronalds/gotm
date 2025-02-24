package services

import "fmt"

type FilesystemReaderWriter interface {
	HasDirectory(directory string) (bool, error)
	CreateDirectory(directory string) error
}

// Service for hanlding initialising new projects
type InitialiserService struct {
	filesystem FilesystemReaderWriter
}

func NewInitialiserService(filesystem FilesystemReaderWriter) InitialiserService {
	return InitialiserService{filesystem}
}

func (s InitialiserService) InitProject(projectName *string) error {
	if projectName != nil {
		if hasDir, err := s.filesystem.HasDirectory(*projectName); err != nil || hasDir {
			return fmt.Errorf("%v already exists in this directory", projectName)
		}

		s.filesystem.CreateDirectory(*projectName)
	}

	// TODO: Write this functionality

	return nil
}

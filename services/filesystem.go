package services

import (
	"errors"
	"io/fs"
	"os"
)

// Service for handling filesystem operations
type FilesystemService struct{}

func NewFilesystemService() FilesystemService {
	return FilesystemService{}
}

func (s FilesystemService) HasDirectory(directory string) (bool, error) {
	_, err := os.Stat(directory)

	if err == nil { // No error means the file/directory exists
		return true, nil
	}

	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}

	return false, err

}

func (s FilesystemService) CreateDirectory(directory string) error {
	if hasDir, err := s.HasDirectory(directory); err != nil || hasDir {
		return errors.New("file with that name already exists")
	}

	// 0755 is the default permission bitmask for a directory in linux
	if err := os.Mkdir(directory, 0755); err != nil {
		return err
	}

	return nil
}

package repositories

import (
	"errors"
	"io/fs"
	"os"
)

// Service for handling filesystem operations
type FilesystemRepository struct{}

func NewFilesystemRepository() FilesystemRepository {
	return FilesystemRepository{}
}

func (r FilesystemRepository) HasDirectory(directory string) (bool, error) {
	_, err := os.Stat(directory)

	if err == nil { // No error means the file/directory exists
		return true, nil
	}

	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}

	return false, err

}

func (r FilesystemRepository) CreateDirectory(directory string) error {
	if hasDir, err := r.HasDirectory(directory); err != nil || hasDir {
		return errors.New("file with that name already exists")
	}

	// 0755 is the default permission bitmask for a directory in linux
	if err := os.Mkdir(directory, 0755); err != nil {
		return err
	}

	return nil
}

func (r FilesystemRepository) CreateFile(filename string) (*os.File, error) {
	file, err := os.Create(filename)

	if err != nil {
		return nil, err
	}

	return file, nil
}

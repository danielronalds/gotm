package repositories

import (
	"errors"
	"io/fs"
	"os"
)

// Repository for handling filesystem operations
type FilesystemRepository struct{}

func NewFilesystemRepository() FilesystemRepository {
	return FilesystemRepository{}
}

func (r FilesystemRepository) HasDirectoryOrFile(directory string) (bool, error) {
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
	if hasDir, err := r.HasDirectoryOrFile(directory); err != nil || hasDir {
		return errors.New("directory with that name already exists")
	}

	// 0755 is the default permission bitmask for a directory in linux
	return os.MkdirAll(directory, 0755)
}

func (r FilesystemRepository) CreateFile(filename string) (*os.File, error) {
	if hasFile, err := r.HasDirectoryOrFile(filename); err != nil || hasFile {
		return nil, errors.New("file with that name already exists")
	}

	return os.Create(filename)
}

package repositories

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
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

func (r FilesystemRepository) ReadDirRecursive(directory string) ([]string, error) {
	files := make([]string, 0)

	err := filepath.WalkDir(directory, func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			files = append(files, path)
		}
		return nil
	})

	return files, err
}

func (r FilesystemRepository) ReadFile(filename string) (string, error) {
	if hasFile, err := r.HasDirectoryOrFile(filename); err != nil || !hasFile {
		return "", errors.New("file with that name does not exist")
	}

	contents, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}

	return string(contents), nil
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

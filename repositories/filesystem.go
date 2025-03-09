package repositories

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// Repository for handling filesystem operations
type FilesystemRepository struct {
	root string
}

func NewFilesystemRepository() (FilesystemRepository, error) {
	root, err := findProjectRoot(".")
	if err != nil {
		return FilesystemRepository{}, err
	}
	return FilesystemRepository{root}, nil
}

// Recursive algorithim to find the root directory of the project, in order to set the programs context
//
// Note: the project root is considred the first parent directory to contain `main.go`
func findProjectRoot(startDir string) (string, error) {
	entries, err := os.ReadDir(startDir)
	if err != nil {
		return "", errors.New("unable to locate project root, are you sure you're in a GOTM project?")
	}

	for _, entry := range entries {
		if entry.Name() == "main.go" {
			return startDir, nil
		}
	}

	return findProjectRoot(fmt.Sprintf("../%v", startDir))
}

func (r FilesystemRepository) Root() string {
	return r.root
}

func (r FilesystemRepository) FromRoot(path string) string {
	return fmt.Sprintf("%v/%v", r.root, strings.TrimPrefix(path, "/"))
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
		if !d.IsDir() && !strings.Contains(d.Name(), "node_modules") {
			files = append(files, path)
		}
		return nil
	})

	return files, err
}

func (r FilesystemRepository) ReadFile(filename string) (string, error) {
	if hasFile, err := r.HasDirectoryOrFile(filename); err != nil || !hasFile {
		return "", nil // If the file doesn't exist, just return an empty string instead of an error
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

func (r FilesystemRepository) DeleteFileRecursive(filename string) error {
	if hasFile, err := r.HasDirectoryOrFile(filename); err != nil || !hasFile {
		return errors.New("file with that name does not exist")
	}

	return os.RemoveAll(filename)
}

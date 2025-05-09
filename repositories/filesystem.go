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
	ignoredDirs map[string]bool
}

func NewFilesystemRepository(dirsToIgnore []string) FilesystemRepository {
	ignoredDirs := make(map[string]bool, 0)
	for _, dir := range dirsToIgnore {
		ignoredDirs[dir] = true
	}

	return FilesystemRepository{ignoredDirs}
}

// Recursive algorithim to find the root directory of the project, in order to set the programs context
//
// Note: the project root is considred the first parent directory to contain `main.go`
func findProjectRoot(startDir string) string {
	entries, err := os.ReadDir(startDir)
	if err != nil {
		fmt.Fprintln(os.Stderr, "unable to locate project root, are you sure you're in a GOTM project?")
		os.Exit(1)
	}

	for _, entry := range entries {
		if entry.Name() == "main.go" {
			return startDir
		}
	}

	return findProjectRoot(fmt.Sprintf("../%v", startDir))
}

func (r FilesystemRepository) Root() string {
	return findProjectRoot(".")
}

func (r FilesystemRepository) FromRoot(path string) string {
	return fmt.Sprintf("%v/%v", findProjectRoot("."), strings.TrimPrefix(path, "/"))
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

func (r FilesystemRepository) isIgnoredDir(dir string) bool {
	for ignoredDir := range r.ignoredDirs {
		if strings.Contains(dir, ignoredDir) {
			return true
		}
	}

	return false
}

func (r FilesystemRepository) ReadDirRecursive(directory string) ([]string, error) {
	files := make([]string, 0)

	err := filepath.WalkDir(directory, func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() && !r.isIgnoredDir(d.Name()) {
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

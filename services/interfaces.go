package services

import (
	"io"
	"os"
)

type cmdRunner interface {
	RunCmdWithPipedOutput(dir, program string, args ...string) error
}

type TemplatesWriter interface {
	WriteTemplate(wr io.Writer, name string, data any) error
}

type FileReader interface {
	ReadFile(filename string) (string, error)
}

type DirReader interface {
	HasDirectoryOrFile(name string) (bool, error)
	ReadDirRecursive(directory string) ([]string, error)
}

type ProjectRoot interface {
	Root() string
	FromRoot(path string) string
}

type DirCreater interface {
	CreateDirectory(directory string) error
}

type FileCreater interface {
	CreateFile(filename string) (*os.File, error)
}

type FileDeleter interface {
	DeleteFileRecursive(filename string) error
}

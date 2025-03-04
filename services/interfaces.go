package services

import (
	"io"
	"os"
)

type Shell interface {
	ExecuteCmdWithPipedOutput(dir, program string, args ...string) error
	ExecuteCmdWithWithOutput(dir, program string, args ...string) (string, error)
}

type TemplatesWriter interface {
	ExecuteTemplate(wr io.Writer, name string, data any) error
}

type FilesystemReader interface {
	HasDirectoryOrFile(directory string) (bool, error)
	ReadDirRecursive(directory string) ([]string, error)
	ReadFile(filename string) (string, error)
	Root() string
	FromRoot(path string) string
}

type FilesystemWriter interface {
	CreateDirectory(directory string) error
	CreateFile(filename string) (*os.File, error)
	DeleteFileRecursive(filename string) error
}

type FilesystemReaderWriter interface {
	FilesystemReader
	FilesystemWriter
}

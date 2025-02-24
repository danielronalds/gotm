package services

import (
	"io"
	"os"
)

type TemplatesWriter interface {
	ExecuteTemplate(wr io.Writer, name string, data any) error
}

type FilesystemReader interface {
	HasDirectoryOrFile(directory string) (bool, error)
}

type FilesystemWriter interface {
	CreateDirectory(directory string) error
	CreateFile(filename string) (*os.File, error)
}

type FilesystemReaderWriter interface {
	FilesystemReader
	FilesystemWriter
}

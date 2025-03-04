package controllers

type FilesystemRoot interface {
	Root() string
	FromRoot(path string) string
}

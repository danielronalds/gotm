package repositories

import (
	"embed"
	"io/fs"
	"text/template"
)

func getAllFilenames(efs embed.FS, pathFiles string) ([]string, error) {
	files, err := fs.ReadDir(efs, pathFiles)
	if err != nil {
		return nil, err
	}

	// only file name
	// 1131 0001-01-01 00:00:00 foo.gohtml -> foo.gohtml
	arr := make([]string, 0, len(files))
	for _, file := range files {
		arr = append(arr, file.Name())
	}

	return arr, nil
}

//go:embed templates/*
var templateFS embed.FS

type TemplatesRepository struct {
	templates *template.Template
}

func NewTemplatesRepository() TemplatesRepository {
	templates, err := template.New("").ParseFS(templateFS,
		"templates/init/*.tmpl",
		"templates/init/controllers/*.tmpl",
		"templates/init/**/*.tmpl",
		"templates/init/frontend/**/*.tmpl",
		"templates/init/frontend/src/**/*.tmpl",
		"templates/init/frontend/src/views/**/*.tmpl",
	)

	if err != nil {
		panic(err)
	}

	return TemplatesRepository{templates}
}

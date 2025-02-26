package repositories

import (
	"embed"
	"io"
	"text/template"
)

//go:embed all:templates/*
var templateFS embed.FS

type TemplatesRepository struct {
	templates *template.Template
}

func NewTemplatesRepository() TemplatesRepository {
	templates, err := template.New("").ParseFS(templateFS,
		"templates/*.tmpl",
		// Initial project stuff
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

func (r TemplatesRepository) ExecuteTemplate(wr io.Writer, name string, data any) error {
	return r.templates.ExecuteTemplate(wr, name, data)
}

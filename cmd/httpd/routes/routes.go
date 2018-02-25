package routes

import (
	"html/template"

	"github.com/maksadbek/dumbwall/internal/database"
)

type Routes struct {
	db *database.Database

	templates *template.Template
}

func New(templatesDirPath string) (*Routes, error) {
	templates, err := template.ParseFiles(templatesDirPath)
	if err != nil {
		return nil, err
	}

	return &Routes{
		templates: templates,
	}, nil
}

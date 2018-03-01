package routes

import (
	"html/template"

	"github.com/maksadbek/dumbwall/internal/auth"
	"github.com/maksadbek/dumbwall/internal/config"
	"github.com/maksadbek/dumbwall/internal/database"
)

type Routes struct {
	db              *database.Database
	auth            *auth.Auth
	recaptchaSecret string

	templates *template.Template
}

func New(c config.Routes, db *database.Database) (*Routes, error) {
	templates, err := template.ParseFiles(c.Templates...)
	if err != nil {
		return nil, err
	}

	authorizer, err := auth.New(c.Certs, "dumbwall.xyz", 3600)
	if err != nil {
		return nil, err
	}

	return &Routes{
		db:              db,
		templates:       templates,
		auth:            authorizer,
		recaptchaSecret: c.RecaptchaSecret,
	}, nil
}

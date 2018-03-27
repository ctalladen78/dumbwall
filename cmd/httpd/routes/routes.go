package routes

import (
	"html/template"
	"os"
	"path"

	"github.com/maksadbek/dumbwall/internal/auth"
	"github.com/maksadbek/dumbwall/internal/config"
	"github.com/maksadbek/dumbwall/internal/database"
	"go.uber.org/zap"
)

type Routes struct {
	db              *database.Database
	auth            *auth.Auth
	recaptchaSecret string
	templates       *template.Template
	logger          *zap.Logger
}

func New(c config.Routes, db *database.Database) (*Routes, error) {
	f, err := os.Open(c.TemplatesDir)
	if err != nil {
		return nil, err
	}

	files, err := f.Readdir(-1)
	if err != nil {
		return nil, err
	}

	templateFiles := []string{}

	for i := range files {
		if files[i].IsDir() {
			continue
		}

		templateFiles = append(templateFiles, path.Join(c.TemplatesDir, files[i].Name()))
	}

	templates, err := template.ParseFiles(templateFiles...)
	if err != nil {
		return nil, err
	}

	authorizer, err := auth.New(c.Certs, "dumbwall.xyz", 99999)
	if err != nil {
		return nil, err
	}

	logger, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}

	return &Routes{
		db:              db,
		templates:       templates,
		auth:            authorizer,
		recaptchaSecret: c.RecaptchaSecret,
		logger:          logger,
	}, nil
}

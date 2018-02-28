package routes

import (
	"html/template"
	"path"

	"github.com/maksadbek/dumbwall/internal/auth"
	"github.com/maksadbek/dumbwall/internal/database"
)

type Routes struct {
	db   *database.Database
	auth *auth.Auth

	templates *template.Template
}

func New(etcPath string) (*Routes, error) {
	templatesDir := path.Join(etcPath, "templates")
	templates, err := template.ParseFiles(
		templatesDir+"/header.tmpl",
		templatesDir+"/footer.tmpl",
		templatesDir+"/signup.tmpl",
		templatesDir+"/index.tmpl",
		templatesDir+"/list.tmpl",
		templatesDir+"/posts/"+"/new.tmpl",
		templatesDir+"/posts/"+"/view.tmpl",
		templatesDir+"/users/"+"/new.tmpl",
		templatesDir+"/users/"+"/profile.tmpl",
	)
	if err != nil {
		return nil, err
	}

	authorizer, err := auth.New(path.Join(etcPath, "certs"), "dumbwall.xyz", 3600)
	if err != nil {
		return nil, err
	}

	return &Routes{
		templates: templates,
		auth:      authorizer,
	}, nil
}

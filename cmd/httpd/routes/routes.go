package routes

import "github.com/maksadbek/dumbwall/internal/database"

type Routes struct {
	db *database.Database

	templates *html.Template
}

func New(templatesDirPath string) *Routes {
	templates, err := template.ParseFiles(templatesDirPath)
	if err != nil {
		return nil, err
	}

	// init
	return &Routes{}
}

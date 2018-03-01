package main

import (
	"net/http"
	"path"

	"github.com/bmizerany/pat"
	"github.com/maksadbek/dumbwall/cmd/httpd/routes"
	"github.com/maksadbek/dumbwall/internal/config"
	"github.com/maksadbek/dumbwall/internal/database"
)

type httpd struct {
	mux *pat.PatternServeMux
}

func (h *httpd) init(etcPath string) error {
	println("init httpd")

	c, err := config.New(path.Join(etcPath, "dumbwall.toml"))
	if err != nil {
		return err
	}

	db, err := database.New(c.Database)
	if err != nil {
		return err
	}

	r, err := routes.New(c.Routes, db)
	if err != nil {
		return err
	}

	m := pat.New()

	m.Get("/posts/new", http.HandlerFunc(r.NewPost))
	m.Post("/posts", http.HandlerFunc(r.CreatePost))
	m.Get("/posts/:id", http.HandlerFunc(r.Post))
	m.Post("/posts/:id/up", http.HandlerFunc(r.UpPost))
	m.Post("/posts/:id/down", http.HandlerFunc(r.DownPost))
	m.Post("/posts/:id/comment", http.HandlerFunc(r.CreateComment))
	m.Post("/comments/:id/comment", http.HandlerFunc(r.ReplyComment))

	m.Get("/listings/hot", http.HandlerFunc(r.Hot))
	m.Get("/listings/top", http.HandlerFunc(r.Top))
	m.Get("/listings/best", http.HandlerFunc(r.Best))
	m.Get("/listings/controversial", http.HandlerFunc(r.Controversial))
	m.Get("/listings/rising", http.HandlerFunc(r.Rising))
	m.Get("/listings/newest", http.HandlerFunc(r.Newest))

	m.Get("/me", http.HandlerFunc(r.Profile))
	m.Get("/me/edit", http.HandlerFunc(r.EditProfile))
	m.Get("/signout", http.HandlerFunc(r.SignOut))
	m.Get("/signup", http.HandlerFunc(r.SignUp))
	m.Get("/signin", http.HandlerFunc(r.SignIn))

	m.Post("/users/create", http.HandlerFunc(r.CreateUser))
	m.Get("/users/:id", http.HandlerFunc(r.User))

	m.Get("/", http.HandlerFunc(r.Best))

	h.mux = m

	return nil
}

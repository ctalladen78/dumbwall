package routes

import (
	"net/http"

	"github.com/maksadbek/dumbwall/internal/posts"
)

func (r *Routes) Hot(w http.ResponseWriter, req *http.Request) {

}

func (r *Routes) Rising(w http.ResponseWriter, req *http.Request) {

}

func (r *Routes) Controversial(w http.ResponseWriter, req *http.Request) {

}

func (r *Routes) Best(w http.ResponseWriter, req *http.Request) {

}

func (r *Routes) Top(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()

	topPosts, errs := r.db.Top(0, 20)
	if len(errs) > 0 {
		panic(errs)
	}

	ctx := struct {
		Posts []posts.Post
	}{
		Posts: topPosts,
	}

	r.templates.ExecuteTemplate(w, "list", ctx)

}

func (r *Routes) Newest(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()

	newestPosts, errs := r.db.Newest(0, 20)
	if len(errs) > 0 {
		panic(errs)
	}

	ctx := struct {
		Posts []posts.Post
	}{
		Posts: newestPosts,
	}

	r.templates.ExecuteTemplate(w, "list", ctx)
}

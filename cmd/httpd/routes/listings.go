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
	userID, tokenErr := r.validateToken(w, req)

	topPosts, errs := r.db.Top(0, 20)
	if len(errs) > 0 {
		panic(errs)
	}

	if tokenErr == nil {
		for i, v := range topPosts {
			userAction, err := r.db.CheckVotes(userID, v.ID)
			if err != nil {
				panic(err)
				return
			}

			println(v.ID, userAction[0])
			topPosts[i].Meta.Action = userAction[0]
		}
	}

	ctx := struct {
		Posts []posts.Post
	}{
		Posts: topPosts,
	}

	err := r.templates.ExecuteTemplate(w, "list", ctx)
	if err != nil {
		panic(err)
	}
}

func (r *Routes) Newest(w http.ResponseWriter, req *http.Request) {
	userID, tokenErr := r.validateToken(w, req)

	newestPosts, errs := r.db.Newest(0, 20)
	if len(errs) > 0 {
		panic(errs)
	}

	if tokenErr == nil {
		for i, v := range newestPosts {
			userAction, err := r.db.CheckVotes(userID, v.ID)
			if err != nil {
				panic(err)
				return
			}

			println(v.ID, userAction[0])
			newestPosts[i].Meta.Action = userAction[0]
		}
	}

	ctx := struct {
		Posts []posts.Post
	}{
		Posts: newestPosts,
	}

	r.templates.ExecuteTemplate(w, "list", ctx)
}

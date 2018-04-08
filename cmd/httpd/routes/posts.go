package routes

import (
	"net/http"
	"strconv"

	"github.com/maksadbek/dumbwall/internal/actions"
	"github.com/maksadbek/dumbwall/internal/posts"
	"go.uber.org/zap"
)

func (r *Routes) NewPost(w http.ResponseWriter, req *http.Request) {
	context := struct {
		flash flash
	}{
		flash: flash{
			Notice: "you're creating a post",
			Alert:  "be careful!",
			Custom: map[string]string{
				"first_ever_post": "hey, post anything what you want",
			},
		},
	}

	r.templates.ExecuteTemplate(w, "new_post", context)
}

func (r *Routes) CreatePost(w http.ResponseWriter, req *http.Request) {
	userID, err := r.validateToken(w, req)
	if err != nil {
		r.logger.Error("failed to validate token", zap.Error(err))
		http.Redirect(w, req, "/signin", http.StatusSeeOther)
		return
	}

	req.ParseForm()

	title := req.PostForm.Get("title")
	body := req.PostForm.Get("body")

	post, err := r.db.CreatePost(userID, posts.Post{
		Title: title,
		Body:  body,
	})

	if err != nil {
		r.logger.Error("failed to create a post", zap.Error(err))
		return
	}

	http.Redirect(w, req, "/posts/"+strconv.Itoa(post.ID), http.StatusFound)
}

func (r *Routes) UpdatePost(w http.ResponseWriter, req *http.Request) {

}

func (r *Routes) DeletePost(w http.ResponseWriter, req *http.Request) {

}

func (r *Routes) UpPost(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()

	userID, err := r.validateToken(w, req)
	if err != nil {
		r.logger.Error("failed to validate token", zap.Error(err))
		http.Redirect(w, req, "/signin", http.StatusSeeOther)
		return
	}

	postID, err := strconv.Atoi(req.URL.Query().Get(":id"))
	if err != nil {
		r.logger.Error("failed to get id", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = r.db.VotePost(userID, postID, actions.ActionUp)
	if err != nil {
		r.logger.Error("failed to vote post", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (r *Routes) DownPost(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()

	userID, err := r.validateToken(w, req)
	if err != nil {
		r.logger.Error("failed to validate token", zap.Error(err))
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}

	postID, err := strconv.Atoi(req.URL.Query().Get(":id"))
	if err != nil {
		r.logger.Error("failed to get id", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = r.db.VotePost(userID, postID, actions.ActionDown)
	if err != nil {
		r.logger.Error("failed to vote post", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (r *Routes) Post(w http.ResponseWriter, req *http.Request) {
	id, err := strconv.Atoi(req.URL.Query().Get(":id"))
	if err != nil {
		r.logger.Error("invalid post id", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	post, err := r.db.GetPost(id)
	if err != nil {
		r.logger.Error("failed to get post by id", zap.Error(err), zap.Int("id", id))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = r.templates.ExecuteTemplate(w, "view_post", post)
	if err != nil {
		r.logger.Error("failed to render template", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
	}
}

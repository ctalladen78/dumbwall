package routes

import (
	"net/http"

	"strconv"

	"github.com/maksadbek/dumbwall/internal/posts"
	"go.uber.org/zap"
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

	var page, _ = strconv.Atoi(req.Form.Get("page"))
	if page == 0 {
		page = 1
	}

	userID, tokenErr := r.validateToken(w, req)

	topPosts, errs := r.db.Top((page*20)-20, page*20)
	if len(errs) > 0 {
		r.logger.Error("failed to get top posts", zap.Errors("errors", errs))
	}

	if tokenErr == nil {
		for i, v := range topPosts {
			userAction, err := r.db.CheckVotes(userID, v.ID)
			if err != nil {
				r.logger.Error("failed check vote", zap.Error(err), zap.Int("post_id", v.ID))
				continue
			}

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
		r.logger.Error("failed to render template", zap.Error(err))
	}
}

func (r *Routes) Newest(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()

	var page, _ = strconv.Atoi(req.Form.Get("page"))
	if page == 0 {
		page = 1
	}

	userID, tokenErr := r.validateToken(w, req)

	newestPosts, errs := r.db.Newest((page*20)-20, page*20)
	if len(errs) > 0 {
		r.logger.Error("failed to get newewst list", zap.Errors("errors", errs))
		http.Redirect(w, req, "/", http.StatusFound)
		return
	}

	if tokenErr == nil {
		for i, v := range newestPosts {
			userAction, err := r.db.CheckVotes(userID, v.ID)
			if err != nil {
				r.logger.Error("failed to check user votes", zap.Error(err))
				http.Redirect(w, req, "/", http.StatusFound)
				return
			}

			newestPosts[i].Meta.Action = userAction[0]
		}
	}

	ctx := struct {
		Posts []posts.Post
	}{
		Posts: newestPosts,
	}

	err := r.templates.ExecuteTemplate(w, "list", ctx)
	if err != nil {
		http.Redirect(w, req, "/", http.StatusFound)
		r.logger.Error("failed to render template", zap.Error(err))
	}
}

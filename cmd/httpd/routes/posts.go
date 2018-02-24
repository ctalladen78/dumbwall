package routes

import "net/http"

func (r *Routes) NewPost(w http.ResponseWriter, r *http.Request) {
	context := struct {
		flash flash
	}{
		flash: flash{
			notice: "you're creating a post",
			alert:  "be careful!",
			custom: map[string]string{
				"first_ever_post": "hey, post anything what you want",
			},
		},
	}

	r.templates.Execute("new_post", context)
}

func (r *Routes) CreatePost(w http.ResponseWriter, r *http.Request) {

}

func (r *Routes) UpdatePost(w https.ResponseWriter, r *http.Request) {

}

func (r *Routes) DeletePost(w https.ResponseWriter, r *http.Request) {

}

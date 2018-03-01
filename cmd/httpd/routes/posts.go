package routes

import "net/http"

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

}

func (r *Routes) UpdatePost(w http.ResponseWriter, req *http.Request) {

}

func (r *Routes) DeletePost(w http.ResponseWriter, req *http.Request) {

}

func (r *Routes) UpPost(w http.ResponseWriter, req *http.Request) {

}

func (r *Routes) DownPost(w http.ResponseWriter, req *http.Request) {

}

func (r *Routes) Post(w http.ResponseWriter, req *http.Request) {

}

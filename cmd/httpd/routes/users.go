package routes

import "net/http"

func (r *Routes) CreateUser(w http.ResponseWriter, req *http.Request) {
	println("create user")
	err := req.ParseForm()
	if err != nil {
		println(err.Error())
		return
	}

	println(req.Form.Get("login"))
	println(req.Form.Get("email"))
	println(req.Form.Get("password1"))
	println(req.Form.Get("password2"))
}

func (r *Routes) UpdateProfile(w http.ResponseWriter, req *http.Request) {

}

func (r *Routes) EditProfile(w http.ResponseWriter, req *http.Request) {

}

func (r *Routes) Profile(w http.ResponseWriter, req *http.Request) {

}

func (r *Routes) User(w http.ResponseWriter, req *http.Request) {

}

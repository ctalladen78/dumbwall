package routes

import "net/http"

func (r *Routes) SignOut(w http.ResponseWriter, req *http.Request) {

}

func (r *Routes) SignIn(w http.ResponseWriter, req *http.Request) {

}

func (r *Routes) SignUp(w http.ResponseWriter, req *http.Request) {
	println("signup motherfucker")

	err := r.templates.ExecuteTemplate(w, "signup", nil)
	if err != nil {
		println(err.Error())
	}
}

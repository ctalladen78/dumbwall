package routes

import (
	"net/http"

	"go.uber.org/zap"
)

func (r *Routes) SignOut(w http.ResponseWriter, req *http.Request) {

}

func (r *Routes) SignIn(w http.ResponseWriter, req *http.Request) {
	println("signin motherfucker")

	err := r.templates.ExecuteTemplate(w, "signin", nil)
	if err != nil {
		r.logger.Error("failed render template", zap.Error(err))
		http.Redirect(w, req, "/", http.StatusFound)
		return
	}
}

func (r *Routes) SignUp(w http.ResponseWriter, req *http.Request) {
	println("signup motherfucker")

	err := r.templates.ExecuteTemplate(w, "signup", nil)
	if err != nil {
		println(err.Error())
	}
}

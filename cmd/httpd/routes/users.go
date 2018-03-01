package routes

import (
	"encoding/json"
	"net/http"
	"net/url"

	valid "github.com/asaskevich/govalidator"
)

func (r *Routes) CreateUser(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		println(err.Error())
		return
	}

	var recaptchaValues = make(url.Values)

	recaptchaValues.Set("secret", r.recaptchaSecret)
	recaptchaValues.Set("response", req.Form.Get("g-recaptcha-response"))
	recaptchaValues.Set("remoteip", req.RemoteAddr)

	resp, err := http.PostForm("https://www.google.com/recaptcha/api/siteverify", recaptchaValues)
	if err != nil {
		println(err.Error())
		return
	}

	defer resp.Body.Close()

	captchaResponse := struct {
		Success bool `json:"success"`
	}{}

	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&captchaResponse)
	if err != nil {
		println(err.Error())
		return
	}

	if !captchaResponse.Success {
		http.Redirect(w, req, "/signup", http.StatusForbidden)
		return
	}

	f := flash{}

	if !valid.IsEmail(req.Form.Get("email")) {
		f.Errors = append(f.Errors, "invalid email")
	}

	if !valid.ByteLength(req.Form.Get("password1"), "6", "100") {
		f.Errors = append(f.Errors, "password must have minimum 6 characters")
	}

	if req.Form.Get("password1") != req.Form.Get("password2") {
		f.Errors = append(f.Errors, "passwords do not match")
	}

	if valid.IsNull(req.Form.Get("login")) {
		f.Errors = append(f.Errors, "passwords do not match")
	}

	if len(f.Errors) > 0 {
		r.templates.ExecuteTemplate(w, "signup", f)
		return
	}
}

func (r *Routes) UpdateProfile(w http.ResponseWriter, req *http.Request) {

}

func (r *Routes) EditProfile(w http.ResponseWriter, req *http.Request) {

}

func (r *Routes) Profile(w http.ResponseWriter, req *http.Request) {

}

func (r *Routes) User(w http.ResponseWriter, req *http.Request) {

}

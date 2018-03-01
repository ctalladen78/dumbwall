package routes

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	valid "github.com/asaskevich/govalidator"
	"github.com/maksadbek/dumbwall/internal/users"
)

var errInvalidCaptchaResponse = errors.New("invalid captcha response")

func (r *Routes) validateCaptcha(recaptchaResponse, remoteAddr string) error {
	var recaptchaValues = make(url.Values)

	recaptchaValues.Set("secret", r.recaptchaSecret)
	recaptchaValues.Set("response", recaptchaResponse)
	recaptchaValues.Set("remoteip", remoteAddr)

	resp, err := http.PostForm("https://www.google.com/recaptcha/api/siteverify", recaptchaValues)
	if err != nil {
		println(err.Error())
		return err
	}

	defer resp.Body.Close()

	captchaResponse := struct {
		Success bool `json:"success"`
	}{}

	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&captchaResponse)
	if err != nil {
		println(err.Error())
		return err
	}

	if !captchaResponse.Success {
		return errInvalidCaptchaResponse
	}

	return nil
}

func (r *Routes) CreateUser(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		println(err.Error())
		return
	}

	err = r.validateCaptcha(req.Form.Get("g-recaptcha-response"), req.RemoteAddr)
	if err != nil {
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

	user, err := r.db.CreateUser(users.User{
		Login:    req.Form.Get("login"),
		Email:    req.Form.Get("email"),
		Password: req.Form.Get("password1"),
	})
	if err != nil {
		println(err.Error())
		return
	}

	fmt.Println(user)

	return
}

func (r *Routes) UpdateProfile(w http.ResponseWriter, req *http.Request) {

}

func (r *Routes) EditProfile(w http.ResponseWriter, req *http.Request) {

}

func (r *Routes) Profile(w http.ResponseWriter, req *http.Request) {

}

func (r *Routes) User(w http.ResponseWriter, req *http.Request) {
	id, err := valid.ToInt(req.URL.Query().Get(":id"))
	if err != nil {
		println(err.Error())
		http.Redirect(w, req, "/404", http.StatusNotFound)
		return
	}

	user, err := r.db.GetUser(uint64(id))
	if err != nil {
		println(err.Error())
		io.WriteString(w, "not found")
		return
	}

	err = r.templates.ExecuteTemplate(w, "profile", user)
	if err != nil {
		println(err.Error())
	}
}

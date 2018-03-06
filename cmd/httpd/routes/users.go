package routes

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	valid "github.com/asaskevich/govalidator"
	"github.com/maksadbek/dumbwall/internal/users"
	"go.uber.org/zap"
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
	println("captcha success", captchaResponse.Success)
	if !captchaResponse.Success {
		return errInvalidCaptchaResponse
	}

	return nil
}

func (r *Routes) CreateUser(w http.ResponseWriter, req *http.Request) {
	println("create user")
	err := req.ParseForm()
	if err != nil {
		println(err.Error())
		return
	}

	/*
		err = r.validateCaptcha(req.Form.Get("g-recaptcha-response"), req.RemoteAddr)
		if err != nil {
			r.logger.Error("failed to validate captcha", zap.Error(err))
			http.Redirect(w, req, "/signup", http.StatusInternalServerError)
			return
		}
	*/

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

	token, err := r.auth.CreateJWTToken(map[string]string{
		"user_id": strconv.FormatUint(user.ID, 10),
	})
	if err != nil {
		println(err.Error())
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "user_session",
		Value:   token,
		Domain:  "localhost",
		Path:    "/",
		Expires: time.Now().Add(time.Duration(90) * time.Hour),
	})

	println("redirect")
	http.Redirect(w, req, "/me", http.StatusFound)

	return
}

func (r *Routes) UpdateProfile(w http.ResponseWriter, req *http.Request) {

}

func (r *Routes) EditProfile(w http.ResponseWriter, req *http.Request) {

}

func (r *Routes) validateToken(req *http.Request) (int, error) {
	var id int

	cookie, err := req.Cookie("user_session")
	if err != nil {
		r.logger.Error("failed to get token", zap.Error(err))
		return id, err
	}

	claims, err := r.auth.Validate(cookie.Value)
	if err != nil {
		r.logger.Error("failed to validate token", zap.Error(err))
		return id, err
	}

	id, err = strconv.Atoi(claims["user_id"].(string))
	if err != nil {
		r.logger.Error("blet", zap.Error(err))
		return id, err
	}

	return id, err
}

func (r *Routes) Profile(w http.ResponseWriter, req *http.Request) {
	cookie, err := req.Cookie("user_session")
	if err != nil {
		r.logger.Error("failed to get token", zap.Error(err))
		http.Redirect(w, req, "/", http.StatusOK)
		return
	}

	println("cookie", cookie.Value)
	claims, err := r.auth.Validate(cookie.Value)
	if err != nil {
		r.logger.Error("failed to validate token", zap.Error(err))
		http.Redirect(w, req, "/", http.StatusOK)
		return
	}

	id, err := strconv.ParseInt(claims["user_id"].(string), 10, 64)
	if err != nil {
		r.logger.Error("blet", zap.Error(err))
		return
	}

	user, err := r.db.GetUser(uint64(id))
	if err != nil {
		r.logger.Error("blet", zap.Error(err))
		return
	}

	err = r.templates.ExecuteTemplate(w, "profile", user)
	if err != nil {
		r.logger.Error("blet", zap.Error(err))
		return
	}
}

func (r *Routes) User(w http.ResponseWriter, req *http.Request) {
	id, err := valid.ToInt(req.URL.Query().Get(":id"))
	if err != nil {
		println(err.Error())
		http.Redirect(w, req, "/404", http.StatusSeeOther)
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

func (r *Routes) Authenticate(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()

	login := req.PostForm.Get("login")
	password := req.PostForm.Get("password")

	id, err := r.db.Authenticate(login, password)
	if err != nil {
		r.logger.Error("password do not match", zap.Error(err))
		http.Redirect(w, req, "/signin", http.StatusSeeOther)
		return
	}

	token, err := r.auth.CreateJWTToken(map[string]string{
		"user_id": strconv.FormatInt(id, 10),
	})
	if err != nil {
		println(err.Error())
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "user_session",
		Value:   token,
		Domain:  "localhost",
		Path:    "/",
		Expires: time.Now().Add(time.Duration(90) * time.Hour),
	})

	http.Redirect(w, req, "/me", http.StatusFound)
	return
}

package handler

import (
	//"io"
	//"os"
	//"fmt"
	"log"
	//"time"
	//"io/ioutil"
	"net/http"
	//"crypto/sha1"
	"html/template"

	"negentropia/webserv/cfg"
	"negentropia/webserv/session"
)

type SignupPage struct {
	HomePath		  string
	LoginPath		  string
	LogoutPath		  string
	SignupProcessPath string

	EmailValue        string	
	BadPasswdMsg      string
	BadConfirmMsg     string
	BadSignupMsg      string
	
	Account         string
	ShowNavAccount  bool
	ShowNavHome     bool
	ShowNavLogin    bool
	ShowNavLogout   bool	
}

func sendSignup(w http.ResponseWriter, p SignupPage) error {
	p.HomePath          = cfg.HomePath()
	p.LoginPath         = cfg.LoginPath()
	p.LogoutPath        = cfg.LogoutPath()
	p.SignupProcessPath = cfg.SignupProcessPath()
	
	// FIXME: we're loading template every time
    t, err := template.ParseFiles(TemplatePath("base.tpl"), TemplatePath("signup.tpl"))
	if err != nil {
		return err
	}
	if err = t.Execute(w, p); err != nil {
		return err
	}
	return nil
}

func Signup(w http.ResponseWriter, r *http.Request, s *session.Session) {
	path := r.URL.Path
	log.Printf("handler.Signup url=%s", path)
	
	account := accountLabel(s)
	
	if err := sendSignup(w, SignupPage{Account:account,ShowNavAccount:true,ShowNavHome:true}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func SignupProcess(w http.ResponseWriter, r *http.Request, s *session.Session) {
	path := r.URL.Path
	log.Printf("handler.SignupProcess url=%s", path)

	//account := accountLabel(s)
}

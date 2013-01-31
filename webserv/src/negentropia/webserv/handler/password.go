package handler

import (
	//"io"
	//"os"
	//"fmt"
	"log"
	//"time"
	//"io/ioutil"
	//"strconv"
	"net/http"
	//"crypto/sha1"
	"html/template"

	"negentropia/webserv/cfg"
	//"negentropia/webserv/store"	
	"negentropia/webserv/session"
)

type PasswordPage struct {
	HomePath		   string
	LoginPath		   string
	LogoutPath		   string
	ConfirmProcessPath string

	EmailValue        string
	BadEmailMsg       string
	BadConfirmMsg     string
	ConfirmDoneMsg    string	
	
	Account         string
	ShowNavAccount  bool
	ShowNavHome     bool
	ShowNavLogin    bool
	ShowNavLogout   bool	
}

/*
const (
	FORM_VAR_EMAIL = "Email"
	FORM_VAR_CONFIRM_ID = "ConfirmId"
)
*/

func sendResetPass(w http.ResponseWriter, p PasswordPage) error {
	p.HomePath           = cfg.HomePath()
	p.LoginPath          = cfg.LoginPath()
	p.LogoutPath         = cfg.LogoutPath()
	p.ConfirmProcessPath = cfg.ConfirmProcessPath()
	
	// FIXME: we're loading template every time
    t, err := template.ParseFiles(TemplatePath("base.tpl"), TemplatePath("confirm.tpl"))
	if err != nil {
		return err
	}
	if err = t.Execute(w, p); err != nil {
		return err
	}
	return nil
}

func ResetPass(w http.ResponseWriter, r *http.Request, s *session.Session) {
	log.Printf("handler.ResePass url=%s", r.URL.Path)
	
	account := accountLabel(s)
	
	if err := sendResetPass(w, PasswordPage{Account:account,ShowNavAccount:true,ShowNavHome:true}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ResetPassProcess(w http.ResponseWriter, r *http.Request, s *session.Session) {
	log.Printf("handler.ResePassProcess url=%s", r.URL.Path)

	account := accountLabel(s)
	
	email := formatEmail(r.FormValue(FORM_VAR_EMAIL))

	if email == "" {
		msg := "Please enter email address."
		if err := sendResetPass(w, PasswordPage{Account:account,ShowNavAccount:true,ShowNavHome:true,BadEmailMsg:msg}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
}

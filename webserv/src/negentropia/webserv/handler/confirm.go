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
	"negentropia/webserv/store"	
	"negentropia/webserv/session"
)

type ConfirmPage struct {
	HomePath		   string
	SignupPath		   string	
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
	ShowNavSignup   bool	
	ShowNavLogin    bool
	ShowNavLogout   bool	
}

const (
	FORM_VAR_EMAIL = "Email"
	FORM_VAR_CONFIRM_ID = "ConfirmId"
	FORM_VAR_PASSWD = "Passwd"
	FORM_VAR_CONFIRM = "Confirm"
)

func sendConfirm(w http.ResponseWriter, p ConfirmPage) error {
	p.HomePath           = cfg.HomePath()
	p.SignupPath         = cfg.SignupPath()	
	p.LoginPath          = cfg.LoginPath()
	p.LogoutPath         = cfg.LogoutPath()
	p.ConfirmProcessPath = cfg.ConfirmProcessPath()
	
	p.ShowNavSignup = true
	
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

func Confirm(w http.ResponseWriter, r *http.Request, s *session.Session) {
	log.Printf("handler.Confirm url=%s", r.URL.Path)
	
	account := accountLabel(s)
	
	if err := sendConfirm(w, ConfirmPage{Account:account,ShowNavAccount:true,ShowNavHome:true}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ConfirmProcess(w http.ResponseWriter, r *http.Request, s *session.Session) {
	log.Printf("handler.ConfirmProcess url=%s", r.URL.Path)

	account := accountLabel(s)
	
	email := formatEmail(r.FormValue(FORM_VAR_EMAIL))
	confId := r.FormValue(FORM_VAR_CONFIRM_ID)

	if email == "" {
		msg := "Please enter email address."
		if err := sendConfirm(w, ConfirmPage{Account:account,ShowNavAccount:true,ShowNavHome:true,BadEmailMsg:msg}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if !store.Exists(email) {
		msg := "The address " + email + " is not registered."
		if err := sendConfirm(w, ConfirmPage{Account:account,ShowNavAccount:true,ShowNavHome:true,BadEmailMsg:msg,EmailValue:email}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	confEmail := store.Get(confId)	
	if (confEmail != email) {
		msg := "Incorrect confirmation id for address."
		if err := sendConfirm(w, ConfirmPage{Account:account,ShowNavAccount:true,ShowNavHome:true,BadConfirmMsg:msg,EmailValue:email}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
		
	store.DelField(email, "unconfirmed") // remove lock
	store.Persist(email)                 // remove expire
	store.Del(confId)                    // just clean-up
	
	msg := "The address " + email + " has been enabled. You can login now."
	if err := sendConfirm(w, ConfirmPage{Account:account,ShowNavAccount:true,ShowNavHome:true,ConfirmDoneMsg:msg,EmailValue:email}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}	
}

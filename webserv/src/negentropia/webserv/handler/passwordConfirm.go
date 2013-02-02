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

type ResetPassConfirmPage struct {
	HomePath		   string
	LoginPath		   string
	LogoutPath		   string
	ResetPassConfirmProcessPath string
	
	EmailValue        string
	ConfirmIdValue        string
	BadEmailMsg       string
	BadConfirmIdMsg     string	
	BadPasswdMsg       string	
	BadConfirmMsg       string		
	ResetPassConfirmDoneMsg    string	
	
	Account         string
	ShowNavAccount  bool
	ShowNavHome     bool
	ShowNavLogin    bool
	ShowNavLogout   bool	
}

func sendResetPassConfirm(w http.ResponseWriter, p ResetPassConfirmPage) error {
	p.HomePath           = cfg.HomePath()
	p.LoginPath          = cfg.LoginPath()
	p.LogoutPath         = cfg.LogoutPath()
	p.ResetPassConfirmProcessPath = cfg.ResetPassConfirmProcessPath()
	
	// FIXME: we're loading template every time
    t, err := template.ParseFiles(TemplatePath("base.tpl"), TemplatePath("passwordConfirm.tpl"))
	if err != nil {
		return err
	}
	if err = t.Execute(w, p); err != nil {
		return err
	}
	return nil
}

func ResetPassConfirm(w http.ResponseWriter, r *http.Request, s *session.Session) {
	log.Printf("handler.ResetPassConfirm url=%s", r.URL.Path)
	
	account := accountLabel(s)

	email := formatEmail(r.FormValue(FORM_VAR_EMAIL))
	confId := r.FormValue(FORM_VAR_CONFIRM_ID)
	
	if err := sendResetPassConfirm(w, ResetPassConfirmPage{Account:account,ShowNavAccount:true,ShowNavHome:true,EmailValue:email,ConfirmIdValue:confId}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ResetPassConfirmProcess(w http.ResponseWriter, r *http.Request, s *session.Session) {
	log.Printf("handler.ResetPassConfirmProcess url=%s", r.URL.Path)

	account := accountLabel(s)
	
	email := formatEmail(r.FormValue(FORM_VAR_EMAIL))
	confId := r.FormValue(FORM_VAR_CONFIRM_ID)
	passwd := r.FormValue(FORM_VAR_PASSWD)
	confirm := r.FormValue(FORM_VAR_CONFIRM)

	if email == "" {
		msg := "Please enter email address."
		if err := sendResetPassConfirm(w, ResetPassConfirmPage{Account:account,ShowNavAccount:true,ShowNavHome:true,BadEmailMsg:msg,ConfirmIdValue:confId}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if !store.Exists(email) {
		msg := "The address " + email + " is not registered."
		if err := sendResetPassConfirm(w, ResetPassConfirmPage{Account:account,ShowNavAccount:true,ShowNavHome:true,BadEmailMsg:msg,EmailValue:email,ConfirmIdValue:confId}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if passwd != confirm {
		msg := "Passwords don't match."
		if err := sendResetPassConfirm(w, ResetPassConfirmPage{Account:account,ShowNavAccount:true,ShowNavHome:true,BadConfirmMsg:msg,EmailValue:email,ConfirmIdValue:confId}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	confEmail := store.Get(confId)	
	if (confEmail != email) {
		msg := "Incorrect confirmation id for address."
		if err := sendResetPassConfirm(w, ResetPassConfirmPage{Account:account,ShowNavAccount:true,ShowNavHome:true,BadConfirmIdMsg:msg,EmailValue:email,ConfirmIdValue:confId}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	
	store.SetField(email, "password-sha1-hex", passDigest(passwd)) // set new password
	store.Del(confId) // just clean-up

	msg := "The password for address " + email + " has been changed."
	if err := sendResetPassConfirm(w, ResetPassConfirmPage{Account:account,ShowNavAccount:true,ShowNavHome:true,ResetPassConfirmDoneMsg:msg,EmailValue:email}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}		
}

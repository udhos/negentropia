package handler

import (
	//"io"
	//"os"
	//"fmt"
	"log"
	//"time"
	//"io/ioutil"
	"strconv"
	"net/http"
	//"crypto/sha1"
	"html/template"

	"negentropia/webserv/cfg"
	"negentropia/webserv/store"	
	"negentropia/webserv/session"
)

type PasswordPage struct {
	HomePath		     string
	LoginPath		     string
	LogoutPath		     string
	ResetPassProcessPath string
	ResetPassConfirmPath string

	EmailValue        string
	BadEmailMsg       string
	ResetPassDoneMsg  string	
	
	Account         string
	ShowNavAccount  bool
	ShowNavHome     bool
	ShowNavLogin    bool
	ShowNavLogout   bool	
}

func sendResetPass(w http.ResponseWriter, p PasswordPage) error {
	p.HomePath           = cfg.HomePath()
	p.LoginPath          = cfg.LoginPath()
	p.LogoutPath         = cfg.LogoutPath()
	p.ResetPassProcessPath = cfg.ResetPassProcessPath()
	p.ResetPassConfirmPath = cfg.ResetPassConfirmPath()	
	
	// FIXME: we're loading template every time
    t, err := template.ParseFiles(TemplatePath("base.tpl"), TemplatePath("password.tpl"))
	if err != nil {
		return err
	}
	if err = t.Execute(w, p); err != nil {
		return err
	}
	return nil
}

func ResetPass(w http.ResponseWriter, r *http.Request, s *session.Session) {
	log.Printf("handler.ResetPass url=%s", r.URL.Path)
	
	account := accountLabel(s)
	
	email := formatEmail(r.FormValue("Email"))
	
	if err := sendResetPass(w, PasswordPage{Account:account,ShowNavAccount:true,ShowNavHome:true,EmailValue:email}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func newResetPassConfirmationId() string {
	return "r:" + strconv.FormatInt(store.Incr("i:resetPassConfirmationIdGenerator"), 10)
}

func ResetPassProcess(w http.ResponseWriter, r *http.Request, s *session.Session) {
	log.Printf("handler.ResetPassProcess url=%s", r.URL.Path)

	account := accountLabel(s)
	
	email := formatEmail(r.FormValue(FORM_VAR_EMAIL))

	if email == "" {
		msg := "Please enter email address."
		if err := sendResetPass(w, PasswordPage{Account:account,ShowNavAccount:true,ShowNavHome:true,BadEmailMsg:msg}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	
	if !store.Exists(email) {
		msg := "The address " + email + " does not exist."
		if err := sendResetPass(w, PasswordPage{Account:account,ShowNavAccount:true,ShowNavHome:true,BadEmailMsg:msg,EmailValue:email}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	
	if store.FieldExists(email, "unconfirmed") {
		msg := "The address " + email + " has not been confirmed."
		if err := sendResetPass(w, PasswordPage{Account:account,ShowNavAccount:true,ShowNavHome:true,BadEmailMsg:msg,EmailValue:email}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	
	confId := newResetPassConfirmationId()
	store.Set(confId, email)
	store.Expire(confId, unconfirmedExpire) // Expire confirmation id after 2 days

	log.Printf("handler.ResetPassProcess: FIXME WRITEME sendMail")
	//go sendMail(email, confId)	
	
	msg := "The validation code for password recovery has been sent to " + email + ". Please check your email to change the password."
	if err := sendResetPass(w, PasswordPage{Account:account,ShowNavAccount:true,ShowNavHome:true,ResetPassDoneMsg:msg,EmailValue:email}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}	
}

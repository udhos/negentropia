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

type SignupPage struct {
	HomePath		  string
	LoginPath		  string
	LogoutPath		  string
	SignupProcessPath string
	ConfirmPath       string	

	EmailValue        string
	BadEmailMsg       string
	BadPasswdMsg      string
	BadConfirmMsg     string
	BadSignupMsg      string
	SignupDoneMsg     string	
	
	Account         string
	ShowNavAccount  bool
	ShowNavHome     bool
	ShowNavLogin    bool
	ShowNavLogout   bool	
}

var (
	unconfirmedExpire int64 = 2 * 86400 // expire unconfirmed email after 2 days
)

func sendSignup(w http.ResponseWriter, p SignupPage) error {
	p.HomePath          = cfg.HomePath()
	p.LoginPath         = cfg.LoginPath()
	p.LogoutPath        = cfg.LogoutPath()
	p.SignupProcessPath = cfg.SignupProcessPath()
	p.ConfirmPath       = cfg.ConfirmPath()
	
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

func newConfirmationId() string {
	return "c:" + strconv.FormatInt(store.Incr("i:confirmationIdGenerator"), 10)
}

func SignupProcess(w http.ResponseWriter, r *http.Request, s *session.Session) {
	path := r.URL.Path
	log.Printf("handler.SignupProcess url=%s", path)

	account := accountLabel(s)
	
	name := r.FormValue("Name")
	email := formatEmail(r.FormValue("Email"))
	password := r.FormValue("Passwd")
	confirm := r.FormValue("Confirm")
	
	if email == "" {
		msg := "Please enter email address."
		if err := sendSignup(w, SignupPage{Account:account,ShowNavAccount:true,ShowNavHome:true,BadEmailMsg:msg}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if store.Exists(email) && !store.FieldExists(email, "unconfirmed") {
		msg := "The address " + email + " is already taken."
		if err := sendSignup(w, SignupPage{Account:account,ShowNavAccount:true,ShowNavHome:true,BadEmailMsg:msg,EmailValue:email}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	
	if password != confirm {
		msg := "Passwords don't match."
		if err := sendSignup(w, SignupPage{Account:account,ShowNavAccount:true,ShowNavHome:true,BadConfirmMsg:msg,EmailValue:email}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	
	confId := newConfirmationId()
	store.Set(confId, email)
	store.Expire(confId, unconfirmedExpire) // Expire confirmation id after 2 days

	store.SetField(email, "name", name)
	store.SetField(email, "password-sha1-hex", passDigest(password))
	store.SetField(email, "unconfirmed", confId) // Save confirmation id here only for informational purpose
	store.Expire(email, unconfirmedExpire) // Expire unconfirmed email after 2 days
	
	msg := "The new account has been created, and a confirmation email has been sent to " + email + ". Please check your email to enable the account."
	if err := sendSignup(w, SignupPage{Account:account,ShowNavAccount:true,ShowNavHome:true,SignupDoneMsg:msg,EmailValue:email}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

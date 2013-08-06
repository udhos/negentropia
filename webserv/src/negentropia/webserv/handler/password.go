package handler

import (
	//"io"
	//"os"
	"bytes"
	"fmt"
	"log"
	//"time"
	//"errors"
	//"io/ioutil"
	"net/http"
	"strconv"
	//"crypto/sha1"
	"html/template"

	"negentropia/webserv/cfg"
	"negentropia/webserv/session"
	"negentropia/webserv/store"
	"negentropia/webserv/util"
)

type PasswordPage struct {
	HomePath             string
	SignupPath           string
	LoginPath            string
	LogoutPath           string
	ResetPassProcessPath string
	ResetPassConfirmPath string

	EmailValue       string
	BadEmailMsg      string
	ResetPassDoneMsg string

	Account        string
	ShowNavAccount bool
	ShowNavHome    bool
	ShowNavSignup  bool
	ShowNavLogin   bool
	ShowNavLogout  bool
}

type ResetPassEmail struct {
	Email      string
	ConfirmId  string
	ClickURL   string
	ConfirmURL string
}

func loadEmailTemplate(filename string, e ResetPassEmail) (string, error) {
	// FIXME: we're loading template every time
	t, err := template.ParseFiles(TemplatePath(filename))
	if err != nil {
		return "", err
	}
	buf := &bytes.Buffer{}
	if err = t.Execute(buf, e); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func sendResetPass(w http.ResponseWriter, p PasswordPage) error {
	p.HomePath = cfg.HomePath()
	p.SignupPath = cfg.SignupPath()
	p.LoginPath = cfg.LoginPath()
	p.LogoutPath = cfg.LogoutPath()
	p.ResetPassProcessPath = cfg.ResetPassProcessPath()
	p.ResetPassConfirmPath = cfg.ResetPassConfirmPath()

	p.ShowNavSignup = true

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

	if err := sendResetPass(w, PasswordPage{Account: account, ShowNavAccount: true, ShowNavHome: true, EmailValue: email}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func newResetPassConfirmationId() string {
	return "r:" + strconv.FormatInt(store.Incr("i:resetPassConfirmationIdGenerator"), 10) + util.RandomSuffix()
}

func sendResetPassEmail(email, confId string) {
	confURL := cfg.ResetPassConfirmURL()
	clickURL := fmt.Sprintf("%s?%s=%s&%s=%s", confURL, FORM_VAR_EMAIL, email, FORM_VAR_CONFIRM_ID, confId)
	emailInfo := ResetPassEmail{email, confId, clickURL, confURL}
	var err error
	var msgPlain string
	if msgPlain, err = loadEmailTemplate("resetPassEmailPlain.tpl", emailInfo); err != nil {
		log.Printf("handler.sendResetPassEmail: failure loading PLAIN template for password recovery email: %s", err)
		return
	}
	var msgHtml string
	if msgHtml, err = loadEmailTemplate("resetPassEmailHtml.tpl", emailInfo); err != nil {
		log.Printf("handler.sendResetPassEmail: failure loading HTML template for password recovery email: %s", err)
		return
	}

	sendSmtp(cfg.SmtpAuthUser, cfg.SmtpAuthPass, cfg.SmtpAuthServer, cfg.SmtpHostPort, cfg.SmtpAuthUser, email, "Negentropia password recovery", msgPlain, msgHtml)
}

func ResetPassProcess(w http.ResponseWriter, r *http.Request, s *session.Session) {
	log.Printf("handler.ResetPassProcess url=%s", r.URL.Path)

	account := accountLabel(s)

	email := formatEmail(r.FormValue(FORM_VAR_EMAIL))

	if email == "" {
		msg := "Please enter email address."
		if err := sendResetPass(w, PasswordPage{Account: account, ShowNavAccount: true, ShowNavHome: true, BadEmailMsg: msg}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if !store.Exists(email) {
		msg := "The address " + email + " does not exist."
		if err := sendResetPass(w, PasswordPage{Account: account, ShowNavAccount: true, ShowNavHome: true, BadEmailMsg: msg, EmailValue: email}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if store.FieldExists(email, "unconfirmed") {
		msg := "The address " + email + " has not been confirmed."
		if err := sendResetPass(w, PasswordPage{Account: account, ShowNavAccount: true, ShowNavHome: true, BadEmailMsg: msg, EmailValue: email}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	confId := newResetPassConfirmationId()
	store.Set(confId, email)
	store.Expire(confId, unconfirmedExpire) // Expire confirmation id after 2 days

	go sendResetPassEmail(email, confId)

	msg := "The validation code for password recovery has been sent to " + email + ". Please check your email to change the password."
	if err := sendResetPass(w, PasswordPage{Account: account, ShowNavAccount: true, ShowNavHome: true, ResetPassDoneMsg: msg, EmailValue: email}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

package handler

import (
	//"io"
	//"os"
	"fmt"
	"log"
	//"time"
	//"io/ioutil"
	"net/http"
	"net/smtp"
	"strconv"
	"strings"
	//"crypto/sha1"
	"html/template"

	"negentropia/webserv/cfg"
	"negentropia/webserv/session"
	"negentropia/webserv/store"
	"negentropia/webserv/util"
)

type SignupPage struct {
	HomePath          string
	SignupPath        string
	LoginPath         string
	LogoutPath        string
	SignupProcessPath string
	ConfirmPath       string

	EmailValue    string
	BadEmailMsg   string
	BadPasswdMsg  string
	BadConfirmMsg string
	BadSignupMsg  string
	SignupDoneMsg string

	Account        string
	ShowNavAccount bool
	ShowNavHome    bool
	ShowNavSignup  bool
	ShowNavLogin   bool
	ShowNavLogout  bool
}

var (
	unconfirmedExpire int64 = 2 * 86400 // expire unconfirmed email after 2 days
)

func sendSignup(w http.ResponseWriter, p SignupPage) error {
	p.HomePath = cfg.HomePath()
	p.SignupPath = cfg.SignupPath()
	p.LoginPath = cfg.LoginPath()
	p.LogoutPath = cfg.LogoutPath()
	p.SignupProcessPath = cfg.SignupProcessPath()
	p.ConfirmPath = cfg.ConfirmPath()

	p.ShowNavSignup = false

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

	email := formatEmail(r.FormValue("Email"))

	if err := sendSignup(w, SignupPage{Account: account, ShowNavAccount: true, ShowNavHome: true, EmailValue: email}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func newConfirmationId() string {
	return "c:" + strconv.FormatInt(store.Incr("i:confirmationIdGenerator"), 10) + util.RandomSuffix()
}

func sendSmtp(authUser, authPass, authServer, smtpHostPort, sender, recipient, subject, msgPlain, msgHtml string) {

	var pass string
	if authPass != "" {
		pass = "<hidden>"
	}

	log.Printf("sendSmtp: sub=[%s] auth=[%s] pass=[%s] sender=[%s] recipient=[%s] sending...", subject, authUser, pass, sender, recipient)

	auth := smtp.PlainAuth(
		"",
		authUser,
		authPass,
		authServer,
	)

	mime := "MIME-version: 1.0;\r\n"
	boundary := "20cf307d051035ce0404d47a8e9b"
	sub := fmt.Sprintf("Subject: %s\r\n", subject)
	from := fmt.Sprintf("From: <%s>\r\n", sender)
	to := fmt.Sprintf("To: <%s>\r\n", recipient)
	bodyTemplate := "Content-Type: multipart/alternative; boundary=" + boundary + "\r\n" +
		"\r\n" +
		"--" + boundary + "\r\n" +
		"Content-Type: text/plain; charset=ISO-8859-1\r\n" +
		"\r\n" +
		"%s" +
		"\r\n" +
		"--" + boundary + "\r\n" +
		mime +
		"Content-Type: text/html; charset=ISO-8859-1\r\n" +
		"\r\n" +
		"%s" +
		"\r\n" +
		"--" + boundary + "--" +
		"\r\n"

	body := fmt.Sprintf(bodyTemplate, msgPlain, msgHtml)

	err := smtp.SendMail(
		smtpHostPort,
		auth,
		sender,
		[]string{recipient},
		[]byte(sub+from+to+body),
	)
	var result string
	if err != nil {
		log.Printf("sendSmtp: failure: %q", strings.Split(err.Error(), "\n"))
		result = "FAIL"
	} else {
		result = "SENT"
	}

	log.Printf("sendSmtp: sub=[%s] auth=[%s] pass=[%s] sender=[%s] recipient=[%s] %s", subject, authUser, pass, sender, recipient, result)
}

func sendSignupMail(email, confId string) {
	confURL := cfg.ConfirmURL()
	clickURL := fmt.Sprintf("%s?%s=%s&%s=%s", cfg.ConfirmProcessURL(), FORM_VAR_EMAIL, email, FORM_VAR_CONFIRM_ID, confId)
	emailInfo := EmailFields{email, confId, clickURL, confURL}

	var err error
	var msgPlain string
	if msgPlain, err = loadEmailTemplate("signupEmailPlain.tpl", emailInfo); err != nil {
		log.Printf("handler.sendSignupEmail: failure loading PLAIN template for signup confirmation email: %s", err)
		return
	}

	var msgHtml string
	if msgHtml, err = loadEmailTemplate("signupEmailHtml.tpl", emailInfo); err != nil {
		log.Printf("handler.sendSignupEmail: failure loading HTML template for signup confirmation email: %s", err)
		return
	}

	sendSmtp(cfg.SmtpAuthUser, cfg.SmtpAuthPass, cfg.SmtpAuthServer, cfg.SmtpHostPort, cfg.SmtpAuthUser, email, "Negentropia signup confirmation", msgPlain, msgHtml)
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
		if err := sendSignup(w, SignupPage{Account: account, ShowNavAccount: true, ShowNavHome: true, BadEmailMsg: msg}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if store.Exists(email) && !store.FieldExists(email, "unconfirmed") {
		msg := "The address " + email + " is already taken."
		if err := sendSignup(w, SignupPage{Account: account, ShowNavAccount: true, ShowNavHome: true, BadEmailMsg: msg, EmailValue: email}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if password != confirm {
		msg := "Passwords don't match."
		if err := sendSignup(w, SignupPage{Account: account, ShowNavAccount: true, ShowNavHome: true, BadConfirmMsg: msg, EmailValue: email}); err != nil {
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
	store.Expire(email, unconfirmedExpire)       // Expire unconfirmed email after 2 days

	go sendSignupMail(email, confId)

	msg := "The new account has been created, and a confirmation email has been sent to " + email + ". Please check your email to enable the account."
	if err := sendSignup(w, SignupPage{Account: account, ShowNavAccount: true, ShowNavHome: true, SignupDoneMsg: msg, EmailValue: email}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

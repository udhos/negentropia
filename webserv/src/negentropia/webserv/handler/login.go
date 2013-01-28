package handler

import (
	"io"
	//"os"
	"fmt"
	"log"
	//"time"
	//"io/ioutil"
	"net/http"
	"crypto/sha1"
	"html/template"

	//"code.google.com/p/goauth2/oauth" // facebook broken
	//"github.com/robfig/goauth2/oauth" // google broken
	//"code.google.com/r/jasonmcvetta-goauth2" // go get broken
	"github.com/HairyMezican/goauth2/oauth"

	"negentropia/webserv/cfg"
	"negentropia/webserv/store"	
	"negentropia/webserv/session"
)

type Page struct {
	HomePath		string
	LoginPath		string
	LogoutPath		string
	LoginAuthPath	string
	SignupPath		string
	EmailValue      string
	
	PasswdBadAuth   string
	GoogleAuthMsg   string
	FacebookAuthMsg string	

	Account         string
	ShowNavAccount  bool
	ShowNavHome     bool
	ShowNavLogin    bool
	ShowNavLogout   bool	
}

func sendLogin(w http.ResponseWriter, p Page) error {
	p.HomePath      = cfg.HomePath()
	p.LoginPath     = cfg.LoginPath()
	p.LogoutPath    = cfg.LogoutPath()
	p.LoginAuthPath = cfg.LoginAuthPath()
	p.SignupPath    = cfg.SignupPath()
	
	// FIXME: we're loading template every time
    t, err := template.ParseFiles(TemplatePath("base.tpl"), TemplatePath("login.tpl"))
	if err != nil {
		return err
	}
	if err = t.Execute(w, p); err != nil {
		return err
	}
	return nil
}

func Login(w http.ResponseWriter, r *http.Request, s *session.Session) {
	path := r.URL.Path
	log.Printf("handler.Login url=%s", path)
	
	account := accountLabel(s)
	
	if err := sendLogin(w, Page{Account:account,ShowNavAccount:true,ShowNavHome:true}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func passDigest(pass string) string {
	h := sha1.New()
	io.WriteString(h, pass)
	return fmt.Sprintf("%0x", h.Sum(nil))
}

func passwordAuth(email string, pass string) bool {
	passHash := passDigest(pass)

	//dbHash := session.RedisQueryField(email, "password-sha1-hex")
	dbHash := store.QueryField(email, "password-sha1-hex")
	
	log.Printf("login.auth: email=%s auth=%s provided=%s", email, dbHash, passHash)	
	
	return passHash == dbHash
}

func googleOauth2Config(host, port string) *oauth.Config {

	redirect := "http://" + host + port + cfg.GoogleCallbackPath()
	
	log.Printf("handler.googleOauth2Config: redirect=%s", redirect)

	return &oauth.Config{
		ClientId:     *GoogleId,
		ClientSecret: *GoogleSecret,
		Scope:        "https://www.googleapis.com/auth/userinfo.profile https://www.googleapis.com/auth/userinfo.email",
		AuthURL:      "https://accounts.google.com/o/oauth2/auth",
		TokenURL:     "https://accounts.google.com/o/oauth2/token",
		RedirectURL:  redirect,
	}
}

func facebookOauth2Config(host, port string) *oauth.Config {

	redirect := "http://" + host + port + cfg.FacebookCallbackPath()
	
	log.Printf("handler.facebookOauth2Config: redirect=%s", redirect)

	return &oauth.Config{
		ClientId:     *FacebookId,
		ClientSecret: *FacebookSecret,
		Scope:        "email",
		AuthURL:      "https://www.facebook.com/dialog/oauth",
		TokenURL:     "https://graph.facebook.com/oauth/access_token",
		RedirectURL:  redirect,
	}
}

func googleOauth2(w http.ResponseWriter, r *http.Request) {
	log.Printf("handler.LoginAuth: google url=%s", r.URL)

	config := googleOauth2Config(RedirectHost, RedirectPort)
	
	// Step one, get an authorization code from the data provider.
	
	url := config.AuthCodeURL("")
	
	http.Redirect(w, r, url, http.StatusFound)
	
	// See next steps under googleCallback handler
}

func facebookOauth2(w http.ResponseWriter, r *http.Request) {
	log.Printf("handler.LoginAuth: facebook url=%s", r.URL)

	config := facebookOauth2Config(RedirectHost, RedirectPort)
	
	// Step one, get an authorization code from the data provider.
	
	url := config.AuthCodeURL("")
	
	http.Redirect(w, r, url, http.StatusFound)
	
	// See next steps under facebookCallback handler
}

func LoginAuth(w http.ResponseWriter, r *http.Request, s *session.Session) {
	path := r.URL.Path

	account := accountLabel(s)

	login := r.FormValue("LoginButton");
	//email := r.FormValue("Email");
	//password := r.FormValue("Passwd");

	google := r.FormValue("GoogleButton");
	facebook := r.FormValue("FacebookButton");	
	
	//debug := "handler.LoginAuth url=%s email=%s pass=%s login=%s google=%s facebook=%s"
	//log.Printf(debug, path, email, password, login, google, facebook)
	//fmt.Fprintf(w, debug, path, email, password, login, google, facebook)
	
	switch {
		case login != "":
			email := formatEmail(r.FormValue("Email"))		
			if store.FieldExists(email, "unconfirmed") {
				msg := "The address " + email + " has not been confirmed."
				if err := sendLogin(w, Page{Account:account,ShowNavAccount:true,ShowNavHome:true,PasswdBadAuth:msg,EmailValue:email}); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				return
			}
				
			password := r.FormValue("Passwd")
			if passwordAuth(email, password) {
				// auth ok

				if s != nil {
					session.Delete(w, s)
				}
				name := session.RedisQueryField(email, "name")
				s = session.Set(w, session.AUTH_PROV_PASSWORD, email, name, email)
				if s == nil {
					log.Printf("login.LoginAuth url=%s could not establish session", path)	
					http.Error(w, "login.LoginAuth could not establish session", http.StatusInternalServerError)
					return
				}
				
				http.Redirect(w, r, cfg.HomePath(), http.StatusFound)
			} else {
				// bad auth
				if err := sendLogin(w, Page{Account:account,ShowNavAccount:true,ShowNavHome:true,PasswdBadAuth:"Invalid email/password. Please try again.",EmailValue:email}); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			}
		case google != "":
			googleOauth2(w, r)
		case facebook != "":
			facebookOauth2(w, r)		
		default:
			log.Printf("handler.LoginAuth: missing button")
			http.Redirect(w, r, cfg.LoginPath(), http.StatusFound)
	}
}


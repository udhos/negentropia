package handler

import (
	//"os"
	"fmt"
	"log"
	//"time"
	//"io/ioutil"
	"net/http"
	"html/template"

	"code.google.com/p/goauth2/oauth"
	
	"negentropia/webserv/session"
)

type Page struct {
	PasswdBadAuth string
	GoogleAuthMsg string

	Account       string
	ShowNavAccount bool
	ShowNavHome    bool
	ShowNavLogin   bool
	ShowNavLogout  bool	
}

func sendLogin(w http.ResponseWriter, p Page) error {
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

func auth(email string, auth string) bool {
	return false
}

func googleOauth2Config() *oauth.Config {
	return &oauth.Config{
		ClientId:     *GoogleId,
		ClientSecret: *GoogleSecret,
		Scope:        "https://www.googleapis.com/auth/userinfo.profile https://www.googleapis.com/auth/userinfo.email",
		AuthURL:      "https://accounts.google.com/o/oauth2/auth",
		TokenURL:     "https://accounts.google.com/o/oauth2/token",
		RedirectURL:  "http://localhost:8080/n/googleCallback",
	}
}

func googleOauth2(w http.ResponseWriter, r *http.Request) {
	log.Printf("handler.LoginAuth: google")
	
	config := googleOauth2Config()
	
	// Step one, get an authorization code from the data provider.
	
	url := config.AuthCodeURL("")
	
	http.Redirect(w, r, url, http.StatusFound)
	
	// See next steps under googleCallback handler
}

func LoginAuth(w http.ResponseWriter, r *http.Request, s *session.Session) {
	//path := r.URL.Path

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
			email := r.FormValue("Email");
			password := r.FormValue("Passwd");
			if auth(email, password) {
				// auth ok
				http.Redirect(w, r, "/n/", http.StatusFound)
			} else {
				// bad auth
				if err := sendLogin(w, Page{Account:account,ShowNavAccount:true,ShowNavHome:true,PasswdBadAuth: "Invalid email/password. Please try again."}); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			}
		case google != "":
			//fmt.Fprintf(w, "handler.LoginAuth: google")
			googleOauth2(w, r)
		case facebook != "":
			log.Printf("handler.LoginAuth: facebook")
			fmt.Fprintf(w, "handler.LoginAuth: facebook")
		default:
			log.Printf("handler.LoginAuth: missing button")
			http.Redirect(w, r, "/n/login", http.StatusFound)
	}
}


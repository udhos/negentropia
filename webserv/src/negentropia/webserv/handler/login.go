package handler

import (
	//"os"
	"fmt"
	"log"
	//"time"
	//"io/ioutil"
	"net/http"
	"html/template"
)

type Page struct {
    //Title string
	LoginBadAuth string
}

func sendLogin(w http.ResponseWriter, p Page) {
	// FIXME: we're loading template every time
    t, _ := template.ParseFiles(TemplatePath("login.tpl"))
    t.Execute(w, p)
}

func Login(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	
	log.Printf("handler.login url=%s", path)

	sendLogin(w, Page{})
}

func auth(email string, auth string) bool {
	return false
}

func LoginAuth(w http.ResponseWriter, r *http.Request) {
	//path := r.URL.Path

	login := r.FormValue("LoginButton");
	//email := r.FormValue("Email");
	//password := r.FormValue("Passwd");

	google := r.FormValue("GoogleButton");
	facebook := r.FormValue("FacebookButton");	
	
	//debug := "handler.loginAuth url=%s email=%s pass=%s login=%s google=%s facebook=%s"
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
				sendLogin(w, Page{"Invalid email/password. Please try again."})
			}
		case google != "":
			fmt.Fprintf(w, "handler.loginAuth: google")
		case facebook != "":
			fmt.Fprintf(w, "handler.loginAuth: facebook")
		default:
			fmt.Fprintf(w, "handler.loginAuth: missing button")
	}
}


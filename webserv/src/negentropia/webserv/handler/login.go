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
    Title string
}

func Login(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	
	log.Printf("handler.login url=%s", path)
	
	p := Page{"Login Test"}
	
    t, _ := template.ParseFiles(TemplatePath("login.tpl"))
    t.Execute(w, p)	
}

func LoginAuth(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	
	log.Printf("handler.loginAuth url=%s", path)
	
	fmt.Fprintf(w, "handler.loginAuth")
}


package handler

import (
	//"os"
	//"fmt"
	"log"
	//"time"
	//"io/ioutil"
	"net/http"
	"html/template"	

	"negentropia/webserv/cfg"	
	"negentropia/webserv/session"
)

type HomePage struct {
	HomePath		string
	LoginPath		string
	LogoutPath		string
	
	Account        string
	
	ShowNavAccount bool
	ShowNavHome    bool
	ShowNavLogin   bool
	ShowNavLogout  bool
}

func sendHome(w http.ResponseWriter, p HomePage) error {
	p.HomePath   = cfg.HomePath()
	p.LoginPath  = cfg.LoginPath()
	p.LogoutPath = cfg.LogoutPath()

	// FIXME: we're loading template every time
    t, err := template.ParseFiles(TemplatePath("base.tpl"), TemplatePath("home.tpl"))
	if err != nil {
		return err
	}
	if err = t.Execute(w, p); err != nil {
		return err
	}
	return nil
}

func Home(w http.ResponseWriter, r *http.Request, s *session.Session) {
	path := r.URL.Path
	
	log.Printf("handler.home url=%s", path)
	
	account := accountLabel(s)
	
	if err := sendHome(w, HomePage{Account:account,ShowNavAccount:true}); err != nil {
		log.Printf("handler.home url=%s %s", path, err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}	
}

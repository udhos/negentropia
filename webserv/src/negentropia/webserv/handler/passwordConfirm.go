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
	//"negentropia/webserv/store"	
	"negentropia/webserv/session"
)

type ResetPassConfirmPage struct {
	HomePath		   string
	LoginPath		   string
	LogoutPath		   string
	ResetPassConfirmProcessPath string
	
	EmailValue        string
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
	
	if err := sendResetPassConfirm(w, ResetPassConfirmPage{Account:account,ShowNavAccount:true,ShowNavHome:true}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ResetPassConfirmProcess(w http.ResponseWriter, r *http.Request, s *session.Session) {
	log.Printf("handler.ResetPassConfirmProcess url=%s", r.URL.Path)

}

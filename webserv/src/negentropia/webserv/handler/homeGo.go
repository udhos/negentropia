package handler

import (
	//"os"
	//"fmt"
	"log"
	//"time"
	//"io/ioutil"
	"html/template"
	"net/http"

	"negentropia/webserv/cfg"
	"negentropia/webserv/session"
	"negentropia/webserv/share"
	"negentropia/webserv/store"
)

type HomeGoPage struct {
	HomePath   string
	HomeJSPath string
	HomeGoPath string
	SignupPath string
	LoginPath  string
	LogoutPath string

	Account string

	ShowNavAccount bool
	ShowNavHome    bool
	ShowNavSignup  bool
	ShowNavLogin   bool
	ShowNavLogout  bool

	Websocket string
}

func sendHomeGo(w http.ResponseWriter, p HomeGoPage) error {
	p.HomePath = cfg.HomePath()
	p.HomeJSPath = cfg.HomeJSPath()
	p.HomeGoPath = cfg.HomeGoPath()
	p.SignupPath = cfg.SignupPath()
	p.LoginPath = cfg.LoginPath()
	p.LogoutPath = cfg.LogoutPath()

	p.ShowNavSignup = true

	p.Websocket = store.Get(share.WORLD_WEBSOCKET)

	// FIXME: we're loading template every time
	t, err := template.ParseFiles(TemplatePath("base.tpl"), TemplatePath("homeGo.tpl"))
	if err != nil {
		return err
	}
	if err = t.Execute(w, p); err != nil {
		return err
	}
	return nil
}

func HomeGo(w http.ResponseWriter, r *http.Request, s *session.Session) {
	path := r.URL.Path

	log.Printf("handler.HomeGo url=%s", path)

	account := accountLabel(s)

	if err := sendHomeGo(w, HomeGoPage{Account: account, ShowNavAccount: true}); err != nil {
		log.Printf("handler.HomeGo url=%s %s", path, err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

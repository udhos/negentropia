package handler

import (
	//"io"
	//"os"
	//"fmt"
	"log"
	//"errors"
	//"time"
	//"io/ioutil"
	"net/http"
	//"encoding/json"
	
	"negentropia/webserv/cfg"
	"negentropia/webserv/session"
)

func Logout(w http.ResponseWriter, r *http.Request, s *session.Session) {
	path := r.URL.Path
	log.Printf("handler.Logout url=%s", path)

	if s != nil {
		session.Delete(w, s)
	}
	
	http.Redirect(w, r, cfg.LoginPath(), http.StatusFound)
}

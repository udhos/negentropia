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
	
	"negentropia/webserv/session"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	
	log.Printf("handler.Logout url=%s", path)

	s := session.Get(r)
	if s != nil {
		session.Delete(w, s)
	}
	
	http.Redirect(w, r, "/n/login", http.StatusFound)
}

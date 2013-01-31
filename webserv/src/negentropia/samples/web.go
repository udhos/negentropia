package main

import (
	"io"
	"fmt"
	"log"
	"net/http"	
)

type StaticHandler struct {
	innerHandler http.Handler
}

func (handler StaticHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("StaticHandler.ServeHTTP url=%s", r.URL.Path)
	handler.innerHandler.ServeHTTP(w, r)
}

func dynamic1(w http.ResponseWriter, r *http.Request) {
	msg := fmt.Sprintf("dynamic1: URL=%s", r.URL.Path)
	log.Printf(msg)
	io.WriteString(w, msg)
	io.WriteString(w, "\n")	
}
func dynamic2(w http.ResponseWriter, r *http.Request) {
	msg := fmt.Sprintf("dynamic2: URL=%s", r.URL.Path)
	log.Printf(msg)
	io.WriteString(w, msg)
	io.WriteString(w, "\n")	
}

func main() {
	
	
	http.Handle("/", StaticHandler{http.StripPrefix("/", http.FileServer(http.Dir("/tmp/www/root/")))})
	http.Handle("/s/", StaticHandler{http.StripPrefix("/s", http.FileServer(http.Dir("/tmp/www/s/")))})
	http.Handle("/s/b/", StaticHandler{http.StripPrefix("/s/b", http.FileServer(http.Dir("/tmp/www/sb/s/b/")))})		
		
	http.HandleFunc("/a/", dynamic1)
	http.HandleFunc("/a/b", dynamic2)

	addr := ":8080"
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Panicf("ListenAndServe: %s: %s", addr, err)
	}		
}

package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

var home string

func root(w http.ResponseWriter, r *http.Request) {
	msg := fmt.Sprintf("root: URL=%s", r.URL.Path)
	log.Printf(msg)

	if r.URL.Path != "/" {
		log.Printf("root: URL=%s refusing to serve", r.URL.Path)
		http.NotFound(w, r)
		return
	}

	io.WriteString(w, msg)
	io.WriteString(w, "\n")
}

func login(w http.ResponseWriter, r *http.Request) {
	msg := fmt.Sprintf("login: URL=%s", r.URL.Path)
	log.Printf(msg)
	io.WriteString(w, msg)
	io.WriteString(w, "\n")
}

func callback(w http.ResponseWriter, r *http.Request) {
	msg := fmt.Sprintf("callback: URL=%s", r.URL.Path)
	log.Printf(msg)
	io.WriteString(w, msg)
	io.WriteString(w, "\n")
}

func main() {

	http.HandleFunc("/", root)
	http.HandleFunc("/login", login)
	http.HandleFunc("/callback", callback)

	addr := "localhost:8080"
	home := "http://" + addr

	log.Printf("serving at %s", home)

	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Panicf("ListenAndServe: %s: %s", addr, err)
	}
}

package main

import (
	//"fmt"
	//"io"
	"log"
	"net/http"
	"os"
)

type StaticHandler struct {
	innerHandler http.Handler
}

func (handler StaticHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("StaticHandler.ServeHTTP url=%s", r.URL.Path)
	handler.innerHandler.ServeHTTP(w, r)
}

/*
func home(w http.ResponseWriter, r *http.Request) {
	msg := fmt.Sprintf("home: URL=%s", r.URL.Path)
	log.Printf(msg)
	homeStr :=
		`<!DOCTYPE html>

<html>
  <head>
    <meta charset="utf-8">
    <title>Launch dartium url</title>
    <link rel="stylesheet" href="launch_dartium_url.css">
  </head>
  <body>
    <h1>Launch dartium url</h1>

    <script type="application/dart" src="/www/launch_dartium_url.dart"></script>
    <script src="/www/packages/browser/dart.js"></script>
  </body>
</html>
`
	io.WriteString(w, homeStr)
}
*/

func notFound(w http.ResponseWriter, r *http.Request) {
	log.Printf("notFound: url=%s", r.URL.Path)
	http.NotFound(w, r) // default not-found handler
}

func main() {

	dir, err := os.Getwd()
	if err != nil {
		log.Panicf("Getwd: %s", err)
	}

	addr := ":8080"

	http.HandleFunc("/", notFound) // trap not-found handler

	/*
		dynamic := "/home"
		log.Printf("serving dynamic path %s", dynamic)
		http.HandleFunc(dynamic, home)
	*/

	path := "/www/"
	log.Printf("serving static directory %s as www path %s", dir, path)
	http.Handle(path, StaticHandler{http.StripPrefix(path, http.FileServer(http.Dir(dir)))})

	log.Printf("serving on port TCP %s", addr)

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Panicf("ListenAndServe: %s: %s", addr, err)
	}
}

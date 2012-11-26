package main

import (
	//"os"
	//"fmt"
	"log"
	//"time"
	//"io/ioutil"
	"net/http"

	"negentropia/webserv/handler"
)

var (
	rootPath     string = "C:\\tmp\\devel\\negentropia\\wwwroot"
	templatePath string = "C:\\tmp\\devel\\negentropia\\template"
)

// Initialize package main
func init() {
	handler.SetRootPath(templatePath)
}

/*
func absPath(path string) string {
	return rootPath + path
}

func static(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	fullPath := absPath(path)

	http.ServeFile(w, r, fullPath)	
	log.Printf("served static url=%s fullPath=%s", path, fullPath)
}
*/

// Wrapper type for Handler
type StaticHandler struct {
	innerHandler http.Handler // save trapped/wrapped Handler
}

func (handler StaticHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	log.Printf("StaticHandler.ServeHTTP url=%s", path)
	handler.innerHandler.ServeHTTP(w, r) // call trapped/wrapped Handler

	/*
		var delay time.Duration = 20 
		log.Printf("blocking for %d secs", delay)
		time.Sleep(delay * time.Second)
	*/
}

func serve(addr string) {
	log.Printf("server starting on " + addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func main() {
	//http.HandleFunc("/", static)
	//http.Handle("/", http.FileServer(http.Dir(rootPath)))

	http.Handle("/", StaticHandler{http.FileServer(http.Dir(rootPath))})
	http.HandleFunc("/n/", handler.Home)
	http.HandleFunc("/n/login", handler.Login)
	http.HandleFunc("/n/loginAuth", handler.LoginAuth)
	http.HandleFunc("/n/callback", handler.Callback)

	go serve(":8080")
	serve(":8000")
}

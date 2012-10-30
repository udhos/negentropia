package main

import (
	//"os"
	//"fmt"
	"log"
	"time"
	//"io/ioutil"
	"net/http"
)

var (
	rootPath string = "../wwwroot"
)

func absPath(path string) string {
	return rootPath + path
}

func static(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	fullPath := absPath(path)
	//log.Printf("static url=%s fullPath=%s", path, fullPath)

	http.ServeFile(w, r, fullPath)	
	log.Printf("served static url=%s fullPath=%s", path, fullPath)

	return
	var delay time.Duration = 20
	log.Printf("sleeping %d secs", delay)
	time.Sleep(delay * time.Second)
}

func serve(addr string) {
	log.Printf("server starting on " + addr)	
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func main() {
	http.HandleFunc("/", static)
	go serve(":8080")
	serve(":8000")
}

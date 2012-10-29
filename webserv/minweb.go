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
	rootPath string = "c:\\wwwroot"
)

func absPath(path string) string {
	return rootPath + path
}

func handler(w http.ResponseWriter, r *http.Request) {

	path := r.URL.Path

	log.Printf("handler path=%s\n", path)	
	
	fullPath := absPath(path)
	
	log.Printf("handler url=%s fullPath=%s\n", path, fullPath)

	http.ServeFile(w, r, fullPath)	

	var delay time.Duration = 20
	log.Printf("served url=%s fullPath=%s sleeping %d secs", path, fullPath, delay)
	time.Sleep(delay * time.Second)
}

func main() {
	log.Printf("server starting\n")	

	http.HandleFunc("/", handler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

package handler

import (
	//"os"
	//"fmt"
	"log"
	//"time"
	//"io/ioutil"
	"net/http"
)

func Callback(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	
	log.Printf("handler.callback url=%s", path)
}

package handler

import (
	//"os"
	//"fmt"
	"log"
	//"time"
	//"io/ioutil"
	"net/http"
)

func Home(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	
	log.Printf("handler.home url=%s", path)
}

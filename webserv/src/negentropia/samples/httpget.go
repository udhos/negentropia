package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	resp, err := http.Get("http://localhost:8080/obj/airship.obj")
	if err != nil {
		log.Fatalf("http get error: %v", err)
	}
	defer resp.Body.Close()

	var info []byte
	info, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("read body error: %v", err)
	}

	fmt.Printf("info=[%v]", string(info))
}

package main

import (
	//"fmt"
	"log"
	"net/http"
	"code.google.com/p/go.net/websocket"
)

func Echo(ws *websocket.Conn) {
	log.Println("Echoing")
	
	sid := "sid"
	cook , err := ws.Request().Cookie(sid)
	if err != nil {
		log.Printf("cookie '%s' NOT FOUND: %s", sid, err)
	} else {
		log.Printf("cookie FOUND: '%s'=[%s]", sid, cook.Value)
	}

	for n := 0; n < 10; n++ {
		msg := "Hello  " + string(n+48)
		log.Println("Sending to client: " + msg)
		err := websocket.Message.Send(ws, msg)
		if err != nil {
			log.Println("Echo: Can't send: %s", err)
			break
		}

		var reply string
		err = websocket.Message.Receive(ws, &reply)
		if err != nil {
			log.Println("Echo: Can't receive: %s", err)
			break
		}
		log.Println("Echo: Received back from client: " + reply)
	}
}

func serve(addr string) {
	if addr == "" {
		log.Printf("server starting on :http (empty address)")
	} else {
		log.Printf("server starting on " + addr)
	}

	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Panicf("ListenAndServe: %s: %s", addr, err)
	}
}

func main() {
	log.Printf("world booting")
	
	http.Handle("/", websocket.Handler(Echo))

	log.Printf("world boot complete")
	
	serve("127.0.0.2:8000")
}

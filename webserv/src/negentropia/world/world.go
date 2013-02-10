package main

import (
	"os"
	//"fmt"
	"log"
	"flag"
	"net/http"
	"code.google.com/p/go.net/websocket"
	
	"negentropia/webserv/session"
	"negentropia/webserv/store"
	"negentropia/webserv/configflag"
)

const (
	CM_CODE_FATAL = 0
	CM_CODE_INFO  = 1
	CM_CODE_AUTH  = 2
)

type ClientMsg struct {
	Code	int
	Data	string
}

var (
	configFlags  *flag.FlagSet = flag.NewFlagSet("config flags", flag.ExitOnError)
	configFile   string
	listenAddr   string
	redisAddr    string
)

// Initialize package main
func init() {
	configFlags.StringVar(&configFile, "config", "", "load config flags from this file")
	configFlags.StringVar(&listenAddr, "listenOn", "127.0.0.2:8000", "listen address [addr]:port")
	configFlags.StringVar(&redisAddr, "redisAddr", "localhost:6379", "redis server address")
}

func Dispatch(ws *websocket.Conn) {

	defer ws.Close()
	
	var sid string
	
	var msg ClientMsg
	err := websocket.JSON.Receive(ws, &msg)
	if err != nil {
		log.Printf("Dispatch: Auth: failure: %s", err)
		return
	}
	if (msg.Code != CM_CODE_AUTH) {
		log.Printf("Dispatch: Auth: non-auth code: %d", msg.Code)
		websocket.JSON.Send(ws, ClientMsg{CM_CODE_FATAL, "auth required"})
		return
	}
	sid = msg.Data
	log.Printf("Dispatch: Auth: sid=%s", sid)
	session := session.Load(sid)
	if (session == nil) {
		log.Printf("Dispatch: Auth: invalid session id sid=%s", sid)
		websocket.JSON.Send(ws, ClientMsg{CM_CODE_FATAL, "bad auth"})
		return
	}

	websocket.JSON.Send(ws, ClientMsg{CM_CODE_INFO, "welcome " + session.ProfileEmail})

	log.Printf("Dispatch: Entering receive loop")
	for {
		err = websocket.JSON.Receive(ws, &msg)
		if err != nil {
			log.Printf("Dispatch: Receive loop: failure: %s", err)
			break
		}
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

	// Parse flags from command-line
	configFlags.Parse(os.Args[1:])

	// Parse flags from file
	if configFile != "" {
		err := configflag.Load(configFlags, configFile)
		if err != nil {
			log.Printf("failure loading config flags: %s", err)
			return
		}
	}
	
	store.Init(redisAddr)
	
	http.Handle("/", websocket.Handler(Dispatch))

	log.Printf("world boot complete")
	
	serve(listenAddr)
}

package main

import (
	"io"
	"os"
	//"fmt"
	"log"
	"flag"
	"net/http"
	
	"code.google.com/p/go.net/websocket"
	
	"negentropia/webserv/session"
	"negentropia/webserv/store"
	"negentropia/webserv/configflag"
	"negentropia/webserv/util"
	"negentropia/webserv/share"
	"negentropia/world/server"
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
	configFlags  	*flag.FlagSet = flag.NewFlagSet("config flags", flag.ExitOnError)
	configList   	configflag.FileList
	listenAddr   	string
	redisAddr    	string
	websocketHost	string
)

// Initialize package main
func init() {
	configFlags.Var(&configList, "config", "load config flags from this file")
	configFlags.StringVar(&listenAddr, "listenOn", "127.0.0.2:8000", "websocket listen address [addr]:port")
	configFlags.StringVar(&redisAddr, "redisAddr", "localhost:6379", "redis server address")
	configFlags.StringVar(&websocketHost, "websocketHost", "127.0.0.2", "host part of websocket uri: ws://host:port/path")
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
		log.Printf("Dispatch: Auth: sid=%s: invalid session id", sid)
		websocket.JSON.Send(ws, ClientMsg{CM_CODE_FATAL, "bad auth"})
		return
	}
		
	server.PlayerAdd(session.SessionId, session.ProfileEmail, ws)

	websocket.JSON.Send(ws, ClientMsg{CM_CODE_INFO, "welcome " + session.ProfileEmail})
	
	log.Printf("Dispatch: Entering receive loop: sid=%s %s", sid, session.ProfileEmail)
	for {
		err = websocket.JSON.Receive(ws, &msg)
		if err == io.EOF {
			log.Printf("Recv: %s %s: disconnected", sid, session.ProfileEmail)
			break
		}
		if err != nil {
			log.Printf("Recv: %s %s: failure: %s", sid, session.ProfileEmail, err)
			break
		}
	}
	
	server.PlayerDel(session.ProfileEmail)
}

func serve(addr string) {
	if addr == "" {
		log.Printf("server starting on :http (empty address)")
	} else {
		log.Printf("server starting on " + addr)
	}
	
	wsAddr := "ws://" + websocketHost + util.GetPort(addr)
	log.Printf("saving websocket address: %s=%s", share.WORLD_WEBSOCKET, wsAddr)
	store.Set(share.WORLD_WEBSOCKET, wsAddr)

	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Printf("ListenAndServe: %s: %s", addr, err)
	}
}

func main() {
	log.Printf("world booting")

	// Parse flags from command-line
	configFlags.Parse(os.Args[1:])

	// Parse flags from files
	log.Printf("config files: %d", len(configList))
	if err := configflag.Load(configFlags, configList); err != nil {
		log.Printf("failure loading config flags: %s", err)
		return
	}
	
	store.Init(redisAddr)
	
	http.Handle("/", websocket.Handler(Dispatch))

	log.Printf("world boot complete")
	
	serve(listenAddr)
}

package main

import (
	//"io"
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

func auth(ws *websocket.Conn) *server.Player {
	
	var msg server.ClientMsg
	err := websocket.JSON.Receive(ws, &msg)
	if err != nil {
		log.Printf("auth: failure: %s", err)
		return nil
	}
	
	if (msg.Code != server.CM_CODE_AUTH) {
		log.Printf("auth: non-auth code: %d", msg.Code)
		websocket.JSON.Send(ws, server.ClientMsg{server.CM_CODE_FATAL, "auth required"})
		return nil
	}
	
	sid := msg.Data
	log.Printf("auth: sid=%s", sid)
	session := session.Load(sid)
	if (session == nil) {
		log.Printf("Dispatch: Auth: sid=%s: invalid session id", sid)
		websocket.JSON.Send(ws, server.ClientMsg{server.CM_CODE_FATAL, "bad auth"})
		return nil
	}
	
	return &server.Player{sid, session.ProfileEmail, ws, make(chan *server.ClientMsg)}
}

func sender(p *server.Player) {
	
	for msg := range p.SendToPlayer {
		err := websocket.JSON.Send(p.Websocket, msg)
		if err != nil {
			log.Printf("sender: %s %s: failure: %s", p.Sid, p.Email, err)
			break
		}
		//log.Printf("sender: %s %s: %q", p.Sid, p.Email, msg)
	}
	
	p.Websocket.Close()
}

func receiver(p *server.Player) {
	for {
		msg := new(server.ClientMsg)
		err := websocket.JSON.Receive(p.Websocket, msg)
		if err != nil {
			log.Printf("receiver: %s %s: failure: %s", p.Sid, p.Email, err)
			break
		}
		//log.Printf("receiver: %s %s: %q", p.Sid, p.Email, msg)
		server.InputCh <- &server.PlayerMsg{p, msg}
	}
	
	p.Websocket.Close()
}

func dispatch(ws *websocket.Conn) {

	defer ws.Close()
	
	newPlayer := auth(ws)
	if newPlayer == nil {
		return
	}

	websocket.JSON.Send(ws, server.ClientMsg{server.CM_CODE_INFO, "welcome " + newPlayer.Email})
	
	server.PlayerAddCh <- newPlayer
	defer func() { server.PlayerDelCh <- newPlayer }()
	
	go sender(newPlayer)
	receiver(newPlayer)
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
	
	http.Handle("/", websocket.Handler(dispatch))

	log.Printf("world boot complete")
	
	serve(listenAddr)
}

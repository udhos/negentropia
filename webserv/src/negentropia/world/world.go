package main

import (
	//"io"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	//"code.google.com/p/go.net/websocket"
	"golang.org/x/net/websocket"

	"negentropia/webserv/configflag"
	"negentropia/webserv/session"
	"negentropia/webserv/share"
	"negentropia/webserv/store"
	"negentropia/webserv/util"
	"negentropia/world/server"
)

var (
	configFlags   *flag.FlagSet = flag.NewFlagSet("config flags", flag.ExitOnError)
	configList    configflag.FileList
	listenAddr    string
	redisAddr     string
	websocketHost string
)

// Initialize package main
func init() {
	configFlags.Var(&configList, "config", "load config flags from this file")
	configFlags.StringVar(&listenAddr, "listenOn", "127.0.0.2:8000", "websocket listen address [addr]:port")
	configFlags.StringVar(&redisAddr, "redisAddr", "localhost:6379", "redis server address")
	configFlags.StringVar(&websocketHost, "websocketHost", "127.0.0.2", "host part of websocket uri: ws://host:port/path")
	configFlags.StringVar(&server.ObjBaseURL, "objBaseURL", "http://localhost:8080", "base URL to get model .OBJ files")
}

func auth(ws *websocket.Conn) *server.Player {

	var msg server.ClientMsg
	if err := websocket.JSON.Receive(ws, &msg); err != nil {
		log.Printf("auth: failure: %s", err)
		return nil
	}

	if msg.Code != server.CM_CODE_AUTH {
		log.Printf("auth: non-auth code: %d", msg.Code)
		websocket.JSON.Send(ws, server.ClientMsg{Code: server.CM_CODE_FATAL, Data: "auth required"})
		return nil
	}

	sid := msg.Data
	//log.Printf("auth: sid=%s", sid)

	session := session.Load(sid)
	if session == nil {
		msg := fmt.Sprintf("auth: sid=%s: invalid session id", sid)
		log.Printf(msg)
		websocket.JSON.Send(ws, server.ClientMsg{Code: server.CM_CODE_FATAL, Data: fmt.Sprintf("bad auth: %s", msg)})
		return nil
	}

	return &server.Player{Sid: sid,
		Email:        session.ProfileEmail,
		Websocket:    ws,
		SendToPlayer: make(chan *server.ClientMsg),
		Quit:         make(chan int)}
}

// read from player channels and write to player socket
func sender(p *server.Player) {

LOOP:
	for {
		select {
		case <-p.Quit:
			log.Printf("sender: %s %s: quit request", p.Sid, p.Email)

			// destroys the session -- otherwise a rogue client could reconnect
			store.Del(p.Sid)

			// tells the client the session is destroyed -- otherwise a sane client would hopelessy retry
			websocket.JSON.Send(p.Websocket, server.ClientMsg{Code: server.CM_CODE_KILL, Data: "session destroyed due to newer login"})

			break LOOP
		case msg := <-p.SendToPlayer:
			if err := websocket.JSON.Send(p.Websocket, msg); err != nil {
				log.Printf("sender: %s %s: failure: %s", p.Sid, p.Email, err)
				break LOOP
			}
		}
		//log.Printf("sender: %s %s: %q", p.Sid, p.Email, msg)
	}

	log.Printf("sender: %s %s: exiting", p.Sid, p.Email)

	p.Websocket.Close()
}

// read from player socket and write to player input channel
func receiver(p *server.Player) {

	for {
		msg := &server.ClientMsg{} // new(server.ClientMsg)
		if err := websocket.JSON.Receive(p.Websocket, msg); err != nil {
			log.Printf("receiver: %s %s: failure: %s", p.Sid, p.Email, err)
			break
		}
		//log.Printf("receiver: %s %s: %q", p.Sid, p.Email, msg)
		server.InputCh <- &server.PlayerMsg{Player: p, Msg: msg}
	}

	p.Websocket.Close()
}

func dispatch(ws *websocket.Conn) {

	defer ws.Close()

	newPlayer := auth(ws)
	if newPlayer == nil {
		return
	}

	websocket.JSON.Send(ws, server.ClientMsg{Code: server.CM_CODE_INFO, Data: "welcome " + newPlayer.Email})

	server.PlayerAddCh <- newPlayer
	defer func() { server.PlayerDelCh <- newPlayer }()

	go sender(newPlayer) // read from player channels and write to player socket
	receiver(newPlayer)  // read from player socket and write to player input channel
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

	if err := http.ListenAndServe(addr, nil); err != nil {
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

	//session.Init()

	http.Handle("/", websocket.Handler(dispatch))

	log.Printf("world boot complete")

	serve(listenAddr)
}

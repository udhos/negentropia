package server

import (
	"log"
	"fmt"

	"code.google.com/p/go.net/websocket"
	
	"negentropia/webserv/store"
)

const (
	CM_CODE_FATAL = 0
	CM_CODE_INFO  = 1
	CM_CODE_AUTH  = 2 // client->server: let me in
	CM_CODE_ECHO  = 3 // client->server: please echo this
	CM_CODE_KILL  = 4 // server->client: do not attempt reconnect on same session
    CM_CODE_REQZ  = 5 // client->server: please send current zone	
    CM_CODE_ZONE  = 6 // server->client: current zone
)

type ClientMsg struct {
	Code	int
	Data	string
}

type Player struct {
	Sid	          string
	Email         string
	Websocket    *websocket.Conn
	SendToPlayer  chan *ClientMsg
	Quit          chan int
}

func (p *Player) getLocation() string {
	var location string
	if location = store.QueryField(p.Email, "location"); location == "" {
		return "demo"
	}
	return location
}

type PlayerMsg struct {
	Player *Player
	Msg    *ClientMsg
}

var (
	playerTable	                = map[string]*Player {}
	PlayerAddCh	chan *Player    = make(chan *Player)
	PlayerDelCh	chan *Player    = make(chan *Player)
	InputCh     chan *PlayerMsg = make(chan *PlayerMsg)
)

func serve() {
	log.Printf("world server.serve: goroutine started")
	
	for {
		select {
			case p := <- PlayerAddCh:
				playerAdd(p)
			case p := <- PlayerDelCh:
				playerDel(p)
			case m := <- InputCh:
				input(m.Player, m.Msg)
		}
	}
}

func input(p *Player, m *ClientMsg) {
	log.Printf("server.input: %s: %q", p.Email, m)
	
	switch m.Code {
	case CM_CODE_ECHO:
		p.SendToPlayer <- &ClientMsg{CM_CODE_INFO, "echo: " + m.Data}
	case CM_CODE_REQZ:
		log.Printf("server.input: CM_CODE_REQZ FIXME WRITEME");
		p.SendToPlayer <- &ClientMsg{CM_CODE_ZONE, p.getLocation()}
	default:
		log.Printf("server.input: unknown code=%d", m.Code);
		p.SendToPlayer <- &ClientMsg{CM_CODE_INFO, fmt.Sprintf("unknown code=%d", m.Code)}
	}
}

func init() {
	go serve()
}

func playerAdd(newPlayer *Player) {
	if p, ok := playerTable[newPlayer.Email]; ok {
		log.Printf("server.playerAdd: sending quit to existing %s", p.Email)
		p.Quit <- 1
	}
	
	// notice this immediately unregisters the previous player
	playerTable[newPlayer.Email] = newPlayer
}

func playerDel(oldPlayer *Player) {
	log.Printf("server.playerDel: %s %s", oldPlayer.Email, oldPlayer.Sid)
	
	if p, ok := playerTable[oldPlayer.Email]; ok && p.Sid == oldPlayer.Sid {
		// do not unregister wrong player
		delete(playerTable, oldPlayer.Email)
	}
}

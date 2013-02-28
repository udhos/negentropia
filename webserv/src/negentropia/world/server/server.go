package server

import (
	"log"

	"code.google.com/p/go.net/websocket"
	
	//"negentropia/webserv/store"
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

type Player struct {
	Sid	          string
	Email         string
	Websocket    *websocket.Conn
	SendToPlayer  chan *ClientMsg
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
}

func init() {
	go serve()
}

func playerAdd(p *Player) {
	playerTable[p.Email] = p
}

func playerDel(p *Player) {
	delete(playerTable, p.Email)
}

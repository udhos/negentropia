package main

import (
	"fmt"
	"github.com/gopherjs/websocket"
	"honnef.co/go/js/dom"
)

type Websocket struct {
	uri  string
	conn *websocket.Conn
}

func (ws *Websocket) open(uri string) {
	ws.uri = uri

	log(fmt.Sprintf("websocket open: opening: %s", ws.uri))

	c, err := websocket.Dial(ws.uri)
	if err != nil {
		log(fmt.Sprintf("websocket open: could not connect: %s: error=%v", ws.uri, err))
		return
	}

	ws.conn = c

	log(fmt.Sprintf("websocket open: connected: %s", ws.uri))
}

func handleWebsocket(wsUri string) {

	ws := new(Websocket)

	go ws.open(wsUri)

	log(fmt.Sprintf("handleWebsocket: spawned websocket handling: %s", ws.uri))
}

func initWebSocket() bool {

	query := "#wsUri"

	el := dom.GetWindow().Document().QuerySelector(query)
	if el == nil {
		log(fmt.Sprintf("initWebSocket: could not find element: %s", query))
		return true // error
	}
	//span := el.(dom.HTMLSpanElement)
	log(fmt.Sprintf("initWebSocket: %s el=%v", query, el))
	wsUri := el.TextContent()
	if wsUri == "" {
		log(fmt.Sprintf("initWebSocket: empty text for element: %s", query))
		return true // error
	}

	log(fmt.Sprintf("initWebSocket: %s wsUri=%v", query, wsUri))

	go handleWebsocket(wsUri)

	return false // ok
}

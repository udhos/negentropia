package main

import (
	"fmt"
	jsws "github.com/gopherjs/websocket"
	//"golang.org/x/net/websocket"
	"honnef.co/go/js/dom"
	"time"
)

/*
type ClientMsg struct {
	Code int
	Data string
	Tab  map[string]string
}
*/

type Websocket struct {
	uri  string
	conn *jsws.Conn
}

func (ws *Websocket) open(uri string) {
	ws.uri = uri

	log(fmt.Sprintf("websocket open: opening: %s", ws.uri))

	c, err := jsws.Dial(ws.uri)
	if err != nil {
		log(fmt.Sprintf("websocket open: could not connect: %s: error=%v", ws.uri, err))
		return
	}

	ws.conn = c

	log(fmt.Sprintf("websocket open: connected: %s", ws.uri))
}

func handleWebsocket(wsUri string) {

	//ws := new(Websocket)
	ws := &Websocket{}

	ws.open(wsUri)

	for {
		if ws.conn == nil {
			var connectDelay time.Duration = 10
			log(fmt.Sprintf("handleWebsocket: reconnect: %s waiting: %d seconds", ws.uri, connectDelay))
			time.Sleep(time.Second * connectDelay)
			ws.open(wsUri)
			continue
		}

		/*
			msg := &ClientMsg{} // new(server.ClientMsg)
			if err := websocket.JSON.Receive(ws.conn, msg); err != nil {
				log(fmt.Sprintf("handleWebsocket: Receive: %s", err))
				break
			}
		*/

		/*
			var delay time.Duration = 10
			log(fmt.Sprintf("handleWebsocket: %s for loop: waiting %d seconds", ws.uri, delay))
			time.Sleep(time.Second * delay)
		*/
	}
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

	log(fmt.Sprintf("initWebSocket: spawned websocket handling: %s", wsUri))

	return false // ok
}

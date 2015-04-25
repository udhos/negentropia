package main

import (
	"fmt"
	"github.com/gopherjs/websocket"
	"honnef.co/go/js/dom"
)

func handleWebsocket(wsUri string) {

	log(fmt.Sprintf("handleWebSocket: opening: %s", wsUri))

	c, err := websocket.Dial(wsUri)
	if err != nil {
		log(fmt.Sprintf("handleWebSocket: could not connect: %s: error=%v", wsUri, err))
		return
	}

	log(fmt.Sprintf("handleWebSocket: connected: %s", wsUri))

	defer c.Close()

	log(fmt.Sprintf("handleWebSocket: disconnecting: %s", wsUri))
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

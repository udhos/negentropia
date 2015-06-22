package main

import (
	"fmt"
	"github.com/gopherjs/websocket"
	//"golang.org/x/net/websocket"
	"encoding/json"
	"github.com/udhos/cookie"
	"honnef.co/go/js/dom"
	"negentropia/world/server"
	"time"
)

type gameWebsocket struct {
	uri     string
	conn    *websocket.Conn
	status  dom.Element
	encoder *json.Encoder
}

func switchZone(sock *gameWebsocket) {
	sock.write(&server.ClientMsg{Code: server.CM_CODE_SWITCH_ZONE})
}

func requestZone(sock *gameWebsocket) {
	sock.write(&server.ClientMsg{Code: server.CM_CODE_REQZ})
}

func (ws *gameWebsocket) write(msg *server.ClientMsg) error {
	log(fmt.Sprintf("websocket write: writing: %v", msg))

	if ws.conn == nil {
		info := "websocket write: write on broken connection"
		log(info)
		return fmt.Errorf(info)
	}

	if err := ws.encoder.Encode(&msg); err != nil {
		log(fmt.Sprintf("websocket write: error: %s", err))
		ws.conn.Close()
		ws.conn = nil
		ws.status.SetTextContent("disconnected")
		return err
	}

	return nil
}

func (ws *gameWebsocket) open(uri, sid string, status dom.Element) {
	ws.uri = uri
	ws.status = status

	info := fmt.Sprintf("opening: %s", ws.uri)
	log(fmt.Sprintf("websocket open: %s", info))
	ws.status.SetTextContent(info)

	c, err := websocket.Dial(ws.uri)
	if err != nil {
		log(fmt.Sprintf("websocket open: could not connect: %s: error=%v", ws.uri, err))
		ws.conn = nil
		ws.status.SetTextContent("disconnected")
		return
	}

	ws.conn = c
	ws.encoder = json.NewEncoder(ws.conn)

	info = fmt.Sprintf("connected: %s", ws.uri)
	log(fmt.Sprintf("websocket open: %s", info))
	ws.status.SetTextContent(info)

	msg := &server.ClientMsg{Code: server.CM_CODE_AUTH, Data: sid}

	if err := ws.write(msg); err != nil {
		log(fmt.Sprintf("websocket open: JSON encoding error: %s", err))
		return
	}

	log(fmt.Sprintf("websocket open: sent=[%v]", msg))
}

func handleWebsocket(gameInfo *gameState, wsUri, sid string, status dom.Element) {

	log(fmt.Sprintf("handleWebsocket: entering read loop: %s", wsUri))

	defer func() {
		log("handleWebsocket: exiting (goroutine finishing)")
	}()

	gameInfo.sock = &gameWebsocket{}
	gameInfo.sock.open(wsUri, sid, status)

	// reconnect loop
	for {
		if gameInfo.sock.conn == nil {
			var connectDelay time.Duration = 10

			log(fmt.Sprintf("handleWebsocket: reconnect: %s waiting: %d seconds", gameInfo.sock.uri, connectDelay))
			gameInfo.sock.status.SetTextContent("waiting")

			time.Sleep(time.Second * connectDelay)
			gameInfo.sock.open(wsUri, sid, status)
			continue
		}

		msg := &server.ClientMsg{}

		// read loop
		for {
			decoder := json.NewDecoder(gameInfo.sock.conn)

			if err := decoder.Decode(&msg); err != nil {
				log(fmt.Sprintf("handleWebsocket: JSON decoding error: %s", err))
				gameInfo.sock.conn = nil
				gameInfo.sock.status.SetTextContent("disconnected")
				break // reconnect
			}

			//log(fmt.Sprintf("handleWebsocket: received=[%v]", msg))

			if msg.Code == server.CM_CODE_KILL {
				info := fmt.Sprintf("server killed our session: %s", msg.Data)
				log(fmt.Sprintf("handleWebsocket: %s", info))
				gameInfo.sock.status.SetTextContent(info)
				return // stop
			}

			dispatch(gameInfo, msg.Code, msg.Data, msg.Tab)
		}

		/*
			var delay time.Duration = 10
			log(fmt.Sprintf("handleWebsocket: %s for loop: waiting %d seconds", ws.uri, delay))
			time.Sleep(time.Second * delay)
		*/
	}
}

func initWebSocket(gameInfo *gameState) bool {

	sidCookie := "sid"
	sid, ok := cookie.Get(sidCookie)
	if !ok {
		log(fmt.Sprintf("initWebSocket: could not find cookie: %s", sidCookie))
		return true // error
	}

	log(fmt.Sprintf("initWebSocket: found cookie %s=%s", sidCookie, sid))

	//
	// websocket URI
	//

	query := "#wsUri"
	//el := dom.GetWindow().Document().QuerySelector(query)
	el := docQuery(query)
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

	//
	// websocket status
	//

	statusQuery := "#ws_status"
	statusEl := docQuery(statusQuery)
	if statusEl == nil {
		log(fmt.Sprintf("initWebSocket: could not find element: %s", statusQuery))
		return true // error
	}

	// spawn websocket handler
	log(fmt.Sprintf("initWebSocket: spawning websocket handler: %s", wsUri))
	go handleWebsocket(gameInfo, wsUri, sid, statusEl)

	return false // ok
}

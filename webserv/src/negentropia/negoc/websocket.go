package main

import (
	"fmt"
	jsws "github.com/gopherjs/websocket"
	//"golang.org/x/net/websocket"
	"encoding/json"
	"github.com/udhos/cookie"
	"honnef.co/go/js/dom"
	"time"
)

// dup from world/server/server.go
const (
	CM_CODE_FATAL           = 0
	CM_CODE_INFO            = 1
	CM_CODE_AUTH            = 2  // client->server: let me in
	CM_CODE_ECHO            = 3  // client->server: please echo this
	CM_CODE_KILL            = 4  // server->client: do not attempt reconnect on same session
	CM_CODE_REQZ            = 5  // client->server: please send current zone
	CM_CODE_ZONE            = 6  // server->client: reset client zone info
	CM_CODE_SKYBOX          = 7  // server->client: set full skybox
	CM_CODE_PROGRAM         = 8  // server->client: set shader program
	CM_CODE_INSTANCE        = 9  // server->client: set instance
	CM_CODE_INSTANCE_UPDATE = 10 // server->client: update instance
	CM_CODE_MESSAGE         = 11 // server->client: message for user
	CM_CODE_MISSION_NEXT    = 12 // client->server: switch mission
	CM_CODE_SWITCH_ZONE     = 13 // client->server: switch zone
)

// dup from world/server/server.go
type ClientMsg struct {
	Code int
	Data string
	Tab  map[string]string
}

type Websocket struct {
	uri    string
	conn   *jsws.Conn
	status dom.Element
}

func (ws *Websocket) open(uri, sid string, status dom.Element) {
	ws.uri = uri
	ws.status = status

	info := fmt.Sprintf("opening: %s", ws.uri)
	log(fmt.Sprintf("websocket open: %s", info))
	ws.status.SetTextContent(info)

	c, err := jsws.Dial(ws.uri)
	if err != nil {
		log(fmt.Sprintf("websocket open: could not connect: %s: error=%v", ws.uri, err))
		ws.conn = nil
		ws.status.SetTextContent("disconnected")
		return
	}

	ws.conn = c

	info = fmt.Sprintf("connected: %s", ws.uri)
	log(fmt.Sprintf("websocket open: %s", info))
	ws.status.SetTextContent(info)

	msg := &ClientMsg{Code: CM_CODE_AUTH, Data: sid}

	encoder := json.NewEncoder(ws.conn)

	if err := encoder.Encode(&msg); err != nil {
		log(fmt.Sprintf("websocket open: JSON encoding error: %s", err))
		ws.conn = nil
		ws.status.SetTextContent("disconnected")
		return
	}

	log(fmt.Sprintf("websocket open: sent=[%v]", msg))
}

func dispatch(code int, data string, tab map[string]string) {
	log(fmt.Sprintf("dispatch: code=%v data=%v tab=%v", code, data, tab))
}

func handleWebsocket(wsUri, sid string, status dom.Element) {

	log(fmt.Sprintf("handleWebsocket: entering read loop: %s", wsUri))

	defer func() {
		log("handleWebsocket: exiting (goroutine finishing)")
	}()

	ws := &Websocket{}

	ws.open(wsUri, sid, status)

	// reconnect loop
	for {
		if ws.conn == nil {
			var connectDelay time.Duration = 10

			log(fmt.Sprintf("handleWebsocket: reconnect: %s waiting: %d seconds", ws.uri, connectDelay))
			ws.status.SetTextContent("waiting")

			time.Sleep(time.Second * connectDelay)
			ws.open(wsUri, sid, status)
			continue
		}

		msg := &ClientMsg{}

		// read loop
		for {
			decoder := json.NewDecoder(ws.conn)

			if err := decoder.Decode(&msg); err != nil {
				log(fmt.Sprintf("handleWebsocket: JSON decoding error: %s", err))
				ws.conn = nil
				ws.status.SetTextContent("disconnected")
				break // reconnect
			}

			log(fmt.Sprintf("handleWebsocket: received=[%v]", msg))

			if msg.Code == CM_CODE_KILL {
				info := fmt.Sprintf("server killed our session: %s", msg.Data)
				log(fmt.Sprintf("handleWebsocket: %s", info))
				ws.status.SetTextContent(info)
				return // stop
			}

			dispatch(msg.Code, msg.Data, msg.Tab)
		}

		/*
			var delay time.Duration = 10
			log(fmt.Sprintf("handleWebsocket: %s for loop: waiting %d seconds", ws.uri, delay))
			time.Sleep(time.Second * delay)
		*/
	}
}

func initWebSocket() bool {

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
	go handleWebsocket(wsUri, sid, statusEl)

	return false // ok
}

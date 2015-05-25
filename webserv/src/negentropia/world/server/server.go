package server

import (
	"fmt"
	"log"
	"strings"
	"time"

	//"code.google.com/p/go.net/websocket"
	"golang.org/x/net/websocket"

	"negentropia/webserv/store"
)

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

type ClientMsg struct {
	Code int
	Data string
	Tab  map[string]string
}

type Player struct {
	Sid          string
	Email        string
	Websocket    *websocket.Conn
	SendToPlayer chan *ClientMsg
	Quit         chan int
}

type PlayerMsg struct {
	Player *Player
	Msg    *ClientMsg
}

var (
	playerTable                  = map[string]*Player{}
	PlayerAddCh  chan *Player    = make(chan *Player)
	PlayerDelCh  chan *Player    = make(chan *Player)
	InputCh      chan *PlayerMsg = make(chan *PlayerMsg)
	previousTick time.Time
)

func tick(currentTick time.Time, tickInterval time.Duration) {
	elapsed := currentTick.Sub(previousTick)
	if elapsed < tickInterval/2 || elapsed > 2*tickInterval {
		log.Printf("server: tick: ugh: elapsed=%d ms far from interval=%d ms", elapsed/time.Millisecond, tickInterval/time.Millisecond)
	}
	updateAllZones(elapsed)
	previousTick = currentTick
}

func serve() {
	log.Printf("world server.serve: goroutine started")

	var tickPeriod time.Duration = 1000 * time.Millisecond
	log.Printf("world server.serve: ticker at %d ms", tickPeriod/time.Millisecond)
	previousTick = time.Now()
	ticker := time.NewTicker(tickPeriod)

	for {
		select {
		case p := <-PlayerAddCh:
			playerAdd(p)
		case p := <-PlayerDelCh:
			playerDel(p)
		case m := <-InputCh:
			input(m.Player, m.Msg)
		case t := <-ticker.C:
			tick(t, tickPeriod)
		}
	}
}

func sendZoneStatic(p *Player) {
	p.SendToPlayer <- &ClientMsg{
		Code: CM_CODE_ZONE,
		Tab: map[string]string{
			"backfaceCulling": "true",
		},
	}
	p.SendToPlayer <- &ClientMsg{
		Code: CM_CODE_SKYBOX,
		Tab: map[string]string{
			"skyboxURL": "/skybox/skybox_galaxy.json",
		},
	}
	p.SendToPlayer <- &ClientMsg{
		Code: CM_CODE_PROGRAM,
		Tab: map[string]string{
			"programName":    "simpleTexturizer",
			"vertexShader":   "/shader/simpleTex_vs.txt",
			"fragmentShader": "/shader/simpleTex_fs.txt",
		},
	}

	front := []float32{1.0, 0.0, 0.0}
	up := []float32{0.0, 1.0, 0.0}
	coord := []float32{0.0, 0.0, 0.0}
	scale := 1.0

	frontStr := fmt.Sprintf("%f,%f,%f", front[0], front[1], front[2])
	upStr := fmt.Sprintf("%f,%f,%f", up[0], up[1], up[2])
	coordStr := fmt.Sprintf("%f,%f,%f", coord[0], coord[1], coord[2])
	scaleStr := fmt.Sprintf("%f", scale)

	p.SendToPlayer <- &ClientMsg{
		Code: CM_CODE_INSTANCE,
		Tab: map[string]string{
			"objURL":         "/obj/airship.obj",
			"programName":    "simpleTexturizer",
			"directionFront": frontStr,
			"directionUp":    upStr,
			"coord":          coordStr,
			"scale":          scaleStr,
		},
	}

}

func sendZoneDynamic(p *Player, loc string) {

	zoneMsg := ClientMsg{
		Code: CM_CODE_ZONE,
		Tab:  map[string]string{},
	}

	if culling := store.QueryField(loc, "backfaceCulling"); culling != "" {
		zoneMsg.Tab["backfaceCulling"] = culling
	}

	if camCoord := store.QueryField(loc, "cameraCoord"); camCoord != "" {
		zoneMsg.Tab["cameraCoord"] = camCoord
	}

	if len(zoneMsg.Tab) > 0 {
		p.SendToPlayer <- &zoneMsg
	}

	if skybox := store.QueryField(loc, "skyboxURL"); skybox != "" {
		p.SendToPlayer <- &ClientMsg{
			Code: CM_CODE_SKYBOX,
			Tab: map[string]string{
				"skyboxURL": skybox,
			},
		}
	}

	if program := store.QueryField(loc, "programName"); program != "" {
		vertex := store.QueryField(program, "vertexShader")
		fragment := store.QueryField(program, "fragmentShader")

		p.SendToPlayer <- &ClientMsg{
			Code: CM_CODE_PROGRAM,
			Tab: map[string]string{
				"programName":    program,
				"vertexShader":   vertex,
				"fragmentShader": fragment,
			},
		}
	}

	if instanceList := store.QueryField(loc, "instanceList"); instanceList != "" {
		instances := store.QuerySet(instanceList)

		for _, inst := range instances {

			obj := store.QueryField(inst, "obj")
			coord := store.QueryField(inst, "coord")
			scale := store.QueryField(inst, "scale")
			mission := store.QueryField(inst, "mission")

			url := store.QueryField(obj, "objURL")
			globeRadius := store.QueryField(obj, "globeRadius")
			globeTextureURL := store.QueryField(obj, "globeTextureURL")
			program := store.QueryField(obj, "programName")
			front := store.QueryField(obj, "modelFront")
			up := store.QueryField(obj, "modelUp")
			repeat := store.QueryField(obj, "repeatTexture")

			log.Printf("sendZoneDynamic: id=%s obj=%s objURL=%s", inst, obj, url)

			m := map[string]string{
				"id":            inst,
				"obj":           obj,
				"programName":   program,
				"modelFront":    front,
				"modelUp":       up,
				"repeatTexture": repeat,
				"coord":         coord,
				"scale":         scale,
				"mission":       mission,
			}

			if url != "" {
				m["objURL"] = url
			} else {
				m["globeRadius"] = globeRadius
				m["globeTextureURL"] = globeTextureURL
			}

			p.SendToPlayer <- &ClientMsg{
				Code: CM_CODE_INSTANCE,
				Tab:  m,
			}

		}
	}

	max := 3
	for i := 1; i <= max; i++ {
		msgPlayer(p, fmt.Sprintf("world server: line %d of %d", i, max))
	}
}

func msgPlayer(p *Player, msg string) {
	p.SendToPlayer <- &ClientMsg{Code: CM_CODE_MESSAGE, Data: msg}
}

func sendZone(p *Player, loc string) {
	if strings.HasPrefix(loc, "z:") {
		sendZoneDynamic(p, loc)
		return
	}

	sendZoneStatic(p)
}

func input(p *Player, m *ClientMsg) {
	log.Printf("server.input: %s: %q", p.Email, m)

	switch m.Code {
	case CM_CODE_ECHO:
		p.SendToPlayer <- &ClientMsg{Code: CM_CODE_INFO, Data: "echo: " + m.Data}
	case CM_CODE_REQZ:
		/*
			location:
				""    --> send "demo"       --> client will load hard-coded demo zone
				"z:*" --> send dynamic zone (loaded from redis db)
				"*"   --> send static zone (hard-coded at server)
		*/
		if loc := store.QueryField(p.Email, "location"); loc == "" {
			p.SendToPlayer <- &ClientMsg{Code: CM_CODE_ZONE, Data: "demo"}
		} else {
			sendZone(p, loc)
		}
	case CM_CODE_MISSION_NEXT:
		for id := range m.Tab {
			missionNext(p, id)
		}
	case CM_CODE_SWITCH_ZONE:
		switchZone(p)
	default:
		log.Printf("server.input: unknown code=%d", m.Code)
		p.SendToPlayer <- &ClientMsg{Code: CM_CODE_INFO, Data: fmt.Sprintf("unknown code=%d data=%v tab=%v", m.Code, m.Data, m.Tab)}
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

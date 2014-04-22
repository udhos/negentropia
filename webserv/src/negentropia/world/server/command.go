package server

import (
	"fmt"
	"log"

	"negentropia/webserv/store"
)

func missionNext(p *Player, unitId string) {
	owner := store.QueryField(unitId, "owner")

	if owner != p.Email {
		msgPlayer(p, fmt.Sprintf("unit %v: mission command: refused", unitId))
		return
	}

	msgPlayer(p, fmt.Sprintf("unit %v: mission command: accepted", unitId))

	log.Printf("FIXME missionNext: player=%v unit=%v owner=%v", p.Email, unitId, owner)
}

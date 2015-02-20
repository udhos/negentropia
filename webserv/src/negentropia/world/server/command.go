package server

import (
	"fmt"
	"log"
	"strings"

	"negentropia/webserv/store"
)

func missionNext(p *Player, unitId string) {

	zoneId := store.QueryField(p.Email, "location")
	if !strings.HasPrefix(zoneId, "z:") {
		log.Printf("missionNext: BAD ZONE ID zone=%v player=%v unit=%v", zoneId, p.Email, unitId)
		return
	}

	zone, zok := zoneTable[zoneId]
	if !zok {
		log.Printf("missionNext: ZONE NOT FOUND zone=%v player=%v unit=%v", zoneId, p.Email, unitId)
		return
	}

	owner := store.QueryField(unitId, "owner")

	if owner != p.Email {
		// mission silently refused
		return
	}

	unit := zone.unitTable[unitId]

	switch unit.mission {
	case "": // no mission
		unit.mission = "rotateYaw"
	case "rotateYaw":
		unit.mission = "hunt"
	case "hunt":
		unit.mission = ""
	default:
		log.Printf("missionNext: UNKNOWN MISSION zone=%s unit=%s mission=%s", zone.zid, unit.uid, unit.mission)
	}

	msgPlayer(p, fmt.Sprintf("unit %v: new mission: [%v]", unitId, unit.mission))
}

func switchZone(p *Player) {

	//
	// 1. Scan zones
	//
	zones := store.QueryKeys("z:*") // FIXME: replace redis KEYS with redis SCAN
	if len(zones) < 1 {
		log.Printf("switchZone: no zone found")
		return
	}

	log.Printf("switchZone: FIXME WRITEME")

	// 2. Sort zone list

	// 3. Switch to next zone

}

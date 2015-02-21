package server

import (
	"fmt"
	"log"
	"sort"
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

func playerToZone(p *Player, zid string) {
	store.SetField(p.Email, "location", zid)
	sendZone(p, zid)
}

func switchZone(p *Player) {

	// 1. Get zone list
	zones := store.QueryKeys("z:*") // FIXME: replace redis KEYS with redis SCAN
	if len(zones) < 1 {
		log.Printf("switchZone: no zone found")
		return
	}

	// no zone!!
	loc := store.QueryField(p.Email, "location")
	if loc == "" {
		playerToZone(p, zones[0]) // pick any zone
		return
	}

	// 2. Sort zone list

	sort.Strings(zones)

	// 3. Switch to next zone

	for i, zid := range zones {
		if zid == loc {
			// found current zone

			var newZid string

			if i < len(zones)-1 {
				newZid = zones[i+1]
			} else {
				newZid = zones[0]
			}

			playerToZone(p, newZid)

			msgPlayer(p, fmt.Sprintf("switchZone: old=%s new=%s", loc, newZid))

			return
		}
	}

	msg := fmt.Sprintf("switchZone: zone %s not found", loc)
	log.Printf(msg)
	msgPlayer(p, msg)
}

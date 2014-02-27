package server

import (
	"log"

	"github.com/spate/vectormath"

	"negentropia/webserv/store"
)

// Non-persitent instance data
type Unit struct {
	uid         string
	front       vectormath.Vector3
	up          vectormath.Vector3
	linearSpeed float32 // m/s
	yawSpeed    float32 // rad/s
	pitchSpeed  float32 // rad/s
}

type Zone struct {
	zid       string
	unitTable map[string]*Unit
}

var (
	zoneTable = map[string]*Zone{}
)

func newZone(zid string) *Zone {
	return &Zone{zid: zid, unitTable: make(map[string]*Unit)}
}

func newUnit(uid string) *Unit {
	return &Unit{uid: uid}
}

func updateAllZones() {
	//
	// Scan zones
	//
	zones := store.QueryKeys("z:*")
	if len(zones) < 1 {
		log.Printf("updateAllZones: no zone found")
		return
	}

	for _, zid := range zones {
		zone, zok := zoneTable[zid]
		if !zok {
			zone = newZone(zid)
			zoneTable[zid] = zone
		}

		if zone == nil {
			log.Printf("updateAllZones: failure creating zone=%s", zid)
			continue
		}

		instanceList := store.QueryField(zid, "instanceList")
		if instanceList == "" {
			log.Printf("updateAllZones: zone=%s: no instanceList", zid)
			continue
		}

		instances := store.QuerySet(instanceList)
		if len(instances) < 1 {
			log.Printf("updateAllZones: zone=%s: empty instanceList", zid)
			continue
		}

		for _, uid := range instances {
			unit, uok := zone.unitTable[uid]
			if !uok {
				unit = newUnit(uid)
				zone.unitTable[uid] = unit
			}

			if unit == nil {
				log.Printf("updateAllZones: failure creating unit=%s", zid)
				continue
			}

			coord := store.QueryField(uid, "coord")
			mission := store.QueryField(uid, "mission")

			log.Printf("updateAllZones: zone=%s unit=%s coord=%s mission=%s", zid, uid, coord, mission)
		}
	}
}

package server

import (
	"fmt"
	"log"
	"math"
	"time"

	"github.com/spate/vectormath"

	"negentropia/webserv/store"
)

// Non-persitent instance data
type Unit struct {
	uid         string
	coord       vectormath.Vector3
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

func newUnit(uid, coord, front, up string) (*Unit, error) {

	unit := &Unit{uid: uid}

	if err := parseVector3(&unit.coord, coord); err != nil {
		e := fmt.Errorf("newUnit: unit=%s coord=[%s] parse failure: %s", uid, coord, err)
		log.Print(e)
		return nil, e
	}
	if err := parseVector3(&unit.front, front); err != nil {
		e := fmt.Errorf("newUnit: unit=%s front=[%s] parse failure: %s", uid, front, err)
		log.Print(e)
		return nil, e
	}
	if err := parseVector3(&unit.up, up); err != nil {
		e := fmt.Errorf("newUnit: unit=%s up=[%s] parse failure: %s", uid, up, err)
		log.Print(e)
		return nil, e
	}

	return unit, nil
}

func updateUnit(elapsed time.Duration, zone *Zone, unit *Unit, mission string) {
	log.Printf("updateUnit: zone=%s unit=%s mission=%s", zone.zid, unit.uid, mission)

	switch mission {
	case "": // no mission
	case "rotateYaw":
		unit.linearSpeed = 0.0
		unit.yawSpeed = 10.0 * math.Pi / 180.0 // 10 degrees/s
		unit.pitchSpeed = 0.0
	default:
		log.Printf("updateUnit: UNKNOWN MISSION zone=%s unit=%s mission=%s", zone.zid, unit.uid, mission)
	}
}

func updateAllZones(elapsed time.Duration) {
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

				coord := store.QueryField(uid, "coord")

				// Fetch from model
				model := store.QueryField(uid, "obj")
				modelFront := store.QueryField(model, "modelFront")
				modelUp := store.QueryField(model, "modelUp")

				var err error
				unit, err = newUnit(uid, coord, modelFront, modelUp)
				if err != nil {
					log.Printf("new unit=%s failure: %s", uid, err)
					continue
				}
				zone.unitTable[uid] = unit
			}

			if unit == nil {
				log.Printf("updateAllZones: failure creating unit=%s", uid)
				continue
			}

			mission := store.QueryField(uid, "mission")

			updateUnit(elapsed, zone, unit, mission)
		}
	}
}

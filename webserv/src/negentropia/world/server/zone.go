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

func (unit *Unit) rightDirection() vectormath.Vector3 {

	if !vector3Unit(unit.front) {
		log.Printf("unit.rightDirection: NOT UNITARY: front=%s length=%f", vector3String(unit.front), unit.front.Length())
	}

	if !vector3Unit(unit.up) {
		log.Printf("unit.rightDirection: NOT UNITARY: up=%s length=%f", vector3String(unit.up), unit.up.Length())
	}

	if !vector3Orthogonal(unit.front, unit.up) {
		log.Printf("unit.rightDirection: NOT ORTHOGONAL: front=%s up=%s: dot=%f",
			vector3String(unit.front), vector3String(unit.up), vectormath.V3Dot(&unit.front, &unit.up))
	}

	var right vectormath.Vector3
	vectormath.V3Cross(&right, &unit.front, &unit.up)
	vectormath.V3Normalize(&right, &right)

	if !vector3Unit(right) {
		log.Printf("unit.rightDirection: NOT UNITARY: right=%s length=%f", vector3String(right), right.Length())
	}

	return right
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

	vectormath.V3Normalize(&unit.front, &unit.front)
	vectormath.V3Normalize(&unit.up, &unit.up)

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

		// angle to rotate
		rad := unit.yawSpeed * float32(elapsed) / float32(time.Second)

		// axis to rotate around
		//var rightDirection vectormath.Vector3
		//vectormath.V3Cross(&rightDirection, &unit.front, &unit.up)
		rightDirection := unit.rightDirection()

		// quaternion representing rotation
		var quat vectormath.Quat
		vectormath.QMakeRotationAxis(&quat, rad, &rightDirection)

		// apply quaternion rotation to front direction
		oldFront := unit.front
		vectormath.QRotate(&unit.front, &quat, &oldFront)
		vectormath.V3Normalize(&unit.front, &unit.front)

		if !vector3Unit(unit.front) {
			log.Printf("updateUnit: NOT UNITARY: front=%s length=%f", vector3String(unit.front), unit.front.Length())
		}

		if !vector3Orthogonal(unit.front, rightDirection) {
			log.Printf("updateUnit: NOT ORTHOGONAL: front=%s right=%s: dot=%f",
				vector3String(unit.front), vector3String(rightDirection), vectormath.V3Dot(&unit.front, &rightDirection))
		}

		// calculate new up direction
		vectormath.V3Cross(&unit.up, &rightDirection, &unit.front)
		vectormath.V3Normalize(&unit.up, &unit.up)

		if !vector3Unit(unit.up) {
			log.Printf("updateUnit: NOT UNITARY: up=%s length=%f", vector3String(unit.up), unit.up.Length())
		}

		log.Printf("rotateYaw: front=%s up=%s right=%s",
			vector3String(unit.front), vector3String(unit.up), vector3String(rightDirection))
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

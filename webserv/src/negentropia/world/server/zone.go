package server

import (
	"fmt"
	"log"
	"math"
	//"strconv"
	"time"

	"github.com/udhos/vectormath"

	"negentropia/webserv/store"
)

type Unit struct {
	uid         string
	coord       vectormath.Vector3
	front       vectormath.Vector3
	up          vectormath.Vector3
	linearSpeed float32 // m/s
	yawSpeed    float32 // rad/s
	pitchSpeed  float32 // rad/s
	rollSpeed   float32 // rad/s
	mission     string
	delete      bool
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
	delete    bool
}

var (
	zoneTable = map[string]*Zone{}
)

func newZone(zid string) *Zone {
	return &Zone{zid: zid, unitTable: make(map[string]*Unit)}
}

func newUnit(uid, coord, front, up, mission string) (*Unit, error) {

	unit := &Unit{uid: uid, mission: mission}

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

func sendUnitUpdate(unit *Unit, zid string) {

	// scan world's player table
	for email, p := range playerTable {

		// get player zone
		loc := store.QueryField(email, "location")
		if loc == "" {
			log.Printf("sendUnitUpdate: empty player location: unit=%s zone=%s player=%s",
				unit.uid, zid, email)
			continue
		}
		// skip players belonging to other zones
		if loc != zid {
			continue
		}

		log.Printf("sendUnitUpdate: unit=%s zone=%s player=%s front=%s up=%s coord=%s mission=%s",
			unit.uid, zid, email,
			vector3String(unit.front), vector3String(unit.up), vector3String(unit.coord),
			unit.mission)

		// send unit update to player
		p.SendToPlayer <- &ClientMsg{
			Code: CM_CODE_INSTANCE_UPDATE,
			Tab: map[string]string{
				"id":      unit.uid,
				"front":   vector3String(unit.front),
				"up":      vector3String(unit.up),
				"coord":   vector3String(unit.coord),
				"mission": unit.mission,
			},
		}

	}

}

func rotateYaw(elapsed time.Duration, zone *Zone, unit *Unit) {
	unit.linearSpeed = 0.0
	unit.yawSpeed = 20.0 * math.Pi / 180.0 // 20 degrees/s
	unit.pitchSpeed = 0.0
	unit.rollSpeed = 0.0

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
		log.Printf("rotateYaw: NOT UNITARY: front=%s length=%f", vector3String(unit.front), unit.front.Length())
	}

	if !vector3Orthogonal(unit.front, rightDirection) {
		log.Printf("rotateYaw: NOT ORTHOGONAL: front=%s right=%s: dot=%f",
			vector3String(unit.front), vector3String(rightDirection), vectormath.V3Dot(&unit.front, &rightDirection))
	}

	// calculate new up direction
	vectormath.V3Cross(&unit.up, &rightDirection, &unit.front)
	vectormath.V3Normalize(&unit.up, &unit.up)

	if !vector3Unit(unit.up) {
		log.Printf("rotateYaw: NOT UNITARY: up=%s length=%f", vector3String(unit.up), unit.up.Length())
	}

	/*
		log.Printf("rotateYaw: front=%s up=%s right=%s",
			vector3String(unit.front), vector3String(unit.up), vector3String(rightDirection))
	*/

	sendUnitUpdate(unit, zone.zid)
}

func hunt(elapsed time.Duration, zone *Zone, unit *Unit) {
	// FIXME WRITEME
	//
	// future bounding sphere intersection => future collision likely
	// if future collision is likely, maneuver to avoid it, then finish
	// fire if there is enemy's bounding sphere available in attack cone
	// maneuver to put attack cone around nearest enemy
	//
	// nearest enemy: bruteforce x kdtree ?
	// kdtree: http://godoc.org/code.google.com/p/biogo.store/kdtree
	// kdtree: http://godoc.org/code.google.com/p/eaburns/kdtree
}

func updateUnit(elapsed time.Duration, zone *Zone, unit *Unit) {
	//log.Printf("updateUnit: zone=%s unit=%s mission=%s", zone.zid, unit.uid, mission)

	switch unit.mission {
	case "": // no mission
	case "rotateYaw":
		rotateYaw(elapsed, zone, unit)
	case "hunt":
		hunt(elapsed, zone, unit)
	default:
		log.Printf("updateUnit: UNKNOWN MISSION zone=%s unit=%s mission=%s", zone.zid, unit.uid, unit.mission)
	}
}

func updateAllZones(elapsed time.Duration) {
	//
	// Scan zones
	//
	zones := store.QueryKeys("z:*") // FIXME: replace redis KEYS with redis SCAN
	if len(zones) < 1 {
		log.Printf("updateAllZones: no zone found")
		return
	}

	// mark all zones for deletion
	for _, zone := range zoneTable {
		zone.delete = true
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

		zone.delete = false // do not delete this zone

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

		// mark units for deletion
		for _, unit := range zone.unitTable {
			unit.delete = true
		}

		for _, uid := range instances {
			unit, uok := zone.unitTable[uid]
			if !uok {

				coord := store.QueryField(uid, "coord")
				mission := store.QueryField(uid, "mission")

				// Fetch from model
				model := store.QueryField(uid, "obj")
				modelFront := store.QueryField(model, "modelFront")
				modelUp := store.QueryField(model, "modelUp")

				var err error
				unit, err = newUnit(uid, coord, modelFront, modelUp, mission)
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

			unit.delete = false // do not delete this unit
		}

		// delete not found units
		for uid, unit := range zone.unitTable {
			if unit.delete {
				log.Printf("deleting unit: zone=%s unit=%s", zid, uid)
				delete(zone.unitTable, uid)
			}
		}

		// from here all units from zone are available in unitTable

		// update all zone units
		for _, unit := range zone.unitTable {
			updateUnit(elapsed, zone, unit)
		}

	} // range zones

	// delete not found zones
	for zid, zone := range zoneTable {
		if zone.delete {
			log.Printf("deleting zone: %s", zid)
			delete(zoneTable, zid)
		}
	}

}

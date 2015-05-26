package server

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/udhos/vectormath"

	"negentropia/webserv/store"
	"negentropia/world/obj"
)

type Unit struct {
	uid            string
	coord          vectormath.Vector3
	front          vectormath.Vector3
	up             vectormath.Vector3
	linearSpeed    float64 // m/s
	yawSpeed       float64 // rad/s
	pitchSpeed     float64 // rad/s
	rollSpeed      float64 // rad/s
	boundingRadius float64 // bounding sphere radius (meter)
	mission        string
	delete         bool
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

func newUnit(uid, coord, front, up, mission string, radius float64) (*Unit, error) {

	unit := &Unit{uid: uid, mission: mission, boundingRadius: radius}

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

		/*
			log.Printf("sendUnitUpdate: unit=%s zone=%s player=%s front=%s up=%s coord=%s mission=%s",
				unit.uid, zid, email,
				vector3String(unit.front), vector3String(unit.up), vector3String(unit.coord),
				unit.mission)
		*/

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

func missionRotateYaw(elapsed time.Duration, zone *Zone, unit *Unit) {
	unit.linearSpeed = 0.0
	unit.yawSpeed = 20.0 * math.Pi / 180.0 // 20 degrees/s
	unit.pitchSpeed = 0.0
	unit.rollSpeed = 0.0

	unitMove(unit, elapsed)

	sendUnitUpdate(unit, zone.zid)
}

func missionHunt(elapsed time.Duration, zone *Zone, unit *Unit) {
	// FIXME WRITEME
	//
	// 1. if future collision is likely, maneuver to avoid it
	// use capsule collision detection
	//
	// 2. fire if there is enemy's bounding sphere available in attack cone
	// search enemy in attack cone
	// collision: attack cone vs enemy bouding sphere
	//
	// 3. maneuver to put attack cone around nearest enemy
	// find nearest enemy
	// maneuver towards nearest enemy
	//
	// nearest enemy: bruteforce x kdtree ?
	// kdtree: http://godoc.org/code.google.com/p/biogo.store/kdtree
	// kdtree: http://godoc.org/code.google.com/p/eaburns/kdtree

	unit.linearSpeed = 0.1 // 0.1  m/s
	unit.yawSpeed = 0.0
	unit.pitchSpeed = 0.0
	unit.rollSpeed = 0.0

	unitMove(unit, elapsed)

	sendUnitUpdate(unit, zone.zid)
}

func updateUnit(elapsed time.Duration, zone *Zone, unit *Unit) {
	//log.Printf("updateUnit: zone=%s unit=%s mission=%s", zone.zid, unit.uid, unit.mission)

	switch unit.mission {
	case "": // no mission
	case "rotateYaw":
		missionRotateYaw(elapsed, zone, unit)
	case "hunt":
		missionHunt(elapsed, zone, unit)
	default:
		log.Printf("updateUnit: UNKNOWN MISSION zone=%s unit=%s mission=%s", zone.zid, unit.uid, unit.mission)
	}
}

var ObjBaseURL string

func httpFetch(url string) ([]byte, error) {
	fullURL := ObjBaseURL + url
	resp, err := http.Get(fullURL)
	if err != nil {
		return nil, fmt.Errorf("httpFetch: get url=%v: %v", fullURL, err)
	}
	defer resp.Body.Close()

	var info []byte
	info, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("httpFetch: read all: url=%v: %v", fullURL, err)
	}

	return info, nil
}

func loadModelRadius(model, objURL string) float64 {
	//buf := []byte("###test\n\nv 2.133923 2.037626 -0.070400\neraseme\nv 2.062405 2.303306 -0.070400\nv 2.062404 2.303306 0.060598")
	buf, err := httpFetch(objURL)
	if err != nil {
		log.Printf("loadModelRadius: fetch model=%v objURL=%v error: %v", model, objURL, err)
		return 1.0
	}
	var o *obj.Obj
	if o, err = obj.NewObjFromBuf(buf, func(msg string) { log.Printf("loadModelRadius: %s", msg) }); err != nil {
		log.Printf("loadModelRadius: parse model=%v objURL=%v error: %v", model, objURL, err)
		return 1.0
	}
	size := len(o.Coord)
	if size < 3 {
		log.Printf("loadModelRadius: model=%v objURL=%v short vertex buffer size=%v", model, objURL, size)
		return 1.0
	}
	minX, minY, minZ := o.Coord[0], o.Coord[1], o.Coord[2]
	maxX, maxY, maxZ := o.Coord[0], o.Coord[1], o.Coord[2]
	for i := 3; i < size; i += 3 {
		x, y, z := o.Coord[i], o.Coord[i+1], o.Coord[i+2]
		if x < minX {
			minX = x
		}
		if y < minY {
			minY = y
		}
		if z < minZ {
			minZ = z
		}
		if x > maxX {
			maxX = x
		}
		if y > maxY {
			maxY = y
		}
		if z > maxZ {
			maxZ = z
		}
	}
	dx := maxX - minX
	dy := maxY - minY
	dz := maxZ - minZ
	radius := math.Sqrt(dx*dx+dy*dy+dz*dz) / 2.0
	return radius
}

var modelRadiusCache = map[string]float64{}

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
			//log.Printf("updateAllZones: zone=%s: no instanceList", zid)
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
				scaleStr := store.QueryField(uid, "scale")

				var err error
				var scale float64
				if scale, err = strconv.ParseFloat(scaleStr, 64); err != nil {
					log.Print(fmt.Errorf("failure parsing unit scale: [%v]: %v", scaleStr, err))
				}

				// Fetch from model
				model := store.QueryField(uid, "obj")
				modelFront := store.QueryField(model, "modelFront")
				modelUp := store.QueryField(model, "modelUp")
				objURL := store.QueryField(model, "objURL")

				modelRadius, rok := modelRadiusCache[objURL]
				if !rok {
					modelRadius = loadModelRadius(model, objURL)
					modelRadiusCache[objURL] = modelRadius
					log.Printf("modelRadiusCache miss: model=%v url=%v radius=%v", model, objURL, modelRadius)
				}

				radius := scale * modelRadius

				if unit, err = newUnit(uid, coord, modelFront, modelUp, mission, radius); err != nil {
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

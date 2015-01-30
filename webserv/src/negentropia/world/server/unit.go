package server

import (
	"log"
	"time"

	"github.com/udhos/vectormath"

	"negentropia/world/util"
)

func unitRotateYaw(unit *Unit, elapsed time.Duration) {

	// angle to rotate
	rad := unit.yawSpeed * float64(elapsed) / float64(time.Second)

	// axis to rotate around
	//var rightDirection vectormath.Vector3
	//vectormath.V3Cross(&rightDirection, &unit.front, &unit.up)
	rightDirection := unit.rightDirection()

	// quaternion representing rotation
	var quat vectormath.Quat
	vectormath.QMakeRotationAxis(&quat, float32(rad), &rightDirection)

	// apply quaternion rotation to front direction
	oldFront := unit.front
	vectormath.QRotate(&unit.front, &quat, &oldFront)
	vectormath.V3Normalize(&unit.front, &unit.front)

	if !vector3Unit(unit.front) {
		log.Printf("rotateYaw: NOT UNITARY: front=%s length=%f", vector3String(unit.front), unit.front.Length())
		return
	}

	if !vector3Orthogonal(unit.front, rightDirection) {
		log.Printf("rotateYaw: NOT ORTHOGONAL: front=%s right=%s: dot=%f",
			vector3String(unit.front), vector3String(rightDirection), vectormath.V3Dot(&unit.front, &rightDirection))
		return
	}

	// calculate new up direction
	vectormath.V3Cross(&unit.up, &rightDirection, &unit.front)
	vectormath.V3Normalize(&unit.up, &unit.up)

	if !vector3Unit(unit.up) {
		log.Printf("rotateYaw: NOT UNITARY: up=%s length=%f", vector3String(unit.up), unit.up.Length())
	}
}

func unitForward(unit *Unit, elapsed time.Duration) {
	if !vector3Unit(unit.front) {
		log.Printf("unitForward: NOT UNITARY: front=%v length=%v", vector3String(unit.front), unit.front.Length())
		return
	}

	/*
		var speed vectormath.Vector3
		vectormath.V3ScalarMul(&speed, &unit.front, float32(unit.linearSpeed*float64(elapsed)/float64(time.Second)))
	*/

	var speed [3]float64
	speed[0] = float64(unit.front.X)
	speed[1] = float64(unit.front.Y)
	speed[2] = float64(unit.front.Z)
	scale := unit.linearSpeed * float64(elapsed) / float64(time.Second)
	//log.Printf("before mul: %v scale=%v", speed, scale)
	v3scalarMul(speed[:], scale)
	//log.Printf("after mul: %v scale=%v", speed, scale)

	/*
		if diff := unit.linearSpeed - float64(speed.Length()); !util.CloseToZero(diff) {
			log.Printf("unitForward: MISMATCH: unit=%v forward=%v speed=%v linearSpeed=%v speed=%v diff=%v", unit.uid, vector3String(unit.front), vector3String(speed), unit.linearSpeed, speed.Length(), diff)
		}
	*/
	speedLen := v3len(speed[0], speed[1], speed[2])
	if diff := unit.linearSpeed - speedLen; !util.CloseToZeroEpsilon(diff, 0.001) {
		log.Printf("unitForward: MISMATCH: unit=%v forward=%v speed=%v linearSpeed=%v speedLen=%v diff=%v", unit.uid, vector3String(unit.front), speed, unit.linearSpeed, speedLen, diff)
	}

	//vectormath.V3Add(&unit.coord, &unit.coord, &speed)
	var coord [3]float64
	coord[0] = float64(unit.coord.X)
	coord[1] = float64(unit.coord.Y)
	coord[2] = float64(unit.coord.Z)
	v3add(coord[:], coord[0], coord[1], coord[2], speed[0], speed[1], speed[2])
	unit.coord.X = float32(coord[0])
	unit.coord.Y = float32(coord[1])
	unit.coord.Z = float32(coord[2])
}

/*
	linearSpeed    float64 // m/s
	yawSpeed       float64 // rad/s
	pitchSpeed     float64 // rad/s
	rollSpeed      float64 // rad/s
*/
func unitMove(unit *Unit, elapsed time.Duration) {
	if !util.CloseToZero(unit.linearSpeed) {
		unitForward(unit, elapsed)
	}

	if !util.CloseToZero(unit.yawSpeed) {
		unitRotateYaw(unit, elapsed)
	}
}

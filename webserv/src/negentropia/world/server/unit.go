package server

import (
	"log"
	"time"

	"github.com/udhos/vectormath"
)

func unitRotateYaw(elapsed time.Duration, unit *Unit) {

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
}

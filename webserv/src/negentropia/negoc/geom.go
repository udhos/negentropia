package main

import (
	"math"
	//"fmt"
)

type ray struct {
	originX, originY, originZ          float64
	directionX, directionY, directionZ float64
}

func (r ray) getPoint(t float64) (float64, float64, float64) {
	return r.originX + r.directionX*t, r.originY + r.directionY*t, r.originZ + r.directionZ*t
}

type sphere struct {
	centerX, centerY, centerZ float64
	radius                    float64
}

/*
http://www.csee.umbc.edu/~olano/435f02/ray-sphere.html

output:
hit: is there intersection?
t1: intersection point1 = r.origin + r.direction * t1
t2: intersection point2 = r.origin + r.direction * t2
*/
func intersectRaySphere(r ray, s sphere) (hit bool, t1, t2 float64) {

	coX, coY, coZ := r.originX-s.centerX, r.originY-s.centerY, r.originZ-s.centerZ

	a := lengthSquared3(r.directionX, r.directionY, r.directionZ)
	b := 2.0 * dot3(r.directionX, r.directionY, r.directionZ, coX, coY, coZ)
	c := lengthSquared3(coX, coY, coZ) - s.radius*s.radius

	delta := b*b - 4*a*c
	if delta < 0.0 {
		return
	}

	hit = true
	deltaRoot := math.Sqrt(delta)
	a2 := 2.0 * a
	t1 = (-b - deltaRoot) / a2
	t2 = (-b + deltaRoot) / a2

	return
}

package main

import (
	"math"
	//"fmt"
)

type ray struct {
	originX, originY, originZ          float64
	directionX, directionY, directionZ float64
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

	a := lengthSquared3(r.directionX, r.directionY, r.directionZ)
	coX, coY, coZ := r.originX-s.centerX, r.originY-s.centerY, r.originZ-s.centerZ
	b := 2 * dot3(r.directionX, r.directionY, r.directionZ, coX, coY, coZ)
	c := lengthSquared3(coX, coY, coZ) - s.radius*s.radius

	delta := b*b - 4*a*c
	if delta < 0.0 {
		return
	}

	hit = true
	deltaRoot := math.Sqrt(delta)
	t1 = (-b - deltaRoot) / 2 * a
	t2 = (-b + deltaRoot) / 2 * a

	return
}

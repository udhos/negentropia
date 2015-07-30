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
bool intersect(Ray* r, Sphere* s, float* t1, float *t2)
{
	//solve for tc
	float L = s->center - r->origin;
	float tc = dot(L, r->direction);

	if ( tc < 0.0 ) return false;
	float d2 = (tc*tc) - (L*L);

	float radius2 = s->radius * s->radius;
	if ( d2 > radius2) return false;

	//solve for t1c
	float t1c = sqrt( radius2 - d2 );

	//solve for intersection points
	*t1 = tc - t1c;
	*t2 = tc + t1c;

	return true;
}

http://kylehalladay.com/blog/tutorial/math/2013/12/24/Ray-Sphere-Intersection.html

output:
hit: is there intersection?
t1: intersection point1 = r.origin + r.direction * t1
t2: intersection point2 = r.origin + r.direction * t2
*/
func intersectRaySphere(r ray, s sphere) (hit bool, t1, t2 float64) {

	Lx, Ly, Lz := s.centerX-r.originX, s.centerY-r.originY, s.centerZ-r.originZ
	tc := dot3(Lx, Ly, Lz, r.directionX, r.directionY, r.directionZ)
	if tc < 0.0 {
		return
	}

	d2 := tc*tc - dot3(Lx, Ly, Lz, Lx, Ly, Lz)
	radius2 := s.radius * s.radius
	if d2 > radius2 {
		return
	}

	t1c := math.Sqrt(radius2 - d2)
	t1 = tc - t1c
	t2 = tc + t1c

	hit = true

	return
}

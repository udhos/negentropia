package server

import (
	"errors"
	"math"

	"github.com/udhos/vectormath"

	"negentropia/world/util"
)

type Sphere struct {
	center    vectormath.Vector3
	radius    float64
	radiusSqr float64
}

type Cone struct {
	vertex        vectormath.Vector3
	axis          vectormath.Vector3
	sinReciprocal float64
	cosSqr        float64
	sinSqr        float64
}

func setSphere(s *Sphere, centerX, centerY, centerZ, radius float64) {
	vectormath.V3MakeFromElems(&s.center, float32(centerX), float32(centerY), float32(centerZ))
	s.radius = radius
	s.radiusSqr = radius * radius
}

func setCone(k *Cone, vertexX, vertexY, vertexZ, axisX, axisY, axisZ, angleRadians float64) error {
	vectormath.V3MakeFromElems(&k.vertex, float32(vertexX), float32(vertexY), float32(vertexZ))
	vectormath.V3MakeFromElems(&k.axis, float32(axisX), float32(axisY), float32(axisZ))
	sin := math.Sin(angleRadians)
	if util.CloseToZero(sin) {
		return errors.New("setCone: invalid zero sin")
	}
	k.sinSqr = sin * sin
	k.cosSqr = 1 - k.sinSqr
	k.sinReciprocal = 1 / sin
	return nil
}

/*
bool SphereIntersectsCone (Sphere S, Cone K) {
	U = K.vertex - (Sphere.radius*K.sinReciprocal)*K.axis;
	D = S.center - U;
	dsqr = Dot(D,D);
	e = Dot(K.axis,D);
	if ( e > 0 and e*e >= dsqr*K.cosSqr ) {
		D = S.center - K.vertex;
		dsqr = Dot(D,D);
		e = -Dot(K.axis,D);
		if ( e > 0 and e*e >= dsqr*K.sinSqr )
			return dsqr <= S.radiusSqr;
		else
			return true;
	}
	return false;
}

K.sinReciprocal = precomputed cone sin reciprocal
K.cosSqr        = precomputed cone cos squared
K.sinSqr        = precomputed cone sin squared
S.radiusSqr     = precomputed sphere radius squared

http://www.geometrictools.com/Documentation/IntersectionSphereCone.pdf
*/

func sphereIntersectsCone(s *Sphere, k *Cone) bool {
	var U, D vectormath.Vector3

	vectormath.V3ScalarMul(&U, &k.axis, float32(s.radius*k.sinReciprocal))
	vectormath.V3Sub(&U, &k.vertex, &U)
	vectormath.V3Sub(&D, &s.center, &U)
	dsqr := float64(vectormath.V3Dot(&D, &D))
	e := float64(vectormath.V3Dot(&k.axis, &D))
	if e > 0 && e*e > dsqr*k.cosSqr {
		vectormath.V3Sub(&D, &s.center, &k.vertex)
		dsqr = float64(vectormath.V3Dot(&D, &D))
		e = -float64(vectormath.V3Dot(&k.axis, &D))
		if e > 0 && e*e > dsqr*k.sinSqr {
			return dsqr <= s.radiusSqr
		} else {
			return true
		}
	}

	return false
}

/*
http://mathematica.stackexchange.com/questions/45265/distance-between-two-line-segments-in-3-space

Another way of doing this (http://mathforum.org/library/drmath/view/51980.html) is to find the mutual perpendicular between the two lines using the cross product, converting this to a unit vector, and then using the dot product between that cross product, and any vector going between the two lines. Like this:

newMinDist[{p1_, p2_}, {q1_, q2_}] :=
 Module[
  {u, v, n, w},
  u = p2 - p1;
  v = q2 - q1;
  n = Normalize[Cross[u,v]];
  w = q1 - p1;
  Dot[w,n]
  ]
Again, u and v are vectors headed along the lines. The vector n is a unit vector in the direction of the cross product of u and v (the unit vector normal to both u and v). The w is any vector between the ps and the qs. You could swap any value of 1 or 2 here (for instance q2-p1). Then w.n gives the shortest distance between the lines.

Gives the same result, and Timing shows this method is better than an order of magnitude faster than the previous one.
*/

// segment 1: point p1 to point p2
// segment 2: point q1 to point q2
func distanceBetweenSegments(p1x, p1y, p1z,
	p2x, p2y, p2z,
	q1x, q1y, q1z,
	q2x, q2y, q2z float64) float64 {
	var p1, p2, q1, q2, u, v, n, w vectormath.Vector3
	
	vectormath.V3MakeFromElems(&p1, float32(p1x), float32(p1y), float32(p1z))
	vectormath.V3MakeFromElems(&p2, float32(p2x), float32(p2y), float32(p2z))
	vectormath.V3MakeFromElems(&q1, float32(q1x), float32(q1y), float32(q1z))
	vectormath.V3MakeFromElems(&q2, float32(q2x), float32(q2y), float32(q2z))

	vectormath.V3Sub(&u, &p2, &p1)
	vectormath.V3Sub(&v, &q2, &q1)

	vectormath.V3Cross(&n, &v, &u)
	vectormath.V3Normalize(&n, &n)

	vectormath.V3Sub(&w, &q1, &p1)

	return float64(vectormath.V3Dot(&w, &n))
}

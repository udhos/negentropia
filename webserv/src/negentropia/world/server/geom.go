package server

import (
	"errors"
	"math"

	"github.com/udhos/vectormath"
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
	if CloseToZero(sin) {
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

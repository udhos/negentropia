package server

import (
	"math"
	"testing"
)

func expectHit(t *testing.T, centerX, centerY, centerZ, radius, vertexX, vertexY, vertexZ, axisX, axisY, axisZ, angle float64) {
	var s Sphere
	var k Cone

	setSphere(&s, centerX, centerY, centerZ, radius)
	if err := setCone(&k, vertexX, vertexY, vertexZ, axisX, axisY, axisZ, angle); err != nil {
		t.Errorf("expectHit: error: %v", err)
	}

	if hit := sphereIntersectsCone(&s, &k); !hit {
		t.Errorf("expectHit: miss")
	}
}

func expectMiss(t *testing.T, centerX, centerY, centerZ, radius, vertexX, vertexY, vertexZ, axisX, axisY, axisZ, angle float64) {
	var s Sphere
	var k Cone

	setSphere(&s, centerX, centerY, centerZ, radius)
	if err := setCone(&k, vertexX, vertexY, vertexZ, axisX, axisY, axisZ, angle); err != nil {
		t.Errorf("expectMiss: error: %v", err)
	}

	if hit := sphereIntersectsCone(&s, &k); hit {
		t.Errorf("expectMiss: hit")
	}
}

func TestSphereIntersectsCone(t *testing.T) {
	expectHit(t, 0.0, 100.0, 0, 5.0, 0.0, 0.0, 0.0, 0.0, 1.0, 0.0, 15.0*math.Pi/180.0)
	expectMiss(t, 100.0, 0.0, 0, 5.0, 0.0, 0.0, 0.0, 0.0, 1.0, 0.0, 15.0*math.Pi/180.0)
}

func TestDistanceBetweenLineSegments(t *testing.T) {
	expected := 2.0
	dist := distanceBetweenLineSegments(-1.0, 0.0, 0.0,
		1.0, 0.0, 0.0,
		0.0, 2.0, -1.0,
		0.0, 2.0, 1.0)
	if dist != expected {
		t.Errorf("TestDistanceBetweenSegments: distance=%v expected=%v", dist, expected)
	}
}

func TestCapsuleCapsuleIntersection(t *testing.T) {

	cap1 := Capsule{0.9, -1.0, 0.0, 0.0, 1.0, 0.0, 0.0}
	cap2 := Capsule{0.9, 0.0, 2.0, -1.0, 0.0, 2.0, 1.0}

	if capsuleCapsuleIntersection(cap1, cap2) {
		t.Errorf("TestCapsuleCapsuleIntersection: unexpected intersection")
	}

	cap1 = Capsule{1.1, -1.0, 0.0, 0.0, 1.0, 0.0, 0.0}
	cap2 = Capsule{1.1, 0.0, 2.0, -1.0, 0.0, 2.0, 1.0}

	if !capsuleCapsuleIntersection(cap1, cap2) {
		t.Errorf("TestCapsuleCapsuleIntersection: unexpected separation")
	}

}

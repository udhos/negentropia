package main

import (
	"math"
	//"fmt"
	//"negentropia/world/parser"
	//"strings"
)

type Matrix4 struct {
	data []float32
}

func (m *Matrix4) scale(x, y, z, w float64) {
	x1 := float32(x)
	y1 := float32(y)
	z1 := float32(z)
	w1 := float32(w)

	m.data[0] *= x1
	m.data[1] *= x1
	m.data[2] *= x1
	m.data[3] *= x1

	m.data[4] *= y1
	m.data[5] *= y1
	m.data[6] *= y1
	m.data[7] *= y1

	m.data[8] *= z1
	m.data[9] *= z1
	m.data[10] *= z1
	m.data[11] *= z1

	m.data[12] *= w1
	m.data[13] *= w1
	m.data[14] *= w1
	m.data[15] *= w1
}

func setNullMatrix(perspectiveMatrix *Matrix4) {
	perspectiveMatrix.data = []float32{
		0, 0, 0, 0, // c0
		0, 0, 0, 0, // c1
		0, 0, 0, 0, // c2
		0, 0, 0, 0, // c3
	}
}

func setIdentityMatrix(perspectiveMatrix *Matrix4) {
	perspectiveMatrix.data = []float32{
		1, 0, 0, 0, // c0
		0, 1, 0, 0, // c1
		0, 0, 1, 0, // c2
		0, 0, 0, 1, // c3
	}
}

func setPerspectiveMatrix(perspectiveMatrix *Matrix4, fieldOfViewYRadians, aspectRatio, zNear, zFar float64) {
	height := math.Tan(fieldOfViewYRadians*0.5) * zNear
	width := height * aspectRatio
	setFrustumMatrix(perspectiveMatrix, -width, width, -height, height, zNear, zFar)
}

func setFrustumMatrix(perspectiveMatrix *Matrix4, left, right, bottom, top, near, far float64) {
	two_near := 2.0 * near
	right_minus_left := right - left
	top_minus_bottom := top - bottom
	far_minus_near := far - near

	r0c0 := float32(two_near / right_minus_left)
	r1c1 := float32(two_near / top_minus_bottom)
	r0c2 := float32((right + left) / right_minus_left)
	r1c2 := float32((top + bottom) / top_minus_bottom)
	r2c2 := float32(-(far + near) / far_minus_near)
	r3c2 := float32(-1.0)
	r2c3 := float32(-(two_near * far) / far_minus_near)

	perspectiveMatrix.data = []float32{
		r0c0, 0, 0, 0, // c0
		0, r1c1, 0, 0, // c1
		r0c2, r1c2, r2c2, r3c2, // c2
		0, 0, r2c3, 0, // c3
	}
}

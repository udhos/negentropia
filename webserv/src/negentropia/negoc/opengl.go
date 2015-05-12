package main

import (
	"math"
	//"fmt"
	//"negentropia/world/parser"
	//"strings"
)

type Matrix4 struct {
	data [16]float64
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

	perspectiveMatrix.data = [16]float64{
		two_near / right_minus_left, 0, 0, 0, // c0
		0, two_near / top_minus_bottom, 0, 0, // c1
		(right + left) / right_minus_left, (top + bottom) / top_minus_bottom, -(far + near) / far_minus_near, -1.0, // c2
		0, 0, -(two_near * far) / far_minus_near, 0, // c3
	}
}

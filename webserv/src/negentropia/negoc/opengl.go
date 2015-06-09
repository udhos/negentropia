package main

import (
	"math"
	//"fmt"
	//"negentropia/world/parser"
	//"strings"
	"negentropia/world/util"
)

type Matrix4 struct {
	data []float32
}

// usually set w to 1.0
func (m *Matrix4) translate(tx, ty, tz, tw float64) {
	x := float32(tx)
	y := float32(ty)
	z := float32(tz)
	w := float32(tw)
	t1 := m.data[0]*x + m.data[4]*y + m.data[8]*z + m.data[12]*w
	t2 := m.data[1]*x + m.data[5]*y + m.data[9]*z + m.data[13]*w
	t3 := m.data[2]*x + m.data[6]*y + m.data[10]*z + m.data[14]*w
	t4 := m.data[3]*x + m.data[7]*y + m.data[11]*z + m.data[15]*w
	m.data[12] = t1
	m.data[13] = t2
	m.data[14] = t3
	m.data[15] = t4
}

// usually set w to 1.0
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

func (m *Matrix4) multiply(n *Matrix4) {
	m00 := m.data[0]
	m01 := m.data[4]
	m02 := m.data[8]
	m03 := m.data[12]
	m10 := m.data[1]
	m11 := m.data[5]
	m12 := m.data[9]
	m13 := m.data[13]
	m20 := m.data[2]
	m21 := m.data[6]
	m22 := m.data[10]
	m23 := m.data[14]
	m30 := m.data[3]
	m31 := m.data[7]
	m32 := m.data[11]
	m33 := m.data[15]

	n00 := n.data[0]
	n01 := n.data[4]
	n02 := n.data[8]
	n03 := n.data[12]
	n10 := n.data[1]
	n11 := n.data[5]
	n12 := n.data[9]
	n13 := n.data[13]
	n20 := n.data[2]
	n21 := n.data[6]
	n22 := n.data[10]
	n23 := n.data[14]
	n30 := n.data[3]
	n31 := n.data[7]
	n32 := n.data[11]
	n33 := n.data[15]

	m.data[0] = (m00 * n00) + (m01 * n10) + (m02 * n20) + (m03 * n30)
	m.data[4] = (m00 * n01) + (m01 * n11) + (m02 * n21) + (m03 * n31)
	m.data[8] = (m00 * n02) + (m01 * n12) + (m02 * n22) + (m03 * n32)
	m.data[12] = (m00 * n03) + (m01 * n13) + (m02 * n23) + (m03 * n33)
	m.data[1] = (m10 * n00) + (m11 * n10) + (m12 * n20) + (m13 * n30)
	m.data[5] = (m10 * n01) + (m11 * n11) + (m12 * n21) + (m13 * n31)
	m.data[9] = (m10 * n02) + (m11 * n12) + (m12 * n22) + (m13 * n32)
	m.data[13] = (m10 * n03) + (m11 * n13) + (m12 * n23) + (m13 * n33)
	m.data[2] = (m20 * n00) + (m21 * n10) + (m22 * n20) + (m23 * n30)
	m.data[6] = (m20 * n01) + (m21 * n11) + (m22 * n21) + (m23 * n31)
	m.data[10] = (m20 * n02) + (m21 * n12) + (m22 * n22) + (m23 * n32)
	m.data[14] = (m20 * n03) + (m21 * n13) + (m22 * n23) + (m23 * n33)
	m.data[3] = (m30 * n00) + (m31 * n10) + (m32 * n20) + (m33 * n30)
	m.data[7] = (m30 * n01) + (m31 * n11) + (m32 * n21) + (m33 * n31)
	m.data[11] = (m30 * n02) + (m31 * n12) + (m32 * n22) + (m33 * n32)
	m.data[15] = (m30 * n03) + (m31 * n13) + (m32 * n23) + (m33 * n33)
}

func setNullMatrix(m *Matrix4) {
	m.data = []float32{
		0, 0, 0, 0, // c0
		0, 0, 0, 0, // c1
		0, 0, 0, 0, // c2
		0, 0, 0, 0, // c3
	}
}

func setIdentityMatrix(m *Matrix4) {
	m.data = []float32{
		1, 0, 0, 0, // c0
		0, 1, 0, 0, // c1
		0, 0, 1, 0, // c2
		0, 0, 0, 1, // c3
	}
}

func ortho3(x1, y1, z1, x2, y2, z2 float64) bool {
	return util.CloseToZero(dot3(x1, y1, z1, x2, y2, z2))
}

func cross3(x1, y1, z1, x2, y2, z2 float64) (float64, float64, float64) {
	return y1*z2 - z1*y2, z1*x2 - x1*z2, x1*y2 - y1*x2
}

func dot3(x1, y1, z1, x2, y2, z2 float64) float64 {
	return x1*x2 + y1*y2 + z1*z2
}

func lengthSquared3(x, y, z float64) float64 {
	return x*x + y*y + z*z
}

func length3(x, y, z float64) float64 {
	return math.Sqrt(lengthSquared3(x, y, z))
}

func normalize3(x, y, z float64) (float64, float64, float64) {
	length := length3(x, y, z)
	if length == 0 {
		return x, y, z // ugh
	}
	return x / length, y / length, z / length
}

/*
	null rotation:
	forward = 0 0 -1 // looking towards -Z
	up = 0 1 0       // up direction is +Y
	setRotationMatrix(&rotation, 0, 0, -1, 0, 1, 0)
*/
func setRotationMatrix(rotationMatrix *Matrix4, forwardX, forwardY, forwardZ, upX, upY, upZ float64) {
	setModelMatrix(rotationMatrix, forwardX, forwardY, forwardZ, upX, upY, upZ, 0, 0, 0)
}

/*
	setModelMatrix builds the model matrix.
	Model transformation is also known as "camera" transformation.
	Model transformation is the inverse of the view transformation.
	Common use is to compute object location/orientation into full transformation matrix.

	null model:
	forward = 0 0 -1    // looking towards -Z
	up = 0 1 0          // up direction is +Y
	translation = 0 0 0 // position at origin
	setModelMatrix(&rotation, 0, 0, -1, 0, 1, 0, 0, 0, 0)
*/
func setModelMatrix(modelMatrix *Matrix4, forwardX, forwardY, forwardZ, upX, upY, upZ, tX, tY, tZ float64) {
	rightX, rightY, rightZ := normalize3(cross3(forwardX, forwardY, forwardZ, upX, upY, upZ))

	rX := float32(rightX)
	rY := float32(rightY)
	rZ := float32(rightZ)

	uX := float32(upX)
	uY := float32(upY)
	uZ := float32(upZ)

	bX := -float32(forwardX)
	bY := -float32(forwardY)
	bZ := -float32(forwardZ)

	oX := float32(tX)
	oY := float32(tY)
	oZ := float32(tZ)

	modelMatrix.data = []float32{
		rX, rY, rZ, 0, // c0
		uX, uY, uZ, 0, // c1
		bX, bY, bZ, 0, // c2
		oX, oY, oZ, 1, // c3
	}
}

/*
	setViewMatrix builds the view matrix.
	View transformation is the inverse of the model transformation.
	Common use is to compute camera location/orientation into full transformation matrix.

	null view matrix:
	pos   = 0 0 0
	focus = 0 0 -1
	up    = 0 1 0
	setViewMatrix(&V, 0, 0, 0, 0, 0, -1, 0, 1, 0)
*/
func setViewMatrix(viewMatrix *Matrix4, posX, posY, posZ, focusX, focusY, focusZ, upX, upY, upZ float64) {
	backX, backY, backZ := normalize3(posX-focusX, posY-focusY, posZ-focusZ)
	rightX, rightY, rightZ := normalize3(cross3(upX, upY, upZ, backX, backY, backZ))
	newUpX, newUpY, newUpZ := normalize3(cross3(backX, backY, backZ, rightX, rightY, rightZ))

	rotatedEyeX := -dot3(rightX, rightY, rightZ, posX, posY, posZ)
	rotatedEyeY := -dot3(newUpX, newUpY, newUpZ, posX, posY, posZ)
	rotatedEyeZ := -dot3(backX, backY, backZ, posX, posY, posZ)

	rX := float32(rightX)
	rY := float32(rightY)
	rZ := float32(rightZ)

	uX := float32(newUpX)
	uY := float32(newUpY)
	uZ := float32(newUpZ)

	bX := float32(backX)
	bY := float32(backY)
	bZ := float32(backZ)

	eX := float32(rotatedEyeX)
	eY := float32(rotatedEyeY)
	eZ := float32(rotatedEyeZ)

	viewMatrix.data = []float32{
		rX, uX, bX, 0, // c0
		rY, uY, bY, 0, // c1
		rZ, uZ, bZ, 0, // c2
		eX, eY, eZ, 1, // c3
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

	// row x col in the representation below
	r0c0 := float32(two_near / right_minus_left)
	r1c1 := float32(two_near / top_minus_bottom)
	r0c2 := float32((right + left) / right_minus_left)
	r1c2 := float32((top + bottom) / top_minus_bottom)
	r2c2 := float32(-(far + near) / far_minus_near)
	r2c3 := float32(-(two_near * far) / far_minus_near)

	perspectiveMatrix.data = []float32{
		r0c0, 0, r0c2, 0, // "r0"
		0, r1c1, r1c2, 0, // "r1"
		0, 0, r2c2, r2c3, // "r2"
		0, 0, -1, 0, // "r3"
	}
}

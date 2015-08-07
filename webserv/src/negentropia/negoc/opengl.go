package main

import (
	"math"
	//"fmt"
	//"strings"
	"errors"

	//"negentropia/world/parser"
	"negentropia/world/util"
)

type Matrix4 struct {
	data []float32
}

func (m *Matrix4) malloc() {
	if len(m.data) != 16 {
		m.data = make([]float32, 16, 16)
	}
}

func (m *Matrix4) copyFrom(src *Matrix4) {
	m.malloc()
	copy(m.data, src.data)
}

func (m *Matrix4) invert() {
	m.copyInverseFrom(m)
}

func (m *Matrix4) copyInverseFrom(src *Matrix4) error {

	a00 := src.data[0]
	a01 := src.data[1]
	a02 := src.data[2]
	a03 := src.data[3]
	a10 := src.data[4]
	a11 := src.data[5]
	a12 := src.data[6]
	a13 := src.data[7]
	a20 := src.data[8]
	a21 := src.data[9]
	a22 := src.data[10]
	a23 := src.data[11]
	a30 := src.data[12]
	a31 := src.data[13]
	a32 := src.data[14]
	a33 := src.data[15]

	b00 := a00*a11 - a01*a10
	b01 := a00*a12 - a02*a10
	b02 := a00*a13 - a03*a10
	b03 := a01*a12 - a02*a11
	b04 := a01*a13 - a03*a11
	b05 := a02*a13 - a03*a12
	b06 := a20*a31 - a21*a30
	b07 := a20*a32 - a22*a30
	b08 := a20*a33 - a23*a30
	b09 := a21*a32 - a22*a31
	b10 := a21*a33 - a23*a31
	b11 := a22*a33 - a23*a32

	det := b00*b11 - b01*b10 + b02*b09 + b03*b08 - b04*b07 + b05*b06
	if det == 0.0 {
		m.copyFrom(src)
		return errors.New("copyInverseFrom: null determinant")
	}
	invDet := 1.0 / det

	m.malloc()

	m.data[0] = (a11*b11 - a12*b10 + a13*b09) * invDet
	m.data[1] = (-a01*b11 + a02*b10 - a03*b09) * invDet
	m.data[2] = (a31*b05 - a32*b04 + a33*b03) * invDet
	m.data[3] = (-a21*b05 + a22*b04 - a23*b03) * invDet
	m.data[4] = (-a10*b11 + a12*b08 - a13*b07) * invDet
	m.data[5] = (a00*b11 - a02*b08 + a03*b07) * invDet
	m.data[6] = (-a30*b05 + a32*b02 - a33*b01) * invDet
	m.data[7] = (a20*b05 - a22*b02 + a23*b01) * invDet
	m.data[8] = (a10*b10 - a11*b08 + a13*b06) * invDet
	m.data[9] = (-a00*b10 + a01*b08 - a03*b06) * invDet
	m.data[10] = (a30*b04 - a31*b02 + a33*b00) * invDet
	m.data[11] = (-a20*b04 + a21*b02 - a23*b00) * invDet
	m.data[12] = (-a10*b09 + a11*b07 - a12*b06) * invDet
	m.data[13] = (a00*b09 - a01*b07 + a02*b06) * invDet
	m.data[14] = (-a30*b03 + a31*b01 - a32*b00) * invDet
	m.data[15] = (a20*b03 - a21*b01 + a22*b00) * invDet

	return nil
}

// transform: multiply this matrix [m] by vector [x,y,z,w]
func (m *Matrix4) transform(x, y, z, w float64) (tx, ty, tz, tw float64) {

	m0 := float64(m.data[0])
	m1 := float64(m.data[1])
	m2 := float64(m.data[2])
	m3 := float64(m.data[3])
	m4 := float64(m.data[4])
	m5 := float64(m.data[5])
	m6 := float64(m.data[6])
	m7 := float64(m.data[7])
	m8 := float64(m.data[8])
	m9 := float64(m.data[9])
	m10 := float64(m.data[10])
	m11 := float64(m.data[11])
	m12 := float64(m.data[12])
	m13 := float64(m.data[13])
	m14 := float64(m.data[14])
	m15 := float64(m.data[15])

	tx = m0*x + m4*y + m8*z + m12*w
	ty = m1*x + m5*y + m9*z + m13*w
	tz = m2*x + m6*y + m10*z + m14*w
	tw = m3*x + m7*y + m11*z + m15*w

	return
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

func distance3(x1, y1, z1, x2, y2, z2 float64) float64 {
	return length3(x2-x1, y2-y1, z2-z1)
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
	return x*x + y*y + z*z // dot3(x,y,z,x,y,z)
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

	obj.coord. -> P*V*T*R*U*S -> clip coord -> divide by w -> NDC coord -> viewport transform -> window coord
	P*V*T*R*U*S = full transformation
	P = Perspective
	V = View (inverse of camera) built by setViewMatrix
	T*R = model transformation built by THIS setModelMatrix
	T = Translation
	R = Rotation
	U = Undo Model Local Rotation
	S = Scaling

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
	View transformation is also known as "lookAt" transformation.
	View transformation is the inverse of the model transformation.
	Common use is to compute camera location/orientation into full transformation matrix.

	obj.coord. -> P*V*T*R*U*S -> clip coord -> divide by w -> NDC coord -> viewport transform -> window coord
	P*V*T*R*U*S = full transformation
	P = Perspective
	V = View (inverse of camera) built by THIS setViewMatrix
	T*R = model transformation built by setModelMatrix
	T = Translation
	R = Rotation
	U = Undo Model Local Rotation
	S = Scaling

	null view matrix:
	focus = 0 0 -1
	up    = 0 1 0
	pos   = 0 0 0
	setViewMatrix(&V, 0, 0, -1, 0, 1, 0, 0, 0, 0)
*/
func setViewMatrix(viewMatrix *Matrix4, focusX, focusY, focusZ, upX, upY, upZ, posX, posY, posZ float64) {
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

func setPerspectiveMatrix1(perspectiveMatrix *Matrix4, fieldOfViewYRadians, aspectRatio, zNear, zFar float64) {
	height := math.Tan(fieldOfViewYRadians*0.5) * zNear
	width := height * aspectRatio
	setFrustumMatrix(perspectiveMatrix, -width, width, -height, height, zNear, zFar)
}

func setFrustumMatrix(frustumMatrix *Matrix4, left, right, bottom, top, near, far float64) {
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

	frustumMatrix.data = []float32{
		r0c0, 0, r0c2, 0, // "r0"
		0, r1c1, r1c2, 0, // "r1"
		0, 0, r2c2, r2c3, // "r2"
		0, 0, -1, 0, // "r3"
	}
}

func setPerspectiveMatrix2(perspectiveMatrix *Matrix4, fieldOfViewYRadians, aspectRatio, zNear, zFar float64) {
	f := math.Tan(math.Pi*0.5 - fieldOfViewYRadians*0.5) // = cotan(fieldOfViewYRadians/2)
	rangeInv := 1.0 / (zNear - zFar)

	d0 := float32(f / aspectRatio)
	d5 := float32(f)
	d10 := float32((zNear + zFar) * rangeInv)
	d14 := float32(zNear * zFar * rangeInv * 2.0)

	perspectiveMatrix.data = []float32{
		d0, 0, 0, 0,
		0, d5, 0, 0,
		0, 0, d10, -1,
		0, 0, d14, 0,
	}
}

/*
	camera = includes both the perspective and view transforms

	obj.coord. -> P*V*T*R*U*S -> clip coord -> divide by w -> NDC coord -> viewport transform -> window coord
	P*V*T*R*U*S = full transformation
	P = Perspective
	V = View (inverse of camera) built by setViewMatrix
	T*R = model transformation built by setModelMatrix
	T = Translation
	R = Rotation
	U = Undo Model Local Rotation
	S = Scaling
*/
func unproject(camera *Matrix4, viewportX, viewportWidth, viewportY, viewportHeight, pickX, pickY int, depth float64) (worldX, worldY, worldZ float64, err error) {

	// from screen coordinates to clip coordinates
	pX := (2.0 * float64(pickX-viewportX) / float64(viewportWidth)) - 1.0
	pY := (2.0 * float64(pickY-viewportY) / float64(viewportHeight)) - 1.0
	pZ := 2.0*depth - 1.0

	if pX < -1.0 || pX > 1.0 || pY < -1.0 || pY > 1.0 || pZ < -1.0 || pZ > 1.0 {
		err = errors.New("unproject: pick point outside unit cube")
		return
	}

	// invertedCamera: clip coord -> undo perspective -> undo view -> world coord
	var invertedCamera Matrix4
	invertedCamera.copyInverseFrom(camera)
	vx, vy, vz, vw := invertedCamera.transform(pX, pY, pZ, 1.0)
	if vw == 0.0 {
		err = errors.New("unproject: unprojected pick point with W=0")
		return
	}
	invW := 1.0 / vw
	worldX = vx * invW
	worldY = vy * invW
	worldZ = vz * invW

	return
}

/*
	camera = includes both the perspective and view transforms

	obj.coord. -> P*V*T*R*U*S -> clip coord -> divide by w -> NDC coord -> viewport transform -> window coord
	P*V*T*R*U*S = full transformation
	P = Perspective
	V = View (inverse of camera) built by setViewMatrix
	T*R = model transformation built by setModelMatrix
	T = Translation
	R = Rotation
	U = Undo Model Local Rotation
	S = Scaling
*/
func pickRay(camera *Matrix4, viewportX, viewportWidth, viewportY, viewportHeight, pickX, pickY int) (nearX, nearY, nearZ, farX, farY, farZ float64, err error) {

	nearX, nearY, nearZ, err = unproject(camera, viewportX, viewportWidth, viewportY, viewportHeight, pickX, viewportHeight-pickY, 0.0)
	if err != nil {
		return
	}

	farX, farY, farZ, err = unproject(camera, viewportX, viewportWidth, viewportY, viewportHeight, pickX, viewportHeight-pickY, 1.0)

	return
}

/*
	viewportTransform: map NDC coordinates to viewport coordinates
	viewportX, viewportWidth, viewportY, viewportHeight: viewport
	depthNear, depthFar: depthRange
*/
func viewportTransform(viewportX, viewportWidth, viewportY, viewportHeight int, depthNear, depthFar, ndcX, ndcY, ndcZ float64) (int, int, float64) {
	halfWidth := float64(viewportWidth) / 2.0
	halfHeight := float64(viewportHeight) / 2.0
	vx := round(ndcX*halfWidth+halfWidth) + viewportX
	vy := round(ndcY*halfHeight+halfHeight) + viewportY
	depth := (ndcZ*(depthFar-depthNear) + (depthFar + depthNear)) / 2.0

	return vx, viewportHeight - vy, depth
}

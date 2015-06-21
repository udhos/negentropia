package main

import (
	"fmt"
	"math"
	//"negentropia/world/parser"
	//"strings"
)

type camera struct {
	camFocusX, camFocusY, camFocusZ float64
	camUpX, camUpY, camUpZ          float64
	camPosX, camPosY, camPosZ       float64

	orbitRadius float64
}

/*
func newCamera() *camera {
	return &camera{
		0, 0, -1, // focus
		0, 1, 0,  // up
		0, 0, 0,  // pos
	}
}
*/

func resetCamera(cam *camera) {
	*cam = camera{
		0, 0, -1, // focus
		0, 1, 0, // up
		0, 0, 0, // pos
		0,
	}

	cam.orbitRadius = distance3(cam.camPosX, cam.camPosY, cam.camPosZ, cam.camFocusX, cam.camFocusY, cam.camFocusZ)
}

var camRad = 0.0

func incRad(r, delta float64) float64 {
	const pi2 = 2 * math.Pi
	r += delta
	if r > .999*pi2 {
		r = 0
	}
	return r
}

func loadCameraViewMatrixInto(cam *camera, V *Matrix4) {

	delta := math.Pi / 20
	camRad = incRad(camRad, delta)

	//log(fmt.Sprintf("camera: angle=%v delta=%v", camRad*180/math.Pi, delta*180/math.Pi))

	cos := math.Cos(camRad)
	sin := math.Sin(camRad)

	cam.camPosX, cam.camPosY, cam.camPosZ = cam.orbitRadius*sin, 0, cam.orbitRadius*cos

	setViewMatrix(V, cam.camFocusX, cam.camFocusY, cam.camFocusZ, cam.camUpX, cam.camUpY, cam.camUpZ, cam.camPosX, cam.camPosY, cam.camPosZ)

	//log(fmt.Sprintf("angle=%v delta=%v up=%v,%v,%v view=%v", camUpRad*180/math.Pi, delta*180/math.Pi, camUpX, camUpY, camUpZ, V))
}

func cameraMoveTo(cam *camera, coord []float64) {
	cam.camPosX = coord[0]
	cam.camPosY = coord[1]
	cam.camPosZ = coord[2]

	cam.orbitRadius = distance3(cam.camPosX, cam.camPosY, cam.camPosZ, cam.camFocusX, cam.camFocusY, cam.camFocusZ)

	log(fmt.Sprintf("cameraMoveTo: orbitRadius=%v", cam.orbitRadius))
}

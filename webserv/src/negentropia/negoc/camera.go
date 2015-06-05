package main

import (
	//"fmt"
	"math"
	//"negentropia/world/parser"
	//"strings"
)

type camera struct {
	camPosX, camPosY, camPosZ       float64
	camFocusX, camFocusY, camFocusZ float64
	camUpX, camUpY, camUpZ          float64
}

/*
func newCamera() *camera {
	return &camera{
		0, 0, 0,
		0, 0, -1,
		0, 1, 0,
	}
}
*/

func resetCamera(cam *camera) {
	*cam = camera{
		0, 0, 0,
		0, 0, -1,
		0, 1, 0,
	}
}

var camUpRad = 0.0

func incRad(r, delta float64) float64 {
	const pi2 = 2 * math.Pi
	r += delta
	if r > pi2 {
		r -= pi2
	}
	return r
}

func loadCameraViewMatrixInto(cam *camera, V *Matrix4) {

	delta := 0.0 // math.Pi / 5
	camUpRad = incRad(camUpRad, delta)

	cam.camUpX, cam.camUpY, cam.camUpZ = normalize3(math.Sin(camUpRad), math.Cos(camUpRad), 0)

	setViewMatrix(V, cam.camPosX, cam.camPosY, cam.camPosZ, cam.camFocusX, cam.camFocusY, cam.camFocusZ, cam.camUpX, cam.camUpY, cam.camUpZ)

	//log(fmt.Sprintf("angle=%v delta=%v up=%v,%v,%v view=%v", camUpRad*180/math.Pi, delta*180/math.Pi, camUpX, camUpY, camUpZ, V))
}

func cameraMoveTo(cam *camera, coord []float64) {
	cam.camPosX = coord[0]
	cam.camPosY = coord[1]
	cam.camPosZ = coord[2]
}

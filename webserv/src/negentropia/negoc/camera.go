package main

import (
	"fmt"
	"math"
	//"strings"
	"time"
	//"negentropia/world/parser"
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

	cameraOrbitRadius(cam)
}

func cameraOrbitRadius(cam *camera) {
	cam.orbitRadius = distance3(cam.camPosX, cam.camPosY, cam.camPosZ, cam.camFocusX, cam.camFocusY, cam.camFocusZ)
}

func cameraOrbitFrom(cam *camera, x, y, z float64) {
	cameraMoveTo(cam, []float64{x, y, z})
	cameraOrbitRadius(cam)
	log(fmt.Sprintf("cameraOrbitFrom: %v,%v,%v radius=%v", x, y, z, cam.orbitRadius))
}

const pi2 = 2 * math.Pi

func cameraUpdate(gameInfo *gameState, t time.Time) {
	var camPosX, camPosY, camPosZ float64

	autoRotate := true

	if autoRotate {
		sec := float64(t.Second()) + float64(t.Nanosecond())/1000000000

		turnsPerSec := .1
		camRad := sec * pi2 * turnsPerSec

		cos := math.Cos(camRad)
		sin := math.Sin(camRad)

		camPosX, camPosY, camPosZ = gameInfo.cam.orbitRadius*sin, 0.0, gameInfo.cam.orbitRadius*cos
	} else {
		camPosX, camPosY, camPosZ = gameInfo.cam.camPosX, gameInfo.cam.camPosY, gameInfo.cam.camPosZ
	}

	cameraControlMoveTo(gameInfo, []float64{camPosX, camPosY, camPosZ})
}

func loadCameraViewMatrixInto(gameInfo *gameState, cam *camera, V *Matrix4) {
	setViewMatrix(V, cam.camFocusX, cam.camFocusY, cam.camFocusZ, cam.camUpX, cam.camUpY, cam.camUpZ, cam.camPosX, cam.camPosY, cam.camPosZ)
}

func cameraMoveTo(cam *camera, coord []float64) {
	cam.camPosX = coord[0]
	cam.camPosY = coord[1]
	cam.camPosZ = coord[2]

	//cam.orbitRadius = distance3(cam.camPosX, cam.camPosY, cam.camPosZ, cam.camFocusX, cam.camFocusY, cam.camFocusZ)

	//log(fmt.Sprintf("cameraMoveTo: orbitRadius=%v", cam.orbitRadius))
}

func cameraControlMoveTo(gameInfo *gameState, coord []float64) {
	cameraMoveTo(&gameInfo.cam, coord)
	skyboxFollowCamera(gameInfo)
}

func skyboxFollowCamera(gameInfo *gameState) {
	skyboxMoveTo(gameInfo, []float64{gameInfo.cam.camPosX, gameInfo.cam.camPosY, gameInfo.cam.camPosZ})
}

func skyboxMoveTo(gameInfo *gameState, coord []float64) {
	if gameInfo.debugDraw {
		log("skyboxMoveTo: begin")
	}
	if gameInfo.skybox == nil {
		return
	}
	if len(gameInfo.skybox.modelList) < 1 {
		return
	}
	m, ok := gameInfo.skybox.modelList[0].(*skyboxModel)
	if !ok {
		return
	}
	if len(m.instanceList) < 1 {
		return
	}
	m.instanceList[0].setTranslation(coord[0], coord[1], coord[2])
	if gameInfo.debugDraw {
		log("skyboxMoveTo: end")
	}
}

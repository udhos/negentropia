package main

import (
	"fmt"
)

func pick(gameInfo *gameState, canvasX, canvasY int) {
	var cameraMatrix Matrix4
	cameraMatrix.copyFrom(&gameInfo.pMatrix) // camera = P

	var V Matrix4
	loadCameraViewMatrixInto(gameInfo, &gameInfo.cam, &V)

	cameraMatrix.multiply(&V) // camera = P * V

	nearX, nearY, nearZ, farX, farY, farZ, err := pickRay(&cameraMatrix, 0, gameInfo.viewportWidth, 0, gameInfo.viewportHeight, canvasX, canvasY)

	log(fmt.Sprintf("pick: canvas=%v,%v near=%v,%v,%v far=%v,%v,%v error=%v", canvasX, canvasY, nearX, nearY, nearZ, farX, farY, farZ, err))
}

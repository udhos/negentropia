package main

import (
	"fmt"
)

func debugPick(gameInfo *gameState, cameraMatrix *Matrix4, nearX, nearY, nearZ, farX, farY, farZ float64) {

	cam := &gameInfo.cam

	log(fmt.Sprintf("debugPick: camera: focus=%v,%v,%v up=%v,%v,%v pos=%v,%v,%v",
		cam.camFocusX, cam.camFocusY, cam.camFocusZ,
		cam.camUpX, cam.camUpY, cam.camUpZ,
		cam.camPosX, cam.camPosY, cam.camPosZ))

	clipNearX, clipNearY, clipNearZ, clipNearW := cameraMatrix.transform(nearX, nearY, nearZ, 1)
	ndcNearX, ndcNearY, ndcNearZ := clipNearX/clipNearW, clipNearY/clipNearW, clipNearZ/clipNearW
	log(fmt.Sprintf("debugPick: projected ndcNear=%v,%v,%v", ndcNearX, ndcNearY, ndcNearZ))

	clipFarX, clipFarY, clipFarZ, clipFarW := cameraMatrix.transform(farX, farY, farZ, 1)
	ndcFarX, ndcFarY, ndcFarZ := clipFarX/clipFarW, clipFarY/clipFarW, clipFarZ/clipFarW
	log(fmt.Sprintf("debugPick: projected ndcFar=%v,%v,%v", ndcFarX, ndcFarY, ndcFarZ))

	screenNearX, screenNearY, screenNearDepth := viewportTransform(0, gameInfo.viewportWidth, 0, gameInfo.viewportHeight, 0.0, 1.0, ndcNearX, ndcNearY, ndcNearZ)
	log(fmt.Sprintf("debugPick: screenNear=%v,%v,%v", screenNearX, screenNearY, screenNearDepth))

	screenFarX, screenFarY, screenFarDepth := viewportTransform(0, gameInfo.viewportWidth, 0, gameInfo.viewportHeight, 0.0, 1.0, ndcFarX, ndcFarY, ndcFarZ)
	log(fmt.Sprintf("debugPick: screenFar=%v,%v,%v", screenFarX, screenFarY, screenFarDepth))
}

func pick(gameInfo *gameState, canvasX, canvasY int) {
	var cameraMatrix Matrix4
	cameraMatrix.copyFrom(&gameInfo.pMatrix) // camera = P

	var V Matrix4
	loadCameraViewMatrixInto(gameInfo, &gameInfo.cam, &V)

	cameraMatrix.multiply(&V) // camera = P * V

	nearX, nearY, nearZ, farX, farY, farZ, err := pickRay(&cameraMatrix, 0, gameInfo.viewportWidth, 0, gameInfo.viewportHeight, canvasX, canvasY)

	log(fmt.Sprintf("pick: canvas=%v,%v near=%v,%v,%v far=%v,%v,%v error=%v", canvasX, canvasY, nearX, nearY, nearZ, farX, farY, farZ, err))

	debugPick(gameInfo, &cameraMatrix, nearX, nearY, nearZ, farX, farY, farZ)

	ray := ray{nearX, nearY, nearZ, farX - nearX, farY - nearY, farZ - nearZ}

	pickInstance(gameInfo.shaderList, ray)
}

func pickInstance(shaderList []shader, r ray) {
	for _, s := range shaderList {
		s.pickInstance(r)
	}
}

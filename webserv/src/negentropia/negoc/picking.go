package main

import (
	"fmt"
	"math"

	"github.com/udhos/goglmath"
)

func debugPick(gameInfo *gameState, cameraMatrix *goglmath.Matrix4, nearX, nearY, nearZ, farX, farY, farZ float64) {

	cam := &gameInfo.cam

	log(fmt.Sprintf("debugPick: camera: focus=%v,%v,%v up=%v,%v,%v pos=%v,%v,%v",
		cam.camFocusX, cam.camFocusY, cam.camFocusZ,
		cam.camUpX, cam.camUpY, cam.camUpZ,
		cam.camPosX, cam.camPosY, cam.camPosZ))

	clipNearX, clipNearY, clipNearZ, clipNearW := cameraMatrix.Transform(nearX, nearY, nearZ, 1)
	ndcNearX, ndcNearY, ndcNearZ := clipNearX/clipNearW, clipNearY/clipNearW, clipNearZ/clipNearW
	log(fmt.Sprintf("debugPick: projected ndcNear=%v,%v,%v", ndcNearX, ndcNearY, ndcNearZ))

	clipFarX, clipFarY, clipFarZ, clipFarW := cameraMatrix.Transform(farX, farY, farZ, 1)
	ndcFarX, ndcFarY, ndcFarZ := clipFarX/clipFarW, clipFarY/clipFarW, clipFarZ/clipFarW
	log(fmt.Sprintf("debugPick: projected ndcFar=%v,%v,%v", ndcFarX, ndcFarY, ndcFarZ))

	screenNearX, screenNearY, screenNearDepth := goglmath.ViewportTransform(0, gameInfo.viewportWidth, 0, gameInfo.viewportHeight, 0.0, 1.0, ndcNearX, ndcNearY, ndcNearZ)
	log(fmt.Sprintf("debugPick: screenNear=%v,%v,%v", screenNearX, screenNearY, screenNearDepth))

	screenFarX, screenFarY, screenFarDepth := goglmath.ViewportTransform(0, gameInfo.viewportWidth, 0, gameInfo.viewportHeight, 0.0, 1.0, ndcFarX, ndcFarY, ndcFarZ)
	log(fmt.Sprintf("debugPick: screenFar=%v,%v,%v", screenFarX, screenFarY, screenFarDepth))
}

func pick(gameInfo *gameState, canvasX, canvasY int) {
	var cameraMatrix goglmath.Matrix4
	cameraMatrix.CopyFrom(&gameInfo.pMatrix) // camera = P

	var V goglmath.Matrix4
	loadCameraViewMatrixInto(gameInfo, &gameInfo.cam, &V)

	cameraMatrix.Multiply(&V) // camera = P * V

	nearX, nearY, nearZ, farX, farY, farZ, err := goglmath.PickRay(&cameraMatrix, 0, gameInfo.viewportWidth, 0, gameInfo.viewportHeight, canvasX, canvasY)

	log(fmt.Sprintf("pick: canvas=%v,%v near=%v,%v,%v far=%v,%v,%v error=%v", canvasX, canvasY, nearX, nearY, nearZ, farX, farY, farZ, err))

	//debugPick(gameInfo, &cameraMatrix, nearX, nearY, nearZ, farX, farY, farZ)

	ray := ray{nearX, nearY, nearZ, farX - nearX, farY - nearY, farZ - nearZ}

	i := pickInstance(gameInfo.shaderList, ray, gameInfo.cam.camPosX, gameInfo.cam.camPosY, gameInfo.cam.camPosZ)

	log(fmt.Sprintf("pick: found=%v", i))
}

type bestPick struct {
	i               *instance
	distanceSquared float64
}

func pickInstance(shaderList []shader, r ray, camPosX, camPosY, camPosZ float64) *instance {
	closest := &bestPick{nil, math.MaxFloat64}
	for _, s := range shaderList {
		s.pickInstance(r, camPosX, camPosY, camPosZ, closest)
	}
	return closest.i
}

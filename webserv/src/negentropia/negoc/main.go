package main

import (
	"fmt"
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/webgl"
	"math"
	"time"
)

func log(msg string) {
	m := fmt.Sprintf("negoc: %s", msg)
	println(m)
}

var scale = 1.0
var rad = 0.0
var camUpRad = 0.0
var tx = 0.0

const pi2 = 2 * math.Pi

func incRad(r, delta float64) float64 {
	r += delta
	if r > pi2 {
		r -= pi2
	}
	return r
}

const VERTEX_POSITION_ITEM_SIZE = 3 // x,y,z

func draw(gameInfo *gameState, t time.Time, a_Position, vertexIndexSize int, prog, vertexPositionBuffer, vertexIndexBuffer *js.Object) {

	gl := gameInfo.gl

	gl.Clear(gl.COLOR_BUFFER_BIT)

	u_P := gl.GetUniformLocation(prog, "u_P")
	u_MV := gl.GetUniformLocation(prog, "u_MV")

	// scan programs

	gl.UseProgram(prog)
	gl.EnableVertexAttribArray(a_Position)

	uploadPerspective(gl, u_P, &gameInfo.pMatrix)

	// scan models

	gl.BindBuffer(gl.ARRAY_BUFFER, vertexPositionBuffer)
	gl.VertexAttribPointer(a_Position, VERTEX_POSITION_ITEM_SIZE, gl.FLOAT, false, 0, 0)

	// scan instances

	uploadModelView(gl, u_MV, &gameInfo.cam)

	vertexIndexOffset := 0
	vertexIndexElementSize := 2 // uint16

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, vertexIndexBuffer)
	gl.DrawElements(gl.TRIANGLES, vertexIndexSize, gl.UNSIGNED_SHORT, vertexIndexOffset*vertexIndexElementSize)
}

const FRAME_RATE = 2                     // frames per second
const FRAME_INTERVAL = 1000 / FRAME_RATE // msec

func gameLoop(gameInfo *gameState, a_Position, vertexIndexSize int, prog, vertexPositionBuffer, vertexIndexBuffer *js.Object) {
	log(fmt.Sprintf("gameLoop: frame_rate=%v fps frame_interval=%v msec", FRAME_RATE, FRAME_INTERVAL))

	ticker := time.NewTicker(time.Millisecond * FRAME_INTERVAL)
	go func() {
		for t := range ticker.C {
			draw(gameInfo, t, a_Position, vertexIndexSize, prog, vertexPositionBuffer, vertexIndexBuffer)
		}
	}()
}

func testModelView() {
	//func setViewMatrix(viewMatrix *Matrix4, posX, posY, posZ, focusX, focusY, focusZ, upX, upY, upZ float64) {
	//func setModelMatrix(modelMatrix *Matrix4, forwardX, forwardY, forwardZ, upX, upY, upZ, tX, tY, tZ float64) {

	pos := []float64{1, 1, 1}
	focus := []float64{0, 0, -1}
	up := []float64{0, 1, 0}
	var V Matrix4
	setViewMatrix(&V, pos[0], pos[1], pos[2], focus[0], focus[1], focus[2], up[0], up[1], up[2])
	log(fmt.Sprintf("testModelView: view = %v", V))

	forward := []float64{focus[0] - pos[0], focus[1] - pos[1], focus[2] - pos[2]}
	forward[0], forward[1], forward[2] = normalize3(forward[0], forward[1], forward[2])
	rightX, rightY, rightZ := normalize3(cross3(forward[0], forward[1], forward[2], up[0], up[1], up[2]))
	uX, uY, uZ := normalize3(cross3(rightX, rightY, rightZ, forward[0], forward[1], forward[2]))
	var M Matrix4
	setModelMatrix(&M, forward[0], forward[1], forward[2], uX, uY, uZ, pos[0], pos[1], pos[2])
	log(fmt.Sprintf("testModelView: model = %v", M))

	V.multiply(&M)
	log(fmt.Sprintf("testModelView: model x view = %v", V))
}

func testRotation() {
	fx := 0.0 //math.Sin(rad)
	fy := 0.0 //math.Cos(rad)
	fz := -1.0
	ux := 0.0 //math.Sin(up)
	uy := 1.0 //math.Cos(up)
	uz := 0.0
	log(fmt.Sprintf("forward=%v,%v,%v up=%v,%v,%v", fx, fy, fz, ux, uy, uz))

	var rotation Matrix4
	setRotationMatrix(&rotation, fx, fy, fz, ux, uy, uz)
	log(fmt.Sprintf("rotation = %v", rotation))
}

func testView() {
	var V Matrix4
	setViewMatrix(&V, 0, 0, 0, 0, 0, -1, 0, 1, 0)
	log(fmt.Sprintf("testView: view = %v", V))
}

type gameState struct {
	gl                 *webgl.Context
	sock               *gameWebsocket
	defaultTextureUnit int
	pMatrix            Matrix4 // perspective matrix
	canvasAspect       float64
	cam                camera
}

var gameInfo *gameState = &gameState{defaultTextureUnit: 0}

func main() {
	log("main: begin")

	gameInfo.gl = initGL()
	if gameInfo.gl == nil {
		log("main: no webgl context, exiting")
		return
	}

	gl := gameInfo.gl // shortcut

	log("main: WebGL context initialized")

	resetCamera(&gameInfo.cam)

	if initWebSocket(gameInfo) {
		log("main: could not initalize web socket, exiting")
		return
	}

	prog := newShaderProgram(gl)

	attr := "a_Position"
	a_Position := gl.GetAttribLocation(prog, attr)
	if a_Position < 0 {
		log(fmt.Sprintf("main: could not get attribute location: %s", attr))
		return
	}

	log(fmt.Sprintf("main: attribute %s=%v", attr, a_Position))

	// create buffer
	vertexIndexBuffer := gl.CreateBuffer()
	vertexPositionBuffer := gl.CreateBuffer()

	// fill buffer

	indices := []uint16{0, 1, 2} // 3 vertices
	vertexIndexSize := len(indices)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, vertexIndexBuffer)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, indices, gl.STATIC_DRAW)

	// triangle vertices
	vertexPositionData := []float32{
		.9, -.9, -.9, // v0
		0, .9, -.9, // v1
		-.9, -.9, -.9, // v2
	}
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexPositionBuffer)
	gl.BufferData(gl.ARRAY_BUFFER, vertexPositionData, gl.STATIC_DRAW)

	initContext(gameInfo) // set aspectRatio

	setPerspective(gameInfo) // requires aspectRatio

	gameLoop(gameInfo, a_Position, vertexIndexSize, prog, vertexPositionBuffer, vertexIndexBuffer)

	log("main: end")

	//testModelView()
	//testRotation()
	//testView()
}

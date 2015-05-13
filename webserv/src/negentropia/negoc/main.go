package main

import (
	"fmt"
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/webgl"
	"honnef.co/go/js/dom"
	"math"
	"time"
)

func log(msg string) {
	m := fmt.Sprintf("negoc: %s", msg)
	println(m)
}

func setViewport(gl *webgl.Context, w, h int) float64 {
	canvas := gl.Get("canvas")

	/*
	   canvas.width, canvas.height = size you requested the canvas's drawingBuffer to be
	   gl.drawingBufferWidth, gl.drawingBufferHeight = size you actually got.
	   canvas.clientWidth, canvas.clientHeight = size the browser is displaying your canvas.
	*/

	canvas.Set("width", w)
	canvas.Set("height", h)

	style := canvas.Get("style")
	style.Set("width", "${w}px")
	style.Set("height", "${h}px")

	drawingBufferWidth := gl.Object.Get("drawingBufferWidth").Int()
	drawingBufferHeight := gl.Object.Get("drawingBufferHeight").Int()

	gl.BindFramebuffer(gl.FRAMEBUFFER, nil) // on-screen framebuffer
	gl.Viewport(0, 0, drawingBufferWidth, drawingBufferHeight)

	// canvasAspect: save aspect for render loop perspective matrix
	canvasAspect := float64(drawingBufferWidth) / float64(drawingBufferHeight)

	log(fmt.Sprintf("setViewport: %v x %v aspect=%v", drawingBufferWidth, drawingBufferHeight, canvasAspect))

	return canvasAspect
}

func initGL() *webgl.Context {

	document := js.Global.Get("document")
	//body := document.Get("body")

	el := dom.GetWindow().Document().QuerySelector("#canvasbox")
	log(fmt.Sprintf("initGL: #canvasbox el=%v", el))
	canvasbox := el.Underlying()

	canvas := document.Call("createElement", "canvas")

	canvasbox.Call("appendChild", canvas)

	attrs := webgl.DefaultAttributes()
	//attrs.Alpha = false
	gl, err := webgl.NewContext(canvas, attrs)
	if err != nil {
		log(err.Error())
		return nil
	}

	return gl
}

func uploadPerspective(gl *webgl.Context, u_P *js.Object, P *Matrix4) {
	// send perspective projection matrix uniform
	//gl.UniformMatrix4fv(u_P, false, P.data)
	var tmp Matrix4
	setIdentityMatrix(&tmp)
	gl.UniformMatrix4fv(u_P, false, tmp.data)
}

var scale = 1.0

func uploadModelView(gl *webgl.Context, u_MV *js.Object) {

	/*
	   V = View (inverse of camera matrix -- translation and rotation)
	   T = Translation
	   R = Rotation
	   U = Undo Model Local Rotation
	   S = Scaling
	*/
	/*
	   cam.loadViewMatrixInto(MV); // MV = V

	   MV.translate(_center[0], _center[1], _center[2]); // MV = V*T

	   MV.multiply(_rotation); // MV = V*T*R*U

	   MV.scale(rescale, rescale, rescale, 1.0); // MV = V*T*R*U*S
	*/

	var MV Matrix4
	setIdentityMatrix(&MV)

	if scale -= .1; scale < 0 {
		scale = 1.0
	}

	MV.scale(scale, scale, scale, 1.0) // MV = V*T*R*U*S

	gl.UniformMatrix4fv(u_MV, false, MV.data)
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

	uploadModelView(gl, u_MV)

	vertexIndexOffset := 0
	vertexIndexElementSize := 2 // uint16

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, vertexIndexBuffer)
	gl.DrawElements(gl.TRIANGLES, vertexIndexSize, gl.UNSIGNED_SHORT, vertexIndexOffset*vertexIndexElementSize)
}

const FRAME_RATE = 1                     // frames per second
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

func requestZone(sock *gameWebsocket) {
	sock.write(&ClientMsg{Code: CM_CODE_REQZ})
}

func updateCulling(gl *webgl.Context, backfaceCulling bool) {
	if backfaceCulling {
		log("backface culling: ON")
		gl.FrontFace(gl.CCW)
		gl.CullFace(gl.BACK)
		gl.Enable(gl.CULL_FACE)
		return
	}

	log("backface culling: OFF")
	gl.Disable(gl.CULL_FACE)
}

func updateViewport(gameInfo *gameState) {
	gameInfo.canvasAspect = setViewport(gameInfo.gl, 600, 400)
}

func enableExtensions() {
	log("enableExtensions: WRITEME detect and enable WebGL extensions")
}

func initContext(gameInfo *gameState) {
	enableExtensions()

	gl := gameInfo.gl

	gl.ClearColor(0.8, 0.3, 0.01, 1)

	gl.Enable(gl.DEPTH_TEST) // enable depth testing
	gl.DepthFunc(gl.LESS)    // gl.LESS is default depth test
	gl.DepthRange(0.0, 1.0)  // default

	updateViewport(gameInfo)

	updateCulling(gl, false)

	// set default texture unit
	gl.ActiveTexture(gl.TEXTURE0 + gameInfo.defaultTextureUnit)

	requestZone(gameInfo.sock)
}

func setPerspective(gameInfo *gameState) {

	fieldOfViewYRadians := 45 * math.Pi / 180
	planeNear := 2.0   // 2m
	planeFar := 5000.0 // 5km

	setPerspectiveMatrix(&gameInfo.pMatrix, fieldOfViewYRadians, gameInfo.canvasAspect, planeNear, planeFar)
}

type gameState struct {
	gl                 *webgl.Context
	sock               *gameWebsocket
	defaultTextureUnit int
	pMatrix            Matrix4 // perspective matrix
	canvasAspect       float64
}

var gameInfo *gameState = &gameState{defaultTextureUnit: 0}

func main() {
	log("main: begin")

	gl := initGL()
	if gl == nil {
		log("main: no webgl context, exiting")
		return
	}

	gameInfo.gl = gl

	log("main: WebGL context initialized")

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
		.7, -.7, 0, // v0
		0, .7, 0, // v1
		-.7, -.7, 0, // v2
	}
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexPositionBuffer)
	gl.BufferData(gl.ARRAY_BUFFER, vertexPositionData, gl.STATIC_DRAW)

	initContext(gameInfo) // set aspectRatio

	setPerspective(gameInfo) // requires aspectRatio

	gameLoop(gameInfo, a_Position, vertexIndexSize, prog, vertexPositionBuffer, vertexIndexBuffer)

	log("main: end")
}

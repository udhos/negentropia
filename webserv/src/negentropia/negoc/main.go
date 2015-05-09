package main

import (
	"fmt"
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/webgl"
	"honnef.co/go/js/dom"
	"time"
)

func log(msg string) {
	m := fmt.Sprintf("negoc: %s", msg)
	println(m)
}

func setViewport(gl *webgl.Context, w, h int) float32 {
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
	canvasAspect := float32(drawingBufferWidth) / float32(drawingBufferHeight)

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

const VERTEX_POSITION_ITEM_SIZE = 3 // x,y,z

func draw(gl *webgl.Context, t time.Time, a_Position, vertexIndexSize int, prog, vertexPositionBuffer, vertexIndexBuffer *js.Object) {
	gl.Clear(gl.COLOR_BUFFER_BIT)

	gl.UseProgram(prog)
	gl.EnableVertexAttribArray(a_Position)

	gl.BindBuffer(gl.ARRAY_BUFFER, vertexPositionBuffer)
	gl.VertexAttribPointer(a_Position, VERTEX_POSITION_ITEM_SIZE, gl.FLOAT, false, 0, 0)

	vertexIndexOffset := 0
	vertexIndexElementSize := 2 // uint16

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, vertexIndexBuffer)
	gl.DrawElements(gl.TRIANGLES, vertexIndexSize, gl.UNSIGNED_SHORT, vertexIndexOffset*vertexIndexElementSize)
}

const FRAME_RATE = 1                     // frames per second
const FRAME_INTERVAL = 1000 / FRAME_RATE // msec

func gameLoop(gl *webgl.Context, a_Position, vertexIndexSize int, prog, vertexPositionBuffer, vertexIndexBuffer *js.Object) {
	log(fmt.Sprintf("entering game loop frame_rate=%v frame_interval=%v", FRAME_RATE, FRAME_INTERVAL))

	ticker := time.NewTicker(time.Millisecond * FRAME_INTERVAL)
	go func() {
		for t := range ticker.C {
			draw(gl, t, a_Position, vertexIndexSize, prog, vertexPositionBuffer, vertexIndexBuffer)
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

func initContext(gameInfo *gameState) {
	log("initContext: WRITEME")

	/*
	   enable_extensions(gl);

	   clearColor(gl, 0.5, 0.5, 0.5, 1.0);
	   gl.enable(RenderingContext.DEPTH_TEST); // enable depth testing
	   gl.depthFunc(RenderingContext.LESS); // gl.LESS is default depth test
	   gl.depthRange(0.0, 1.0); // default

	   setViewport(gl, gl.canvas.width, gl.canvas.height);

	   updateCulling(gl);

	   // set default texture unit
	   gl.activeTexture(RenderingContext.TEXTURE0 + defaultTextureUnit);
	*/

	gl := gameInfo.gl

	gl.ClearColor(0.8, 0.3, 0.01, 1)

	gl.Enable(gl.DEPTH_TEST) // enable depth testing
	gl.DepthFunc(gl.LESS)    // gl.LESS is default depth test
	gl.DepthRange(0.0, 1.0)  // default

	setViewport(gl, 700, 400)

	updateCulling(gl, gameInfo.backfaceCulling)

	// set default texture unit
	gl.ActiveTexture(gl.TEXTURE0 + gameInfo.defaultTextureUnit)

	requestZone(gameInfo.sock)
}

func setPerspective() {
	// aspect = canvas.width / canvas.height
	//setPerspectiveMatrix(pMatrix, fieldOfViewYRadians, canvasAspect, planeNear, planeFar)

	log("setPerspective: WRITEME")
}

type gameState struct {
	gl                 *webgl.Context
	sock               *gameWebsocket
	backfaceCulling    bool
	defaultTextureUnit int
}

var gameInfo *gameState = &gameState{backfaceCulling: true, defaultTextureUnit: 0}

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

	setPerspective() // requires aspectRatio

	gameLoop(gl, a_Position, vertexIndexSize, prog, vertexPositionBuffer, vertexIndexBuffer)

	log("main: end")
}

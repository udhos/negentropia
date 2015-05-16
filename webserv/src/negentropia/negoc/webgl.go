package main

import (
	"fmt"
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/webgl"
	"honnef.co/go/js/dom"
	"math"
)

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

func uploadPerspective(gl *webgl.Context, u_P *js.Object, P *Matrix4) {
	gl.UniformMatrix4fv(u_P, false, P.data)
}

func loadCameraViewMatrixInto(V *Matrix4) {

	delta := 0.0 // math.Pi / 5
	camUpRad = incRad(camUpRad, delta)
	camUpX, camUpY, camUpZ := normalize3(math.Sin(camUpRad), math.Cos(camUpRad), 0)

	setViewMatrix(V, 0, 0, 0, 0, 0, -1, camUpX, camUpY, camUpZ)

	//log(fmt.Sprintf("angle=%v delta=%v up=%v,%v,%v view=%v", camUpRad*180/math.Pi, delta*180/math.Pi, camUpX, camUpY, camUpZ, V))
}

func uploadModelView(gl *webgl.Context, u_MV *js.Object) {

	/*
	   V = View (inverse of camera matrix -- translation and rotation)
	   T = Translation
	   R = Rotation
	   U = Undo Model Local Rotation
	   S = Scaling

	   MV = V*T*R*U*S
	*/

	// cam.loadViewMatrixInto(MV); // MV = V
	var MV Matrix4
	loadCameraViewMatrixInto(&MV)

	tx += 0.02
	if tx > .5 {
		tx = 0
	}
	MV.translate(tx, 0, 0, 1.0) // MV = V*T

	//rad = incRad(rad, math.Pi/5)
	upX, upY, upZ := normalize3(math.Sin(rad), math.Cos(rad), 0)
	var rotation Matrix4
	setRotationMatrix(&rotation, 0, 0, -1, upX, upY, upZ)
	MV.multiply(&rotation) // MV = V*T*R*U

	//scale -= .1
	if scale < 0 {
		scale = 1.0
	}
	MV.scale(scale, scale, scale, 1.0) // MV = V*T*R*U*S

	gl.UniformMatrix4fv(u_MV, false, MV.data)
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

func setPerspective(gameInfo *gameState) {

	fieldOfViewYRadians := 45 * math.Pi / 180
	planeNear := 2.0   // 2m
	planeFar := 5000.0 // 5km

	setPerspectiveMatrix(&gameInfo.pMatrix, fieldOfViewYRadians, gameInfo.canvasAspect, planeNear, planeFar)

	//log(fmt.Sprintf("perspective: %v", gameInfo.pMatrix))
}

func enableExtensions() {
	log("enableExtensions: WRITEME detect and enable WebGL extensions")
}

func updateViewport(gameInfo *gameState) {
	gameInfo.canvasAspect = setViewport(gameInfo.gl, 600, 400)
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

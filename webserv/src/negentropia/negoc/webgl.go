package main

import (
	"fmt"
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/webgl"
	//"honnef.co/go/js/dom"
	"math"
	"time"
)

func initGL() (*webgl.Context, *js.Object) {

	document := js.Global.Get("document")
	//body := document.Get("body")

	//el := dom.GetWindow().Document().QuerySelector("#canvasbox")
	canvasid := "#canvasbox"
	el := docQuery(canvasid)
	log(fmt.Sprintf("initGL: %s el=%v", canvasid, el))
	canvasbox := el.Underlying()

	canvas := document.Call("createElement", "canvas")

	canvasbox.Call("appendChild", canvas)

	attrs := webgl.DefaultAttributes()
	//attrs.Alpha = false
	gl, err := webgl.NewContext(canvas, attrs)
	if err != nil {
		log(err.Error())
		return nil, nil
	}

	return gl, canvas
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

func updateCulling(gl *webgl.Context, backfaceCulling bool) {

	log(fmt.Sprintf("updateCulling: backfaceCulling=%v", backfaceCulling))

	if backfaceCulling {
		gl.FrontFace(gl.CCW)
		gl.CullFace(gl.BACK)
		gl.Enable(gl.CULL_FACE)
		return
	}

	gl.Disable(gl.CULL_FACE)
}

func setPerspective(gameInfo *gameState) {

	fieldOfViewYRadians := 30 * math.Pi / 180
	planeNear := 2.0   // 2m
	planeFar := 5000.0 // 5km

	setPerspectiveMatrix(&gameInfo.pMatrix, fieldOfViewYRadians, gameInfo.canvasAspect, planeNear, planeFar)

	//log(fmt.Sprintf("perspective: %v", gameInfo.pMatrix))
}

func enableExtensionUintIndex(gameInfo *gameState) {
	extName := "OES_element_index_uint"
	ext := gameInfo.gl.GetExtension(extName)
	gameInfo.extensionUintIndexEnabled = ext != nil
	log(fmt.Sprintf("enableExtensionUintIndex: support for extension [%s] is [%v]", extName, gameInfo.extensionUintIndexEnabled))
}

func enableExtensions(gameInfo *gameState) {
	enableExtensionUintIndex(gameInfo)
}

func updateViewport(gameInfo *gameState) {
	gameInfo.canvasAspect = setViewport(gameInfo.gl, 600, 400)
}

func initContext(gameInfo *gameState) {
	enableExtensions(gameInfo)

	gl := gameInfo.gl

	gl.ClearColor(0.8, 0.3, 0.01, 1)

	gl.Enable(gl.DEPTH_TEST) // enable depth testing
	gl.DepthFunc(gl.LEQUAL)  // gl.LESS is default depth test
	gl.DepthRange(0.0, 1.0)  // default

	updateViewport(gameInfo)

	updateCulling(gl, false)

	// set default texture unit
	gl.ActiveTexture(gl.TEXTURE0 + gameInfo.defaultTextureUnit)

	for {
		if gameInfo.sock != nil {
			requestZone(gameInfo.sock)
			break
		}
		log("initContext: websocket is not ready")
		time.Sleep(time.Second * 1)
	}
}

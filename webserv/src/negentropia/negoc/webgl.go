package main

import (
	"fmt"
	"math"
	"time"

	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/webgl"
	"honnef.co/go/js/dom"
)

func initGL() (*webgl.Context, *js.Object) {

	document := js.Global.Get("document")
	//body := document.Get("body")

	//el := dom.GetWindow().Document().QuerySelector("#canvasbox")
	boxid := "#canvasbox"
	el := docQuery(boxid)
	log(fmt.Sprintf("initGL: %s el=%v", boxid, el))
	canvasbox := el.Underlying()

	canvas := document.Call("createElement", "canvas")
	canvas.Set("id", "main_canvas") // CSS style is attached to id #main_canvas

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

func getCanvasSize(gl *webgl.Context) (int, int) {
	canvas := gl.Get("canvas")
	w := canvas.Get("width").Int()
	h := canvas.Get("height").Int()

	cw := canvas.Get("clientWidth").Int()
	ch := canvas.Get("clientHeight").Int()

	sw := canvas.Get("scrollWidth").Int()
	sh := canvas.Get("scrollHeight").Int()

	el := dom.WrapElement(canvas)

	rect := el.GetBoundingClientRect()
	rl := rect.Left
	rt := rect.Top
	rw := rect.Width
	rh := rect.Height

	style := windowGetComputedStyle(el)
	stw := style.Get("width").Int()
	sth := style.Get("height").Int()

	log(fmt.Sprintf("getCanvasSize: canvas=%dx%d client=%dx%d scroll=%dx%d rect=(%fx%f)%fx%f style=%dx%d", w, h, cw, ch, sw, sh, rl, rt, rw, rh, stw, sth))

	return cw, ch
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

	getCanvasSize(gl)

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

	setPerspectiveMatrix1(&gameInfo.pMatrix, fieldOfViewYRadians, gameInfo.canvasAspect, planeNear, planeFar)
	setPerspectiveMatrix2(&gameInfo.pMatrix, fieldOfViewYRadians, gameInfo.canvasAspect, planeNear, planeFar)

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
	gameInfo.canvasAspect = setViewport(gameInfo.gl, gameInfo.viewportWidth, gameInfo.viewportHeight)
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

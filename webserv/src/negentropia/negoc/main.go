package main

import (
	"fmt"
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/webgl"
	"honnef.co/go/js/dom"
	"io/ioutil"
	"net/http"
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

	setViewport(gl, 700, 400)

	ticker := time.NewTicker(time.Millisecond * FRAME_INTERVAL)
	go func() {
		for t := range ticker.C {
			draw(gl, t, a_Position, vertexIndexSize, prog, vertexPositionBuffer, vertexIndexBuffer)
		}
	}()
}

var vertShaderSrc = `
attribute vec3 a_Position;
 
void main(void) {
	gl_Position = vec4(a_Position, 1.0);
}
`

var fragShaderSrc = `
precision mediump float; // required

void main(void) {
	gl_FragColor = vec4(0.95, 0.95, .95, 1.0); // white opaque
}
`

func compileShader(gl *webgl.Context, shaderSource string, shaderType int) *js.Object {
	shader := gl.CreateShader(shaderType)

	gl.ShaderSource(shader, shaderSource)
	gl.CompileShader(shader)
	parameter := gl.GetShaderParameterb(shader, gl.COMPILE_STATUS)
	//log(fmt.Sprintf("shader parameter=%v", parameter))
	if !parameter {
		infoLog := gl.GetShaderInfoLog(shader)
		log(fmt.Sprintf("compileShader error: infoLog=%v", infoLog))
		return nil
	}

	return shader
}

func httpFetch(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("httpFetch: get url=%v: %v", url, err)
	}
	defer resp.Body.Close()

	var info []byte
	info, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("httpFetch: read all: url=%v: %v", url, err)
	}

	return info, nil
}

func newShaderProgram(gl *webgl.Context) *js.Object {

	vertShaderURL := "/shader/clip_vs.txt"
	fragShaderURL := "/shader/clip_fs.txt"

	if buf, err := httpFetch(vertShaderURL); err != nil {
		log(fmt.Sprintf("newShaderProgram: fetch url=%v error: %v", vertShaderURL, err))
	} else {
		vertShaderSrc = string(buf[:])
		log(fmt.Sprintf("newShaderProgram: url=%v loaded: %d bytes", vertShaderURL, len(vertShaderSrc)))
	}

	if buf, err := httpFetch(fragShaderURL); err != nil {
		log(fmt.Sprintf("newShaderProgram: fetch url=%v error: %v", fragShaderURL, err))
	} else {
		fragShaderSrc = string(buf[:])
		log(fmt.Sprintf("newShaderProgram: url=%v loaded: %d bytes", fragShaderURL, len(fragShaderSrc)))
	}

	vertShader := compileShader(gl, vertShaderSrc, gl.VERTEX_SHADER)
	if vertShader == nil {
		log("newShaderProgram: failure compiling vertex shader")
		return nil
	}
	fragShader := compileShader(gl, fragShaderSrc, gl.FRAGMENT_SHADER)
	if fragShader == nil {
		log("newShaderProgram: failure compiling fragment shader")
		return nil
	}

	program := gl.CreateProgram()
	gl.AttachShader(program, vertShader)
	gl.AttachShader(program, fragShader)
	gl.LinkProgram(program)
	progParameter := gl.GetProgramParameterb(program, gl.LINK_STATUS)
	//log(fmt.Sprintf("program parameter=%v", progParameter))
	if !progParameter {
		infoLog := gl.GetProgramInfoLog(program)
		log(fmt.Sprintf("newShaderProgram: infoLog=%v", infoLog))
		return nil
	}

	log("newShaderProgram: done")

	return program
}

func initWebSocket() bool {

	query := "#wsUri"

	el := dom.GetWindow().Document().QuerySelector(query)
	if el == nil {
		log(fmt.Sprintf("initWebSocket: could not find element: %s", query))
		return true // error
	}
	//span := el.(dom.HTMLSpanElement)
	log(fmt.Sprintf("initWebSocket: %s el=%v", query, el))
	wsUri := el.TextContent()
	if wsUri == "" {
		log(fmt.Sprintf("initWebSocket: empty text for element: %s", query))
		return true // error
	}

	log(fmt.Sprintf("initWebSocket: %s wsUri=%v", query, wsUri))

	return false // ok
}

func main() {
	log("main: Hello world, console")
	gl := initGL()
	if gl == nil {
		log("main: no webgl context, exiting")
		return
	}

	if initWebSocket() {
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

	gl.ClearColor(0.8, 0.3, 0.01, 1)

	gameLoop(gl, a_Position, vertexIndexSize, prog, vertexPositionBuffer, vertexIndexBuffer)
}

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

func initGL() *webgl.Context {

	document := js.Global.Get("document")
	//body := document.Get("body")

	el := dom.GetWindow().Document().QuerySelector("#canvasbox")
	log(fmt.Sprintf("el=%v", el))
	canvasbox := el.Underlying()

	canvas := document.Call("createElement", "canvas")

	canvasbox.Call("appendChild", canvas)

	attrs := webgl.DefaultAttributes()
	attrs.Alpha = false

	gl, err := webgl.NewContext(canvas, attrs)
	if err != nil {
		log(err.Error())
		return nil
	}

	return gl
}

func draw(gl *webgl.Context, t time.Time) {
	gl.Clear(gl.COLOR_BUFFER_BIT)
}

const FRAME_RATE = .5                    // frames per second
const FRAME_INTERVAL = 1000 / FRAME_RATE // msec

func gameLoop(gl *webgl.Context) {
	log(fmt.Sprintf("entering game loop frame_rate=%v frame_interval=%v", FRAME_RATE, FRAME_INTERVAL))

	log("entering game loop")

	ticker := time.NewTicker(time.Millisecond * FRAME_INTERVAL)
	go func() {
		for t := range ticker.C {
			draw(gl, t)
		}
	}()
}

const vertShaderSrc = `
attribute vec3 a_Position;
attribute vec3 a_Position_test; // eraseme
 
void main(void) {
	gl_Position = vec4(a_Position + a_Position_test, 1.0);
}
`

const fragShaderSrc = `
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

func newShaderProgram(gl *webgl.Context) *js.Object {
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

func main() {
	log("main: Hello world, console")
	gl := initGL()
	if gl == nil {
		log("main: no webgl context, exiting")
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

	gl.ClearColor(0.8, 0.3, 0.01, 1)

	gameLoop(gl)
}

package main

import (
	"fmt"
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/webgl"
)

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

type Shader interface {
	ShaderName() string
}

type simpleTexturizer struct {
	program  *js.Object
	progName string
}

func (s *simpleTexturizer) ShaderName() string {
	return s.progName
}

func findShader(shaderList []Shader, name string) Shader {
	for _, s := range shaderList {
		if name == s.ShaderName() {
			return s
		}
	}
	return nil
}

func fetchShaderProgram(gameInfo *gameState, programName, vertShader, fragShader string) {

	log(fmt.Sprintf("fetchShaderProgram: prog=%v vert=%v frag=%v", programName, vertShader, fragShader))

	var s Shader
	s = findShader(gameInfo.shaderList, programName)
	if s != nil {
		log(fmt.Sprintf("fetchShaderProgram: shader FOUND: prog=%v vert=%v frag=%v", programName, vertShader, fragShader))
		return
	}

	log(fmt.Sprintf("fetchShaderProgram: shader NOT found: prog=%v vert=%v frag=%v", programName, vertShader, fragShader))
	prog := newShaderProgram(gameInfo.gl, vertShader, fragShader)
	if prog == nil {
		log(fmt.Sprintf("fetchShaderProgram: failure creating shader: prog=%v vert=%v frag=%v", programName, vertShader, fragShader))
		return
	}
	s = &simpleTexturizer{program: prog, progName: programName}
	gameInfo.shaderList = append(gameInfo.shaderList, s)
}

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

func newShaderProgram(gl *webgl.Context, vertShaderURL, fragShaderURL string) *js.Object {

	/*
		vertShaderURL := "/shader/simple_vs.txt"
		fragShaderURL := "/shader/simple_fs.txt"
	*/

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

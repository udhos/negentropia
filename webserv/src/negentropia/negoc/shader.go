package main

import (
	"fmt"
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/webgl"
	//"negentropia/world/obj"
)

/*
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
*/

type shader interface {
	name() string
	init(gl *webgl.Context)
	draw(gameInfo *gameState)
	findModel(name string) *model
	addModel(model *model)
	unif_MV() *js.Object
	unif_Sampler() *js.Object
	attrLoc_Position() int
	attrLoc_TextureCoord() int
}

type simpleTexturizer struct {
	program        *js.Object
	progName       string
	u_P            *js.Object
	u_MV           *js.Object
	u_Sampler      *js.Object
	a_Position     int
	a_TextureCoord int
	modelList      []*model
}

func (s *simpleTexturizer) unif_MV() *js.Object {
	return s.u_MV
}

func (s *simpleTexturizer) unif_Sampler() *js.Object {
	return s.u_Sampler
}

func (s *simpleTexturizer) attrLoc_Position() int {
	return s.a_Position
}

func (s *simpleTexturizer) attrLoc_TextureCoord() int {
	return s.a_TextureCoord
}

func (s *simpleTexturizer) addModel(m *model) {
	s.modelList = append(s.modelList, m)
}

func (s *simpleTexturizer) findModel(name string) *model {
	for _, m := range s.modelList {
		if name == m.name() {
			return m
		}
	}
	return nil
}

func (s *simpleTexturizer) name() string {
	return s.progName
}

func (s *simpleTexturizer) getUniform(gl *webgl.Context, uniform string) *js.Object {
	u := gl.GetUniformLocation(s.program, uniform)
	if u == nil {
		log(fmt.Sprintf("simpleTexturizer.getUniform: could not get uniform location: %s", uniform))
	}
	return u
}

func (s *simpleTexturizer) getAttrib(gl *webgl.Context, attr string) int {
	a := gl.GetAttribLocation(s.program, attr)
	if a < 0 {
		log(fmt.Sprintf("simpleTexturizer.getAttrib: could not get attrib location: %s", attr))
	}
	return a
}

func (s *simpleTexturizer) init(gl *webgl.Context) {
	s.a_Position = s.getAttrib(gl, "a_Position")
	s.a_TextureCoord = s.getAttrib(gl, "a_TextureCoord")

	s.u_P = s.getUniform(gl, "u_P")
	s.u_MV = s.getUniform(gl, "u_MV")
	s.u_Sampler = s.getUniform(gl, "u_Sampler")
}

func (s *simpleTexturizer) draw(gameInfo *gameState) {
	gl := gameInfo.gl // shortcut

	gl.UseProgram(s.program)
	gl.EnableVertexAttribArray(s.a_Position)
	gl.EnableVertexAttribArray(s.a_TextureCoord)

	uploadPerspective(gl, s.u_P, &gameInfo.pMatrix)

	// draw every model
	for _, m := range s.modelList {
		m.draw(gameInfo, s)
	}
}

func findShader(shaderList []shader, name string) shader {
	for _, s := range shaderList {
		if name == s.name() {
			return s
		}
	}
	return nil
}

func fetchShaderProgram(gameInfo *gameState, programName, vertShader, fragShader string) {

	log(fmt.Sprintf("fetchShaderProgram: prog=%v vert=%v frag=%v", programName, vertShader, fragShader))

	s := findShader(gameInfo.shaderList, programName)
	if s != nil {
		log(fmt.Sprintf("fetchShaderProgram: existing shader FOUND: prog=%v vert=%v frag=%v", programName, vertShader, fragShader))
		return
	}

	log(fmt.Sprintf("fetchShaderProgram: will create new shader: prog=%v vert=%v frag=%v", programName, vertShader, fragShader))
	prog := newShaderProgram(gameInfo.gl, vertShader, fragShader)
	if prog == nil {
		log(fmt.Sprintf("fetchShaderProgram: failure creating shader: prog=%v vert=%v frag=%v", programName, vertShader, fragShader))
		return
	}

	// create new shader
	s = &simpleTexturizer{program: prog, progName: programName}
	s.init(gameInfo.gl)
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

	var vertShaderSrc, fragShaderSrc string

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

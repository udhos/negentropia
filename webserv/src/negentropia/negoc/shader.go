package main

import (
	"fmt"
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/webgl"
	"negentropia/world/obj"
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
}

type instance struct {
	instanceName string
}

func (i *instance) name() string {
	return i.instanceName
}

func (i *instance) draw(gameInfo *gameState) {
	// draw instance
}

type model struct {
	modelName    string
	instanceList []*instance
	ready        bool // mesh and textures loaded
}

func newModel(s shader, modelName string, gl *webgl.Context, objURL string,
	front, up []float64, assetPath asset, textureTable map[string]texture, repeatTexture bool) *model {

	// allocate new model
	mod := &model{modelName: modelName, ready: false}

	// fetch model from objURL

	var buf []byte
	var err error

	if buf, err = httpFetch(objURL); err != nil {
		log(fmt.Sprintf("newModel: fetch model from objURL=%s error: %v", objURL, err))
		return nil
	}

	var o *obj.Obj
	if o, err = obj.NewObjFromBuf(buf, func(msg string) { log(fmt.Sprintf("newModel: %s", msg)) }); err != nil {
		log(fmt.Sprintf("newModel: parse error objURL=%s error: %v", objURL, err))
		return nil
	}

	log(fmt.Sprintf("newModel: objURL=%s coordinates: %d", objURL, len(o.Coord)))

	// push new model into shader.modelList
	s.addModel(mod)

	// return new model
	mod.ready = false // FIXME

	return mod
}

func (m *model) draw(gameInfo *gameState) {
	// draw every instance
	for _, i := range m.instanceList {
		i.draw(gameInfo)
	}
}

func (m *model) name() string {
	return m.modelName
}

func (m *model) findInstance(name string) *instance {
	for _, i := range m.instanceList {
		if name == i.name() {
			return i
		}
	}
	return nil
}

type simpleTexturizer struct {
	program    *js.Object
	progName   string
	u_P        *js.Object
	u_MV       *js.Object
	a_Position int
	modelList  []*model
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

func (s *simpleTexturizer) init(gl *webgl.Context) {
	attr := "a_Position"
	s.a_Position = gl.GetAttribLocation(s.program, attr)
	if s.a_Position < 0 {
		log(fmt.Sprintf("simpleTexturizer.init: could not get attribute location: %s", attr))
	}

	s.u_P = s.getUniform(gl, "u_P")
	s.u_MV = s.getUniform(gl, "u_MV")
}

func (s *simpleTexturizer) draw(gameInfo *gameState) {
	gl := gameInfo.gl // shortcut

	gl.UseProgram(s.program)
	gl.EnableVertexAttribArray(s.a_Position)

	// draw every model
	for _, m := range s.modelList {
		if m.ready {
			m.draw(gameInfo)
		}
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

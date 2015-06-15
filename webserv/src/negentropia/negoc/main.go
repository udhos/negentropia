package main

import (
	"fmt"
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/webgl"
	//"math"
	"negentropia/world/obj"
	"time"
)

func log(msg string) {
	m := fmt.Sprintf("negoc: %s", msg)
	println(m)
}

/*
var scale = 1.0
var rad = 0.0
var camUpRad = 0.0
var tx = 0.0

const pi2 = 2 * math.Pi

func incRad(r, delta float64) float64 {
	r += delta
	if r > pi2 {
		r -= pi2
	}
	return r
}
*/

const VERTEX_POSITION_ITEM_SIZE = 3 // x,y,z

func draw(gameInfo *gameState, t time.Time) {

	gl := gameInfo.gl

	gl.BindFramebuffer(gl.FRAMEBUFFER, nil) // select default framebuffer attached to on-screen canvas
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	// scan programs

	for _, s := range gameInfo.shaderList {
		s.draw(gameInfo)
	}

	/*
		gl.UseProgram(prog)
		gl.EnableVertexAttribArray(a_Position)

		u_P := gl.GetUniformLocation(prog, "u_P")
		u_MV := gl.GetUniformLocation(prog, "u_MV")

		uploadPerspective(gl, u_P, &gameInfo.pMatrix)

		// scan models

		gl.BindBuffer(gl.ARRAY_BUFFER, vertexPositionBuffer)
		gl.VertexAttribPointer(a_Position, VERTEX_POSITION_ITEM_SIZE, gl.FLOAT, false, 0, 0)

		// scan instances

		// put triangle at pos[0 0 1] front[0 0 -1] up[0 1 0]
		i := newInstance("builtin", 0, 0, -1, 0, 1, 0, 0, 0, 1, 10)
		i.uploadModelView(gl, u_MV, &gameInfo.cam)

		vertexIndexOffset := 0
		vertexIndexElementSize := 2 // uint16

		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, vertexIndexBuffer)
		gl.DrawElements(gl.TRIANGLES, vertexIndexSize, gl.UNSIGNED_SHORT, vertexIndexOffset*vertexIndexElementSize)
	*/
}

const FRAME_RATE = 2                     // frames per second
const FRAME_INTERVAL = 1000 / FRAME_RATE // msec

func gameLoop(gameInfo *gameState) {
	log(fmt.Sprintf("gameLoop: frame_rate=%v fps frame_interval=%v msec", FRAME_RATE, FRAME_INTERVAL))

	ticker := time.NewTicker(time.Millisecond * FRAME_INTERVAL)
	go func() {
		for t := range ticker.C {
			draw(gameInfo, t)
		}
	}()
}

func testModelView() {

	pos := []float64{1, 1, 1}
	focus := []float64{0, 0, -1}
	up := []float64{0, 1, 0}
	var V Matrix4
	setViewMatrix(&V, focus[0], focus[1], focus[2], up[0], up[1], up[2], pos[0], pos[1], pos[2])
	log(fmt.Sprintf("testModelView: view = %v", V))

	forward := []float64{focus[0] - pos[0], focus[1] - pos[1], focus[2] - pos[2]}
	forward[0], forward[1], forward[2] = normalize3(forward[0], forward[1], forward[2])
	rightX, rightY, rightZ := normalize3(cross3(forward[0], forward[1], forward[2], up[0], up[1], up[2]))
	uX, uY, uZ := normalize3(cross3(rightX, rightY, rightZ, forward[0], forward[1], forward[2]))
	var M Matrix4
	setModelMatrix(&M, forward[0], forward[1], forward[2], uX, uY, uZ, pos[0], pos[1], pos[2])
	log(fmt.Sprintf("testModelView: model = %v", M))

	V.multiply(&M)
	log(fmt.Sprintf("testModelView: model x view = %v", V))
}

func testRotation() {
	fx := 0.0 //math.Sin(rad)
	fy := 0.0 //math.Cos(rad)
	fz := -1.0
	ux := 0.0 //math.Sin(up)
	uy := 1.0 //math.Cos(up)
	uz := 0.0
	log(fmt.Sprintf("forward=%v,%v,%v up=%v,%v,%v", fx, fy, fz, ux, uy, uz))

	var rotation Matrix4
	setRotationMatrix(&rotation, fx, fy, fz, ux, uy, uz)
	log(fmt.Sprintf("rotation = %v", rotation))
}

func testView() {
	var V Matrix4
	setViewMatrix(&V, 0, 0, -1, 0, 1, 0, 0, 0, 0)
	log(fmt.Sprintf("testView: view = %v", V))
}

func testModelTRU() {
	ufx := 0.0
	ufy := 0.0
	ufz := -1.0
	uux := 0.0
	uuy := 1.0
	uuz := 0.0

	fx := 1.0
	fy := 0.0
	fz := 0.0
	ux := 0.0
	uy := 0.0
	uz := -1.0
	tx := 1.1
	ty := 2.2
	tz := 3.3

	i1 := instance{}
	i1.undoModelRotationFrom(ufx, ufy, ufz, uux, uuy, uuz) // rotation = U
	i1.setRotationFrom(fx, fy, fz, ux, uy, uz)             // rotation = R*U
	var M1 Matrix4
	setIdentityMatrix(&M1)
	M1.translate(tx, ty, tz, 1) // M1 = T
	M1.multiply(&i1.rotation)   // M1 = T*R*U

	log(fmt.Sprintf("testModelTRU: M1 = %v", M1))

	i2 := instance{}
	i2.undoModelRotationFrom(ufx, ufy, ufz, uux, uuy, uuz) // rotation = U
	i2.setRotation(fx, fy, fz, ux, uy, uz)                 // rotation = T*R*U
	i2.setTranslation(tx, ty, tz)                          // rotation = T*R*U
	var M2 Matrix4
	setIdentityMatrix(&M2)
	M2.multiply(&i2.rotation) // M2 = T*R*U

	log(fmt.Sprintf("testModelTRU: M2 = %v", M2))
}

type gameState struct {
	gl                 *webgl.Context
	sock               *gameWebsocket
	defaultTextureUnit int
	pMatrix            Matrix4 // perspective matrix
	canvasAspect       float64
	cam                camera
	shaderList         []shader
	textureTable       map[string]*texture
	assetPath          asset
	materialLib        obj.MaterialLib
}

var gameInfo *gameState = &gameState{defaultTextureUnit: 0}

func main() {
	log("main: begin")

	gameInfo.gl = initGL()
	if gameInfo.gl == nil {
		log("main: no webgl context, exiting")
		return
	}

	log("main: WebGL context initialized")

	resetCamera(&gameInfo.cam)

	gameInfo.assetPath.setRoot("/")

	if initWebSocket(gameInfo) {
		log("main: could not initalize web socket, exiting")
		return
	}

	gameInfo.textureTable = map[string]*texture{}
	gameInfo.materialLib = obj.NewMaterialLib()

	// BEGIN: useless code
	eraseme(gameInfo.gl, "/shader/simple_vs.txt", "/shader/simple_fs.txt")
	// END: useless code

	initContext(gameInfo) // set aspectRatio

	setPerspective(gameInfo) // requires aspectRatio

	gameLoop(gameInfo)

	log("main: end")

	//testModelView()
	//testRotation()
	//testView()
	//testModelTRU()

	var eraseme *js.Object = nil
	if eraseme == nil {
		// y u do dis spoderman?
	}
}

func eraseme(gl *webgl.Context, vertShaderURL, fragShaderURL string) *js.Object {

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
	if !progParameter {
		infoLog := gl.GetProgramInfoLog(program)
		log(fmt.Sprintf("newShaderProgram: infoLog=%v", infoLog))
		return nil
	}

	log("newShaderProgram: done")

	return program
}

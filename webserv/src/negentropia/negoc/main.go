package main

import (
	"fmt"
	//"math"
	"time"

	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/webgl"

	"negentropia/world/obj"
)

func log(msg string) {
	m := fmt.Sprintf("negoc: %s", msg)
	println(m)
}

func update(gameInfo *gameState, t time.Time) {
	cameraUpdate(gameInfo, t)
}

func draw(gameInfo *gameState, t time.Time) {

	gl := gameInfo.gl

	gl.BindFramebuffer(gl.FRAMEBUFFER, nil) // select default framebuffer attached to on-screen canvas
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	for _, s := range gameInfo.shaderList {
		s.draw(gameInfo)
	}

	if box := gameInfo.box; box != nil {
		box.draw(gameInfo)
	}

	if skybox := gameInfo.skybox; skybox != nil {
		skybox.draw(gameInfo)
	}
}

/*
const FRAME_RATE = 10                    // frames per second
const FRAME_INTERVAL = 1000 / FRAME_RATE // msec
*/

/*
func gameLoop(gameInfo *gameState) {
	log(fmt.Sprintf("gameLoop: frame_rate=%v fps frame_interval=%v msec", FRAME_RATE, FRAME_INTERVAL))

	ticker := time.NewTicker(time.Millisecond * FRAME_INTERVAL)
	go func() {
		for t := range ticker.C {
			if gameInfo.debugDraw {
				log(fmt.Sprintf("gameLoop: exiting due debugDraw=%v", gameInfo.debugDraw))
				break
			}
			update(gameInfo, t)
			draw(gameInfo, t)
		}
	}()
}
*/

func gameLoopStart(gameInfo *gameState) {
	var tick func(timestamp float32)

	tick = func(timestamp float32) {
		if !gameInfo.debugDraw {
			gameInfo.animFrameId = requestAnimationFrame(tick) // re-schedule
		}
		t := time.Now()
		update(gameInfo, t)
		draw(gameInfo, t)
	}

	gameInfo.animFrameId = requestAnimationFrame(tick) // schedule
}

func gameLoopStop(gameInfo *gameState) {
	cancelAnimationFrame(gameInfo.animFrameId)
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

func testIntersectRaySphere() {
	s := sphere{0, 0, 0, 1}
	r1 := ray{5, -1, 0, 0, 1, 0}
	r2 := ray{1, -1, 0, 0, 1, 0}
	r3 := ray{.5, -1, 0, 0, 1, 0}
	hit1, t1a, t1b := intersectRaySphere(r1, s)
	hit2, t2a, t2b := intersectRaySphere(r2, s)
	hit3, t3a, t3b := intersectRaySphere(r3, s)
	p1ax, p1ay, p1az := r1.getPoint(t1a)
	p1bx, p1by, p1bz := r1.getPoint(t1b)
	p2ax, p2ay, p2az := r2.getPoint(t2a)
	p2bx, p2by, p2bz := r2.getPoint(t2b)
	p3ax, p3ay, p3az := r3.getPoint(t3a)
	p3bx, p3by, p3bz := r3.getPoint(t3b)
	log(fmt.Sprintf("testIntersectRaySphere: ray 1: expected=MISS hit=%v A=%v,%v,%v B=%v,%v,%v", hit1, p1ax, p1ay, p1az, p1bx, p1by, p1bz))
	log(fmt.Sprintf("testIntersectRaySphere: ray 2: expected=HIT  hit=%v A=%v,%v,%v B=%v,%v,%v", hit2, p2ax, p2ay, p2az, p2bx, p2by, p2bz))
	log(fmt.Sprintf("testIntersectRaySphere: ray 3: expected=HIT  hit=%v A=%v,%v,%v B=%v,%v,%v", hit3, p3ax, p3ay, p3az, p3bx, p3by, p3bz))
}

type gameState struct {
	gl                        *webgl.Context
	sock                      *gameWebsocket
	defaultTextureUnit        int
	pMatrix                   Matrix4 // perspective matrix
	canvasAspect              float64
	cam                       camera
	shaderList                []shader
	textureTable              map[string]*texture
	assetPath                 asset
	materialLib               obj.MaterialLib
	kb                        keyboard
	extensionUintIndexEnabled bool
	skybox                    *skyboxShader
	box                       *boxdemo
	debugDraw                 bool
	animFrameId               int
	canvas                    *js.Object
	viewportWidth             int
	viewportHeight            int
}

func resetGame(gameInfo *gameState) {
	gameInfo.materialLib = obj.NewMaterialLib()
	gameInfo.shaderList = []shader{}              // drop existing shaders
	gameInfo.textureTable = map[string]*texture{} // drop existing texture table
}

func main() {
	log("main: begin")

	gameInfo := &gameState{}

	gameInfo.gl, gameInfo.canvas = initGL()
	if gameInfo.gl == nil {
		log("main: no webgl context, exiting")
		return
	}

	log("main: WebGL context initialized")

	w, h := getCanvasSize(gameInfo.gl)
	gameInfo.viewportWidth = w
	gameInfo.viewportHeight = h
	log(fmt.Sprintf("main: canvas size: %d x %d", gameInfo.viewportWidth, gameInfo.viewportHeight))

	resetCamera(&gameInfo.cam)

	gameInfo.assetPath.setRoot("/")

	if initWebSocket(gameInfo) {
		log("main: could not initalize web socket, exiting")
		return
	}

	resetGame(gameInfo)

	initContext(gameInfo) // set aspectRatio

	setPerspective(gameInfo) // requires aspectRatio

	trapKeyboard(gameInfo)
	trapMouse(gameInfo)

	gameLoopStart(gameInfo)

	//gameInfo.box = newBoxdemo(gameInfo)

	log("main: end")

	//testModelView()
	//testRotation()
	//testView()
	//testModelTRU()
	testIntersectRaySphere()
}

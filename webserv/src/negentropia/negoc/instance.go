package main

import (
	"fmt"
	//"math"
	//"negentropia/world/parser"
	//"strings"
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/webgl"
)

type instance struct {
	id                           string
	posX, posY, posZ             float64
	forwardX, forwardY, forwardZ float64
	upX, upY, upZ                float64
	scale                        float64
	undoModelRotation            Matrix4 // U
	rotation                     Matrix4 // R * U
}

// used only when instance is initialized
func (i *instance) undoModelRotationFrom(forwardX, forwardY, forwardZ, upX, upY, upZ float64) {
	setViewMatrix(&i.undoModelRotation, 0, 0, 0, forwardX, forwardY, forwardZ, upX, upY, upZ)
}

func (i *instance) setRotationFrom(forwardX, forwardY, forwardZ, upX, upY, upZ float64) {
	setRotationMatrix(&i.rotation, forwardX, forwardY, forwardZ, upX, upY, upZ) // rotation = R
	i.rotation.multiply(&i.undoModelRotation)                                   // rotation = R * U
}

// update T*R*U
func (i *instance) updateModelMatrix() {
	setModelMatrix(&i.rotation, i.forwardX, i.forwardY, i.forwardZ, i.upX, i.upY, i.upZ, i.posX, i.posY, i.posZ) // rotation = T*R
	i.rotation.multiply(&i.undoModelRotation)                                                                    // rotation = T*R*U
}

func (i *instance) setRotation(forwardX, forwardY, forwardZ, upX, upY, upZ float64) {
	i.forwardX, i.forwardY, i.forwardZ, i.upX, i.upY, i.upZ = forwardX, forwardY, forwardZ, upX, upY, upZ
	i.updateModelMatrix() // rotation = T*R*U
}

func (i *instance) setTranslation(x, y, z float64) {
	i.posX, i.posY, i.posZ = x, y, z
	i.updateModelMatrix() // rotation = T*R*U
}

func (i *instance) draw(gameInfo *gameState, mod *model) {
	// scan model groups
	for i, g := range mod.mesh.Groups {
		t := mod.textures[i]
		if t.texture == nil {
			continue // texture not ready
		}

		// draw group here

		if g.IndexBegin > g.IndexCount {
			// bogus usage of g to make go compiler happy
		}
	}
}

func (i *instance) uploadModelView(gl *webgl.Context, u_MV *js.Object, cam *camera) {

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
	loadCameraViewMatrixInto(cam, &MV)

	/*
		tx += 0.02
		if tx > .5 {
			tx = 0
		}
	*/
	MV.translate(i.posX, i.posY, i.posZ, 1) // MV = V*T

	//rad = incRad(rad, math.Pi/5)
	/*
		upX, upY, upZ := normalize3(math.Sin(rad), math.Cos(rad), 0)
		var rotation Matrix4
		setRotationMatrix(&rotation, 0, 0, -1, upX, upY, upZ)
		MV.multiply(&rotation) // MV = V*T*R*U
	*/
	MV.multiply(&i.rotation) // MV = V*T*R*U

	/*
		//scale -= .1
		if scale < 0 {
			scale = 1.0
		}
	*/
	MV.scale(i.scale, i.scale, i.scale, 1.0) // MV = V*T*R*U*S

	gl.UniformMatrix4fv(u_MV, false, MV.data)
}

func createInstance(gameInfo *gameState, tab map[string]string) {

	var ok bool
	var err error
	var id string

	if id, ok = tab["id"]; !ok {
		log(fmt.Sprintf("createInstance: missing id"))
		return
	}

	var front string

	if front, ok = tab["modelFront"]; !ok {
		log(fmt.Sprintf("createInstance: id=%s missing modelFront", id))
		return
	}

	var f []float64

	if f, err = parseVector3(front); err != nil {
		log(fmt.Sprintf("createInstance: id=%s bad modelFront=%v: error: %v", id, front, err))
		return
	}

	var up string

	if up, ok = tab["modelUp"]; !ok {
		log(fmt.Sprintf("createInstance: id=%s missing modelUp", id))
		return
	}

	var u []float64

	if u, err = parseVector3(up); err != nil {
		log(fmt.Sprintf("createInstance: id=%s bad modelUp=%v: error: %v", id, up, err))
		return
	}

	if !ortho3(f[0], f[1], f[2], u[0], u[1], u[2]) {
		log(fmt.Sprintf("createInstance: id=%s NOT ORTHOGONAL f=%v u=%v: dot=%f", id, f, u, dot3(f[0], f[1], f[2], u[0], u[1], u[2])))
	}

	var coord string

	if coord, ok = tab["coord"]; !ok {
		log(fmt.Sprintf("createInstance: id=%s missing coord", id))
		return
	}

	var c []float64

	if c, err = parseVector3(coord); err != nil {
		log(fmt.Sprintf("createInstance: id=%s bad coord=%v: error: %v", id, coord, err))
		return
	}

	var programName string

	if programName, ok = tab["programName"]; !ok {
		log(fmt.Sprintf("createInstance: id=%s missing program name", id))
		return
	}

	shader := findShader(gameInfo.shaderList, programName)
	if shader == nil {
		log(fmt.Sprintf("createInstance: id=%s shader programName=%s not found", id, programName))
		return
	}

	var modelName string

	if modelName, ok = tab["obj"]; !ok {
		log(fmt.Sprintf("createInstance: id=%s missing obj", id))
		return
	}

	var objURL string

	if objURL, ok = tab["objURL"]; !ok {
		log(fmt.Sprintf("createInstance: id=%s missing objURL", id))
		return
	}

	repeatTexture := tab["repeatTexture"]
	repeat := stringIsTrue(repeatTexture)

	mod := shader.findModel(modelName)
	if mod == nil {
		log(fmt.Sprintf("createInstance: id=%s program=%s model=%s not found", id, programName, modelName))
		if mod = newModel(shader, modelName, gameInfo.gl, objURL, f, u, gameInfo.assetPath, gameInfo.textureTable, repeat, gameInfo.materialLib); mod == nil {
			log(fmt.Sprintf("createInstance: id=%s program=%s failure creating model=%s", id, programName, modelName))
			return
		}
	}

	inst := mod.findInstance(id)
	if inst != nil {
		log(fmt.Sprintf("createInstance: id=%s model=%s prog=%s ignoring instance redefinition", id, modelName, programName))
		return
	}

	// WRITEME: create instance of model

	log(fmt.Sprintf("createInstance: id=%s prog=%s coord=%v f=%v u=%v WRITEME", id, programName, c, f, u))

}

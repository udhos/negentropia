package main

import (
	"fmt"
	//"math"
	"strconv"
	//"strings"

	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/webgl"
	"github.com/udhos/goglmath"
)

type instance struct {
	id                           string
	posX, posY, posZ             float64
	forwardX, forwardY, forwardZ float64
	upX, upY, upZ                float64
	scale                        float64
	picking                      bool
	undoModelRotation            goglmath.Matrix4 // U
	rotation                     goglmath.Matrix4 // R * U
	modelBoundingRadius          float64
}

func (i *instance) boundingRadius() float64 {
	return i.scale * i.modelBoundingRadius
}

func newInstanceNull(id string) *instance {
	return newInstance(id, 0, 0, 0, -1, 0, 1, 0, 0, 0, 0, 1, false)
}

func newInstance(id string, modelRadius, modelForwardX, modelForwardY, modelForwardZ, modelUpX, modelUpY, modelUpZ, posX, posY, posZ, scale float64, picking bool) *instance {
	i := &instance{id: id, scale: scale, picking: picking}

	i.forwardX, i.forwardY, i.forwardZ = goglmath.Normalize3(modelForwardX, modelForwardY, modelForwardZ)
	i.upX, i.upY, i.upZ = goglmath.Normalize3(modelUpX, modelUpY, modelUpZ)
	i.posX, i.posY, i.posZ = posX, posY, posZ

	// U: undo model implicit rotation
	// R: apply instance-specific rotation
	// Initially, before R is modified by instance specific rotation, U=inverse(R), R*U=I
	i.undoModelRotationFrom(modelForwardX, modelForwardY, modelForwardZ, modelUpX, modelUpY, modelUpZ) // setup U
	i.updateModelMatrix()                                                                              // rotation = T*R*U

	i.modelBoundingRadius = modelRadius

	log(fmt.Sprintf("newInstance: instance=%s boundingRadius=%v picking=%v", id, i.boundingRadius(), i.picking))

	return i
}

// called only when instance is initialized
func (i *instance) undoModelRotationFrom(forwardX, forwardY, forwardZ, upX, upY, upZ float64) {
	// focus = translation + forward
	// in object-space coordinates: translation = 0, focus = forward
	goglmath.SetViewMatrix(&i.undoModelRotation, forwardX, forwardY, forwardZ, upX, upY, upZ, 0, 0, 0)
}

func (i *instance) setRotationFrom(forwardX, forwardY, forwardZ, upX, upY, upZ float64) {
	goglmath.SetRotationMatrix(&i.rotation, forwardX, forwardY, forwardZ, upX, upY, upZ) // rotation = R
	i.rotation.Multiply(&i.undoModelRotation)                                            // rotation = R * U
}

// update T*R*U
func (i *instance) updateModelMatrix() {
	goglmath.SetModelMatrix(&i.rotation, i.forwardX, i.forwardY, i.forwardZ, i.upX, i.upY, i.upZ, i.posX, i.posY, i.posZ) // rotation = T*R
	i.rotation.Multiply(&i.undoModelRotation)                                                                             // rotation = T*R*U
}

func (i *instance) setRotation(forwardX, forwardY, forwardZ, upX, upY, upZ float64) {
	i.forwardX, i.forwardY, i.forwardZ, i.upX, i.upY, i.upZ = forwardX, forwardY, forwardZ, upX, upY, upZ
	i.updateModelMatrix() // rotation = T*R*U
}

func (i *instance) setTranslation(x, y, z float64) {
	i.posX, i.posY, i.posZ = x, y, z
	i.updateModelMatrix() // rotation = T*R*U
}

/*
func (i *instance) draw(gameInfo *gameState, mod *model, u_MV, u_Sampler *js.Object) {

	gl := gameInfo.gl

	i.uploadModelView(gl, u_MV, &gameInfo.cam)

	// scan model groups
	for i, g := range mod.mesh.Groups {
		t := mod.textures[i]
		if t == nil {
			continue // skip group because texture is not ready
		}

		gl.BindTexture(gl.TEXTURE_2D, t.texture)

		// set sampler to use texture assigned to unit
		gl.Uniform1i(u_Sampler, gameInfo.defaultTextureUnit)

		gl.DrawElements(gl.TRIANGLES, g.IndexCount,
			mod.vertexIndexElementType,
			g.IndexBegin*mod.vertexIndexElementSize)
	}
}
*/

func (i *instance) uploadModelView(gameInfo *gameState, gl *webgl.Context, u_MV *js.Object, cam *camera) {

	/*
	   P = Perspective
	   V = View (inverse of camera matrix)
	   T = Translation
	   R = Rotation
	   U = Undo Model Local Rotation
	   S = Scaling

	   MV = V*T*R*U*S
	   Full transform: P * MV
	*/

	var MV goglmath.Matrix4
	loadCameraViewMatrixInto(gameInfo, cam, &MV) // MV = V

	//MV.translate(i.posX, i.posY, i.posZ, 1) // MV = V*T

	MV.Multiply(&i.rotation) // MV = V*T*R*U (rotation = T*R*U)

	MV.Scale(i.scale, i.scale, i.scale, 1.0) // MV = V*T*R*U*S

	gl.UniformMatrix4fv(u_MV, false, MV.Data())
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

	if !goglmath.Ortho3(f[0], f[1], f[2], u[0], u[1], u[2]) {
		log(fmt.Sprintf("createInstance: id=%s NOT ORTHOGONAL f=%v u=%v: dot=%f", id, f, u, goglmath.Dot3(f[0], f[1], f[2], u[0], u[1], u[2])))
	}

	var coord string

	if coord, ok = tab["coord"]; !ok {
		log(fmt.Sprintf("createInstance: id=%s missing coord", id))
		return
	}

	s := 1.0

	if scale, scaleFound := tab["scale"]; scaleFound {
		if v, parseFloatErr := strconv.ParseFloat(scale, 64); parseFloatErr == nil {
			s = v
		} else {
			log(fmt.Sprintf("createInstance: id=%s bad parse float scale=%s: %v", id, scale, parseFloatErr))
		}
	} else {
		log(fmt.Sprintf("createInstance: id=%s missing scale", id))
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

	picking := !stringIsFalse(tab["picking"])

	mod := shader.findModel(modelName)
	if mod == nil {
		log(fmt.Sprintf("createInstance: id=%s program=%s model=%s not found", id, programName, modelName))
		mod = newModel(shader, modelName, gameInfo.gl, objURL, f, u, gameInfo.assetPath, gameInfo.textureTable, repeat, gameInfo.materialLib, gameInfo.extensionUintIndexEnabled)
		if mod == nil {
			log(fmt.Sprintf("createInstance: id=%s program=%s failure creating model=%s", id, programName, modelName))
			return
		}
	}

	var trueNil model
	log(fmt.Sprintf("createInstance: id=%s program=%s model=%s newModel=%v newModelIsNil=%v (trueNil=%v trueNilIsNil=%v)", id, programName, modelName, mod, mod == nil, trueNil, trueNil == nil))

	inst := mod.findInstance(id)
	if inst != nil {
		log(fmt.Sprintf("createInstance: id=%s model=%s prog=%s ignoring instance redefinition", id, modelName, programName))
		return
	}

	inst = newInstance(id, mod.getBoundingRadius(), f[0], f[1], f[2], u[0], u[1], u[2], c[0], c[1], c[2], s, picking)

	log(fmt.Sprintf("createInstance: id=%s prog=%s coord=%v f=%v u=%v scale=%f inst=%v", id, programName, c, f, u, s, inst))

	mod.addInstance(inst)
}

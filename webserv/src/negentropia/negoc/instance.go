package main

import (
	"fmt"
	//"negentropia/world/parser"
	//"strings"
)

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
		if mod = newModel(shader, modelName, gameInfo.gl, objURL, f, u, gameInfo.assetPath, gameInfo.textureTable, repeat); mod == nil {
			log(fmt.Sprintf("createInstance: id=%s program=%s failure creating model=%s", id, programName, modelName))
			return
		}
	}

	// WRITEME: create instance of model

	log(fmt.Sprintf("createInstance: id=%s prog=%s coord=%v f=%v u=%v WRITEME", id, programName, c, f, u))
}

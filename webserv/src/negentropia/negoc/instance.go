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
		log(fmt.Sprintf("createInstance: id=%s bad modelFront=%v", id, front))
		return
	}

	log(fmt.Sprintf("createInstance: id=%s f=%v WRITEME", id, f))

}

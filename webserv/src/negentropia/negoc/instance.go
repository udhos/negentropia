package main

import (
	"fmt"
	//"negentropia/world/parser"
	//"strings"
)

func createInstance(gameInfo *gameState, tab map[string]string) {

	var ok bool
	var id string

	if id, ok = tab["id"]; !ok {
		log(fmt.Sprintf("createInstance: missing id"))
		return
	}

	log(fmt.Sprintf("createInstance: id=%v WRITEME", id))

}

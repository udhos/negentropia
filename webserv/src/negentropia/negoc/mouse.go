package main

import (
	"fmt"

	"honnef.co/go/js/dom"
)

func trapMouse(gameInfo *gameState) {

	wheelHandler := func(ev dom.Event) {
		wheelEv := ev.(*dom.WheelEvent)
		log(fmt.Sprintf("wheelHandler: event=%v", wheelEv))
	}

	mouseUpHandler := func(ev dom.Event) {
		mouseEv := ev.(*dom.MouseEvent)
		log(fmt.Sprintf("mouseUpHandler: event=%v", mouseEv))
	}

	mouseDownHandler := func(ev dom.Event) {
		mouseEv := ev.(*dom.MouseEvent)
		log(fmt.Sprintf("mouseDownHandler: event=%v", mouseEv))
	}

	el := dom.WrapElement(gameInfo.canvas)

	el.AddEventListener("wheel", false, wheelHandler)
	el.AddEventListener("mouseup", false, mouseUpHandler)
	el.AddEventListener("mousedown", false, mouseDownHandler)
}

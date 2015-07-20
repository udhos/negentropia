package main

import (
	"fmt"

	"honnef.co/go/js/dom"
)

func getCanvasCoord(canvas dom.Element, clientX, clientY int) (int, int) {
	rect := canvas.GetBoundingClientRect()

	canvasX := clientX - rect.Left
	canvasY := clientY - rect.Top

	return canvasX, canvasY
}

func trapMouse(gameInfo *gameState) {

	wheelHandler := func(ev dom.Event) {
		wheelEv := ev.(*dom.WheelEvent)
		log(fmt.Sprintf("wheelHandler: event=%v", wheelEv))
	}

	mouseUpHandler := func(ev dom.Event) {
		m := ev.(*dom.MouseEvent)
		el := dom.WrapElement(gameInfo.canvas)
		canvasX, canvasY := getCanvasCoord(el, m.ClientX, m.ClientY)
		log(fmt.Sprintf("mouseUpHandler: %v,%v", canvasX, canvasY))
	}

	mouseDownHandler := func(ev dom.Event) {
		m := ev.(*dom.MouseEvent)
		el := dom.WrapElement(gameInfo.canvas)
		canvasX, canvasY := getCanvasCoord(el, m.ClientX, m.ClientY)
		log(fmt.Sprintf("mouseDownHandler: %v,%v", canvasX, canvasY))
	}

	el := dom.WrapElement(gameInfo.canvas)

	el.AddEventListener("wheel", false, wheelHandler)
	el.AddEventListener("mouseup", false, mouseUpHandler)
	el.AddEventListener("mousedown", false, mouseDownHandler)
}

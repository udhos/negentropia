package main

import (
	"fmt"

	"honnef.co/go/js/dom"
)

func getCanvasCoord(canvas dom.Element, clientX, clientY int) (int, int) {
	rect := canvas.GetBoundingClientRect()

	canvasX := clientX - roundToInt(rect.Left)
	canvasY := clientY - roundToInt(rect.Top)

	return canvasX, canvasY
}

func normalizeWheel(delta float64) int {
	if int(delta*100)%5625 == 0 {
		// IE11: Nx56.25
		return int(delta*10000) / 5625
	}
	d := int(delta)
	if intAbs(d) < 100 {
		// Firefox: Nx3
		return d * 100 / 3
	}
	if d%120 == 0 {
		// IE: Nx120
		return d * 100 / 120
	}
	if d%100 == 0 {
		// Chrome, Opera: Nx100
		return d
	}
	return d // unknown browser
}

func trapMouse(gameInfo *gameState) {

	canvas := dom.WrapElement(gameInfo.canvas)

	wheelHandler := func(ev dom.Event) {
		w := ev.(*dom.WheelEvent)
		w.PreventDefault()
		delta := normalizeWheel(w.DeltaY)
		log(fmt.Sprintf("wheelHandler: deltaY=%v normalized=%v", w.DeltaY, delta))
	}

	mouseUpHandler := func(ev dom.Event) {
		m := ev.(*dom.MouseEvent)
		m.PreventDefault()
		canvasX, canvasY := getCanvasCoord(canvas, m.ClientX, m.ClientY)
		log(fmt.Sprintf("mouseUpHandler: %v,%v", canvasX, canvasY))
	}

	mouseDownHandler := func(ev dom.Event) {
		m := ev.(*dom.MouseEvent)
		m.PreventDefault()
		canvasX, canvasY := getCanvasCoord(canvas, m.ClientX, m.ClientY)
		log(fmt.Sprintf("mouseDownHandler: %v,%v", canvasX, canvasY))
		pick(gameInfo, canvasX, canvasY)
	}

	canvas.AddEventListener("wheel", false, wheelHandler)
	canvas.AddEventListener("mouseup", false, mouseUpHandler)
	canvas.AddEventListener("mousedown", false, mouseDownHandler)
}

package main

import (
	"fmt"
	"honnef.co/go/js/dom"
)

func handleKeyDown(ev dom.Event) {
	kbev := ev.(*dom.KeyboardEvent)

	switch kbev.KeyCode {
	case 90:
		log("handleKeyDown: Z key hit")
	default:
		log(fmt.Sprintf("handleKeyDown: keyCode=%d", kbev.KeyCode))
	}
}

func trapKeyboard() {
	docAddEventListener("keydown", false, handleKeyDown)
}

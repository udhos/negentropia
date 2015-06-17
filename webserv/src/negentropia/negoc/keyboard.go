package main

import (
	"fmt"
	"honnef.co/go/js/dom"
)

func trapKeyboard(gameInfo *gameState) {

	keyHandler := func(ev dom.Event) {
		kbev := ev.(*dom.KeyboardEvent)

		switch kbev.KeyCode {
		case 90:
			log("handleKeyDown: Z key hit: requesting zone switch")
			switchZone(gameInfo.sock)
		default:
			log(fmt.Sprintf("handleKeyDown: keyCode=%d", kbev.KeyCode))
		}

	}

	docAddEventListener("keydown", false, keyHandler)
}

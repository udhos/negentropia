package main

import (
	//"fmt"

	"honnef.co/go/js/dom"
)

type keyboard struct {
	keyDownZ bool // Z key state is down
}

// keyPressedZ is called only once when key state changes from UP to DOWN
func keyPressedZ(gameInfo *gameState) {
	log("keyPressedZ: requesting zone switch")
	switchZone(gameInfo.sock)
}

func trapKeyboard(gameInfo *gameState) {

	keyDownHandler := func(ev dom.Event) {
		kbev := ev.(*dom.KeyboardEvent)

		switch kbev.KeyCode {
		case 90:
			if !gameInfo.kb.keyDownZ {
				gameInfo.kb.keyDownZ = true
				keyPressedZ(gameInfo)
			}
		default:
			//log(fmt.Sprintf("keyDownHandler: keyCode=%d", kbev.KeyCode))
		}

	}

	keyUpHandler := func(ev dom.Event) {
		kbev := ev.(*dom.KeyboardEvent)

		switch kbev.KeyCode {
		case 90:
			if gameInfo.kb.keyDownZ {
				gameInfo.kb.keyDownZ = false
				//keyReleasedZ(gameInfo)
			}
		default:
			//log(fmt.Sprintf("keyUpHandler: keyCode=%d", kbev.KeyCode))
		}

	}

	docAddEventListener("keydown", false, keyDownHandler)
	docAddEventListener("keyup", false, keyUpHandler)
}

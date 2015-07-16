package main

import (
	"fmt"
	"time"

	"honnef.co/go/js/dom"
)

type keyboard struct {
	keyDownD bool // D key state is down
	keyDownZ bool // Z key state is down
}

// keyPressedD is called only once when key state changes from UP to DOWN
func keyPressedD(gameInfo *gameState) {
	if gameInfo.debugDraw {
		now := time.Now()
		log(fmt.Sprintf("keyPressedD: drawing once now=%v", now))
		update(gameInfo, now)
		draw(gameInfo, now)
	} else {
		log("keyPressedD: disabling draw loop -- hit again to draw")
		gameInfo.debugDraw = true
	}
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
		case 68:
			if !gameInfo.kb.keyDownD {
				gameInfo.kb.keyDownD = true
				keyPressedD(gameInfo)
			}
		case 90:
			if !gameInfo.kb.keyDownZ {
				gameInfo.kb.keyDownZ = true
				keyPressedZ(gameInfo)
			}
		default:
			log(fmt.Sprintf("keyDownHandler: keyCode=%d", kbev.KeyCode))
		}

	}

	keyUpHandler := func(ev dom.Event) {
		kbev := ev.(*dom.KeyboardEvent)

		switch kbev.KeyCode {
		case 68:
			if gameInfo.kb.keyDownD {
				gameInfo.kb.keyDownD = false
				//keyReleasedD(gameInfo)
			}
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

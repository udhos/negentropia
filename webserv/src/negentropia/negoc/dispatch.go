package main

import (
	"fmt"
)

func dispatch(gameInfo *gameState, code int, data string, tab map[string]string) {
	log(fmt.Sprintf("dispatch: code=%v data=%v tab=%v", code, data, tab))
}

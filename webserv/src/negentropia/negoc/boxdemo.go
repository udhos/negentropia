package main

import (
	"fmt"
)

type boxdemo struct {
}

func newBoxdemo(gameInfo *gameState) *boxdemo {
	log(fmt.Sprintf("newBoxdemo: FIXME WRITEME"))

	box := &boxdemo{}

	box = nil
	return box
}

func (b *boxdemo) draw(gameInfo *gameState) {
}

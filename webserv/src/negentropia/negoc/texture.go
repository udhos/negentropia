package main

import (
	"fmt"
	//"negentropia/world/parser"
	//"strings"
)

type texture struct {
}

func fetchTexture(textureURL string) (*texture, error) {
	log(fmt.Sprintf("fetchTexture: texture=%s", textureURL))
	return &texture{}, nil
}

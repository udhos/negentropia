package main

import (
	"fmt"
	//"negentropia/world/parser"
	//"strings"
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/webgl"
)

type texture struct {
	texture *js.Object
}

func onLoad(gl *webgl.Context, t *texture, textureURL string) {
	log(fmt.Sprintf("onLoad: texture=%s LOADED", textureURL))

	log(fmt.Sprintf("onLoad: texture=%s upload texture to GPU", textureURL))

	t.texture = gl.CreateTexture() // this marks model.group.texture as done
}

func newImage() *js.Object {
	return js.Global.Call("eval", "(new Image())")
}

func fetchTexture(gl *webgl.Context, textureURL string) (*texture, error) {
	log(fmt.Sprintf("fetchTexture: texture=%s", textureURL))

	image := newImage()

	t := &texture{}

	image.Set("onload", func() {
		go onLoad(gl, t, textureURL)
	})

	image.Set("src", textureURL)

	return t, nil
}

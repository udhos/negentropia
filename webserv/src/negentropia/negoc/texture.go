package main

import (
	"fmt"
	//"negentropia/world/parser"
	//"strings"
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/webgl"
)

type texture struct {
	URL     string
	texture *js.Object
}

func onLoad(gl *webgl.Context, t *texture, textureURL string) {
	log(fmt.Sprintf("onLoad: texture=%s image LOADED", textureURL))

	log(fmt.Sprintf("onLoad: texture=%s FIXME WRITEME upload texture to GPU", textureURL))
}

func newImage() *js.Object {
	return js.Global.Call("eval", "(new Image())")
}

func loadSolidColor(gl *webgl.Context, texture *js.Object, rgba []byte) {

	log(fmt.Sprintf("loadSolidColor: %v", rgba))

	gl.BindTexture(gl.TEXTURE_2D, texture)

	//gl.texImage2D(gl.TEXTURE_2D, 0, gl.RGBA, 1, 1, 0, gl.RGBA, gl.UNSIGNED_BYTE, data);
	gl.Call("texImage2D", gl.TEXTURE_2D, 0, gl.RGBA, 1, 1, 0, gl.RGBA, gl.UNSIGNED_BYTE, rgba)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)

	gl.BindTexture(gl.TEXTURE_2D, nil)
}

func fetchTexture(gl *webgl.Context, textureURL string, tempColor []byte) (*texture, error) {
	log(fmt.Sprintf("fetchTexture: texture=%s", textureURL))

	tex := gl.CreateTexture()
	loadSolidColor(gl, tex, tempColor)
	t := &texture{URL: textureURL}
	t.texture = tex // on return this will mark model.mesh.group.texture as done

	image := newImage()

	image.Set("onload", func() {
		go onLoad(gl, t, textureURL)
	})

	image.Set("src", textureURL)

	return t, nil
}

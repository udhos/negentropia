package main

import (
	"fmt"
	//"strings"

	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/webgl"
	//"negentropia/world/parser"
)

type texture struct {
	URL     string
	texture *js.Object
	image   *js.Object
	wrap    int
}

func isPowerOfTwo(v int) bool {
	return v != 0 && (v&(v-1)) == 0
}

func onLoad(gl *webgl.Context, t *texture, textureURL string) {
	log(fmt.Sprintf("onLoad: texture=%s image LOADED", textureURL))

	gl.BindTexture(gl.TEXTURE_2D, t.texture)

	gl.PixelStorei(gl.UNPACK_FLIP_Y_WEBGL, 1)

	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, gl.RGBA, gl.UNSIGNED_BYTE, t.image)

	// undo flip Y otherwise it could affect other texImage calls
	gl.PixelStorei(gl.UNPACK_FLIP_Y_WEBGL, 0)

	width := t.image.Get("width").Int()
	height := t.image.Get("height").Int()

	mipmap := isPowerOfTwo(width) && isPowerOfTwo(height)

	if mipmap {
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_NEAREST)
		gl.GenerateMipmap(gl.TEXTURE_2D)
	} else {
		log(fmt.Sprintf("onLoad: texture=%s w=%d x h=%d can't enable MIPMAP for NPOT texture", textureURL, width, height))
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	}

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, t.wrap)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, t.wrap)

	//int anisotropy = anisotropic_filtering_enable(gl, textureName)
	log(fmt.Sprintf("onLoad: texture=%s FIXEME WRITEME enable anisotropic filtering", textureURL))

	gl.BindTexture(gl.TEXTURE_2D, nil)
}

func newImage() *js.Object {
	//return js.Global.Call("eval", "(new Image())")
	return js.Global.Get("Image").New()
}

func loadSolidColor(gl *webgl.Context, texture *js.Object, rgba []byte) {

	//log(fmt.Sprintf("loadSolidColor: %v", rgba))

	gl.BindTexture(gl.TEXTURE_2D, texture)

	//gl.texImage2D(gl.TEXTURE_2D, 0, gl.RGBA, 1, 1, 0, gl.RGBA, gl.UNSIGNED_BYTE, data);
	gl.Call("texImage2D", gl.TEXTURE_2D, 0, gl.RGBA, 1, 1, 0, gl.RGBA, gl.UNSIGNED_BYTE, rgba)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)

	gl.BindTexture(gl.TEXTURE_2D, nil)
}

func fetchTexture(gl *webgl.Context, textureURL string, tempColor []byte, wrap int) (*texture, error) {
	log(fmt.Sprintf("fetchTexture: texture=%s", textureURL))

	tex := gl.CreateTexture()
	loadSolidColor(gl, tex, tempColor)
	t := &texture{URL: textureURL,
		wrap:    wrap,
		texture: tex, // on return this will mark model.mesh.group.texture as done
	}

	t.image = newImage()

	t.image.Set("onload", func() {
		go onLoad(gl, t, textureURL)
	})

	t.image.Set("src", textureURL)

	return t, nil
}

package main

import (
	"fmt"
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/webgl"
	"honnef.co/go/js/dom"
	"time"
)

func log(msg string) {
	m := fmt.Sprintf("negoc: %s", msg)
	println(m)
}

func initGL() *webgl.Context {

	document := js.Global.Get("document")
	//body := document.Get("body")

	el := dom.GetWindow().Document().QuerySelector("#canvasbox")
	log(fmt.Sprintf("el=%v", el))
	canvasbox := el.Underlying()

	canvas := document.Call("createElement", "canvas")

	canvasbox.Call("appendChild", canvas)

	attrs := webgl.DefaultAttributes()
	attrs.Alpha = false

	gl, err := webgl.NewContext(canvas, attrs)
	if err != nil {
		log(err.Error())
		return nil
	}

	return gl
}

func draw(gl *webgl.Context, t time.Time) {
	//log(fmt.Sprintf("draw: %v", t))
	gl.ClearColor(0.8, 0.3, 0.01, 1)
	gl.Clear(gl.COLOR_BUFFER_BIT)
}

const FRAME_RATE = .5                    // frames per second
const FRAME_INTERVAL = 1000 / FRAME_RATE // msec

func gameLoop(gl *webgl.Context) {
	log(fmt.Sprintf("entering game loop frame_rate=%v frame_interval=%v", FRAME_RATE, FRAME_INTERVAL))

	log("entering game loop")

	ticker := time.NewTicker(time.Millisecond * FRAME_INTERVAL)
	go func() {
		for t := range ticker.C {
			draw(gl, t)
		}
	}()
}

func main() {
	log("main: Hello world, console")
	gl := initGL()
	if gl == nil {
		log("main: no webgl context, exiting")
		return
	}
	gameLoop(gl)
}

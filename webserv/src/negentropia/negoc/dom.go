package main

import (
	"github.com/gopherjs/gopherjs/js"
	"honnef.co/go/js/dom"
)

func windowGetComputedStyle(el dom.Element) *dom.CSSStyleDeclaration {
	return dom.GetWindow().GetComputedStyle(el, "")
}

func docQuery(query string) dom.Element {
	return dom.GetWindow().Document().QuerySelector(query)
}

func docAddEventListener(event string, useCapture bool, listener func(dom.Event)) {
	dom.GetWindow().Document().AddEventListener(event, useCapture, listener)
}

func requestAnimationFrame(callback func(timestamp float32)) int {
	return js.Global.Call("requestAnimationFrame", callback).Int()
}

func cancelAnimationFrame(id int) {
	js.Global.Call("cancelAnimationFrame", id)
}

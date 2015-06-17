package main

import (
	"honnef.co/go/js/dom"
)

func docQuery(query string) dom.Element {
	return dom.GetWindow().Document().QuerySelector(query)
}

func docAddEventListener(event string, useCapture bool, listener func(dom.Event)) {
	dom.GetWindow().Document().AddEventListener(event, useCapture, listener)
}

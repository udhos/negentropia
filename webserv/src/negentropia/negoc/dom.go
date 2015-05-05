package main

import (
	"honnef.co/go/js/dom"
)

func docQuery(query string) dom.Element {
	return dom.GetWindow().Document().QuerySelector(query)
}

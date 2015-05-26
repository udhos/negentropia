package main

import (
	"fmt"
	//"negentropia/world/parser"
	//"strings"
)

type asset struct {
	mesh    string
	mtl     string
	obj     string
	shader  string
	texture string
}

func (a *asset) setRoot(path string) {
	log(fmt.Sprintf("asset.setRoot: setting asset root to path=%s", path))

	a.mesh = fmt.Sprintf("%smesh", path)
	a.mtl = fmt.Sprintf("%smtl", path)
	a.obj = fmt.Sprintf("%sobj", path)
	a.shader = fmt.Sprintf("%sshader", path)
	a.texture = fmt.Sprintf("%stexture", path)
}

package main

import (
	"fmt"
	//"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/webgl"
	"negentropia/world/obj"
)

type model struct {
	modelName    string
	instanceList []*instance
	ready        bool // mesh and textures loaded
}

func newModel(s shader, modelName string, gl *webgl.Context, objURL string,
	front, up []float64, assetPath asset, textureTable map[string]texture, repeatTexture bool) *model {

	// allocate new model
	mod := &model{modelName: modelName, ready: false}

	// fetch model from objURL

	var buf []byte
	var err error

	// push new model into shader.modelList
	s.addModel(mod)

	if buf, err = httpFetch(objURL); err != nil {
		log(fmt.Sprintf("newModel: fetch model from objURL=%s error: %v", objURL, err))
		return nil
	}

	var o *obj.Obj
	if o, err = obj.NewObjFromBuf(buf, func(msg string) { log(fmt.Sprintf("newModel: %s", msg)) }, nil); err != nil {
		log(fmt.Sprintf("newModel: parse error objURL=%s error: %v", objURL, err))
		return nil
	}

	log(fmt.Sprintf("newModel: objURL=%s elements=%d bigIndex=%v texCoord=%v normCoord=%v", objURL, o.NumberOfElements(), o.BigIndexFound, o.TextCoordFound, o.NormCoordFound))

	log(fmt.Sprintf("newModel: objURL=%s FIXME load OBJ textures", objURL))

	mod.ready = false // FIXME when all model data is loaded (mesh, textures)

	return mod
}

func (m *model) draw(gameInfo *gameState) {
	// draw every instance
	for _, i := range m.instanceList {
		i.draw(gameInfo)
	}
}

func (m *model) name() string {
	return m.modelName
}

func (m *model) findInstance(name string) *instance {
	for _, i := range m.instanceList {
		if name == i.name() {
			return i
		}
	}
	return nil
}

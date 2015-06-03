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

func fetchMaterialLib(materialLib map[string]obj.Material, libURL string) error {
	//var buf []byte

	buf, err := httpFetch(libURL)
	if err != nil {
		return fmt.Errorf("fetchMaterialLib: URL=%s failure: %v", libURL, err)
	}

	opt := &obj.ObjParserOptions{Logger: func(msg string) { log(fmt.Sprintf("fetchMaterialLib: %s", msg)) }}

	var lib map[string]obj.Material
	if lib, err = obj.ReadMaterialLibFromBuf(buf, opt); err != nil {
		return err
	}

	// save new material into lib
	for k, v := range lib {
		if _, found := materialLib[k]; found {
			log(fmt.Sprintf("fetchMaterialLib: mtllib=%s REWRITING material=%s", libURL, k))
		}
		materialLib[k] = v
	}

	return nil
}

func newModel(s shader, modelName string, gl *webgl.Context, objURL string,
	front, up []float64, assetPath asset, textureTable map[string]texture,
	repeatTexture bool, materialLib map[string]obj.Material) *model {

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

	opt := &obj.ObjParserOptions{Logger: func(msg string) { log(fmt.Sprintf("newModel: %s", msg)) }}
	var o *obj.Obj
	if o, err = obj.NewObjFromBuf(buf, opt); err != nil {
		log(fmt.Sprintf("newModel: parse error objURL=%s error: %v", objURL, err))
		return nil
	}

	log(fmt.Sprintf("newModel: objURL=%s elements=%d bigIndex=%v texCoord=%v normCoord=%v", objURL, o.NumberOfElements(), o.BigIndexFound, o.TextCoordFound, o.NormCoordFound))

	libURL := fmt.Sprintf("%s/%s", assetPath.mtl, o.Mtllib)

	for _, g := range o.Groups {
		log(fmt.Sprintf("newModel: objURL=%s group=%s size=%d mtllib=%s consider material=%s", objURL, g.Name, g.IndexCount, o.Mtllib, g.Usemtl))

		var mat obj.Material
		var matOk bool
		if mat, matOk = materialLib[g.Usemtl]; !matOk {
			// material not found -- fetch lib

			log(fmt.Sprintf("newModel: objURL=%s group=%s load mtllib=%s", objURL, g.Name, o.Mtllib))

			if libErr := fetchMaterialLib(materialLib, libURL); libErr == nil {
				log(fmt.Sprintf("newModel: objURL=%s group=%s size=%d mtllib=%s material=%s LIB FAILURE: %v", objURL, g.Name, g.IndexCount, libURL, g.Usemtl, libErr))
				continue // ugh
			}

			if mat, matOk = materialLib[g.Usemtl]; !matOk {
				log(fmt.Sprintf("newModel: objURL=%s group=%s size=%d mtllib=%s MISSING material=%s", objURL, g.Name, g.IndexCount, libURL, g.Usemtl))
				continue // ugh
			}
		}

		log(fmt.Sprintf("newModel: objURL=%s group=%s mtllib=%s usemtl=%s load texture=%s", objURL, g.Name, o.Mtllib, g.Usemtl, mat.Map_Kd))
	}

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

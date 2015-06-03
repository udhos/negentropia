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

func fetchMaterialLib(materialLib obj.MaterialLib, libURL string) error {
	//var buf []byte

	buf, err := httpFetch(libURL)
	if err != nil {
		return fmt.Errorf("fetchMaterialLib: URL=%s failure: %v", libURL, err)
	}

	opt := &obj.ObjParserOptions{Logger: func(msg string) { log(fmt.Sprintf("fetchMaterialLib: %s", msg)) }}

	var lib obj.MaterialLib
	if lib, err = obj.ReadMaterialLibFromBuf(buf, opt); err != nil {
		return err
	}

	// save new material into lib
	for k, v := range lib.Lib {
		if _, found := materialLib.Lib[k]; found {
			log(fmt.Sprintf("fetchMaterialLib: mtllib=%s REWRITING material=%s", libURL, k))
		}
		materialLib.Lib[k] = v
	}

	return nil
}

func newModel(s shader, modelName string, gl *webgl.Context, objURL string,
	front, up []float64, assetPath asset, textureTable map[string]texture,
	repeatTexture bool, materialLib obj.MaterialLib) *model {

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

	// Load textures for groups
	for _, g := range o.Groups {
		//log(fmt.Sprintf("newModel: objURL=%s group=%s size=%d mtllib=%s consider material=%s", objURL, g.Name, g.IndexCount, o.Mtllib, g.Usemtl))

		if g.IndexCount < 3 {
			log(fmt.Sprintf("newModel: objURL=%s group=%s size=%d mtllib=%s bad index list size", objURL, g.Name, g.IndexCount, o.Mtllib))
			continue // skip group missing index list
		}

		if g.Usemtl == "" {
			log(fmt.Sprintf("newModel: objURL=%s group=%s size=%d mtllib=%s missing material name", objURL, g.Name, g.IndexCount, o.Mtllib))
			continue // skip group missing material name
		}

		var mat *obj.Material
		var matOk bool
		if mat, matOk = materialLib.Lib[g.Usemtl]; !matOk {
			// material not found -- fetch lib

			//log(fmt.Sprintf("newModel: objURL=%s group=%s load mtllib=%s", objURL, g.Name, o.Mtllib))

			if libErr := fetchMaterialLib(materialLib, libURL); libErr != nil {
				log(fmt.Sprintf("newModel: objURL=%s group=%s size=%d mtllib=%s material=%s LIB FAILURE: %v", objURL, g.Name, g.IndexCount, libURL, g.Usemtl, libErr))
				continue // ugh
			}

			log(fmt.Sprintf("newModel: objURL=%s group=%s size=%d mtllib=%s LIB LOADED for material=%s", objURL, g.Name, g.IndexCount, libURL, g.Usemtl))

			if mat, matOk = materialLib.Lib[g.Usemtl]; !matOk {
				log(fmt.Sprintf("newModel: objURL=%s group=%s size=%d mtllib=%s MISSING material=%s", objURL, g.Name, g.IndexCount, libURL, g.Usemtl))
				continue // ugh
			}

			//log(fmt.Sprintf("newModel: objURL=%s group=%s size=%d mtllib=%s material=%s MATERIAL LOADED", objURL, g.Name, g.IndexCount, o.Mtllib, g.Usemtl))
		}

		//log(fmt.Sprintf("newModel: objURL=%s group=%s size=%d mtllib=%s material=%s MATERIAL OK", objURL, g.Name, g.IndexCount, o.Mtllib, g.Usemtl))

		log(fmt.Sprintf("newModel: objURL=%s group=%s mtllib=%s usemtl=%s load texture=%s", objURL, g.Name, o.Mtllib, g.Usemtl, mat.Map_Kd))
	}

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

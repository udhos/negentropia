package main

import (
	"fmt"
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/webgl"
	"negentropia/world/obj"
	"sort"
)

type model struct {
	modelName         string
	instanceList      []*instance
	mesh              *obj.Obj
	textures          []*texture
	vertexBuffer      *js.Object
	vertexIndexBuffer *js.Object
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

func addGroupTexture(mod *model, gl *webgl.Context, textureTable map[string]*texture, groupListSize, i int, textureName, textureURL string, tempColor []byte) error {

	//log(fmt.Sprintf("addGroupTexture: index=%d texture=%s", i, textureURL))

	texSize := len(mod.textures)
	if i != texSize {
		err := fmt.Errorf("addGroupTexture: model=%s index=%d texture=%s currentTextureListSize=%d not last texture index", mod.modelName, i, textureURL, texSize)
		log(fmt.Sprintf("%s", err))
		return err
	}

	if i >= groupListSize {
		err := fmt.Errorf("addGroupTexture: model=%s index=%d texture=%s currentGroupListSize=%d texture index beyond last group", mod.modelName, i, textureURL, groupListSize)
		log(fmt.Sprintf("%s", err))
		return err
	}

	var t *texture
	if textureName != "" {
		var ok bool
		if t, ok = textureTable[textureURL]; !ok {
			// texture not found - load it
			var err error
			if t, err = fetchTexture(gl, textureURL, tempColor); err != nil {
				log(fmt.Sprintf("addGroupTexture: %s", err)) // warning only
			}
			textureTable[textureURL] = t
		}
	}
	mod.textures = append(mod.textures, t)

	return nil
}

func addGroupTextureNull(mod *model, groupListSize, i int) error {
	return addGroupTexture(mod, nil, nil, groupListSize, i, "", "", []byte{})
}

type GroupByTextureName struct {
	m *model
}

func textureNameLess(t1, t2 *texture) bool {
	if t1 == nil {
		return t2 != nil
	}
	if t2 == nil {
		return false
	}
	return t1.URL < t2.URL
}

func (m GroupByTextureName) Len() int { return len(m.m.textures) }
func (m GroupByTextureName) Swap(i, j int) {
	g := m.m.mesh.Groups
	t := m.m.textures
	g[i], g[j], t[i], t[j] = g[j], g[i], t[j], t[i]
}
func (m GroupByTextureName) Less(i, j int) bool {
	return textureNameLess(m.m.textures[i], m.m.textures[j])
}

func showGroups(m *model) {
	for i, g := range m.mesh.Groups {
		var textureURL string
		if m.textures[i] != nil {
			textureURL = m.textures[i].URL
		}
		log(fmt.Sprintf("showGroups: %d group=%s texture=%s", i, g.Name, textureURL))
	}
}

func newModel(s shader, modelName string, gl *webgl.Context, objURL string,
	front, up []float64, assetPath asset, textureTable map[string]*texture,
	repeatTexture bool, materialLib obj.MaterialLib) *model {

	// allocate new model
	mod := &model{modelName: modelName}

	// fetch model from objURL

	var buf []byte
	var err error

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

	groupListSize := len(o.Groups)

	// Load textures for groups
	for i, g := range o.Groups {
		//log(fmt.Sprintf("newModel: objURL=%s group=%s size=%d mtllib=%s consider material=%s", objURL, g.Name, g.IndexCount, o.Mtllib, g.Usemtl))

		if g.IndexCount < 3 {
			log(fmt.Sprintf("newModel: objURL=%s group=%s size=%d mtllib=%s bad index list size", objURL, g.Name, g.IndexCount, o.Mtllib))
			if addGroupTextureNull(mod, groupListSize, i) != nil {
				return nil
			}
			continue // skip group missing index list
		}

		if g.Usemtl == "" {
			log(fmt.Sprintf("newModel: objURL=%s group=%s size=%d mtllib=%s missing material name", objURL, g.Name, g.IndexCount, o.Mtllib))
			if addGroupTextureNull(mod, groupListSize, i) != nil {
				return nil
			}
			continue // skip group missing material name
		}

		var mat *obj.Material
		var matOk bool
		if mat, matOk = materialLib.Lib[g.Usemtl]; !matOk {
			// material not found -- fetch lib

			//log(fmt.Sprintf("newModel: objURL=%s group=%s load mtllib=%s", objURL, g.Name, o.Mtllib))

			if libErr := fetchMaterialLib(materialLib, libURL); libErr != nil {
				log(fmt.Sprintf("newModel: objURL=%s group=%s size=%d mtllib=%s material=%s LIB FAILURE: %v", objURL, g.Name, g.IndexCount, libURL, g.Usemtl, libErr))
				if addGroupTextureNull(mod, groupListSize, i) != nil {
					return nil
				}
				continue // ugh
			}

			log(fmt.Sprintf("newModel: objURL=%s group=%s size=%d mtllib=%s LIB LOADED for material=%s", objURL, g.Name, g.IndexCount, libURL, g.Usemtl))

			if mat, matOk = materialLib.Lib[g.Usemtl]; !matOk {
				log(fmt.Sprintf("newModel: objURL=%s group=%s size=%d mtllib=%s MISSING material=%s", objURL, g.Name, g.IndexCount, libURL, g.Usemtl))
				if addGroupTextureNull(mod, groupListSize, i) != nil {
					return nil
				}
				continue // ugh
			}

			//log(fmt.Sprintf("newModel: objURL=%s group=%s size=%d mtllib=%s material=%s MATERIAL LOADED", objURL, g.Name, g.IndexCount, o.Mtllib, g.Usemtl))
		}

		//log(fmt.Sprintf("newModel: objURL=%s group=%s size=%d mtllib=%s material=%s MATERIAL OK", objURL, g.Name, g.IndexCount, o.Mtllib, g.Usemtl))

		//log(fmt.Sprintf("newModel: objURL=%s group=%s mtllib=%s usemtl=%s load texture=%s", objURL, g.Name, o.Mtllib, g.Usemtl, mat.Map_Kd))

		textureURL := fmt.Sprintf("%s/%s", assetPath.texture, mat.Map_Kd)

		r := byte(mat.Kd[0] * 255.0)
		g := byte(mat.Kd[1] * 255.0)
		b := byte(mat.Kd[2] * 255.0)
		tempColor := []byte{r, g, b, 255}

		if addGroupTexture(mod, gl, textureTable, groupListSize, i, mat.Map_Kd, textureURL, tempColor) != nil {
			return nil
		}
	}

	if len(o.Groups) != len(mod.textures) {
		log(fmt.Sprintf("newModel: objURL=%s BAD group/texture count: groups=%d textures=%d", objURL, len(o.Groups), len(mod.textures)))
		return nil
	}

	mod.mesh = o

	//showGroups(mod)
	sort.Sort(GroupByTextureName(GroupByTextureName{mod}))
	//showGroups(mod)

	// buffer for vertex data
	mod.vertexBuffer = gl.CreateBuffer()
	gl.BindBuffer(gl.ARRAY_BUFFER, mod.vertexBuffer)
	gl.BufferData(gl.ARRAY_BUFFER, o.Coord, gl.STATIC_DRAW)

	// buffer for indices
	mod.vertexIndexBuffer = gl.CreateBuffer()
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, mod.vertexIndexBuffer)
	if o.BigIndexFound {
		list := make([]uint32, len(o.Indices))
		for i, v := range o.Indices {
			list[i] = uint32(v)
		}
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, list, gl.STATIC_DRAW)
	} else {
		list := make([]uint16, len(o.Indices))
		for i, v := range o.Indices {
			list[i] = uint16(v)
		}
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, list, gl.STATIC_DRAW)
	}

	// push new model into shader.modelList
	s.addModel(mod)

	return mod
}

const vertexPositionBufferItemSize = 3 // coord x,y,z
const textureCoordBufferItemSize = 2   // coord s,t

func (m *model) draw(gameInfo *gameState, prog shader) {
	gl := gameInfo.gl // shortcut

	gl.BindBuffer(gl.ARRAY_BUFFER, m.vertexBuffer)

	// vertex coord x,y,z
	gl.VertexAttribPointer(prog.attrLoc_Position(),
		vertexPositionBufferItemSize,
		gl.FLOAT, false, m.mesh.StrideSize,
		m.mesh.StrideOffsetPosition)

	// texture coord s,t
	gl.VertexAttribPointer(prog.attrLoc_TextureCoord(),
		textureCoordBufferItemSize,
		gl.FLOAT, false, m.mesh.StrideSize,
		m.mesh.StrideOffsetTexture)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, m.vertexIndexBuffer)

	for _, inst := range m.instanceList {
		inst.draw(gameInfo, m, prog.unif_MV())
	}
}

func (m *model) name() string {
	return m.modelName
}

func (m *model) findInstance(id string) *instance {
	for _, i := range m.instanceList {
		if id == i.id {
			return i
		}
	}
	return nil
}

func (m *model) addInstance(inst *instance) {
	m.instanceList = append(m.instanceList, inst)
	//log(fmt.Sprintf("model.addInstance: model=%s newInstance=%s instances=%d instanceList=%v", m.modelName, inst.id, len(m.instanceList), m.instanceList))
}

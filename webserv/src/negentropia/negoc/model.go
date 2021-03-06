package main

import (
	"fmt"
	"math"
	"sort"

	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/webgl"
	"github.com/udhos/goglmath"
	"github.com/udhos/gwob"
)

type model interface {
	name() string
	findInstance(id string) (*instance, bool)
	addInstance(inst *instance)
	draw(gameInfo *gameState, prog shader)
	getBoundingRadius() float64
	pickInstance(r ray, camPosX, camPosY, camPosZ float64, closest *bestPick)
}

type simpleModel struct {
	modelName              string
	instanceList           []*instance
	mesh                   *gwob.Obj
	vertexBuffer           *js.Object
	vertexIndexBuffer      *js.Object
	vertexIndexElementType int
	vertexIndexElementSize int
	boundingRadius         float64
}

type texturizedModel struct {
	simpleModel
	textures []*texture
}

func (m *simpleModel) getBoundingRadius() float64 {
	return m.boundingRadius
}

func (m *simpleModel) pickInstance(r ray, camPosX, camPosY, camPosZ float64, closest *bestPick) {
	for _, inst := range m.instanceList {
		log(fmt.Sprintf("pickInstance: model=%s instance=%s pick=%v", m.name(), inst.id, inst.picking))
		if !inst.picking {
			continue // not a pickable instance
		}
		sph := sphere{inst.posX, inst.posY, inst.posZ, inst.boundingRadius()}
		hit, t1, t2 := intersectRaySphere(r, sph)
		log(fmt.Sprintf("pickInstance: model=%s instance=%s ray=%v sphere=%v hit=%v", m.name(), inst.id, r, sph, hit))
		if !hit {
			continue
		}

		// test distance from camera to hit point 1
		p1x, p1y, p1z := r.getPoint(t1)
		dist1 := goglmath.DistanceSquared3(p1x, p1y, p1z, camPosX, camPosY, camPosZ)
		if dist1 < closest.distanceSquared {
			closest.distanceSquared = dist1
			closest.i = inst
		}

		// test distance from camera to hit point 2
		p2x, p2y, p2z := r.getPoint(t2)
		dist2 := goglmath.DistanceSquared3(p2x, p2y, p2z, camPosX, camPosY, camPosZ)
		if dist2 < closest.distanceSquared {
			closest.distanceSquared = dist2
			closest.i = inst
		}
	}
}

func fetchMaterialLib(materialLib gwob.MaterialLib, libURL string) error {
	//var buf []byte

	buf, err := httpFetch(libURL)
	if err != nil {
		return fmt.Errorf("fetchMaterialLib: URL=%s failure: %v", libURL, err)
	}

	opt := &gwob.ObjParserOptions{Logger: func(msg string) { log(fmt.Sprintf("fetchMaterialLib: %s", msg)) }}

	var lib gwob.MaterialLib
	if lib, err = gwob.ReadMaterialLibFromBuf(buf, opt); err != nil {
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

func addGroupTexture(mod *texturizedModel, gl *webgl.Context, textureTable map[string]*texture, groupListSize, i int, textureName, textureURL string, tempColor []byte, wrap int) error {

	//log(fmt.Sprintf("addGroupTexture: index=%d texture=%s", i, textureURL))

	texSize := len(mod.textures)
	if i != texSize {
		err := fmt.Errorf("addGroupTexture: model=%s index=%d texture=%s currentTextureListSize=%d not last texture index", mod.name(), i, textureURL, texSize)
		log(fmt.Sprintf("%s", err))
		return err
	}

	if i >= groupListSize {
		err := fmt.Errorf("addGroupTexture: model=%s index=%d texture=%s currentGroupListSize=%d texture index beyond last group", mod.name(), i, textureURL, groupListSize)
		log(fmt.Sprintf("%s", err))
		return err
	}

	var t *texture
	if textureName != "" {
		var ok bool
		if t, ok = textureTable[textureURL]; !ok {
			// texture not found - load it
			var err error
			if t, err = fetchTexture(gl, textureURL, tempColor, wrap); err != nil {
				log(fmt.Sprintf("addGroupTexture: %s", err)) // warning only
			}
			textureTable[textureURL] = t
		}
	}
	mod.textures = append(mod.textures, t)

	return nil
}

func addGroupTextureNull(mod *texturizedModel, gl *webgl.Context, groupListSize, i int) error {
	return addGroupTexture(mod, nil, nil, groupListSize, i, "", "", []byte{}, gl.CLAMP_TO_EDGE)
}

type GroupByTextureName struct {
	m *texturizedModel
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

func showGroups(m *texturizedModel) {
	for i, g := range m.mesh.Groups {
		var textureURL string
		if m.textures[i] != nil {
			textureURL = m.textures[i].URL
		}
		log(fmt.Sprintf("showGroups: %d group=%s texture=%s", i, g.Name, textureURL))
	}
}

func (m *simpleModel) createBuffers(objURL string, gl *webgl.Context, extensionUintIndexEnabled bool) {
	o := m.mesh

	// buffer for vertex data
	m.vertexBuffer = gl.CreateBuffer()
	gl.BindBuffer(gl.ARRAY_BUFFER, m.vertexBuffer)
	gl.BufferData(gl.ARRAY_BUFFER, o.Coord, gl.STATIC_DRAW)

	// buffer for indices
	m.vertexIndexBuffer = gl.CreateBuffer()
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, m.vertexIndexBuffer)
	if o.BigIndexFound && extensionUintIndexEnabled {

		list := make([]uint32, len(o.Indices))
		for i, v := range o.Indices {
			list[i] = uint32(v)
		}
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, list, gl.STATIC_DRAW)
		m.vertexIndexElementType = gl.UNSIGNED_INT
		m.vertexIndexElementSize = 4
	} else {
		if o.BigIndexFound && !extensionUintIndexEnabled {
			log(fmt.Sprintf("createBuffers: objURL=%s BigIndexFound BUT WebGL extension missing", objURL))
		}

		list := make([]uint16, len(o.Indices))
		for i, v := range o.Indices {
			list[i] = uint16(v)
		}
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, list, gl.STATIC_DRAW)
		m.vertexIndexElementType = gl.UNSIGNED_SHORT
		m.vertexIndexElementSize = 2
	}

	log(fmt.Sprintf("createBuffers: objURL=%s bigIndexFound=%v uintIndexEnabled=%v elemType=%d elemSize=%d",
		objURL, o.BigIndexFound, extensionUintIndexEnabled, m.vertexIndexElementType, m.vertexIndexElementSize))
}

func findBoundingRadius(o *gwob.Obj) float64 {

	farthestSquaredDist := 0.0

	offset := o.StrideOffsetPosition / 4
	step := o.StrideSize / 4
	for i := offset; i < len(o.Coord); i += step {
		x, y, z := float64(o.Coord[i]), float64(o.Coord[i+1]), float64(o.Coord[i+2])
		squaredDist := x*x + y*y + z*z
		if squaredDist > farthestSquaredDist {
			farthestSquaredDist = squaredDist
		}
	}

	return math.Sqrt(farthestSquaredDist)
}

func newModel(s shader, modelName string, gl *webgl.Context, objURL string,
	front, up []float64, assetPath asset, textureTable map[string]*texture,
	repeatTexture bool, materialLib gwob.MaterialLib, extensionUintIndexEnabled bool) (*texturizedModel, error) {

	// allocate new model
	mod := &texturizedModel{simpleModel: simpleModel{modelName: modelName}}

	// fetch model from objURL

	var buf []byte
	var err error

	if buf, err = httpFetch(objURL); err != nil {
		return nil, fmt.Errorf("newModel: fetch model from objURL=%s error: %v", objURL, err)
	}

	opt := &gwob.ObjParserOptions{Logger: func(msg string) { log(fmt.Sprintf("newModel: %s", msg)) }}
	var o *gwob.Obj
	if o, err = gwob.NewObjFromBuf(buf, opt); err != nil {
		return nil, fmt.Errorf("newModel: parse error objURL=%s error: %v", objURL, err)
	}

	log(fmt.Sprintf("newModel: objURL=%s elements=%d bigIndex=%v texCoord=%v normCoord=%v", objURL, o.NumberOfElements(), o.BigIndexFound, o.TextCoordFound, o.NormCoordFound))

	if !o.TextCoordFound {
		return nil, fmt.Errorf("newModel: objURL=%s FIXME texture coordinates required", objURL)
	}

	libURL := fmt.Sprintf("%s/%s", assetPath.mtl, o.Mtllib)

	groupListSize := len(o.Groups)

	// Load textures for groups
	for i, g := range o.Groups {
		//log(fmt.Sprintf("newModel: objURL=%s group=%s size=%d mtllib=%s consider material=%s", objURL, g.Name, g.IndexCount, o.Mtllib, g.Usemtl))

		if g.IndexCount < 3 {
			log(fmt.Sprintf("newModel: objURL=%s group=%s size=%d mtllib=%s bad index list size", objURL, g.Name, g.IndexCount, o.Mtllib))
			if err = addGroupTextureNull(mod, gl, groupListSize, i); err != nil {
				return nil, err
			}
			continue // skip group missing index list
		}

		if g.Usemtl == "" {
			log(fmt.Sprintf("newModel: objURL=%s group=%s size=%d mtllib=%s missing material name", objURL, g.Name, g.IndexCount, o.Mtllib))
			if err = addGroupTextureNull(mod, gl, groupListSize, i); err != nil {
				return nil, err
			}
			continue // skip group missing material name
		}

		var mat *gwob.Material
		var matOk bool
		if mat, matOk = materialLib.Lib[g.Usemtl]; !matOk {
			// material not found -- fetch lib

			//log(fmt.Sprintf("newModel: objURL=%s group=%s load mtllib=%s", objURL, g.Name, o.Mtllib))

			if libErr := fetchMaterialLib(materialLib, libURL); libErr != nil {
				log(fmt.Sprintf("newModel: objURL=%s group=%s size=%d mtllib=%s material=%s LIB FAILURE: %v", objURL, g.Name, g.IndexCount, libURL, g.Usemtl, libErr))
				if err = addGroupTextureNull(mod, gl, groupListSize, i); err != nil {
					return nil, err
				}
				continue // ugh
			}

			log(fmt.Sprintf("newModel: objURL=%s group=%s size=%d mtllib=%s LIB LOADED for material=%s", objURL, g.Name, g.IndexCount, libURL, g.Usemtl))

			if mat, matOk = materialLib.Lib[g.Usemtl]; !matOk {
				log(fmt.Sprintf("newModel: objURL=%s group=%s size=%d mtllib=%s MISSING material=%s", objURL, g.Name, g.IndexCount, libURL, g.Usemtl))
				if err = addGroupTextureNull(mod, gl, groupListSize, i); err != nil {
					return nil, err
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

		var wrap int
		if repeatTexture {
			wrap = gl.REPEAT
		} else {
			wrap = gl.CLAMP_TO_EDGE
		}

		if err = addGroupTexture(mod, gl, textureTable, groupListSize, i, mat.Map_Kd, textureURL, tempColor, wrap); err != nil {
			return nil, err
		}
	}

	if len(o.Groups) != len(mod.textures) {
		return nil, fmt.Errorf("newModel: objURL=%s BAD group/texture count: groups=%d textures=%d", objURL, len(o.Groups), len(mod.textures))
	}

	mod.boundingRadius = findBoundingRadius(o)

	mod.mesh = o

	//showGroups(mod)
	sort.Sort(GroupByTextureName(GroupByTextureName{mod}))
	//showGroups(mod)

	mod.createBuffers(objURL, gl, extensionUintIndexEnabled)

	// push new model into shader.modelList
	s.addModel(mod)

	return mod, nil
}

const vertexPositionBufferItemSize = 3 // coord x,y,z
const textureCoordBufferItemSize = 2   // coord s,t

func (m *texturizedModel) drawGroups(gameInfo *gameState, u_Sampler *js.Object) {
	gl := gameInfo.gl

	// scan model groups
	for i, g := range m.mesh.Groups {
		t := m.textures[i]
		if t == nil {
			continue // skip group because texture is not ready
		}

		gl.BindTexture(gl.TEXTURE_2D, t.texture)

		// set sampler to use texture assigned to unit
		gl.Uniform1i(u_Sampler, gameInfo.defaultTextureUnit)

		gl.DrawElements(gl.TRIANGLES, g.IndexCount,
			m.vertexIndexElementType,
			g.IndexBegin*m.vertexIndexElementSize)
	}
}

func (m *simpleModel) draw(gameInfo *gameState, prog shader) {
	// abstract -- FIXME: Draw non-texturized models ???
}

func (m *texturizedModel) draw(gameInfo *gameState, prog shader) {
	gl := gameInfo.gl // shortcut

	gl.BindBuffer(gl.ARRAY_BUFFER, m.vertexBuffer)

	// vertex coord x,y,z
	gl.VertexAttribPointer(prog.attrLoc_Position(),
		vertexPositionBufferItemSize,
		gl.FLOAT, false, m.mesh.StrideSize,
		m.mesh.StrideOffsetPosition)

	texturizer, isTexturizer := prog.(*simpleTexturizer)

	// texture coord s,t
	if isTexturizer {
		gl.VertexAttribPointer(texturizer.attrLoc_TextureCoord(),
			textureCoordBufferItemSize,
			gl.FLOAT, false, m.mesh.StrideSize,
			m.mesh.StrideOffsetTexture)
	}

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, m.vertexIndexBuffer)

	u_MV := prog.unif_MV()

	if isTexturizer {
		u_Sampler := texturizer.unif_Sampler()
		for _, inst := range m.instanceList {
			inst.uploadModelView(gameInfo, gl, u_MV, &gameInfo.cam)
			m.drawGroups(gameInfo, u_Sampler)
		}
	}
}

func (m *simpleModel) name() string {
	return m.modelName
}

func (m *simpleModel) findInstance(id string) (*instance, bool) {
	for _, i := range m.instanceList {
		if id == i.id {
			return i, true
		}
	}
	return nil, false
}

func (m *simpleModel) addInstance(inst *instance) {
	m.instanceList = append(m.instanceList, inst)
	//log(fmt.Sprintf("model.addInstance: model=%s newInstance=%s instances=%d instanceList=%v", m.name(), inst.id, len(m.instanceList), m.instanceList))
}

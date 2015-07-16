package main

import (
	"encoding/json"
	"fmt"

	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/webgl"

	"negentropia/world/obj"
)

// skybox struct for decoding json
type skybox struct {
	Cube           string
	VertexShader   string
	FragmentShader string
	FaceRight      string
	FaceLeft       string
	FaceUp         string
	FaceDown       string
	FaceFront      string
	FaceBack       string
}

// cube struct for decoding json
type cube struct {
	VertCoord []float32
	TexCoord  []float32
	VertInd   []int
}

type skyboxShader struct {
	simpleShader
	u_Skybox *js.Object
}

type skyboxModel struct {
	simpleModel
	cubemapTexture *js.Object
}

func setCubemapFace(gl *webgl.Context, image, texture *js.Object, face int, faceURL string) {
	log(fmt.Sprintf("setCubemapFace: faceURL=%s", faceURL))

	gl.BindTexture(gl.TEXTURE_CUBE_MAP, texture)

	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MIN_FILTER, gl.NEAREST)

	gl.TexImage2D(face, 0, gl.RGBA, gl.RGBA, gl.UNSIGNED_BYTE, image)

	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)

	gl.BindTexture(gl.TEXTURE_CUBE_MAP, nil)
}

func onCubemapFaceLoad(gl *webgl.Context, s *skyboxModel, image *js.Object, face int, faceURL string) {
	//log(fmt.Sprintf("onCubemapFaceLoad: faceURL=%s", faceURL))
	setCubemapFace(gl, image, s.cubemapTexture, face, faceURL)
}

func (s *skyboxModel) addCubemapFace(gl *webgl.Context, face int, faceURL string) {

	image := newImage()

	image.Set("onload", func() {
		go onCubemapFaceLoad(gl, s, image, face, faceURL)
	})

	image.Set("src", faceURL)
}

func reverse(list []int) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func fetchSkybox(gameInfo *gameState, skyboxURL string) {

	buf, err := httpFetch(skyboxURL)
	if err != nil {
		log(fmt.Sprintf("fetchSkybox: skyboxURL=%s failure: %v", skyboxURL, err))
		return
	}

	box := skybox{}
	if err = json.Unmarshal(buf, &box); err != nil {
		log(fmt.Sprintf("fetchSkybox: skyboxURL=%s JSON=%v: error=%v", skyboxURL, string(buf), err))
		return
	}

	cubeURL := box.Cube
	buf, err = httpFetch(cubeURL)
	if err != nil {
		log(fmt.Sprintf("fetchSkybox: cubeURL=%s failure: %v", cubeURL, err))
		return
	}

	cube := cube{}
	if err = json.Unmarshal(buf, &cube); err != nil {
		log(fmt.Sprintf("fetchSkybox: cubeURL=%s JSON=%v: error=%v", cubeURL, string(buf), err))
		return
	}

	log(fmt.Sprintf("fetchSkybox: cube=%v", cube))

	vertShader := box.VertexShader
	fragShader := box.FragmentShader
	prog := newShaderProgram(gameInfo.gl, vertShader, fragShader)
	if prog == nil {
		log(fmt.Sprintf("fetchSkybox: skyboxURL=%s failure creating shader: vert=%v frag=%v", skyboxURL, vertShader, fragShader))
		return
	}

	skybox := &skyboxShader{simpleShader: simpleShader{program: prog, progName: "skybox"}}

	gl := gameInfo.gl // shortcut

	skybox.init(gl)

	// create model from mesh
	var o *obj.Obj
	o, err = obj.NewObjFromVertex(cube.VertCoord, cube.VertInd)

	log(fmt.Sprintf("fetchSkybox: skyboxURL=%s reversing cube indices", skyboxURL))
	reverse(o.Indices)

	m := &skyboxModel{cubemapTexture: gl.CreateTexture(), simpleModel: simpleModel{modelName: "skybox-model", mesh: o}}

	log(fmt.Sprintf("fetchSkybox: skyboxURL=%s JSON=%v skybox=%v mesh=%v FIXME WRITEME", skyboxURL, string(buf), box, o))

	m.createBuffers(cubeURL, gl, gameInfo.extensionUintIndexEnabled)

	// add cubemap faces to model

	m.addCubemapFace(gl, gl.TEXTURE_CUBE_MAP_POSITIVE_X, box.FaceRight)
	m.addCubemapFace(gl, gl.TEXTURE_CUBE_MAP_NEGATIVE_X, box.FaceLeft)
	m.addCubemapFace(gl, gl.TEXTURE_CUBE_MAP_POSITIVE_Y, box.FaceUp)
	m.addCubemapFace(gl, gl.TEXTURE_CUBE_MAP_NEGATIVE_Y, box.FaceDown)
	m.addCubemapFace(gl, gl.TEXTURE_CUBE_MAP_POSITIVE_Z, box.FaceFront)
	m.addCubemapFace(gl, gl.TEXTURE_CUBE_MAP_NEGATIVE_Z, box.FaceBack)

	i := newInstanceNull("skybox-instance")
	// skyboxScale should not matter much when it is centered on camera
	// skybox faces should however be within view frustum
	i.scale = 30

	m.addInstance(i) // add instance to model

	skybox.addModel(m) // add model to shader

	gameInfo.skybox = skybox
}

func (s *skyboxShader) init(gl *webgl.Context) {
	s.a_Position = s.getAttrib(gl, "a_Position")

	s.u_P = s.getUniform(gl, "u_P")
	s.u_MV = s.getUniform(gl, "u_MV")
	s.u_Skybox = s.getUniform(gl, "u_Skybox")
}

func (s *skyboxShader) draw(gameInfo *gameState) {
	gl := gameInfo.gl // shortcut

	gl.DepthRange(1.0, 1.0) // draw skybox at far plane

	gl.UseProgram(s.program)
	gl.EnableVertexAttribArray(s.a_Position)

	uploadPerspective(gl, s.u_P, &gameInfo.pMatrix)

	// draw every model
	for _, m := range s.modelList {
		m.draw(gameInfo, s)
	}

	gl.DepthRange(0.0, 1.0) // restore default
}

func (m *skyboxModel) draw(gameInfo *gameState, prog shader) {

	skyboxSh, isSkybox := prog.(*skyboxShader)

	if !isSkybox {
		return
	}

	gl := gameInfo.gl // shortcut

	gl.BindBuffer(gl.ARRAY_BUFFER, m.vertexBuffer)

	// vertex coord x,y,z
	gl.VertexAttribPointer(skyboxSh.a_Position, vertexPositionBufferItemSize, gl.FLOAT, false, 0, 0)

	gl.BindTexture(gl.TEXTURE_CUBE_MAP, m.cubemapTexture)

	gl.Uniform1i(skyboxSh.u_Skybox, gameInfo.defaultTextureUnit)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, m.vertexIndexBuffer)

	u_MV := skyboxSh.u_MV

	for i, inst := range m.instanceList {
		inst.uploadModelView(gameInfo, gl, u_MV, &gameInfo.cam)

		for g, grp := range m.mesh.Groups {
			if gameInfo.debugDraw {
				log(fmt.Sprintf("skyboxModel.draw(): model=%v instance[%d]=%v group=%d DrawElements: elemType=%d elemSize=%d count=%d begin=%d ",
					m.name(), i, inst.id, g, m.vertexIndexElementType, m.vertexIndexElementSize, grp.IndexCount, grp.IndexBegin))
			}
			gl.DrawElements(gl.TRIANGLES, grp.IndexCount, m.vertexIndexElementType, grp.IndexBegin*m.vertexIndexElementSize)
		}
	}

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, nil)
	gl.BindTexture(gl.TEXTURE_CUBE_MAP, nil)
	gl.BindBuffer(gl.ARRAY_BUFFER, nil)
}

package main

import (
	"fmt"

	"strconv"
	"strings"

	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/webgl"
)

type boxdemo struct {
	cubemapTexture    *js.Object
	program           *js.Object
	u_P               *js.Object
	u_MV              *js.Object
	u_Skybox          *js.Object
	a_Position        int
	vertexBuffer      *js.Object
	vertexIndexBuffer *js.Object
	cubeIndices       []uint16
	cubeCoord         []float32
	inst              *instance
}

func newBoxdemo(gameInfo *gameState) *boxdemo {

	gl := gameInfo.gl // shortcut

	box := &boxdemo{cubemapTexture: gl.CreateTexture(), vertexBuffer: gl.CreateBuffer(), vertexIndexBuffer: gl.CreateBuffer()}

	box.cubeIndices = []uint16{}
	indices := "22 23 20 21 22 20 19 18 16 18 17 16 15 14 12 14 13 12 11 10 8 10 9 8 6 7 4 5 6 4 3 2 0 2 1 0"
	list := strings.Fields(indices)
	for _, s := range list {
		v, _ := strconv.Atoi(s)
		box.cubeIndices = append(box.cubeIndices, uint16(v))
	}
	log(fmt.Sprintf("newBoxdemo: indices = (%d) %v", len(box.cubeIndices), box.cubeIndices))

	box.cubeCoord = []float32{}
	coord := "1 1 1 -1 1 1 -1 -1 1 1 -1 1 1 1 -1 -1 1 -1 -1 -1 -1 1 -1 -1 -1 1 1 -1 1 -1 -1 -1 -1 -1 -1 1 1 1 1 1 -1 1 1 -1 -1 1 1 -1 1 1 1 1 1 -1 -1 1 -1 -1 1 1 1 -1 1 1 -1 -1 -1 -1 -1 -1 -1 1"
	list = strings.Fields(coord)
	for _, s := range list {
		v, _ := strconv.ParseFloat(s, 32)
		box.cubeCoord = append(box.cubeCoord, float32(v))
	}

	log(fmt.Sprintf("newBoxdemo: coord = (%d) %v", len(box.cubeCoord), box.cubeCoord))

	gl.BindBuffer(gl.ARRAY_BUFFER, box.vertexBuffer)
	gl.BufferData(gl.ARRAY_BUFFER, box.cubeCoord, gl.STATIC_DRAW)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, box.vertexIndexBuffer)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, box.cubeIndices, gl.STATIC_DRAW)

	box.inst = newInstanceNull("boxdemo-instance")
	box.inst.scale = 20

	box.program = newShaderProgram(gl, "/shader/skybox_vs.txt", "/shader/skybox_fs.txt")
	if box.program == nil {
		log(fmt.Sprintf("newBoxdemo: failure creating shader"))
		return nil
	}
	log(fmt.Sprintf("newBoxdemo: shader compiled"))

	box.u_P = gl.GetUniformLocation(box.program, "u_P")
	box.u_MV = gl.GetUniformLocation(box.program, "u_MV")
	box.u_Skybox = gl.GetUniformLocation(box.program, "u_Skybox")
	box.a_Position = gl.GetAttribLocation(box.program, "a_Position")

	box.fetchCubemapFace(gl, gl.TEXTURE_CUBE_MAP_POSITIVE_X, "/texture/space_rt.jpg")
	box.fetchCubemapFace(gl, gl.TEXTURE_CUBE_MAP_NEGATIVE_X, "/texture/space_lf.jpg")
	box.fetchCubemapFace(gl, gl.TEXTURE_CUBE_MAP_POSITIVE_Y, "/texture/space_up.jpg")
	box.fetchCubemapFace(gl, gl.TEXTURE_CUBE_MAP_NEGATIVE_Y, "/texture/space_dn.jpg")
	box.fetchCubemapFace(gl, gl.TEXTURE_CUBE_MAP_POSITIVE_Z, "/texture/space_fr.jpg")
	box.fetchCubemapFace(gl, gl.TEXTURE_CUBE_MAP_NEGATIVE_Z, "/texture/space_bk.jpg")

	return box
}

func (b *boxdemo) fetchCubemapFace(gl *webgl.Context, face int, faceURL string) {

	image := newImage()

	image.Set("onload", func() {
		log(fmt.Sprintf("newBoxdemo: fetchCubemapFace: onload: URL=%s", faceURL))
		go setCubemapFace(gl, image, b.cubemapTexture, face, faceURL)
	})

	image.Set("src", faceURL)
}

func (b *boxdemo) draw(gameInfo *gameState) {

	gl := gameInfo.gl // shortcut

	//gl.DepthRange(1.0, 1.0) // draw skybox at far plane

	gl.UseProgram(b.program)
	gl.EnableVertexAttribArray(b.a_Position)

	uploadPerspective(gl, b.u_P, &gameInfo.pMatrix)

	gl.BindBuffer(gl.ARRAY_BUFFER, b.vertexBuffer)

	// vertex coord x,y,z
	gl.VertexAttribPointer(b.a_Position, vertexPositionBufferItemSize, gl.FLOAT, false, 0, 0)

	gl.BindTexture(gl.TEXTURE_CUBE_MAP, b.cubemapTexture)

	gl.Uniform1i(b.u_Skybox, gameInfo.defaultTextureUnit)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, b.vertexIndexBuffer)

	b.inst.uploadModelView(gameInfo, gl, b.u_MV, &gameInfo.cam)

	gl.DrawElements(gl.TRIANGLES, len(b.cubeIndices), gl.UNSIGNED_SHORT, 0)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, nil)
	gl.BindTexture(gl.TEXTURE_CUBE_MAP, nil)
	gl.BindBuffer(gl.ARRAY_BUFFER, nil)

	gl.DepthRange(0.0, 1.0) // restore default
}

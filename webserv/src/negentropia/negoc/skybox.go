package main

import (
	"encoding/json"
	"fmt"

	"github.com/gopherjs/gopherjs/js"
)

// skybox struct for decoding json
type skybox struct {
	Cube           string
	VertexShader   string
	FragmentShader string
	FaceRight      string
	FaceLeft       string
	FaceDown       string
	FaceFront      string
	FaceBack       string
}

type skyboxShader struct {
	simpleShader
	u_Skybox *js.Object
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

	log(fmt.Sprintf("fetchSkybox: skyboxURL=%s JSON=%v skybox=%v FIXME WRITEME", skyboxURL, string(buf), box))
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

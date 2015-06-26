package main

import (
	"encoding/json"
	"fmt"
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

package main

import (
	"fmt"
	"negentropia/world/parser"
	"strings"
)

func resetPickColor() {
	log("resetPickColor: WRITEME")
}

func resetZone() {
	/*
	   programList = new List<ShaderProgram>(); // drop existing shaders
	   shaderCache = new Map<String, Shader>(); // drop existing compile shader cache
	   textureTable = new Map<String, Texture>(); // drop existing texture table

	   skybox =
	       null; // drop skybox shader (re-created only when new skybox is added to zone)
	   picker =
	       null; // drop picking shader (re-created only when any model instance is added)
	   solidShader =
	       null; // drop axis shader (re-created only when any model instance is added)
	*/

	resetPickColor()
}

func dispatch(gameInfo *gameState, code int, data string, tab map[string]string) {
	//log(fmt.Sprintf("dispatch: code=%v data=%v tab=%v", code, data, tab))

	switch code {
	case CM_CODE_INFO:

		log(fmt.Sprintf("dispatch: server sent info: %s", data))

		if strings.HasPrefix(data, "welcome") {
			// test echo loop thru server
			msg := &ClientMsg{Code: CM_CODE_ECHO, Data: "hi, please echo back this test"}
			gameInfo.sock.write(msg)
		}

	case CM_CODE_ZONE:

		if backfaceCulling, ok := tab["backfaceCulling"]; ok {
			culling := stringIsTrue(backfaceCulling)
			//log(fmt.Sprintf("dispatch: zone: backfaceCulling: recv=%s parsed=%v", backfaceCulling, culling))
			updateCulling(gameInfo.gl, culling)
		}

		if camCoord, ok := tab["cameraCoord"]; ok {
			if coord, err := parser.ParseFloatVector3Comma(camCoord); err != nil {
				log(fmt.Sprintf("dispatch: error parsing Vector3(%s): %v", camCoord, err))
			} else {
				cameraMoveTo(&gameInfo.cam, coord)
			}
		}

		resetZone()

	case CM_CODE_SKYBOX:

		if skyboxURL, ok := tab["skyboxURL"]; ok {
			fetchSkybox(skyboxURL)
		} else {
			log("dispatch: missing skybox URL")
		}

	case CM_CODE_PROGRAM:

		var nameOk, vertOk, fragOk bool
		var programName, vertShader, fragShader string

		if programName, nameOk = tab["programName"]; !nameOk {
		}
		if vertShader, vertOk = tab["vertexShader"]; !vertOk {
		}
		if fragShader, fragOk = tab["fragmentShader"]; !fragOk {
		}

		if nameOk && vertOk && fragOk {
			fetchShaderProgram(programName, vertShader, fragShader)
		}

	default:
		log(fmt.Sprintf("dispatch: unknown code=%v data=%v tab=%v", code, data, tab))
	}
}

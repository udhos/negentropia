package main

import (
	"fmt"
	"strings"

	"negentropia/ipc"
	"negentropia/world/parser"
)

func resetPickColor() {
	log("resetPickColor: WRITEME")
}

func resetZone(gameInfo *gameState) {
	resetGame(gameInfo)
	resetPickColor()
}

func dispatch(gameInfo *gameState, code int, data string, tab map[string]string) {
	//log(fmt.Sprintf("dispatch: code=%v data=%v tab=%v", code, data, tab))

	switch code {
	case ipc.CM_CODE_INFO:

		log(fmt.Sprintf("dispatch: server sent info: %s", data))

		if strings.HasPrefix(data, "welcome") {
			// test echo loop thru server
			msg := &ipc.ClientMsg{Code: ipc.CM_CODE_ECHO, Data: "hi, please echo back this test"}
			gameInfo.sock.write(msg)
		}

	case ipc.CM_CODE_ZONE:

		log("dispatch: server sending NEW ZONE")

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

		resetZone(gameInfo)

	case ipc.CM_CODE_SKYBOX:

		if skyboxURL, ok := tab["skyboxURL"]; ok {
			fetchSkybox(skyboxURL)
		} else {
			log("dispatch: missing skybox URL")
		}

	case ipc.CM_CODE_PROGRAM:

		var nameOk, vertOk, fragOk bool
		var programName, vertShader, fragShader string

		if programName, nameOk = tab["programName"]; !nameOk {
			log(fmt.Sprintf("dispatch: program: missing name"))
		}
		if vertShader, vertOk = tab["vertexShader"]; !vertOk {
			log(fmt.Sprintf("dispatch: program: missing vertex shader URL"))
		}
		if fragShader, fragOk = tab["fragmentShader"]; !fragOk {
			log(fmt.Sprintf("dispatch: program: missing fragment shader URL"))
		}

		if nameOk && vertOk && fragOk {
			fetchShaderProgram(gameInfo, programName, vertShader, fragShader)
		}

	case ipc.CM_CODE_INSTANCE:

		createInstance(gameInfo, tab)

		//countInstances(gameInfo)

	case ipc.CM_CODE_INSTANCE_UPDATE:
		log(fmt.Sprintf("dispatch: instance update: WRITEME"))

	case ipc.CM_CODE_MESSAGE:
		log(fmt.Sprintf("dispatch: message: WRITEME"))

	default:
		log(fmt.Sprintf("dispatch: unknown code=%v data=%v tab=%v", code, data, tab))
	}
}

func countInstances(gameInfo *gameState) {
	log(fmt.Sprintf("countInstances: shaderList=%v size=%d", &gameInfo.shaderList, len(gameInfo.shaderList)))

	for _, s := range gameInfo.shaderList {
		t := s.(*simpleTexturizer)

		log(fmt.Sprintf("countInstances: shader=%v models=%d", t.name(), len(t.modelList)))

		for _, m := range t.modelList {
			log(fmt.Sprintf("countInstances: shader=%v model=%s instances=%d", t.name(), m.modelName, len(m.instanceList)))
		}
	}
}

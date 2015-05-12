package main

import (
	"fmt"
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

		resetZone()

	default:
		log(fmt.Sprintf("dispatch: unknown code=%v data=%v tab=%v", code, data, tab))
	}
}

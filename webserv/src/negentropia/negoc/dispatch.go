package main

import (
	"fmt"
	"strings"
)

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

	default:
		log(fmt.Sprintf("dispatch: uknown code=%v data=%v tab=%v", code, data, tab))
	}
}

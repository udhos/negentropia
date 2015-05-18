package main

import (
	"fmt"
)

func fetchSkybox(skyboxURL string) {

	log(fmt.Sprintf("fetchSkybox: WRITEME skyboxURL=%v", skyboxURL))

	/*
	   case CM_CODE_SKYBOX:
	     String skyboxURL = tab['skyboxURL'];

	     void handleResponse(String response) {
	       Map<String, String> skybox = JSON.decode(response);
	       addSkybox(gl, skybox);
	     }

	     HttpRequest.getString(skyboxURL).then(handleResponse).catchError((e) {
	       err("dispatcher: failure fetching skyboxURL=$skyboxURL: $e");
	     });

	     break;
	*/

}

function SkyboxProgram(vertexShaderURL, fragmentShaderURL) {
	console.log("new skybox program: vsURL=" + vertexShaderURL + " fsURL=" + fragmentShaderURL);
	this.modelList = [];
	
	// Async request for shader program
	var p = this; // don't put 'this' inside the closure below
	fetchProgramFromURL(vertexShaderURL, fragmentShaderURL, function (prog) { skyboxShaderLoaded(p, prog); });
}

function skyboxShaderLoaded(p, prog) {

	if (!('shaderProgram' in prog)) {
		console.log("skybox: shader program load failure");
		return;
	}

	console.log("skybox: shader program loaded");
	p.shaderProgram = prog.shaderProgram;
		
	// save vertex attribute location
	p.aPosition = gl.getAttribLocation(p.shaderProgram, "aPosition");
	if (p.aPosition < 0) {
		console.log("skybox: aPosition: failure querying attribute location");
	}
	
	p.uSkybox = gl.getUniformLocation(p.shaderProgram, "uSkybox");
	if (p.uSkybox == null) {
		console.log("skybox: uSkybox: failure querying uniform location");
	}
}

SkyboxProgram.prototype.addModel = function(m) {
	this.modelList.push(m);
}

SkyboxProgram.prototype.drawModels = function() {
	
    gl.useProgram(this.shaderProgram);
    gl.enableVertexAttribArray(this.aVertexPosition);

	for (var m in this.modelList) {
		this.modelList[m].drawInstances();
	}
	
	// clean up
	gl.bindBuffer(gl.ARRAY_BUFFER, null);
	gl.bindBuffer(gl.ELEMENT_ARRAY_BUFFER, null);
	
    //gl.disableVertexAttribArray(this.aVertexPosition); // needed ??
}

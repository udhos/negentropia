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

    p.uView = gl.getUniformLocation(p.shaderProgram, "uView");
	if (p.uView == null) {
		console.log("skybox: uView: failure querying uniform location");
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
    gl.enableVertexAttribArray(this.aPosition);
	
	/*
	var unit = 0;
	gl.activeTexture(gl.TEXTURE0 + unit);
	gl.uniform1i(this.uSkybox, unit);
	*/
	
	for (var m in this.modelList) {
		this.modelList[m].drawInstances();
	}
	
	// clean up
	gl.bindBuffer(gl.ARRAY_BUFFER, null);
	gl.bindBuffer(gl.ELEMENT_ARRAY_BUFFER, null);
	
    //gl.disableVertexAttribArray(this.aVertexPosition); // needed ??
}

function SkyboxModel(program, URL, reverse, rescale) {
	this.program = program;
	this.URL = URL;
	this.instanceList = [];
	
	this.cubemapTexture = gl.createTexture();
	gl.bindTexture(gl.TEXTURE_CUBE_MAP, this.cubemapTexture);
	gl.texParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MAG_FILTER, gl.LINEAR);
	gl.texParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MIN_FILTER, gl.LINEAR);
	
	// Async request for buffer data (mesh)
	var m = this; // don't put 'this' inside the closure below
	fetchBufferData(this.URL, function (buf) { skyboxModelBufferDataLoaded(m, buf); }, reverse, rescale);	
}

function skyboxModelBufferDataLoaded(model, buf, reverse) {
	model.buffer = buf;
}

SkyboxModel.prototype.addInstance = function(i) {
	this.instanceList.push(i);
}

SkyboxModel.prototype.drawInstances = function() {

	gl.bindTexture(gl.TEXTURE_CUBE_MAP, this.cubemapTexture);

	for (var i in this.instanceList) {
		this.instanceList[i].draw(this.program);
	}
	
	gl.bindTexture(gl.TEXTURE_CUBE_MAP, null);
}

SkyboxModel.prototype.addCubemapFace = function(face, URL) {
	console.log("add cubemap face: " + URL);
	loadCubemapFace(face, this.cubemapTexture, URL);
}

function SkyboxInstance(model) {
	this.model = model;
	this.viewMatrix  = mat4.create();
}

SkyboxInstance.prototype.draw = function(program) {

	var buf = this.model.buffer;

	// view transform
	//mat4.identity(this.viewMatrix);
	//mat4.lookAt([0,0,0], [0,0,-1], [0,1,0], this.viewMatrix);

	mat4.lookAt(neg.eye, neg.center, neg.up, this.viewMatrix);
	mat4.multiply(this.viewMatrix, neg.pMatrix);
	gl.uniformMatrix4fv(program.uView, false, this.viewMatrix);
	
	// vertex coord
    gl.bindBuffer(gl.ARRAY_BUFFER, buf.vertexPositionBuffer);
   	gl.vertexAttribPointer(program.aVertexPosition, buf.vertexPositionBufferItemSize, gl.FLOAT, false, 0, 0);

	// cubemap sampler
	var unit = 0;
	gl.activeTexture(gl.TEXTURE0 + unit);
	gl.bindTexture(gl.TEXTURE_CUBE_MAP, this.model.cubemapTexture);
	gl.uniform1i(program.uSkybox, unit);
	
	// draw
	gl.bindBuffer(gl.ELEMENT_ARRAY_BUFFER, buf.vertexIndexBuffer);
	gl.drawElements(gl.TRIANGLES, buf.vertexIndexLength, gl.UNSIGNED_SHORT, 0 * buf.vertexIndexBufferItemSize);
}

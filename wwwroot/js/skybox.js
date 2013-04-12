function SkyboxProgram(vertexShaderURL, fragmentShaderURL) {
	console.log("new skybox program: vsURL=" + vertexShaderURL + " fsURL=" + fragmentShaderURL);
	this.modelList = [];
	
	// Async request for shader program
	var p = this; // don't put 'this' inside the closure below
	fetchProgramFromURL(vertexShaderURL, fragmentShaderURL, function (prog) { skyboxShaderLoaded(p, prog); });
}

function skyboxGetAttr(p, attr) {
	p[attr] = gl.getAttribLocation(p.shaderProgram, attr);
	if (p[attr] < 0) {
		console.log("skybox: failure querying attribute location: " + attr);
	}
}

function skyboxGetUniform(p, uniform) {
	p[uniform] = gl.getUniformLocation(p.shaderProgram, uniform);
	if (p[uniform] < 0) {
		console.log("skybox: failure querying uniform location: " + uniform);
	}
}

function skyboxShaderLoaded(p, prog) {

	if (!('shaderProgram' in prog)) {
		console.log("skybox: shader program load failure");
		return;
	}

	console.log("skybox: shader program loaded");
	p.shaderProgram = prog.shaderProgram;
		
	// save attribute location
	skyboxGetAttr(p, "a_Position");

	// save uniform location
	skyboxGetUniform(p, "u_MV");
	skyboxGetUniform(p, "u_P");
	skyboxGetUniform(p, "u_Skybox");
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
	gl.uniform1i(this.u_Skybox, unit);
	*/

	gl.uniformMatrix4fv(this.u_P, false, neg.pMatrix);
	
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

function SkyboxInstance(model, center, scale) {
	this.model = model;
	this.center = center;
	this.scale = scale;
}

SkyboxInstance.prototype.draw = function(program) {

	var buf = this.model.buffer;

	// view transform
	//mat4.identity(this.viewMatrix);
	//mat4.lookAt([0,0,0], [0,0,-1], [0,1,0], this.viewMatrix);

	var MV = mat4.create(); // model-view

	// 6/7. camera
	mat4.lookAt(neg.eye, neg.center, neg.up, MV);
	
	/*
	// 5. obj translate
    mat4.translate(MV, this.center);
		
	// 1. obj scale
	mat4.scale(MV, [this.scale, this.scale, this.scale]);
	*/
	gl.uniformMatrix4fv(program.u_MV, false, MV);
	
	// vertex coord
    gl.bindBuffer(gl.ARRAY_BUFFER, buf.vertexPositionBuffer);
   	gl.vertexAttribPointer(program.aVertexPosition, buf.vertexPositionBufferItemSize, gl.FLOAT, false, 0, 0);

	// cubemap sampler
	var unit = 0;
	gl.activeTexture(gl.TEXTURE0 + unit);
	gl.bindTexture(gl.TEXTURE_CUBE_MAP, this.model.cubemapTexture);
	gl.uniform1i(program.u_Skybox, unit);
	
	// draw
	gl.bindBuffer(gl.ELEMENT_ARRAY_BUFFER, buf.vertexIndexBuffer);
	gl.drawElements(gl.TRIANGLES, buf.vertexIndexLength, gl.UNSIGNED_SHORT, 0 * buf.vertexIndexBufferItemSize);
}

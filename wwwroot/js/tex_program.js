function TexProgram(vertexShaderURL, fragmentShaderURL) {
	console.log("new tex program: vsURL=" + vertexShaderURL + " fsURL=" + fragmentShaderURL);
	this.modelList = [];
	this.vsURL = vertexShaderURL;
	this.fsURL = fragmentShaderURL;
}

function texGetAttr(p, attr) {
	p[attr] = gl.getAttribLocation(p.shaderProgram, attr);
	if (p[attr] < 0) {
		console.log("tex program: failure querying attribute location: " + attr);
	}
}

function texGetUniform(p, uniform) {
	p[uniform] = gl.getUniformLocation(p.shaderProgram, uniform);
	if (p[uniform] < 0) {
		console.log("tex program: failure querying uniform location: " + uniform);
	}
}

TexProgram.prototype.fetch = function() {
	// Async request for shader program
	var p = this; // don't put 'this' inside the closure below
	fetchProgramFromURL(this.vsURL, this.fsURL, function (prog) { texShaderProgramLoaded(p, prog); });
}

function texShaderProgramLoaded(p, prog) {

	if (!('shaderProgram' in prog)) {
		console.log("tex program: shader program load failure");
		return;
	}

	console.log("tex program: shader program loaded");
	p.shaderProgram = prog.shaderProgram;
		
	// save attribute location
	texGetAttr(p, "a_Position");
	texGetAttr(p, "a_TextureCoord");

	// save uniform location
	texGetUniform(p, "u_MV");
	texGetUniform(p, "u_P");	
	texGetUniform(p, "u_Sampler");	
	texGetUniform(p, "u_Color");	
}

TexProgram.prototype.addModel = function(m) {
	this.modelList.push(m);
}

TexProgram.prototype.animate = function() {
	for (var m in this.modelList) {
		this.modelList[m].animate();
	}
}

TexProgram.prototype.drawModels = function() {
	
    gl.useProgram(this.shaderProgram);
    gl.enableVertexAttribArray(this.a_Position);
	gl.enableVertexAttribArray(this.a_TextureCoord);

	// perspective projection
	gl.uniformMatrix4fv(this.u_P, false, neg.pMatrix);
	
	// fallback solid color for textured objects
	gl.uniform4fv(this.u_Color, [1.0, 1.0, 1.0, 1.0]);

	for (var m in this.modelList) {
		this.modelList[m].drawInstances();
	}
	
	// clean up
	gl.bindBuffer(gl.ARRAY_BUFFER, null);
	gl.bindBuffer(gl.ELEMENT_ARRAY_BUFFER, null);
	gl.bindTexture(gl.TEXTURE_2D, null);
	
    //gl.disableVertexAttribArray(this.a_Position); // needed ??
}

function TexModel(program, URL, reverse, mesh) {
	this.program = program;
	this.URL = URL;
	this.instanceList = [];
	
	if (mesh !== undefined) {
		this.buffer = texCreateBuffers(mesh.vertices, mesh.textures, mesh.indices);
		return;
	}
	
	// Async request for buffer data (mesh)
	var m = this; // don't put 'this' inside the closure below
	fetchBufferData(this.URL, function (buf) { texModelBufferDataLoaded(m, buf); }, reverse);	
}

function texModelBufferDataLoaded(model, buf, reverse) {
	model.buffer = buf;
}

TexModel.prototype.addInstance = function(i) {
	this.instanceList.push(i);
}

TexModel.prototype.animate = function() {
	for (var i in this.instanceList) {
		this.instanceList[i].animate();
	}
}

TexModel.prototype.drawInstances = function() {

	for (var i in this.instanceList) {
		this.instanceList[i].draw(this.program);
	}
}

function TexInstance(model, center) {
	this.model = model;
	if (center == null) {
		this.center = [0.0, 0.0, 0.0];
	}
	else {
		this.center = center;
	}
}

TexInstance.prototype.animate = function() {
}

TexInstance.prototype.draw = function(program) {

	var buf = this.model.buffer;
	
	var MV = mat4.create(); // model-view

	// 6/7. camera
	mat4.lookAt(neg.eye, neg.center, neg.up, MV);
	
	// 5. obj translate
    mat4.translate(MV, this.center);
		
	// 1. obj scale
	var s = 1;
	mat4.scale(MV, [s, s, s]);
	
	// send model-view matrix uniform
	gl.uniformMatrix4fv(program.u_MV, false, MV);
	
	// set texture sampler
	var unit = 0;
	gl.activeTexture(gl.TEXTURE0 + unit);
	gl.bindTexture(gl.TEXTURE_2D, texture);
	gl.uniform1i(program.u_Sampler, unit);
	
	// vertex coord
    gl.bindBuffer(gl.ARRAY_BUFFER, buf.vertexPositionBuffer);
   	gl.vertexAttribPointer(program.a_Position, buf.vertexPositionBufferItemSize, gl.FLOAT, false, 0, 0);
	
	// texture coord
	gl.bindBuffer(gl.ARRAY_BUFFER, buf.vertexTextureCoordBuffer);
	gl.vertexAttribPointer(program.a_TextureCoord, buf.vertexTextureCoordBufferItemSize, gl.FLOAT, false, 0, 0);
	
	// draw
	gl.bindBuffer(gl.ELEMENT_ARRAY_BUFFER, buf.vertexIndexBuffer);
	gl.drawElements(gl.TRIANGLES, buf.vertexIndexLength, gl.UNSIGNED_SHORT, 0 * buf.vertexIndexBufferItemSize);
}

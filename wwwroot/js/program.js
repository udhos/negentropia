function Program(vertexShaderURL, fragmentShaderURL) {
	console.log("new program: vsURL=" + vertexShaderURL + " fsURL=" + fragmentShaderURL);
	this.modelList = [];
	this.vsURL = vertexShaderURL;
	this.fsURL = fragmentShaderURL;
}

function getAttr(p, attr) {
	p[attr] = gl.getAttribLocation(p.shaderProgram, attr);
	if (p[attr] < 0) {
		console.log("program: failure querying attribute location: " + attr);
	}
}

function getUniform(p, uniform) {
	p[uniform] = gl.getUniformLocation(p.shaderProgram, uniform);
	if (p[uniform] < 0) {
		console.log("program: failure querying uniform location: " + uniform);
	}
}


Program.prototype.fetch = function() {
	// Async request for shader program
	var p = this; // don't put 'this' inside the closure below
	fetchProgramFromURL(this.vsURL, this.fsURL, function (prog) { shaderProgramLoaded(p, prog); });
}

function shaderProgramLoaded(p, prog) {

	if (!('shaderProgram' in prog)) {
		console.log("program: shader program load failure");
		return;
	}

	console.log("program: shader program loaded");
	p.shaderProgram = prog.shaderProgram;
		
	// save attribute location
	getAttr(p, "a_Position");

	// save uniform location
	getUniform(p, "u_MV");
	getUniform(p, "u_P");	
}

Program.prototype.addModel = function(m) {
	this.modelList.push(m);
}

Program.prototype.animate = function() {
	for (var m in this.modelList) {
		this.modelList[m].animate();
	}
}

Program.prototype.drawModels = function() {
	
    gl.useProgram(this.shaderProgram);
    gl.enableVertexAttribArray(this.a_Position);

	// perspective projection
	gl.uniformMatrix4fv(this.u_P, false, neg.pMatrix);

	for (var m in this.modelList) {
		this.modelList[m].drawInstances();
	}
	
	// clean up
	gl.bindBuffer(gl.ARRAY_BUFFER, null);
	gl.bindBuffer(gl.ELEMENT_ARRAY_BUFFER, null);
	
    //gl.disableVertexAttribArray(this.a_Position); // needed ??
}

function Model(program, URL, reverse, mesh) {
	this.program = program;
	this.URL = URL;
	this.instanceList = [];
	
	if (mesh !== undefined) {
		this.buffer = createBuffers(mesh.vertices, mesh.indices);
		return;
	}
	
	// Async request for buffer data (mesh)
	var m = this; // don't put 'this' inside the closure below
	fetchBufferData(this.URL, function (buf) { modelBufferDataLoaded(m, buf); }, reverse);	
}

function modelBufferDataLoaded(model, buf) {
	model.buffer = buf;
}

Model.prototype.addInstance = function(i) {
	this.instanceList.push(i);
}

Model.prototype.animate = function() {
	for (var i in this.instanceList) {
		this.instanceList[i].animate();
	}
}

Model.prototype.drawInstances = function() {
	for (var i in this.instanceList) {
		this.instanceList[i].draw(this.program);
	}
}

function Instance(model, center) {
	this.model  = model;
	if (center == null) {
		this.center = [0.0, 0.0, 0.0];
	}
	else {
		this.center = center;
	}
}

Instance.prototype.animate = function() {
}

Instance.prototype.draw = function(program) {

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
	
	// vertex coord
    gl.bindBuffer(gl.ARRAY_BUFFER, buf.vertexPositionBuffer);
   	gl.vertexAttribPointer(program.a_Position, buf.vertexPositionBufferItemSize, gl.FLOAT, false, 0, 0);
	
	// draw
	gl.bindBuffer(gl.ELEMENT_ARRAY_BUFFER, buf.vertexIndexBuffer);
	gl.drawElements(gl.TRIANGLES, buf.vertexIndexLength, gl.UNSIGNED_SHORT, 0 * buf.vertexIndexBufferItemSize);
}

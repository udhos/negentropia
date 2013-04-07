function Program(vertexShaderURL, fragmentShaderURL) {
	console.log("new program: vsURL=" + vertexShaderURL + " fsURL=" + fragmentShaderURL);
	this.modelList = [];
	
	// Async request for shader program
	var p = this; // don't put 'this' inside the closure below
	fetchProgramFromURL(vertexShaderURL, fragmentShaderURL, function (prog) { shaderProgramLoaded(p, prog); });
}

function shaderProgramLoaded(p, prog) {

	if (!('shaderProgram' in prog)) {
		console.log("program: shader program load failure");
		return;
	}

	console.log("program: shader program loaded");
	p.shaderProgram = prog.shaderProgram;
		
	// save vertex attribute location
	p.aVertexPosition = gl.getAttribLocation(p.shaderProgram, "aVertexPosition");
	if (p.aVertexPosition == -1) {
		console.log("program: aVertextPosition: failure querying attribute location");
	}
}

Program.prototype.addModel = function(m) {
	this.modelList.push(m);
}

Program.prototype.drawModels = function() {
	
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

function Model(program, URL, reverse) {
	this.program = program;
	this.URL = URL;
	this.instanceList = [];
	
	// Async request for buffer data (mesh)
	var m = this; // don't put 'this' inside the closure below
	fetchBufferData(this.URL, function (buf) { modelBufferDataLoaded(m, buf); }, reverse);	
}

function modelBufferDataLoaded(model, buf, reverse) {
	model.buffer = buf;
}

Model.prototype.addInstance = function(i) {
	this.instanceList.push(i);
}

Model.prototype.drawInstances = function() {
	for (var i in this.instanceList) {
		this.instanceList[i].draw(this.program);
	}
}

function Instance(model) {
	this.model = model;
}

Instance.prototype.draw = function(program) {

	var buf = this.model.buffer;
	
    gl.bindBuffer(gl.ARRAY_BUFFER, buf.vertexPositionBuffer);
   	gl.vertexAttribPointer(program.aVertexPosition, buf.vertexPositionBufferItemSize, gl.FLOAT, false, 0, 0);
	
	gl.bindBuffer(gl.ELEMENT_ARRAY_BUFFER, buf.vertexIndexBuffer);
	gl.drawElements(gl.TRIANGLES, buf.vertexIndexLength, gl.UNSIGNED_SHORT, 0 * buf.vertexIndexBufferItemSize);
}

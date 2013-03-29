function Program(vertexShaderURL, fragmentShaderURL) {
	console.log("new program: vsURL=" + vertexShaderURL + " fsURL=" + fragmentShaderURL);
	this.modelList = [];
	
	// Async request for shader program
	var p = this;
	fetchProgramFromURL("/shader/min_vs.txt", "/shader/min2_fs.txt", function (prog) { shaderProgramLoaded(p, prog); });
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

	/*
	if (!('shaderProgram' in this)) {
		// shader program not ready
		return;
	}
	*/
	
    gl.useProgram(this.shaderProgram);
    gl.enableVertexAttribArray(this.aVertexPosition);

	for (var m in this.modelList) {
		this.modelList[m].drawInstances();
	}
	
	// clean up
	gl.bindBuffer(gl.ARRAY_BUFFER, null);
	gl.bindBuffer(gl.ELEMENT_ARRAY_BUFFER, null);
	
    gl.disableVertexAttribArray(this.aVertexPosition);
}

function Model(nam, program, URL) {
	console.log("new model '" + name + "': URL=" + URL);
	this.name = nam;
	this.program = program;
	this.URL = URL;
	this.instanceList = [];
}

Model.prototype.init = function() {	
	// Async request for buffer data (mesh)
	var m = this;
	fetchBufferData(this.URL, function (buf) { modelBufferDataLoaded(m, buf); } );
}

function modelBufferDataLoaded(model, buf) {
	if (buf == null) {
		console.log(model + ": model '" + model.name + "': buffer data load failure");
		return;
	}
	
	console.log(model + ": model '" + model.name + "': buffer data loaded");
	model.buffer = buf;
	
	console.log(model + ": model '" + model.name + "': buffer: " + model.buffer);
}

/*
Model.prototype.bufferDataLoaded = function(buffer) {
	if (buffer == null) {
		console.log("model '" + this.name + "': buffer data load failure");
		return;
	}
	
	console.log("model '" + this.name + "': buffer data loaded");
	this.buffer = buffer;
	
	console.log("model '" + this.name + "': buffer: " + this.buffer);
}
*/

Model.prototype.addInstance = function(i) {
	this.instanceList.push(i);
}

Model.prototype.drawInstances = function() {
	for (var i in this.instanceList) {
		this.instanceList[i].draw(this.program);
	}
}

function Instance(model, name) {
	console.log("new instance: name=" + name);
	this.model = model;
	this.name = name;
}

Instance.prototype.draw = function(program) {

	//console.log("model: " + this.model);
	//console.log("model buffer: " + this.model.buffer);
	var buf = this.model.buffer;
	//console.log("model buffer: " + buf);
	
	/*
	console.log("buf.vertexPositionBuffer = " + buf.vertexPositionBuffer);
	console.log("buf.vertexIndexBuffer = " + buf.vertexIndexBuffer);
	console.log("buf.vertexIndexLength = " + buf.vertexIndexLength);
	console.log("buf.vertexPositionBufferItemSize = " + buf.vertexPositionBufferItemSize);
	console.log("buf.vertexIndexLength = " + buf.vertexIndexLength);
	console.log("program.aVertexPosition = " + program.aVertexPosition);
	*/
	
    gl.bindBuffer(gl.ARRAY_BUFFER, buf.vertexPositionBuffer);
   	gl.vertexAttribPointer(program.aVertexPosition, buf.vertexPositionBufferItemSize, gl.FLOAT, false, 0, 0);
	
	gl.bindBuffer(gl.ELEMENT_ARRAY_BUFFER, buf.vertexIndexBuffer);
	gl.drawElements(gl.TRIANGLES, buf.vertexIndexLength, gl.UNSIGNED_SHORT, 0 * buf.vertexIndexBufferItemSize);
}

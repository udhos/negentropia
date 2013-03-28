function Program(vertexShaderURL, fragmentShaderURL) {
	console.log("new program: vsURL=" + vertexShaderURL + " fsURL=" + fragmentShaderURL);
	this.modelList = [];
}

Program.prototype.addModel = function(modelURL) {
	var m = new Model(modelURL);
	this.modelList.push(m);
	return m;
}

Program.prototype.drawModels = function() {
	for (var m in this.modelList) {
		m.drawInstances();
	}
}

function Model(modelURL) {
	console.log("new model: URL=" + modelURL);
	this.instanceList = [];
}

Model.prototype.addInstance = function(name) {
	var i = new Instance(name);
	this.instanceList.push(i);
	return i;
}

Model.prototype.drawInstances = function() {
	for (var i in this.instanceList) {
		i.draw();
	}
}

function Instance(name) {
	console.log("new instance: name=" + name);
}

Instance.prototype.draw = function() {
}

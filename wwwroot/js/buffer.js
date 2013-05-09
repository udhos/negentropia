
function bufferAlert(msg) {
	console.log(msg);
	alert(msg);
}

function fetchBufferData(bufferURL, callbackOnDone, reverse, rescale) {
	var opaque = {
		URL: bufferURL,
		onDone: callbackOnDone,
		doReverse: reverse,
		doRescale: rescale
	};
	fetchFile(bufferURL, processBufferData, opaque);
}

function texFetchBufferData(bufferURL, callbackOnDone, reverse, rescale) {
	var opaque = {
		URL: bufferURL,
		onDone: callbackOnDone,
		doReverse: reverse,
		doRescale: rescale
	};
	fetchFile(bufferURL, texProcessBufferData, opaque);
}

function createBuffers(vertCoord, indices) {
	var buf = {};

	buf.vertexIndexLength = indices.length;
	
	buf.vertexPositionBuffer = gl.createBuffer();
	gl.bindBuffer(gl.ARRAY_BUFFER, buf.vertexPositionBuffer);
	gl.bufferData(gl.ARRAY_BUFFER, new Float32Array(vertCoord), gl.STATIC_DRAW);
    buf.vertexPositionBufferItemSize = 3; // coord x,y,z
	
	buf.vertexIndexBuffer = gl.createBuffer();
	gl.bindBuffer(gl.ELEMENT_ARRAY_BUFFER, buf.vertexIndexBuffer);
	gl.bufferData(gl.ELEMENT_ARRAY_BUFFER, new Uint16Array(indices), gl.STATIC_DRAW);
	buf.vertexIndexBufferItemSize = 2; // size of Uint16Array
	
	gl.bindBuffer(gl.ARRAY_BUFFER, null);
	gl.bindBuffer(gl.ELEMENT_ARRAY_BUFFER, null);
	
	return buf;
}

function texCreateBuffers(vertCoord, textCoord, indices) {
	var buf = {};

	buf.vertexIndexLength = indices.length;
	
	console.log("texCreateBuffers: vertCoord: " + vertCoord.length);
	console.log("texCreateBuffers: textCoord: " + textCoord.length);
	console.log("texCreateBuffers: indices: " + indices.length);
	
	buf.vertexPositionBuffer = gl.createBuffer();
	gl.bindBuffer(gl.ARRAY_BUFFER, buf.vertexPositionBuffer);
	gl.bufferData(gl.ARRAY_BUFFER, new Float32Array(vertCoord), gl.STATIC_DRAW);
    buf.vertexPositionBufferItemSize = 3; // coord x,y,z
	
	buf.vertexTextureCoordBuffer = gl.createBuffer();
	gl.bindBuffer(gl.ARRAY_BUFFER, buf.vertexTextureCoordBuffer);
	gl.bufferData(gl.ARRAY_BUFFER, new Float32Array(textCoord), gl.STATIC_DRAW);
	buf.vertexTextureCoordBufferItemSize = 2; // coord s,t
	
	buf.vertexIndexBuffer = gl.createBuffer();
	gl.bindBuffer(gl.ELEMENT_ARRAY_BUFFER, buf.vertexIndexBuffer);
	gl.bufferData(gl.ELEMENT_ARRAY_BUFFER, new Uint16Array(indices), gl.STATIC_DRAW);
	buf.vertexIndexBufferItemSize = 2; // size of Uint16Array
	
	gl.bindBuffer(gl.ARRAY_BUFFER, null);
	gl.bindBuffer(gl.ELEMENT_ARRAY_BUFFER, null);
		
	return buf;
}

function processBufferData(opaque, response) {
	//console.log("buffer data: " + opaque.URL + ": [" + response + "]");
	if (response == null) {
		bufferAlert("buffer data: FATAL ERROR: could not load from URL: " + opaque.URL);
		opaque.onDone(null);
		return;
	}
	var bufferData = JSON.parse(response);
	console.log("buffer data: " + opaque.URL + ": json parsed");
	
	if (opaque.doReverse) {
		// reverse vertex indices
		bufferData.vertInd = bufferData.vertInd.reverse();
	}
	
	if (opaque.doRescale) {
		for (var i in bufferData.vertCoord) {
			bufferData.vertCoord[i] *= opaque.doRescale;
		}
	}
	
	var buf = createBuffers(bufferData.vertCoord, bufferData.vertInd);
	
	console.log("buffer data: " + opaque.URL + ": ready: vertexIndexLength=" + buf.vertexIndexLength);
	
	opaque.onDone(buf);
}

function texProcessBufferData(opaque, response) {
	console.log("buffer data: " + opaque.URL + ": [" + response + "]");
	if (response == null) {
		bufferAlert("buffer data: FATAL ERROR: could not load from URL: " + opaque.URL);
		opaque.onDone(null);
		return;
	}
	var bufferData = JSON.parse(response);
	console.log("buffer data: " + opaque.URL + ": json parsed");
	
	if (opaque.doReverse) {
		// reverse vertex indices
		bufferData.vertInd = bufferData.vertInd.reverse();
	}
	
	if (opaque.doRescale) {
		for (var i in bufferData.vertCoord) {
			bufferData.vertCoord[i] *= opaque.doRescale;
		}
	}
	
	var buf = texCreateBuffers(bufferData.vertCoord, bufferData.textCoord, bufferData.vertInd);
	
	console.log("buffer data: " + opaque.URL + ": ready: vertexIndexLength=" + buf.vertexIndexLength);
	
	opaque.onDone(buf);
}

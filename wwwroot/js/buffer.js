
function bufferAlert(msg) {
	console.log(msg);
	alert(msg);
}

function fetchBufferData(bufferURL, callbackOnDone, reverse) {
	var opaque = {
		URL: bufferURL,
		onDone: callbackOnDone,
		doReverse: reverse
	};
	fetchFile(bufferURL, processBufferData, opaque);
}

function processBufferData(opaque, response) {
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
	
	var buf = {};

	buf.vertexIndexLength = bufferData.vertInd.length;
	
	buf.vertexPositionBuffer = gl.createBuffer();
	gl.bindBuffer(gl.ARRAY_BUFFER, buf.vertexPositionBuffer);
	gl.bufferData(gl.ARRAY_BUFFER, new Float32Array(bufferData.vertCoord), gl.STATIC_DRAW);
    buf.vertexPositionBufferItemSize = 3; // coord x,y,z
	
	buf.vertexIndexBuffer = gl.createBuffer();
	gl.bindBuffer(gl.ELEMENT_ARRAY_BUFFER, buf.vertexIndexBuffer);
	gl.bufferData(gl.ELEMENT_ARRAY_BUFFER, new Uint16Array(bufferData.vertInd), gl.STATIC_DRAW);
	buf.vertexIndexBufferItemSize = 2; // size of Uint16Array
	
	gl.bindBuffer(gl.ARRAY_BUFFER, null);
	gl.bindBuffer(gl.ELEMENT_ARRAY_BUFFER, null);
	
	console.log("buffer data: " + opaque.URL + ": ready: vertexIndexLength=" + buf.vertexIndexLength);
	
	opaque.onDone(buf);
}

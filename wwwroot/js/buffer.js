
function bufferAlert(msg) {
	console.log(msg);
	alert(msg);
}

function fetchSquare(squareFile) {
	delete neg.square;
	var opaque = {
		filename: squareFile
	};
	fetchFile(squareFile, processSquareBuffers, opaque);
}

function processSquareBuffers(opaque, response) {
	console.log(opaque.filename + ": buffers: [" + response + "]");
	if (response == null) {
		bufferAlert("buffer: FATAL ERROR: could not load: " + opaque.filename);
		return;
	}
	var square = JSON.parse(response);
	console.log(opaque.filename + ": json parsed");
	
	square.vertexPositionBuffer = gl.createBuffer();
	gl.bindBuffer(gl.ARRAY_BUFFER, square.vertexPositionBuffer);
	gl.bufferData(gl.ARRAY_BUFFER, new Float32Array(square.vertCoord), gl.STATIC_DRAW);
    square.vertexPositionBufferItemSize = 3; // coord x,y,z
	
	square.vertexIndexBuffer = gl.createBuffer();
	gl.bindBuffer(gl.ELEMENT_ARRAY_BUFFER, square.vertexIndexBuffer);
	gl.bufferData(gl.ELEMENT_ARRAY_BUFFER, new Uint16Array(square.vertInd), gl.STATIC_DRAW);
	square.vertexIndexBufferItemSize = 2; // size of Uint16Array
	
	gl.bindBuffer(gl.ARRAY_BUFFER, null);
	gl.bindBuffer(gl.ELEMENT_ARRAY_BUFFER, null);
	
	neg.square = square;
	
	console.log(opaque.filename + ": square buffers ready");
}



// Global variables
var neg = {
	debugLostContext: true,
	debugWebGL:       true,
	drawOnce:         true,
	cullBackface:     true,
	fieldOfViewY:     45
};
var gl = null;
var websocket = null;

function boot() {
	var sid = docCookies.getItem("sid");
	var statusElem = document.getElementById("ws_status");	
	
	var canvasbox = document.getElementById("canvasbox");
	neg.canvas = document.createElement('canvas');
	neg.canvas.id = "main_canvas";
	neg.canvas.width = 780;
	neg.canvas.height = 500;
	canvasbox.appendChild(neg.canvas);
	neg.canvas = document.getElementById("main_canvas");
		
	console.log("main_canvas: width=" + neg.canvas.width + " height=" + neg.canvas.height);
	//console.log("main_canvas: style width=" + neg.canvas.style.width + " height=" + neg.canvas.style.height);

	console.log("WebGL: initializing");
		
	if (neg.debugLostContext) {
		// DEBUG wrapper context
		neg.canvas = WebGLDebugUtils.makeLostContextSimulatingCanvas(neg.canvas);
	}

	gl = initGL(neg.canvas);
	if (gl) {
		console.log("WebGL: initialized");
	}
	else {
		console.log("WebGL: initialization failure");
		return;
	}
	
	var wsUri = document.getElementById("wsUri");
	
	initWebSocket(wsUri.innerHTML, statusElem, sid);
}

function animate() {
}

function render() {

        gl.viewport(0, 0, neg.canvas.width, neg.canvas.height); // define viewport size
		gl.depthRange(0.0, 1.0); // default
		
		// http://www.opengl.org/sdk/docs/man/xhtml/glClear.xml
        gl.clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT);    // clear color buffer and depth buffer

		// set perspective matrix
		// field of view y: 45 degrees
		// width to height ratio
		// view from 1.0 to 1000.0 distance units
		//
		// tan(45/2) = (h/2) / near
		// h = 2 * tan(45/2) * near
		// h = 2 * 0.414 * 1.0
		// h = 0.828
		//
        //mat4.perspective(neg.fieldOfViewY, neg.canvas.width / neg.canvas.height, 1.0, 1000.0, neg.pMatrix);

	drawSquare();
}

function loop() {

	if (neg.drawOnce) {
		console.log("loop: drawOnce ON: will render only one frame")
	}
	else {
		neg.reqId = window.requestAnimationFrame(loop); // from game-shim.js
	}
	
	animate();		// update state
	render();		// draw
}

function backfaceCulling(gl, enable) {
	if (enable) {
		gl.frontFace(gl.CCW);
		gl.cullFace(gl.BACK);
		gl.enable(gl.CULL_FACE);		
	}
	else {
		gl.disable(gl.CULL_FACE);		
	}
}

function initSquare() {

	console.log("square: init buffers");

	var square = {};
	neg.square = square;
	
	var size = 0.5;

	square.vertCoord = [
		size, size, 0.0,
		-size, size, 0.0,
		-size, -size, 0.0,
		size, -size, 0.0,
	];
	
	square.vertInd = [
		0, 1, 2,
		0, 2, 3,
	];
	
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
	
	console.log("square: init buffers: done");
}

function drawSquare() {

	var square = neg.square;
	var aVertexPosition = neg.prog.aVertexPosition;

    gl.bindBuffer(gl.ARRAY_BUFFER, square.vertexPositionBuffer);
   	gl.vertexAttribPointer(aVertexPosition, square.vertexPositionBufferItemSize, gl.FLOAT, false, 0, 0);
    gl.enableVertexAttribArray(aVertexPosition);
	
	gl.bindBuffer(gl.ELEMENT_ARRAY_BUFFER, square.vertexIndexBuffer);

	gl.drawElements(gl.TRIANGLES, square.vertInd.length, gl.UNSIGNED_SHORT, 0 * square.vertexIndexBufferItemSize);

	// clean up
	gl.bindBuffer(gl.ARRAY_BUFFER, null);
	gl.bindBuffer(gl.ELEMENT_ARRAY_BUFFER, null);				
}

function initBuffers() {
	initSquare();
}

function initContext() {

	/*
		Load program shaders
	*/
	console.log("shader program: loading");
	neg.prog = loadProgram(gl, "shader-vs", "shader-fs");
	if (!neg.prog) {
		console.log("shader program: failure");
		return;
	}
	console.log("shader program: loaded");
	
	initBuffers();
	
   	gl.clearColor(0.5, 0.5, 0.5, 1.0);	// clear color
    gl.enable(gl.DEPTH_TEST);			// perform depth testing
	gl.depthFunc(gl.LESS);				// gl.LESS is default depth test
		
	backfaceCulling(gl, neg.cullBackface);
}

function play() {
	boot();
	
	if (!gl) {
		return;
	}

	initContext();
	
	if (!neg.prog) {
		return;
	}
	
	loop();
}

function main() {
	play();
	
	console.log("main: exit: NOT REACHED"); // prevented by drawing loop
}

function initGL(canvas) {
	var ctx = WebGLUtils.setupWebGL(canvas);
	if (!ctx) {
		console.log("initGL: failure");
		return null;
	}
		
	if (neg.debugWebGL) {
		// DEBUG-only wrapper context -- performance PENALTY!
		ctx = WebGLDebugUtils.makeDebugContext(ctx);
	}
	
	return ctx;
}

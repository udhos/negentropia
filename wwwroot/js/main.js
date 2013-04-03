
// Global variables
var neg = {
	debugLostContext:    true,
	debugWebGL:          true,
	drawOnce:            false,
	cullBackface:        true,
	fieldOfViewY:        45,
	ongoingImageLoads:   [],
	programList:         [],
	shaderCache:         {}
};
var gl = null;
var websocket = null;

// Stats.js
function initStats() {
	neg.stats = new Stats();

	neg.stats.setMode(0); // 0: fps, 1: ms

	//neg.stats.domElement.style.position = 'inherit';
		
	var framerate = document.getElementById("framerate");
	if (framerate.appendChild) {
	
		// remove all existing node children
	    while (framerate.childNodes.length > 0) {
			framerate.removeChild(framerate.firstChild);       
		}
			
		// attach child
		framerate.appendChild(neg.stats.domElement);
	}
}
	
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

	if (neg.debugLostContext) {
		initDebugLostContext(neg.canvas);
	}
	
	var wsUri = document.getElementById("wsUri");
	
	initWebSocket(wsUri.innerHTML, statusElem, sid);
	
	initStats();
}

function animate() {
    // TODO: FIXME: WRITEME: update state
}

function render() {
		
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
	// neg.canvasAspect = neg.canvas.width / neg.canvas.height;
	//
    //mat4.perspective(neg.fieldOfViewY, neg.canvasAspect, 1.0, 1000.0, neg.pMatrix);

	//drawSquare();
	
	for (var p in neg.programList) {
		neg.programList[p].drawModels();
	}
}

function loop() {
	neg.stats.update();         // update framerate statistics

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

/*
function drawSquare() {

	if (!('programSquare' in neg)) {
		// square shader program is not loaded yet
		return;
	}

	if (!('square' in neg)) {
		// square buffers are not loaded yet
		return;
	}
	
	gl.useProgram(neg.programSquare.shaderProgram);

	var aVertexPosition = neg.programSquare.aVertexPosition;
	var square = neg.square;

    gl.bindBuffer(gl.ARRAY_BUFFER, square.vertexPositionBuffer);
   	gl.vertexAttribPointer(aVertexPosition, square.vertexPositionBufferItemSize, gl.FLOAT, false, 0, 0);
    gl.enableVertexAttribArray(aVertexPosition);
	
	gl.bindBuffer(gl.ELEMENT_ARRAY_BUFFER, square.vertexIndexBuffer);
	gl.drawElements(gl.TRIANGLES, square.vertInd.length, gl.UNSIGNED_SHORT, 0 * square.vertexIndexBufferItemSize);

    gl.disableVertexAttribArray(aVertexPosition);

	// clean up
	gl.bindBuffer(gl.ARRAY_BUFFER, null);
	gl.bindBuffer(gl.ELEMENT_ARRAY_BUFFER, null);				
}
*/

/*
function initBuffers() {
	fetchSquare("/mesh/square.json");
}
*/

/*
function squareProgramLoaded(prog) {
	if ('shaderProgram' in prog) {
		console.log("main.js: square shader program: ready");

		// save vertex attribute location
		prog.aVertexPosition = gl.getAttribLocation(prog.shaderProgram, "aVertexPosition");

		neg.programSquare = prog;
	}
	else {
		console.log("main.js: square shader program: failure");	
	}
}
*/

function initContext() {

	// Async request for shader program
	//fetchProgramFromURL("/shader/min_vs.txt", "/shader/min_fs.txt", squareProgramLoaded);
	
	//initBuffers();
	
	neg.programList = []; // drop existing full programs
	neg.shaderCache = {}; // drop existing compiled shaders

	var squareProgram = new Program("/shader/min_vs.txt", "/shader/min_fs.txt");
	neg.programList.push(squareProgram);
	var squareModel = new Model(squareProgram, "/mesh/square.json");
	squareProgram.addModel(squareModel);
	var squareInstance = new Instance(squareModel);
	squareModel.addInstance(squareInstance);

	// create 2nd program after 2 secs (time for the first program to populate the shader cache)
	setTimeout(function() {
		var squareProgram2 = new Program("/shader/min_vs.txt", "/shader/min2_fs.txt");
		neg.programList.push(squareProgram2);
		var squareModel2 = new Model(squareProgram2, "/mesh/square2.json");
		squareProgram2.addModel(squareModel2);
		var squareInstance2 = new Instance(squareModel2);
		squareModel2.addInstance(squareInstance2);
	}, 2000);
		
   	gl.clearColor(0.5, 0.5, 0.5, 1.0);	// clear color
    gl.enable(gl.DEPTH_TEST);			// perform depth testing
	gl.depthFunc(gl.LESS);				// gl.LESS is default depth test
	gl.depthRange(0.0, 1.0);            // default
	
	// define viewport size
    gl.viewport(0, 0, neg.canvas.width, neg.canvas.height);
	neg.canvasAspect = neg.canvas.width / neg.canvas.height; // save aspect for render loop mat4.perspective
		
	backfaceCulling(gl, neg.cullBackface);
	
	loop(); // render loop
}

function main() {
	boot();
	
	if (!gl) {
		return;
	}

	initContext(); // calls loop()
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

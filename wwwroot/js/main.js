
// Global variables
var neg = {
	debugLostContext: true,
	debugWebGL:       true
};
var gl = null;

function boot() {
	var sid = docCookies.getItem("sid");
	var statusElem = document.getElementById("ws_status");	
	
	var canvasbox = document.getElementById("canvasbox");
	neg.canvas = document.createElement('canvas');
	neg.canvas.id = "main_canvas";
	canvasbox.appendChild(neg.canvas);
	neg.canvas = document.getElementById("main_canvas");

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
	
	initWebSocket(wsUri.innerText, statusElem, sid);
}

function loop() {
	console.log("loop: FIXME WRITEME");
}

function play() {
	boot();
	
	if (!gl) {
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

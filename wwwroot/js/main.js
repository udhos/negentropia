
// Global variables
var neg = {
	debugLostContext: true,
	debugWebGL:       true
};
function main() {
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

	initGL(neg.canvas);
	
	if (neg.gl) {
		console.log("WebGL: initialized");
	}
	else {
		console.log("WebGL: initialization failure");
		return;
	}
	
	var wsUri = document.getElementById("wsUri");
	
	initWebSocket(wsUri.innerText, statusElem, sid);
}

function initGL(canvas) {
	neg.gl = WebGLUtils.setupWebGL(canvas);
	if (!neg.gl) {
		console.log("initGL: failure");
		return;
	}
		
	if (neg.debugWebGL) {
		// DEBUG-only wrapper context -- performance PENALTY!
		neg.gl = WebGLDebugUtils.makeDebugContext(neg.gl);
	}
}


function main() {
	var sid = docCookies.getItem("sid");
	var statusElem = document.getElementById("ws_status");	
	initWebSocket(neg.wsUri, statusElem, sid);
	
	var canvasbox = document.getElementById("canvasbox");
	neg.canvas = document.createElement('canvas');
	neg.canvas.id = "main_canvas";
	canvasbox.appendChild(neg.canvas);
	neg.canvas = document.getElementById("main_canvas");
	
	if (neg.debugLostContext) {
		// DEBUG wrapper context
		neg.canvas = WebGLDebugUtils.makeLostContextSimulatingCanvas(neg.canvas);
	}

	initGL(neg.canvas);
}

function initGL(canvas) {
	neg.gl = WebGLUtils.setupWebGL(canvas);
		
	if (neg.debugWebGL) {
		// DEBUG-only wrapper context -- performance PENALTY!
		neg.gl = WebGLDebugUtils.makeDebugContext(neg.gl);
	}

	console.log("WebGL initialized");
}


function shaderAlert(msg) {
	console.log(msg);
	alert(msg);
}

function compileShader(gl, shaderString, shaderType) {
	var shader = gl.createShader(shaderType);
    gl.shaderSource(shader, shaderString);
    gl.compileShader(shader);

    if (!gl.getShaderParameter(shader, gl.COMPILE_STATUS) && !gl.isContextLost()) {
        shaderAlert("Error compiling shader: " + gl.getShaderInfoLog(shader));
		gl.deleteShader(shader);
        return null;
    }
	
	return shader;
}

function linkProg(prog, gl, vertexShader, fragmentShader) {

	// link program
    var shaderProgram = gl.createProgram();
    gl.attachShader(shaderProgram, vertexShader);
    gl.attachShader(shaderProgram, fragmentShader);
    gl.linkProgram(shaderProgram);

    if (!gl.getProgramParameter(shaderProgram, gl.LINK_STATUS) && !gl.isContextLost()) {
        shaderAlert("Error linking program: " + gl.getProgramInfoLog(shaderProgram));
		return;
    }

	// save shader program
	prog.shaderProgram = shaderProgram;
}
		
function tryLinkProgram(prog) {
	
	if (!prog.vertexShader || !prog.fragmentShader) {
		// not ready
		return;
	}
	
	console.log("shader program: linking");
	linkProg(prog, gl, prog.vertexShader, prog.fragmentShader);

	shaderOngoingStop(prog);
}

function processVertexShader(opaque, response) {
	var prog = opaque;

	console.log(prog.vsFile + ": vertex shader: [" + response + "]");
	if (response == null) {
		shaderAlert("vertex shader: FATAL ERROR: could not load");
		shaderOngoingStop(prog);
		return;
	}
	prog.vertexShader = compileShader(gl, response, gl.VERTEX_SHADER);
	tryLinkProgram(prog);
}

function processFragmentShader(opaque, response) {
	var prog = opaque;
	
	console.log(prog.fsFile + ": fragment shader: [" + response + "]");
	if (response == null) {
		shaderAlert("fragment shader: FATAL ERROR: could not load");
		shaderOngoingStop(prog);
		return;
	}
	prog.fragmentShader = compileShader(gl, response, gl.FRAGMENT_SHADER);
	tryLinkProgram(prog);
}

function shaderOngoingStart(prog) {
}

function shaderOngoingStop(prog) {
	if ('shaderProgram' in prog) {
		console.log("shader.js: shader program: ready");
	}
	else {
		console.log("shader.js: shader program: failure");	
	}

	//neg.ongoingProgramLoads.splice(neg.ongoingProgramLoads.indexOf(prog), 1);
	
	// callback
	prog.callbackOnDone(prog);	
}

function fetchProgramFromURL(vs, fs, callbackOnDone) {
	var prog = {};
	prog.vsFile = vs;
	prog.fsFile = fs;
	prog.callbackOnDone = callbackOnDone;
	
	//neg.ongoingProgramLoads = [];
	
	shaderOngoingStart(prog);
	
	fetchFile(vs, processVertexShader, prog);
	fetchFile(fs, processFragmentShader, prog);
}

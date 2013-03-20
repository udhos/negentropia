
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

	// use program
    gl.useProgram(prog.shaderProgram);

	// save vertex attribute location
    prog.aVertexPosition = gl.getAttribLocation(prog.shaderProgram, "aVertexPosition");
}
		
function tryLinkProgram() {
	var prog = neg.prog;
	
	if (prog.vertexShader && prog.fragmentShader) {
		console.log("shader program: linking");
		linkProg(prog, gl, prog.vertexShader, prog.fragmentShader);
		if ('shaderProgram' in prog) {
			console.log("shader program: ready");
		}
	}
}

function processVertexShader(opaque, response) {
	console.log(neg.prog.vsFile + ": vertex shader: [" + response + "]");
	if (response == null) {
		shaderAlert("vertex shader: FATAL ERROR: could not load");
		return;
	}
	neg.prog.vertexShader = compileShader(gl, response, gl.VERTEX_SHADER);
	tryLinkProgram();
}

function processFragmentShader(opaque, response) {
	console.log(neg.prog.fsFile + ": fragment shader: [" + response + "]");
	if (response == null) {
		shaderAlert("fragment shader: FATAL ERROR: could not load");
		return;
	}
	neg.prog.fragmentShader = compileShader(gl, response, gl.FRAGMENT_SHADER);
	tryLinkProgram();
}

function fetchProgramFromURL(vs, fs) {
	neg.prog = {};
	neg.prog.vsFile = vs;
	neg.prog.fsFile = fs;
	fetchFile(vs, processVertexShader, null);
	fetchFile(fs, processFragmentShader, null);
}

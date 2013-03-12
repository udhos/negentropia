
function shaderAlert(msg) {
	console.log(msg);
	alert(msg);
}

    function loadShaderFromDOM(gl, id) {
        var shaderScript = document.getElementById(id);
        if (!shaderScript) {
            return null;
        }

        var str = "";
        var k = shaderScript.firstChild;
        while (k) {
            if (k.nodeType == 3) { // 3 corresponds to TEXT_NODE
                str += k.textContent;
            }
            k = k.nextSibling;
        }

        var shader;
        if (shaderScript.type == "x-shader/x-fragment") {
            shader = gl.createShader(gl.FRAGMENT_SHADER);
        } else if (shaderScript.type == "x-shader/x-vertex") {
            shader = gl.createShader(gl.VERTEX_SHADER);
        } else {
            return null;
        }

        gl.shaderSource(shader, str);
        gl.compileShader(shader);

        if (!gl.getShaderParameter(shader, gl.COMPILE_STATUS) && !gl.isContextLost()) {
            shaderAlert("Error compiling shader: " + gl.getShaderInfoLog(shader));
			gl.deleteShader(shader);
            return null;
        }

        return shader;
    }


	// prog = loadProgram(gl, "shader-vs", "shader-fs")
    function loadProgram(gl, vs, fs) {
	
		// load shaders
        var vertexShader = loadShaderFromDOM(gl, vs);
        var fragmentShader = loadShaderFromDOM(gl, fs);
		
		var prog = {};

		// link program
        var shaderProgram = gl.createProgram();
        gl.attachShader(shaderProgram, vertexShader);
        gl.attachShader(shaderProgram, fragmentShader);
        gl.linkProgram(shaderProgram);

        if (!gl.getProgramParameter(shaderProgram, gl.LINK_STATUS) && !gl.isContextLost()) {
            shaderAlert("Error linking program: " + gl.getProgramInfoLog(shaderProgram));
			return null;
        }

		// use program
		prog.shaderProgram = shaderProgram;
        gl.useProgram(prog.shaderProgram);

		// save vertex attribute location
        prog.aVertexPosition = gl.getAttribLocation(prog.shaderProgram, "aVertexPosition");
				
		return prog;
    }

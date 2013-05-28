
// Global variables
var neg = {
	debugLostContext:    true,
	debugWebGL:          true,
	drawOnce:            false,
	cullBackface:        false,
	fieldOfViewY:        45,
	ongoingImageLoads:   [],
	programList:         [],
	shaderCache:         {},
	pMatrix:             mat4.create(),
	
	// camera tests
	angleY:              0,
	deltaY:              1,
	eye:                 [0,0,15],
    center:	             [0,0,0],
	up:                  [0,1,0],
	scale:               1
};
var gl = null;
var websocket = null;
var assetRoot = "/";
var assetPath = {
	mesh:    assetRoot + "mesh/",
	mtl:     assetRoot + "mtl/",
	obj:     assetRoot + "obj/",
	shader:  assetRoot + "shader/",
	texture: assetRoot + "texture/",
}

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
	
	var camOrbitRadius = 15;
	var degreesPerSec  = 30;
	
	neg.angleY += neg.deltaY * degreesPerSec / 60;
	neg.angleY %= 360;
	var radY = neg.angleY * Math.PI / 180;
	
	neg.eye = [ camOrbitRadius * Math.sin(radY), 0, camOrbitRadius * Math.cos(radY) ];
	neg.center = [ 0, 0, 0 ];
	
	neg.scale = 15 * Math.abs(Math.sin(radY)) + 1;

	for (var p in neg.programList) {
		neg.programList[p].animate();
	}
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
    mat4.perspective(neg.fieldOfViewY, neg.canvasAspect, 1, 1000, neg.pMatrix);

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

function initSquareWhite() {
	var squareProgram = new Program(assetPath.shader + "clip_vs.txt", assetPath.shader + "clip_fs.txt");
	neg.programList.push(squareProgram);
	var squareModel = new Model(squareProgram, assetPath.mesh + "square.json");
	squareProgram.addModel(squareModel);
	var squareInstance = new Instance(squareModel, [0.0, 0.0, 0.0]);
	squareModel.addInstance(squareInstance);
	return squareProgram;
}

function initSquareBlue() {
	var squareProgram2 = new Program(assetPath.shader + "clip_vs.txt", assetPath.shader + "clip2_fs.txt");
	neg.programList.push(squareProgram2);		
	var squareModel2 = new Model(squareProgram2, assetPath.mesh + "square2.json");
	squareProgram2.addModel(squareModel2);
	var squareInstance2 = new Instance(squareModel2, [0.0, 0.0, 0.0]);
	squareModel2.addInstance(squareInstance2);
	return squareProgram2;
}

function initSquareRed() {
	var squareProgram3 = new Program(assetPath.shader + "clip_vs.txt", assetPath.shader + "clip3_fs.txt");
	neg.programList.push(squareProgram3);
	var squareModel3 = new Model(squareProgram3, assetPath.mesh + "square3.json");
	squareProgram3.addModel(squareModel3);
	var squareInstance3 = new Instance(squareModel3);
	squareModel3.addInstance(squareInstance3);
	return squareProgram3;
}

function initSquares() {
	var whiteSquareProgram = initSquareWhite();
	neg.programList.push(whiteSquareProgram);
	whiteSquareProgram.fetch();
	
	// create 2nd program after 2 secs (time for the first program to populate the shader cache)
	var blueSquareProgram = initSquareBlue();
	neg.programList.push(blueSquareProgram);
	setTimeout(function() { blueSquareProgram.fetch(); }, 2000);

	var redSquareProgram = initSquareRed();
	neg.programList.push(redSquareProgram);
	redSquareProgram.fetch();
}

function initSkybox() {
	var skyboxProgram = new SkyboxProgram(assetPath.shader + "skybox_vs.txt", assetPath.shader + "skybox_fs.txt");
	neg.programList.push(skyboxProgram);
	var skyboxModel = new SkyboxModel(skyboxProgram, assetPath.mesh + "cube.json", true, 0);
	skyboxModel.addCubemapFace(gl.TEXTURE_CUBE_MAP_POSITIVE_X, assetPath.texture + 'space_rt.jpg');
	skyboxModel.addCubemapFace(gl.TEXTURE_CUBE_MAP_NEGATIVE_X, assetPath.texture + 'space_lf.jpg');
	skyboxModel.addCubemapFace(gl.TEXTURE_CUBE_MAP_POSITIVE_Y, assetPath.texture + 'space_up.jpg');
	skyboxModel.addCubemapFace(gl.TEXTURE_CUBE_MAP_NEGATIVE_Y, assetPath.texture + 'space_dn.jpg');
	skyboxModel.addCubemapFace(gl.TEXTURE_CUBE_MAP_POSITIVE_Z, assetPath.texture + 'space_fr.jpg');
	skyboxModel.addCubemapFace(gl.TEXTURE_CUBE_MAP_NEGATIVE_Z, assetPath.texture + 'space_bk.jpg');	
	skyboxProgram.addModel(skyboxModel);
	var skyboxInstance = new SkyboxInstance(skyboxModel, [0, 0, 0], 1.0);
	skyboxModel.addInstance(skyboxInstance);
}

function onObjDone(opaque, response) {

	if (response == null) {
		console.log("onObjDone: fetch FAILURE: " + opaque.URL);
		return;
	}

	console.log("onObjDone: fetch done: " + opaque.URL);
	
	var airship = new obj_loader.Mesh(response);
	
	console.log("onObjDone: parsing done: " + opaque.URL);
			
	var mod = new Model(opaque.program, null, false, airship);
	opaque.program.addModel(mod);
	var inst = new Instance(mod, opaque.center);
	mod.addInstance(inst);
}

function initAirship() {
	var prog = new Program(assetPath.shader + "simple_vs.txt", assetPath.shader + "simple_fs.txt");
	neg.programList.push(prog);
	prog.fetch();

	var objURL = assetPath.obj + "airship.obj";
	var opaque = {
		URL: objURL,
		program: prog,
		center: [-8.0, 0.0, 0.0]
	};
	console.log("initAirship: loading OBJ from: " + objURL);
	fetchFile(objURL, onObjDone, opaque);
}

function onMaterialDone(opaque, response) {
	if (response == null) {
		console.log("onMaterialDone: fetch FAILURE: " + opaque.URL);
		return;
	}
	
	console.log("onMaterialDone: " + opaque.URL);
	
	var mtllib = mtllib_parse(response);
	
	console.log("onMaterialDone mtllib = " + mtllib);
	
	var usemtl = opaque.airship.usemtl;
	
	console.log("onMaterialDone usemtl = " + usemtl);

	var map_Kd = mtllib[usemtl].map_Kd;
	
	var textureName = assetPath.texture + map_Kd;

	console.log("onMaterialDone textureName = " + textureName);
	
	var tex = new Texture(neg.textureTable, 0, opaque.airship.indices.length, textureName);
	opaque.mod.addTexture(tex);
}

function onTexObjDone(opaque, response) {

	if (response == null) {
		console.log("onTexObjDone: fetch FAILURE: " + opaque.URL);
		return;
	}

	console.log("onTexObjDone: fetch done: " + opaque.URL);
	
	var obj_airship = new obj_loader.Mesh(response);
	
	console.log("onTexObjDone: parsing done: " + opaque.URL);
			
	var mod = new TexModel(opaque.program, null, false, obj_airship);
	opaque.program.addModel(mod);
	
	console.log("onTexObjDone: mod.buffer.vertexPositionBufferItemSize = " + mod.buffer.vertexPositionBufferItemSize);
	console.log("onTexObjDone: mod.buffer.vertexTextureCoordBuffer = " + mod.buffer.vertexTextureCoordBuffer);
	console.log("onTexObjDone: mod.buffer.vertexIndexBuffer = " + mod.buffer.vertexIndexBuffer);

	if (opaque.forceSolidColor) {
		console.log("onTexObjDone: FORCED SOLID COLOR");
		var tex = new Texture(neg.textureTable, 0, obj_airship.indices.length, "FORCED-SOLID-COLOR");
		mod.addTexture(tex);
		
		var inst = new TexInstance(mod, opaque.center);
		mod.addInstance(inst);
		
		return;
	}
	
	if ('mtllib' in obj_airship) {
		console.log("onTexObjDone: mtllib FOUND: " + obj_airship.mtllib);
		var opaqueMaterial = {
			URL: assetPath.mtl + obj_airship.mtllib,
			airship: obj_airship,
			mod: mod
		};
		fetchFile(opaqueMaterial.URL, onMaterialDone, opaqueMaterial);
		
		var inst = new TexInstance(mod, opaque.center);
		mod.addInstance(inst);
		
		return;
	}
	
	console.log("onTexObjDone: mtllib NOT FOUND: " + opaque.URL);
}

function initAirshipTex() {
	var prog = new TexProgram(assetPath.shader + "simpleTex_vs.txt", assetPath.shader + "simpleTex_fs.txt");
	neg.programList.push(prog);
	prog.fetch();

	var objURL = assetPath.obj + "airship.obj"; 
	var opaque = {
		URL: objURL,
		program: prog,
		center: [0.0, 0.0, 0.0]
	};
	console.log("initAirshipTex: loading OBJ from: " + objURL);
	fetchFile(objURL, onTexObjDone, opaque)
	
	var opaque2 = {
		URL: objURL,
		forceSolidColor: true,
		program: prog,
		center: [8.0, 0.0, 0.0]
	};
	console.log("initAirshipTex: loading OBJ from: " + objURL);
	fetchFile(objURL, onTexObjDone, opaque2);
}

function initShips() {
	initAirship();
	initAirshipTex();
}

function initContext() {
	
	neg.programList  = []; // drop existing full programs
	neg.shaderCache  = {}; // drop existing compiled shaders
	neg.textureTable = {}; // drop existing textures

	initSquares();
	initShips();
	initSkybox();
		
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

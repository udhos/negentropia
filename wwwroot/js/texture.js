
function textureFinishedLoading(textureTable, textureName, onTextureLoad, texture, image) {
  console.log("textureFinishedLoading: " + textureName);
  
  gl.bindTexture(gl.TEXTURE_2D, texture);
  gl.pixelStorei(gl.UNPACK_FLIP_Y_WEBGL, true);

  gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST);
  gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST);

	gl.texImage2D(gl.TEXTURE_2D, 0, gl.RGBA, gl.RGBA, gl.UNSIGNED_BYTE, image);

	/*
	While OpenGL 2.0 and later for the desktop offer full support for non-power-of-two (NPOT) textures, OpenGL ES 2.0 and WebGL have only limited NPOT support. The restrictions are defined in Sections 3.8.2, "Shader Execution", and 3.7.11, "Mipmap Generation", of the OpenGL ES 2.0 specification, and are summarized here:
	generateMipmap(target) generates an INVALID_OPERATION error if the level 0 image of the texture currently bound to target has an NPOT width or height.
	Sampling an NPOT texture in a shader will produce the RGBA color (0, 0, 0, 1) if:
	The minification filter is set to anything but NEAREST or LINEAR: in other words, if it uses one of the mipmapped filters.
	The repeat mode is set to anything but CLAMP_TO_EDGE; repeating NPOT textures are not supported.
	If your application doesn't require the REPEAT wrap mode, and can tolerate the lack of mipmaps, then you can simply configure the WebGLTexture object appropriately at creation time:
	var texture = gl.createTexture();
	gl.bindTexture(gl.TEXTURE_2D, texture);
	gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR);
	gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE);
	gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE);

	http://www.khronos.org/webgl/wiki/WebGL_and_OpenGL_Differences
	*/  
	gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE);
	gl.texParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE);  
	
	gl.bindTexture(gl.TEXTURE_2D, null);
	
	onTextureLoad(textureTable, textureName, texture); // callback
}

function textureLoadError(url, image) {
	neg.ongoingImageLoads.splice(neg.ongoingImageLoads.indexOf(image), 1);
	console.log("error loading texture: " + url);
}

function textureLoadAborted(url, image) {
	neg.ongoingImageLoads.splice(neg.ongoingImageLoads.indexOf(image), 1);
	console.log("aborted loading texture: " + url);
}

function loadImageForTexture(textureTable, textureName, onTextureLoad, texture) {
  var image = new Image();
  image.onerror = function() { textureLoadError(textureName, image); }
  image.onabort = function() { textureLoadAborted(textureName, image); }
  image.onload = function() {
    neg.ongoingImageLoads.splice(neg.ongoingImageLoads.indexOf(image), 1);
    textureFinishedLoading(textureTable, textureName, onTextureLoad, texture, image);
  }
  neg.ongoingImageLoads.push(image);
  image.src = textureName;
}

function loadTexture(textureTable, textureName, onTextureLoad) {

	console.log("loadTexture: " + textureName);
	
	var texture = gl.createTexture();
	if (texture == null) {
		console.log("loadTexture: failure creating texture: " + textureName);
		return;
	}
	
	loadImageForTexture(textureTable, textureName, onTextureLoad, texture);
}
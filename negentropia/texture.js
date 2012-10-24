
function textureFinishedLoading(url, image, texture) {
  console.log("loaded: " + url);
  
  neg.gl.bindTexture(neg.gl.TEXTURE_2D, texture);
  neg.gl.pixelStorei(neg.gl.UNPACK_FLIP_Y_WEBGL, true);
  neg.gl.texImage2D(neg.gl.TEXTURE_2D, 0, neg.gl.RGBA, neg.gl.RGBA, neg.gl.UNSIGNED_BYTE, image);

  neg.gl.texParameteri(neg.gl.TEXTURE_2D, neg.gl.TEXTURE_MAG_FILTER, neg.gl.NEAREST);
  neg.gl.texParameteri(neg.gl.TEXTURE_2D, neg.gl.TEXTURE_MIN_FILTER, neg.gl.NEAREST);

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
	neg.gl.texParameteri(neg.gl.TEXTURE_2D, neg.gl.TEXTURE_WRAP_S, neg.gl.CLAMP_TO_EDGE);
	neg.gl.texParameteri(neg.gl.TEXTURE_2D, neg.gl.TEXTURE_WRAP_T, neg.gl.CLAMP_TO_EDGE);  
	
  neg.gl.bindTexture(neg.gl.TEXTURE_2D, null); 
}

function textureLoadError(url, image) {
	neg.ongoingImageLoads.splice(neg.ongoingImageLoads.indexOf(image), 1);
	console.log("error loading: " + url);
}

function textureLoadAborted(url, image) {
	neg.ongoingImageLoads.splice(neg.ongoingImageLoads.indexOf(image), 1);
	console.log("aborted loading: " + url);
}

function loadImageForTexture(url, texture) {
  var image = new Image();
  image.onerror = function() { textureLoadError(url, image); }
  image.onabort = function() { textureLoadAborted(url, image); }
  image.onload = function() {
    neg.ongoingImageLoads.splice(neg.ongoingImageLoads.indexOf(image), 1);
    textureFinishedLoading(url, image, texture);
  }
  neg.ongoingImageLoads.push(image);
  image.src = url;
}

function loadTexture(name) {
	if (!neg.textureTable.hasOwnProperty(name)) {
		// texture cache miss: do fetch
		var texture = neg.gl.createTexture();
		if (texture != null) {
			// got a new texture object
			neg.textureTable[name] = texture;
			loadImageForTexture("texture/" + name, neg.textureTable[name]);
		}
	}
}

function cubemapFaceFinishedLoading(targetFace, url, image, texture) {
	console.log("loaded cubemap face: " + url);

    neg.gl.bindTexture(neg.gl.TEXTURE_CUBE_MAP, texture);
    neg.gl.texImage2D(targetFace, 0, neg.gl.RGBA, neg.gl.RGBA, neg.gl.UNSIGNED_BYTE, image);
	neg.gl.bindTexture(neg.gl.TEXTURE_CUBE_MAP, null);
}

function loadCubemapFace(targetFace, texture, url) {
	var image = new Image();
	image.onload = function() {
		neg.ongoingImageLoads.splice(neg.ongoingImageLoads.indexOf(image), 1);
		cubemapFaceFinishedLoading(targetFace, url, image, texture);
	}
    neg.ongoingImageLoads.push(image);
	image.src = url;
};

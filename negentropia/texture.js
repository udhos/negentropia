
function textureFinishedLoading(url, image, texture) {
  console.log("loaded: " + url);
  
  neg.gl.bindTexture(neg.gl.TEXTURE_2D, texture);
  neg.gl.pixelStorei(neg.gl.UNPACK_FLIP_Y_WEBGL, true);
  neg.gl.texImage2D(neg.gl.TEXTURE_2D, 0, neg.gl.RGBA, neg.gl.RGBA, neg.gl.UNSIGNED_BYTE, image);

  neg.gl.texParameteri(neg.gl.TEXTURE_2D, neg.gl.TEXTURE_MAG_FILTER, neg.gl.NEAREST);
  neg.gl.texParameteri(neg.gl.TEXTURE_2D, neg.gl.TEXTURE_MIN_FILTER, neg.gl.NEAREST);
	
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

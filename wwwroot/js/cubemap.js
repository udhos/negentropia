function cubemapFaceFinishedLoading(targetFace, url, image, texture) {
	console.log("loaded cubemap face: " + url);

    gl.bindTexture(gl.TEXTURE_CUBE_MAP, texture);
	
	gl.texParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MAG_FILTER, gl.NEAREST);
	gl.texParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MIN_FILTER, gl.NEAREST);
	
    gl.texImage2D(targetFace, 0, gl.RGBA, gl.RGBA, gl.UNSIGNED_BYTE, image);

	gl.texParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE);
	gl.texParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE);
	
	gl.bindTexture(gl.TEXTURE_CUBE_MAP, null);
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

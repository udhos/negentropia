function Picker(gl, canvas_width, canvas_height) {
	this.gl = gl;
	this.width = canvas_width;
	this.height = canvas_height;
	this.texture = null;
	this.framebuffer = null;
	this.renderbuffer = null;
    
	this.configure();
}

Picker.prototype.resize = function(canvas_width, canvas_height) {

	var gl = this.gl;
	this.width = canvas_width;
	this.height = canvas_height;
	var width = canvas_width;
	var height = canvas_height;
   
	gl.bindTexture(gl.TEXTURE_2D, this.texture);
    gl.texImage2D(gl.TEXTURE_2D, 0, gl.RGBA, width, height, 0, gl.RGBA, gl.UNSIGNED_BYTE, null);
	
	// 2. Init Render Buffer
    gl.bindRenderbuffer(gl.RENDERBUFFER, this.renderbuffer);
    gl.renderbufferStorage(gl.RENDERBUFFER, gl.DEPTH_COMPONENT16, width, height);
}

Picker.prototype.configure = function() {

	var gl = this.gl;
	var width = this.width;
	var height = this.height;
	
	// 1. Init Picking Texture
	this.texture = gl.createTexture();
	gl.bindTexture(gl.TEXTURE_2D, this.texture);
    gl.texImage2D(gl.TEXTURE_2D, 0, gl.RGBA, width, height, 0, gl.RGBA, gl.UNSIGNED_BYTE, null);
	
	// 2. Init Render Buffer
	this.renderbuffer = gl.createRenderbuffer();
    gl.bindRenderbuffer(gl.RENDERBUFFER, this.renderbuffer);
    gl.renderbufferStorage(gl.RENDERBUFFER, gl.DEPTH_COMPONENT16, width, height);	
    
    // 3. Init Frame Buffer
    this.framebuffer = gl.createFramebuffer();
	gl.bindFramebuffer(gl.FRAMEBUFFER, this.framebuffer);
	gl.framebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_2D, this.texture, 0);
    gl.framebufferRenderbuffer(gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT, gl.RENDERBUFFER, this.renderbuffer);

	// 4. Clean up
	gl.bindTexture(gl.TEXTURE_2D, null);
    gl.bindRenderbuffer(gl.RENDERBUFFER, null);
    gl.bindFramebuffer(gl.FRAMEBUFFER, null);
};

Picker.prototype.compare = function(readout, color) {
    return (Math.abs(Math.round(color[0]*255) - readout[0]) <= 1 &&
			Math.abs(Math.round(color[1]*255) - readout[1]) <= 1 && 
			Math.abs(Math.round(color[2]*255) - readout[2]) <= 1);
}

Picker.prototype.find = function(coords) {

	var gl = this.gl;
	
	// read one pixel from offscreen framebuffer
	var readout = new Uint8Array(1 * 1 * 4);
	gl.bindFramebuffer(gl.FRAMEBUFFER, this.framebuffer);
	gl.readPixels(coords[0], coords[1], 1, 1, gl.RGBA, gl.UNSIGNED_BYTE, readout);
	gl.bindFramebuffer(gl.FRAMEBUFFER, null);
	
	//console.info('offscreen pixel at (' + coords[0] + ',' + coords[1] + ') = readout ('+ readout[0]+','+ readout[1]+','+ readout[2] + ',' + readout[3] + ')');
	
	//return [readout[0], readout[1], readout[2], readout[3]];
	return readout;
};





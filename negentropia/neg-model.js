
function NegFace(vertexIndexOffset, vertexIndexNumber, textureName) {
	this.vertexIndexOffset = vertexIndexOffset;
	this.vertexIndexNumber = vertexIndexNumber;
	this.textureName = textureName;
}

function NegModelInstance(model, center, scale, animation) {
	this.model = model;
	this.center = center;
	this.scale = scale;
	this.animation = animation;
	
	this.mvMatrix  = mat4.create();	
	this.rotationAngle = 0;
	
	this.name = 'Model ' + generateNameId();
	
	//this.rotationMatrix = mat4.create();
	//mat4.identity(this.rotationMatrix);
}

NegModelInstance.prototype.rotate1 = function(elapsedTime) {
	this.rotationAngle += 45 * elapsedTime / 1000.0; // degrees per second
	this.rotationAngle %= 360;
}

NegModelInstance.prototype.coord = function() {
	return this.center.slice(0);
}

NegModelInstance.prototype.draw = function(offscreen) {

	// Won´t draw object to offscreen framebuffer without picking color
	if (offscreen && !this.pickingColor) {
		return;
	}
			
	// load mvMatrix with either identity or lookAt
	//mat4.identity(this.mvMatrix);
	mat4.lookAt(neg.cam_eye, neg.cam_center, neg.cam_up, this.mvMatrix);
    mat4.translate(this.mvMatrix, this.center);
	mat4.rotate(this.mvMatrix, degToRad(this.rotationAngle), [1, 0, 0]);
	mat4.rotate(this.mvMatrix, degToRad(this.rotationAngle), [0, 1, 0]);
	mat4.rotate(this.mvMatrix, degToRad(this.rotationAngle), [0, 0, 1]);
	mat4.scale(this.mvMatrix, [this.scale, this.scale, this.scale]);

	neg.gl.uniformMatrix4fv(neg.mvMatrixUniformLoc, false, this.mvMatrix);

	// Send picking color to fragment shader
	if (offscreen) {
		if (this.pickingColor) {
			neg.gl.uniform4fv(neg.uniformPickingColorLoc, this.pickingColor);
		}
	}
			
	for (var i in this.model.faceList) {
		var face = this.model.faceList[i];
		
		if (this.model.verticesTextureCoord) {
			if (!offscreen) {
				var texture = neg.textureTable[face.textureName];
				neg.gl.enableVertexAttribArray(neg.vertexTextureAttributeLoc);
				neg.gl.bindTexture(neg.gl.TEXTURE_2D, texture);
			}
		}
		else {
			neg.gl.disableVertexAttribArray(neg.vertexTextureAttributeLoc);
		}
		
		neg.gl.drawElements(this.model.drawPrimitive, face.vertexIndexNumber, neg.gl.UNSIGNED_SHORT, face.vertexIndexOffset * this.model.vertexIndexBufferItemSize);
	}
}

function NegModel(verticesCoord, verticesTextureCoord, verticesColors, verticesIndices, primitive, texWeight) {
	this.verticesCoord = verticesCoord;
	this.verticesTextureCoord = verticesTextureCoord;
	this.verticesColors = verticesColors;
	this.verticesIndices = verticesIndices;
	this.faceList = [];
	this.instanceList = [];
	if ((primitive != undefined) && (primitive != null)) {
		this.drawPrimitive = primitive;
	}
	else {
		this.drawPrimitive = neg.gl.TRIANGLES;
	}
	if ((texWeight != undefined) && (texWeight != null)) {
		this.textureWeight = texWeight;
	}
	else {
		this.textureWeight = 0.5; // 50%
	}
}

NegModel.prototype.addFace = function(vertexIndexOffset, vertexIndexNumber, textureName) {
	this.faceList.push(new NegFace(vertexIndexOffset, vertexIndexNumber, textureName));
}

NegModel.prototype.addInstance = function(center, scale, animation, pick) {
	var mi = new NegModelInstance(this, center, scale, animation);
	if (pick) {
		mi.pickingColor = generatePickingColor();
	}
	this.instanceList.push(mi);
	return mi;
}

NegModel.prototype.drawInstances = function(offscreen) {
		
    neg.gl.bindBuffer(neg.gl.ARRAY_BUFFER, this.vertexPositionBuffer);
   	neg.gl.vertexAttribPointer(neg.vertexPositionAttributeLoc, this.vertexPositionBufferItemSize, neg.gl.FLOAT, false, 0, 0);

	if (this.verticesTextureCoord) {
		neg.gl.bindBuffer(neg.gl.ARRAY_BUFFER, this.vertexTextureCoordBuffer);
		neg.gl.vertexAttribPointer(neg.vertexTextureAttributeLoc, this.vertexTextureCoordBufferItemSize, neg.gl.FLOAT, false, 0, 0);
	}

	if (this.verticesColors.constantColor) {
		neg.gl.disableVertexAttribArray(neg.vertexColorAttributeLoc);
		neg.gl.vertexAttrib4f(neg.vertexColorAttributeLoc,
			this.verticesColors.constantColor[0],
			this.verticesColors.constantColor[1],
			this.verticesColors.constantColor[2],
			this.verticesColors.constantColor[3]);
	}
	else {
		neg.gl.enableVertexAttribArray(neg.vertexColorAttributeLoc);	
		neg.gl.bindBuffer(neg.gl.ARRAY_BUFFER, this.vertexColorBuffer);
		neg.gl.vertexAttribPointer(neg.vertexColorAttributeLoc, this.vertexColorBufferItemSize, neg.gl.FLOAT, false, 0, 0);
	}
	
	neg.gl.bindBuffer(neg.gl.ELEMENT_ARRAY_BUFFER, this.vertexIndexBuffer);
	
	neg.gl.uniform1f(neg.texWeightUniformLoc, this.textureWeight);
			
	for (var i in this.instanceList) {
		this.instanceList[i].draw(offscreen);
	}
}

NegModel.prototype.initBuffers = function() {
	
	this.vertexPositionBuffer = neg.gl.createBuffer();
	neg.gl.bindBuffer(neg.gl.ARRAY_BUFFER, this.vertexPositionBuffer);
	neg.gl.bufferData(neg.gl.ARRAY_BUFFER, new Float32Array(this.verticesCoord), neg.gl.STATIC_DRAW);
    this.vertexPositionBufferItemSize = 3; // coord x,y,z
	
	if (this.verticesTextureCoord) {
		this.vertexTextureCoordBuffer = neg.gl.createBuffer();
		neg.gl.bindBuffer(neg.gl.ARRAY_BUFFER, this.vertexTextureCoordBuffer);
		neg.gl.bufferData(neg.gl.ARRAY_BUFFER, new Float32Array(this.verticesTextureCoord), neg.gl.STATIC_DRAW);
		this.vertexTextureCoordBufferItemSize = 2; // coord s,t
	}

	this.vertexColorBuffer = neg.gl.createBuffer();
	neg.gl.bindBuffer(neg.gl.ARRAY_BUFFER, this.vertexColorBuffer);
	neg.gl.bufferData(neg.gl.ARRAY_BUFFER, new Float32Array(this.verticesColors), neg.gl.STATIC_DRAW);
    this.vertexColorBufferItemSize = 4; // color rgba

	this.vertexIndexBuffer = neg.gl.createBuffer();
	neg.gl.bindBuffer(neg.gl.ELEMENT_ARRAY_BUFFER, this.vertexIndexBuffer);
	neg.gl.bufferData(neg.gl.ELEMENT_ARRAY_BUFFER, new Uint16Array(this.verticesIndices), neg.gl.STATIC_DRAW);
	this.vertexIndexBufferItemSize = 2;
	
	neg.gl.bindBuffer(neg.gl.ARRAY_BUFFER, null);
	neg.gl.bindBuffer(neg.gl.ELEMENT_ARRAY_BUFFER, null);
}

NegModel.prototype.animateInstances = function(elapsedTime) {
	for (var i in this.instanceList) {
		if (this.instanceList[i].animation) {
			this.instanceList[i].animation(elapsedTime);
		}
	}
}

NegModel.prototype.initTextures = function() {
	for (var i in this.faceList) {
		var face = this.faceList[i];
		loadTexture(face.textureName);
	}
	
	var unit = 0;
	neg.gl.uniform1i(neg.uniformSamplerLoc, unit); // (default) point sampler to texture unit 0
	neg.gl.activeTexture(neg.gl.TEXTURE0 + unit); // (default) activate texture unit 0
}
	
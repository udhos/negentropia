
function MeshFace(vertexIndexOffset, vertexIndexNumber, textureName) {
	this.vertexIndexOffset = vertexIndexOffset;
	this.vertexIndexNumber = vertexIndexNumber;
	this.textureName = textureName;
}

function MeshInstance(model, center, scale, animation) {
	this.model = model;
	this.center = center;
	this.scale = scale;
	this.animation = animation;
	
	//this.orientationQuaternion = quat4.fromAxes([0, 0, -1], [1, 0, 0], [0, 1, 0], quat4.create());
	this.orientationQuaternion = quat4.identity();
	
	this.mvMatrix  = mat4.create();	
	
	this.name = 'Model ' + generateNameId();
}

MeshInstance.prototype.orbitZ = function(elapsedTime) {

	if (!('orbitQuaternion' in this)) {
		return;
	}
	
	var change = quat4.fromAngleAxis(degToRad(90 * elapsedTime / 1000.0), [0, 0, 1]);
	var orig = quat4.create(this.orbitQuaternion);
	
	quat4.multiply(orig, change, this.orbitQuaternion); // result = orig * change
	
	// FIXME: How often to normalize ??
	quat4.normalize(this.orbitQuaternion);
}

MeshInstance.prototype.rotateX = function(elapsedTime) {
}

MeshInstance.prototype.rotateY = function(elapsedTime) {
}

MeshInstance.prototype.rotateZ = function(elapsedTime) {

	var change = quat4.fromAngleAxis(degToRad(30 * elapsedTime / 1000.0), [0, 0, 1]);
	var orig = quat4.create(this.orientationQuaternion);
	
	quat4.multiply(orig, change, this.orientationQuaternion); // result = orig * change
	
	// FIXME: How often to normalize ??
	quat4.normalize(this.orientationQuaternion);
}

MeshInstance.prototype.coord = function() {
	return this.center.slice(0);
}

MeshInstance.prototype.draw = function(offscreen) {

	// Won´t draw object to offscreen framebuffer without picking color
	if (offscreen && !('pickingColor' in this)) {
		return;
	}
	
	// grand world coordinate system:
	// 1. obj scale
	// 2. obj rotate
	// 3. obj orbit translate
	// 4. obj orbit rotate
	// 5. obj translate
	// 6. camera rotate
	// 7. camera translate

	// camera
	
	if ('cam' in neg) {
		// quaternion camera
		
		mat4.identity(this.mvMatrix);

		// 7. camera translate
		neg.cam.translate(this.mvMatrix);

		// 6. camera rotate
		neg.cam.rotate(this.mvMatrix);	
	}
	else {
		// matrix camera
		
		// load mvMatrix with either identity or lookAt
		//mat4.identity(this.mvMatrix);
		mat4.lookAt(neg.cam_eye, neg.cam_center, neg.cam_up, this.mvMatrix);
	}	

	// 5. obj translate
    mat4.translate(this.mvMatrix, this.coord());
	
	// 4. obj orbit rotate

	if ('orbitQuaternion' in this) {
		var quatMat = mat4.create();
		quat4.toMat4(this.orbitQuaternion, quatMat);
		mat4.multiply(this.mvMatrix, quatMat);
	}
	
	// 3. obj orbit translate
	if ('orbitTranslation' in this) {
	    mat4.translate(this.mvMatrix, this.orbitTranslation);
	}

	// 2. obj rotate

    var quatMat = mat4.create();
    quat4.toMat4(this.orientationQuaternion, quatMat);
    mat4.multiply(this.mvMatrix, quatMat);
	 
	// 1. obj scale
	mat4.scale(this.mvMatrix, [this.scale, this.scale, this.scale]);
	
	neg.gl.uniformMatrix4fv(neg.mvMatrixUniformLoc, false, this.mvMatrix);

	// Send picking color to fragment shader
	if (offscreen) {
		if ('pickingColor' in this) {
			neg.gl.uniform4fv(neg.uniformPickingColorLoc, this.pickingColor);
		}
	}

	
	for (var i in this.model.faceList) {
		var face = this.model.faceList[i];
		

		if ('verticesTextureCoord' in this.model) {
			if (!offscreen) {
				var texture = neg.textureTable[face.textureName];
				var unit = 0;
				neg.gl.activeTexture(neg.gl.TEXTURE0 + unit);
				neg.gl.bindTexture(neg.gl.TEXTURE_2D, texture);
				neg.gl.uniform1i(neg.uniformSamplerLoc, unit);
				
				neg.gl.enableVertexAttribArray(neg.vertexTextureAttributeLoc);
			}
		}
		else {
			neg.gl.disableVertexAttribArray(neg.vertexTextureAttributeLoc);
		}
		
		neg.gl.drawElements(this.model.drawPrimitive, face.vertexIndexNumber, neg.gl.UNSIGNED_SHORT, face.vertexIndexOffset * this.model.vertexIndexBufferItemSize);
	}
}

function Mesh(verticesCoord, verticesTextureCoord, verticesColors, verticesIndices, primitive, texWeight) {
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
	
	//console.info("debug mesh = " + JSON.stringify(this)); // debug dump mesh to console
}

Mesh.prototype.addFace = function(vertexIndexOffset, vertexIndexNumber, textureName) {
	this.faceList.push(new MeshFace(vertexIndexOffset, vertexIndexNumber, textureName));
}

Mesh.prototype.addInstance = function(center, scale, animation, pick) {
	var mi = new MeshInstance(this, center, scale, animation);
	if (pick) {
		mi.pickingColor = generatePickingColor();
	}
	this.instanceList.push(mi);
	return mi;
}

Mesh.prototype.drawInstances = function(offscreen) {
		
    neg.gl.bindBuffer(neg.gl.ARRAY_BUFFER, this.vertexPositionBuffer);
   	neg.gl.vertexAttribPointer(neg.vertexPositionAttributeLoc, this.vertexPositionBufferItemSize, neg.gl.FLOAT, false, 0, 0);

	if (this.verticesTextureCoord) {
		neg.gl.bindBuffer(neg.gl.ARRAY_BUFFER, this.vertexTextureCoordBuffer);
		neg.gl.vertexAttribPointer(neg.vertexTextureAttributeLoc, this.vertexTextureCoordBufferItemSize, neg.gl.FLOAT, false, 0, 0);
	}

	if ('constantColor' in this.verticesColors) {
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

Mesh.prototype.initBuffers = function() {
	
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

Mesh.prototype.animateInstances = function(elapsedTime) {
	for (var i in this.instanceList) {
		if (this.instanceList[i].animation) {
			this.instanceList[i].animation(elapsedTime);
		}
	}
}

Mesh.prototype.initTextures = function() {
	for (var i in this.faceList) {
		var face = this.faceList[i];
		loadTexture(face.textureName);
	}
	
	/*
	var unit = 0;
	neg.gl.uniform1i(neg.uniformSamplerLoc, unit); // (default) point sampler to texture unit 0
	neg.gl.activeTexture(neg.gl.TEXTURE0 + unit); // (default) activate texture unit 0
	++unit;
	//neg.gl.uniform1i(neg.uniformCubeSamplerLoc, unit); // point sampler to texture unit 1
	//neg.gl.activeTexture(neg.gl.TEXTURE0 + unit); // activate texture unit 1
	*/
}

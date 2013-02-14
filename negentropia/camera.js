
function Camera() {
	//this.orientationQuaternion = quat4.fromAxes([0, 0, -1], [1, 0, 0], [0, 1, 0], quat4.create());
	this.orientationQuaternion = quat4.identity();
	this.eyePosition = vec3.create(neg.cameraHomeCoord.slice(0));
	this.tmpMat = mat4.create();
	this.tmpQuat = quat4.create();
	this.tmpEye = vec3.create();
	this.tmpMat3 = mat3.create();
}

// Apply orientationQuaternion rotation to matrix mat
Camera.prototype.rotate = function(mat) {

	//quat4.toMat4(this.orientationQuaternion, this.tmpMat);
	quat4.conjugate(this.orientationQuaternion, this.tmpQuat);
	
	//quat4.normalize(this.tmpQuat); // ???
	
	quat4.toMat4(this.tmpQuat, this.tmpMat);
	mat4.multiply(mat, this.tmpMat);
}

// Apply eyePosition translation to matrix mat
Camera.prototype.translate = function(mat) {
	mat4.translate(mat, vec3.negate(this.eyePosition, this.tmpEye));
}

// Rotate orientationQuaternion around Y axis
Camera.prototype.rotateY = function(angle) {
	//quat4.fromAngleAxis(angle, [0,1,0], this.tmpQuat);
	quat4.fromAngleAxis(angle, this.getUp(), this.tmpQuat);
	quat4.multiply(this.orientationQuaternion, this.tmpQuat); 
	quat4.multiply(this.tmpQuat, this.orientationQuaternion, this.orientationQuaternion); 
}

Camera.prototype.getUp = function() {

	// get matrix 3 from quaternion
	quat4.toMat3(this.orientationQuaternion, this.tmpMat3);
	
	// get vector from column 1 (up vector of matrix)
	return vec3.createFrom(this.tmpMat3[1], this.tmpMat3[4], this.tmpMat3[7]);
}


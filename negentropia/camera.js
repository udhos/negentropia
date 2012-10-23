
function Camera() {
	//this.orientationQuaternion = quat4.fromAxes([0, 0, -1], [1, 0, 0], [0, 1, 0], quat4.create());
	this.orientationQuaternion = quat4.identity();
	this.eyePosition = vec3.create(neg.cameraHomeCoord.slice(0));
	this.tmpMat = mat4.create();
	this.tmpQuat = quat4.create();
	this.tmpEye = vec3.create();
}

Camera.prototype.rotate = function(mat) {

	//quat4.toMat4(this.orientationQuaternion, this.tmpMat);
	quat4.conjugate(this.orientationQuaternion, this.tmpQuat);
	
	//quat4.normalize(this.tmpQuat); // ???
	
	quat4.toMat4(this.tmpQuat, this.tmpMat);
	mat4.multiply(mat, this.tmpMat);
}

Camera.prototype.translate = function(mat) {
	mat4.translate(mat, vec3.negate(this.eyePosition, this.tmpEye));
}

// quatCam.js
//
// WebGL Camera based on Quaternions
//
// http://gamedev.stackexchange.com/questions/46123/how-can-i-create-a-webgl-camera-based-on-quaternions

function Camera()
{
   this.position = vec3.create();
   this.rotation = quat4.create();

   quat4.fromAxes([0, 0, 1], [1, 0, 0], [0, 1, 0], this.rotation);
   //mat4.lookAt([0,0,0], [0,0,1], [0,1,0], this.mat);

   //quat4.normalize(this.rotation);

   this.getMVMatrix = function(mat)
   {
      // this next line doesn't work, it rotates around the origin. useful for an OrbitCamera I guess.
      // mat4.fromRotationTranslation(this.rotation, this.position, mat);

      // doesn't seem to be needed: mat4.identity( mat);
      quat4.toMat4( this.rotation, mat);
      mat4.translate( mat, this.position);
   }

   var mat = mat3.create();
   // get the values from column 2. They represent the direction of this matrx
   this.getDir = function()
   {
      //console.log(quat4.str(this));
      quat4.toMat3( this.rotation, mat);
      //console.log(mat3.str(mat));
      var a02 = mat[2],
          a12 = mat[5],
          a22 = mat[8];
      return vec3.createFrom( a02, a12, a22); // todo: reuse temp var
   }

   // get the values from column 1. They represent the up vector of this matrx
   this.getUp = function()
   {
      quat4.toMat3( this.rotation, mat);
      var a01 = mat[1],
          a11 = mat[4],
          a21 = mat[7];
      return vec3.createFrom( a01, a11, a21);
   }

   // get the values from column 0. They represent the left vector of this matrx
   this.getLeft = function()
   {
      quat4.toMat3( this.rotation, mat);
      var a00 = mat[0],
          a10 = mat[3],
          a20 = mat[6];
      return vec3.createFrom( a00, a10, a20);
   }

   this.moveForward = function(amount)
   {
      var dir = this.getDir();
      vec3.normalize(dir);
      vec3.scale(dir, amount);
      vec3.add( this.position, dir);
      //console.log(vec3.str(this.position));
      //console.log( vec3.str(quat4.multiplyVec3( this.rotation, vec3.create(0.0, 0.0, -1))) )
      //vec3.add( this.position, quat4.multiplyVec3( this.rotation, vec3.create(0.0, 0.0, -1)));
   }

   this.moveBackward = function(amount)
   {
      var dir = this.getDir();
      vec3.normalize(dir);
      vec3.scale(dir, amount);
      vec3.subtract( this.position, dir);
      //vec3.add( this.position, quat4.multiplyVec3( this.rotation, vec3.create(0.0, 0.0, amount)));
   }

   this.moveLeft = function(amount)
   {
      var dir = this.getLeft();
      vec3.normalize(dir);
      vec3.scale(dir, amount);
      vec3.add( this.position, dir);
      //vec3.add( this.position, quat4.multiplyVec3( this.rotation, vec3.create(-amount, 0.0, 0.0)));
   }

   this.moveRight = function(amount)
   {
      var dir = this.getLeft();
      vec3.normalize(dir);
      vec3.scale(dir, amount);
      vec3.subtract( this.position, dir);
      //vec3.add( this.position, quat4.multiplyVec3( this.rotation, vec3.create(amount, 0.0, 0.0)));
   }

   var tempQuat = quat4.create();
   this.lookUp = function(amount)
   {
      quat4.fromAngleAxis( -amount, this.getLeft(), tempQuat);
      quat4.multiply( this.rotation, tempQuat);
   }
   this.lookDown = function(amount)
   {
      quat4.fromAngleAxis( amount, this.getLeft(), tempQuat);
      quat4.multiply( this.rotation, tempQuat);
   }
   this.lookLeft = function(amount)
   {
      quat4.fromAngleAxis( -amount, this.getUp(), tempQuat);
      quat4.multiply( this.rotation, tempQuat);
   }
   this.lookRight = function(amount)
   {
      quat4.fromAngleAxis( amount, this.getUp(), tempQuat);
      quat4.multiply( this.rotation, tempQuat);
   }
}

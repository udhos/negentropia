// makeSphere
//
// g.sphere = makeSphere(gl, 1, 30, 30);
// https://www.khronos.org/registry/webgl/sdk/demos/webkit/resources/J3DI.js
//
// demo:
// https://www.khronos.org/registry/webgl/sdk/demos/webkit/Earth.html

function sphereGenerate(radius, bands, vertexPositionData, textureCoordData, vertexIndexData) {
		
	var latitudeBands = bands;
	var longitudeBands = bands * 2;
	
	//var theta = Math.PI / 2;
	for (var lat = 0; lat <= latitudeBands; ++lat) {
		var theta = lat * Math.PI / latitudeBands;
		var sinTheta = Math.sin(theta);
		var cosTheta = Math.cos(theta);
				
		//var phi = 0;
		for (var lon = 0; lon <= longitudeBands; ++lon) {
			var phi = lon * 2 * Math.PI / longitudeBands;
			var sinPhi = Math.sin(phi);
			var cosPhi = Math.cos(phi);
			
			var x = cosPhi * sinTheta;
            var z = cosTheta;
            var y = sinPhi * sinTheta;
			var u = lon / longitudeBands;
			var v = 1 - lat / latitudeBands;
			
			vertexPositionData.push(radius * x);
			vertexPositionData.push(radius * y);
			vertexPositionData.push(radius * z);			
			textureCoordData.push(u);
			textureCoordData.push(v);
		}
	}
	
	for (var lat = 0; lat < latitudeBands; ++lat) {
		for (var lon = 0; lon < longitudeBands; ++lon) {
			var first = (lat * (longitudeBands + 1)) + lon;
			var second = first + longitudeBands + 1;
			
			vertexIndexData.push(first);
			vertexIndexData.push(second);
			vertexIndexData.push(first + 1);

			vertexIndexData.push(second);
			vertexIndexData.push(second + 1);
			vertexIndexData.push(first + 1);
		}
    }
	
}


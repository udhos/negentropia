library buffer;

import 'dart:html';
import 'dart:async';
import 'dart:json';

class Model {
  WebGLBuffer vertexPositionBuffer;
  WebGLBuffer vertexIndexBuffer;
  int vertexPositionBufferItemSize;
  int vertexIndexBufferItemSize;
  int vertexIndexLength;
  
  Model(WebGLRenderingContext gl, List<num> vertCoord, List<int> vertInd) {
    this.vertexPositionBuffer = gl.createBuffer();
    gl.bindBuffer(WebGLRenderingContext.ARRAY_BUFFER, this.vertexPositionBuffer);
    gl.bufferData(WebGLRenderingContext.ARRAY_BUFFER, new Float32Array.fromList(vertCoord), WebGLRenderingContext.STATIC_DRAW);
    this.vertexPositionBufferItemSize = 3; // coord x,y,z
    
    this.vertexIndexBuffer = gl.createBuffer();
    gl.bindBuffer(WebGLRenderingContext.ELEMENT_ARRAY_BUFFER, this.vertexIndexBuffer);
    gl.bufferData(WebGLRenderingContext.ELEMENT_ARRAY_BUFFER, new Uint16Array.fromList(vertInd), WebGLRenderingContext.STATIC_DRAW);
    this.vertexIndexBufferItemSize = 2; // size of Uint16Array
    
    this.vertexIndexLength = vertInd.length;
    
    print("Model: vertex index length: ${this.vertexIndexLength}");
    
    // clean-up
    gl.bindBuffer(WebGLRenderingContext.ARRAY_BUFFER, null);
    gl.bindBuffer(WebGLRenderingContext.ELEMENT_ARRAY_BUFFER, null);
  }
}

void fetchSquare(WebGLRenderingContext gl, String jsonUrl, void deliverSquare(Model)) {
  
  void handleResponse(String response) {
    print("fetched square JSON from URL: $jsonUrl: [$response]");
    Map square;
    try {
      square = parse(response);
    }
    catch (e) {
      print("failure parsing square JSON: $e");
      return;
    }
    print("square JSON parsed: [$square]");
    
    List<num> vertCoord = square['vertCoord'];
    List<int> vertInd = square['vertInd'];
    
    Model squareModel = new Model(gl, vertCoord, vertInd);
    deliverSquare(squareModel); // callback
  }
  
  void handleError(AsyncError err) {
    print("failure fetching square JSON from URL: $jsonUrl: $err");
  }
  
  // dart magic :-)
  HttpRequest.getString(jsonUrl)
    .then(handleResponse)
    .catchError(handleError);
}




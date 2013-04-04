library buffer;

import 'dart:html';
import 'dart:async';
import 'dart:json';

import 'shader.dart';

class Instance {
  
  Model model;
  
  Instance(Model this.model);
  
  void draw() {
    
    WebGLRenderingContext gl = model.program.gl;
    
    gl.bindBuffer(WebGLRenderingContext.ARRAY_BUFFER, model.vertexPositionBuffer);
    gl.vertexAttribPointer(model.program.aVertexPosition, model.vertexPositionBufferItemSize, WebGLRenderingContext.FLOAT, false, 0, 0);
  
    gl.bindBuffer(WebGLRenderingContext.ELEMENT_ARRAY_BUFFER, model.vertexIndexBuffer);
    gl.drawElements(WebGLRenderingContext.TRIANGLES, model.vertexIndexLength, WebGLRenderingContext.UNSIGNED_SHORT, 0 * model.vertexIndexBufferItemSize);
  }
}

class Model {
    
  WebGLBuffer vertexPositionBuffer;
  WebGLBuffer vertexIndexBuffer;
  int vertexPositionBufferItemSize;
  int vertexIndexBufferItemSize;
  int vertexIndexLength;

  List<Instance> instanceList = new List<Instance>();
  Program program;
  
  void _createBuffers(WebGLRenderingContext gl, List<num> vertCoord, List<int> vertInd) {
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
  
  Model.fromLists(WebGLRenderingContext gl, Program prog, List<num> vertCoord, List<int> vertInd) {
    this.program = prog;
    _createBuffers(gl, vertCoord, vertInd);
  }
  
  Model.fromURL(WebGLRenderingContext gl, Program prog, String URL) {
    this.program = prog;

    /*
    // load JSON from URL
    var req = new HttpRequest();
    req.open("GET", URL);
    req.onLoad.listen((ProgressEvent e) {
      String response = req.responseText;
      if (req.status != 200) {
        print("Model.fromURL: error: [$response]");
        return;
      }
      print("Model.fromURL: loaded: json=[$response]");
      Map m = parse(response);
      List<num> vertCoord = m["vertCoord"];
      List<int> vertInd = m["vertInd"];
      _createBuffers(gl, vertCoord, vertInd);
    });
    req.onError.listen((e) {
      print("Model.fromURL: error: [$e]");
    });
    req.send();
    */

    void handleResponse(String response) {
      print("Model.fromURL: fetched JSON from URL: $URL: [$response]");
      Map m;
      try {
        m = parse(response);
      }
      catch (e) {
        print("Model.fromURL: failure parsing square JSON: $e");
        return;
      }
      print("Model.fromURL: JSON parsed: [$m]");
      
      List<num> vertCoord = m['vertCoord'];
      List<int> vertInd = m['vertInd'];

      _createBuffers(gl, vertCoord, vertInd);
    }
    
    void handleError(AsyncError err) {
      print("Model.fromURL: failure fetching square JSON from URL: $URL: $err");
    }

    HttpRequest.getString(URL)
    .then(handleResponse)
    .catchError(handleError);
  }
  
  void addInstance(Instance i) {
    this.instanceList.add(i);
  }
 
  void drawInstances() {
    this.instanceList.forEach((Instance i) => i.draw());
  }  
  
}

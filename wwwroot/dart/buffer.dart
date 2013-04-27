library buffer;

import 'dart:html';
import 'dart:async';
import 'dart:json';
import 'dart:web_gl';
import 'dart:typed_data';

import 'shader.dart';

class Instance {
  
  Model model;
  List<num> center;
  num scale;
  
  Instance(Model this.model, List<num> this.center, num this.scale);
  
  void draw() {
    
    RenderingContext gl = model.program.gl;
    
    gl.bindBuffer(RenderingContext.ARRAY_BUFFER, model.vertexPositionBuffer);
    gl.vertexAttribPointer(model.program.a_Position, model.vertexPositionBufferItemSize, RenderingContext.FLOAT, false, 0, 0);
  
    gl.bindBuffer(RenderingContext.ELEMENT_ARRAY_BUFFER, model.vertexIndexBuffer);
    gl.drawElements(RenderingContext.TRIANGLES, model.vertexIndexLength, RenderingContext.UNSIGNED_SHORT, 0 * model.vertexIndexBufferItemSize);
  }
}

class Model {
    
  Buffer vertexPositionBuffer;
  Buffer vertexIndexBuffer;
  int vertexPositionBufferItemSize;
  int vertexIndexBufferItemSize;
  int vertexIndexLength;

  List<Instance> instanceList = new List<Instance>();
  ShaderProgram program; // parent program
  
  void _createBuffers(RenderingContext gl, List<num> vertCoord, List<int> vertInd) {
    this.vertexPositionBuffer = gl.createBuffer();
    gl.bindBuffer(RenderingContext.ARRAY_BUFFER, this.vertexPositionBuffer);
    gl.bufferData(RenderingContext.ARRAY_BUFFER, new Float32List.fromList(vertCoord), RenderingContext.STATIC_DRAW);
    this.vertexPositionBufferItemSize = 3; // coord x,y,z
    
    this.vertexIndexBuffer = gl.createBuffer();
    gl.bindBuffer(RenderingContext.ELEMENT_ARRAY_BUFFER, this.vertexIndexBuffer);
    gl.bufferData(RenderingContext.ELEMENT_ARRAY_BUFFER, new Uint16List.fromList(vertInd), RenderingContext.STATIC_DRAW);
    this.vertexIndexBufferItemSize = 2; // size of Uint16Array
    
    this.vertexIndexLength = vertInd.length;
    
    print("Model: vertex index length: ${this.vertexIndexLength}");
    
    // clean-up
    gl.bindBuffer(RenderingContext.ARRAY_BUFFER, null);
    gl.bindBuffer(RenderingContext.ELEMENT_ARRAY_BUFFER, null);
  }
  
  Model.fromLists(RenderingContext gl, ShaderProgram this.program, List<num> vertCoord, List<int> vertInd) {
    _createBuffers(gl, vertCoord, vertInd);
  }
  
  Model.fromURL(RenderingContext gl, ShaderProgram this.program, String URL) {

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
    
    void handleError(Object err) {
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

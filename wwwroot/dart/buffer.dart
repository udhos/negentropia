library buffer;

import 'dart:html';
import 'dart:async';
import 'dart:json';
import 'dart:web_gl';
import 'dart:typed_data';
import 'dart:math' as math;

import 'package:vector_math/vector_math.dart';
import 'package:game_loop/game_loop_html.dart';

import 'shader.dart';
import 'camera.dart';
import 'obj.dart';

class Instance {
  
  Model model;
  Vector3 center;
  double scale;
  Matrix4 MV = new Matrix4.identity(); // model-view matrix
  
  Instance(Model this.model, Vector3 this.center, double this.scale);
  
  void update(GameLoopHtml gameLoop) {
  }
  
  void draw(GameLoopHtml gameLoop, Camera cam) {

    if (model.vertexPositionBufferItemSize == null) {
      // model isn't ready
      return;  
    }
    
    setViewMatrix(MV, cam.eye, cam.center, cam.up);
    
    MV.translate(center[0], center[1], center[2]);
    
    MV.scale(scale, scale, scale);
    
    ShaderProgram prog = model.program;
    RenderingContext gl = prog.gl;

    gl.uniformMatrix4fv(prog.u_MV, false, MV.storage);

    gl.bindBuffer(RenderingContext.ARRAY_BUFFER, model.vertexPositionBuffer);
    gl.vertexAttribPointer(prog.a_Position, model.vertexPositionBufferItemSize, RenderingContext.FLOAT, false, 0, 0);
  
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
  
  void _createBuffers(RenderingContext gl, List<int> indices, List<double> vertCoord, List<double> textCoord, List<double> normCoord) {
        
    vertexPositionBuffer = gl.createBuffer();
    gl.bindBuffer(RenderingContext.ARRAY_BUFFER, vertexPositionBuffer);
    gl.bufferData(RenderingContext.ARRAY_BUFFER, new Float32List.fromList(vertCoord), RenderingContext.STATIC_DRAW);
    vertexPositionBufferItemSize = 3; // coord x,y,z
    
    vertexIndexBuffer = gl.createBuffer();
    gl.bindBuffer(RenderingContext.ELEMENT_ARRAY_BUFFER, vertexIndexBuffer);
    gl.bufferData(RenderingContext.ELEMENT_ARRAY_BUFFER, new Uint16List.fromList(indices), RenderingContext.STATIC_DRAW);
    vertexIndexBufferItemSize = 2; // size of Uint16Array
    
    vertexIndexLength = indices.length;
    
    print("Model._createBuffers: vertex index length: $vertexIndexLength");
    
    // clean-up
    gl.bindBuffer(RenderingContext.ARRAY_BUFFER, null);
    gl.bindBuffer(RenderingContext.ELEMENT_ARRAY_BUFFER, null);
  }
  
  Model();
  
  Model.fromLists(RenderingContext gl, ShaderProgram this.program, List<double> vertCoord, List<int> vertInd) {
    _createBuffers(gl, vertInd, vertCoord, null, null);
  }
  
  Model.fromJson(RenderingContext gl, ShaderProgram this.program, String URL) {

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
      //print("Model.fromJson: fetched JSON from URL: $URL: [$response]");
      Map m;
      try {
        m = parse(response);
      }
      catch (e) {
        print("Model.fromJson: failure parsing JSON: $e");
        return;
      }
      //print("Model.fromJson: JSON parsed: [$m]");
      
      List<int> indices = m['vertInd'];
      List<double> vertCoord = m['vertCoord'];

      _createBuffers(gl, indices, vertCoord, null, null);
    }
    
    void handleError(Object err) {
      print("Model.fromJson: failure fetching JSON from URL: $URL: $err");
    }

    HttpRequest.getString(URL)
    .then(handleResponse)
    .catchError(handleError);
  }

  Model.fromOBJ(RenderingContext gl, ShaderProgram this.program, String URL) {

    void handleResponse(String response) {
      print("Model.fromOBJ: fetched OBJ from URL: $URL");
      
      Obj obj = new Obj.fromString(URL, response);
      
      _createBuffers(gl, obj.indices, obj.vertCoord, obj.textCoord, obj.normCoord);
    }

    void handleError(Object err) {
      print("Model.fromOBJ: failure fetching OBJ from URL: $URL: $err");
    }

    HttpRequest.getString(URL)
    .then(handleResponse)
    .catchError(handleError);    
  }

  void addInstance(Instance i) {
    this.instanceList.add(i);
  }
 
  void drawInstances(GameLoopHtml gameLoop, Camera cam) {
    this.instanceList.forEach((Instance i) => i.draw(gameLoop, cam));
  }  

  void update(GameLoopHtml gameLoop) {
    this.instanceList.forEach((Instance i) => i.update(gameLoop));
  }  

}

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

class Instance {
  
  Model model;
  vec3 center;
  double scale;
  //double _size;
  mat4 MV = new mat4.identity(); // model-view matrix
  
  Instance(Model this.model, vec3 this.center, double this.scale);
  
  void update(GameLoopHtml gameLoop) {
    /*
    double degreesPerSec = 60.0;
    double angle = (gameLoop.gameTime * degreesPerSec) % 360.0; 
    double rad = angle * math.PI / 180.0;
    _size = 10 * math.sin(rad).abs() + 1;
    */
  }
  
  void draw(GameLoopHtml gameLoop, Camera cam) {

    double size = 10 * math.sin(cam.rad).abs() + 1;

    setViewMatrix(MV, cam.eye, cam.center, cam.up);
    
    MV.translate(center[0], center[1], center[2]);
    
    double s = scale * size;
    MV.scale(s, s, s);
    
    ShaderProgram prog = model.program;
    RenderingContext gl = prog.gl;

    // send model-view matrix uniform
    /*
    List<num> MV_tmp = new List<num>(16); 
    MV.copyIntoArray(MV_tmp);
    gl.uniformMatrix4fv(prog.u_MV, false, MV_tmp);
    */
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
 
  void drawInstances(GameLoopHtml gameLoop, Camera cam) {
    this.instanceList.forEach((Instance i) => i.draw(gameLoop, cam));
  }  

  void update(GameLoopHtml gameLoop) {
    this.instanceList.forEach((Instance i) => i.update(gameLoop));
  }  

}

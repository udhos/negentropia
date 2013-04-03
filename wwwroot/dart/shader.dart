library shader;

import 'dart:html';

import 'buffer.dart';

class Program {
  
  WebGLProgram program;
  int aVertexPosition;
  bool ready = false;
  WebGLRenderingContext gl;
  
  List<Model> modelList = new List<Model>();  
  
  Program._load(WebGLRenderingContext gl, String vertexShaderURL, fragmentShaderURL) {
    print("Program._load: vsUrl=$vertexShaderURL fsURL=$fragmentShaderURL");
    
    String vertShaderSrc, fragShaderSrc;
    
    void tryCompileShaders() {
      if (vertShaderSrc == null || fragShaderSrc == null) {
        return;
      }
      
      print("shaders: available to compile");      
      
      WebGLShader vs = gl.createShader(WebGLRenderingContext.VERTEX_SHADER);
      gl.shaderSource(vs, vertShaderSrc);
      gl.compileShader(vs);
      if (!gl.getShaderParameter(vs, WebGLRenderingContext.COMPILE_STATUS) && !gl.isContextLost()) { 
        print(gl.getShaderInfoLog(vs));
        return;
      }
      
      WebGLShader fs = gl.createShader(WebGLRenderingContext.FRAGMENT_SHADER);
      gl.shaderSource(fs, fragShaderSrc);
      gl.compileShader(fs);
      if (!gl.getShaderParameter(fs, WebGLRenderingContext.COMPILE_STATUS) && !gl.isContextLost()) { 
        print(gl.getShaderInfoLog(fs));
        return;
      }

      WebGLProgram p = gl.createProgram();
      gl.attachShader(p, vs);
      gl.attachShader(p, fs);
      gl.linkProgram(p);
      if (!gl.getProgramParameter(p, WebGLRenderingContext.LINK_STATUS) && !gl.isContextLost()) { 
        print(gl.getProgramInfoLog(p));
      }
      
      this.aVertexPosition = gl.getAttribLocation(p, "aVertexPosition");
      this.ready = true;
      
      print("shader program: ready");
      
      this.program = p;
      this.gl = gl;
    }
        
    var requestVert = new HttpRequest();
    requestVert.open("GET", vertexShaderURL);
    requestVert.onLoad.listen((ProgressEvent e) {
      String response = requestVert.responseText;
      if (requestVert.status != 200) {
        print("vertexShader: error: [$response]");
        return;
      }
      print("vertexShader: loaded: [$response]");
      vertShaderSrc = response;
      tryCompileShaders();
    });
    requestVert.onError.listen((e) {
      print("vertexShader: error: [$e]");
    });
    requestVert.send();

    var requestFrag = new HttpRequest();
    requestFrag.open("GET", fragmentShaderURL);
    requestFrag.onLoad.listen((ProgressEvent e) {
      String response = requestFrag.responseText;
      if (requestFrag.status != 200) {
        print("fragmentShader: error: [$response]");
        return;
      }
      print("fragmentShader: loaded: [$response]");
      fragShaderSrc = response;
      tryCompileShaders();      
    });
    requestFrag.onError.listen((e) {
      print("fragmentShader: error: [$e]");
    });
    requestFrag.send();
  }
  
  factory Program(WebGLRenderingContext gl, String vertexShader, fragmentShader) {
    return new Program._load(gl, vertexShader, fragmentShader);
  }
  
  void addModel(Model m) {
    this.modelList.add(m);
  }
  
  void drawModels() {
    // FIXME WRITEME
    
    gl.useProgram(program);
    gl.enableVertexAttribArray(aVertexPosition);
    
    modelList.forEach((Model m) => m.drawInstances());

    // clean up
    gl.bindBuffer(WebGLRenderingContext.ARRAY_BUFFER, null);
    gl.bindBuffer(WebGLRenderingContext.ELEMENT_ARRAY_BUFFER, null);
    
    //gl.disableVertexAttribArray(aVertexPosition); // needed ??
  }
}


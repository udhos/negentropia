library shader;

import 'dart:html';

import 'buffer.dart';

class Program {
  
  WebGLProgram program;
  int aVertexPosition;
  WebGLRenderingContext gl;
  
  List<Model> modelList = new List<Model>();  
  
  Program._load(WebGLRenderingContext this.gl, Map<String,WebGLShader> shaderCache, String vertexShaderURL, fragmentShaderURL) {
    print("Program._load: vsUrl=$vertexShaderURL fsURL=$fragmentShaderURL");
    
    //String vertShaderSrc, fragShaderSrc;
    
    WebGLShader compileShader(String shaderURL, String shaderSource, int shaderType) {
      WebGLShader shader = gl.createShader(shaderType);
      gl.shaderSource(shader, shaderSource);
      gl.compileShader(shader);
      if (!gl.getShaderParameter(shader, WebGLRenderingContext.COMPILE_STATUS) && !gl.isContextLost()) { 
        print("compileShader: compilation FAILURE: " + shaderURL + ": " + gl.getShaderInfoLog(shader));
        return null;
      }
      
      shaderCache[shaderURL] = shader;
      print("compileShader: " + shaderURL + ": compiled and cached");
      
      return shader;
    }

    WebGLShader vertexShader, fragmentShader;

    void tryLink() {
      if (vertexShader == null || fragmentShader == null) {
        // not ready
        return;
      }
      
      WebGLProgram p = gl.createProgram();
      gl.attachShader(p, vertexShader);
      gl.attachShader(p, fragmentShader);
      gl.linkProgram(p);
      if (!gl.getProgramParameter(p, WebGLRenderingContext.LINK_STATUS) && !gl.isContextLost()) { 
        print(gl.getProgramInfoLog(p));
      }
      
      this.aVertexPosition = gl.getAttribLocation(p, "aVertexPosition");
      this.program = p;
      
      print("shader program: ready");      
    }
    
    void fetchVertexShader() {      
      var requestVert = new HttpRequest();
      requestVert.open("GET", vertexShaderURL);
      requestVert.onLoad.listen((ProgressEvent e) {
        String response = requestVert.responseText;
        if (requestVert.status != 200) {
          print("vertexShader: error: [$response]");
          return;
        }
        print("vertexShader: loaded: [$response]");
        vertexShader = compileShader(vertexShaderURL, response, WebGLRenderingContext.VERTEX_SHADER);
        tryLink();
      });
      requestVert.onError.listen((e) {
        print("vertexShader: error: [$e]");
      });
      requestVert.send();
    }

    void fetchFragmentShader() {
      var requestFrag = new HttpRequest();
      requestFrag.open("GET", fragmentShaderURL);
      requestFrag.onLoad.listen((ProgressEvent e) {
        String response = requestFrag.responseText;
        if (requestFrag.status != 200) {
          print("fragmentShader: error: [$response]");
          return;
        }
        print("fragmentShader: loaded: [$response]");
        fragmentShader = compileShader(fragmentShaderURL, response, WebGLRenderingContext.FRAGMENT_SHADER);
        tryLink();      
      });
      requestFrag.onError.listen((e) {
        print("fragmentShader: error: [$e]");
      });
      requestFrag.send();
    }
    
    vertexShader = shaderCache[vertexShaderURL];
    if (vertexShader == null) {
      print("vertexShader: " + vertexShaderURL + ": cache MISS");
      fetchVertexShader();
    }
    else {
      print("vertexShader: " + vertexShaderURL + ": cache HIT");      
    }
    if (fragmentShader == null) {
      print("fragmentShader: " + fragmentShaderURL + ": cache MISS");
      fetchFragmentShader();
    }
    else {
      print("fragmentShader: " + fragmentShaderURL + ": cache HIT");
    }
    tryLink(); // needed when both vertexShader and fragmentShader are found in cache
  }
 
  
  factory Program(WebGLRenderingContext gl, Map<String,WebGLShader> shaderCache, String vertexShaderURL, fragmentShaderURL) {
    return new Program._load(gl, shaderCache, vertexShaderURL, fragmentShaderURL);
  }
  
  void addModel(Model m) {
    this.modelList.add(m);
  }
  
  void drawModels() {
    
    gl.useProgram(program);
    gl.enableVertexAttribArray(aVertexPosition);
    
    modelList.forEach((Model m) => m.drawInstances());

    // clean up
    gl.bindBuffer(WebGLRenderingContext.ARRAY_BUFFER, null);
    gl.bindBuffer(WebGLRenderingContext.ELEMENT_ARRAY_BUFFER, null);
    
    //gl.disableVertexAttribArray(aVertexPosition); // needed ??
  }
}



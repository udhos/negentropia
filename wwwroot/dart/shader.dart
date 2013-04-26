library shader;

import 'dart:html';
import 'dart:web_gl';

import 'package:vector_math/vector_math.dart';

import 'buffer.dart';

class ShaderProgram {
  
  Program program;
  int a_Position;
  RenderingContext gl;
  
  List<Model> modelList = new List<Model>();  
 
  ShaderProgram(RenderingContext this.gl);

  void fetch(Map<String,Shader> shaderCache, String vertexShaderURL, fragmentShaderURL) {
    print("Program.fetch: vsUrl=$vertexShaderURL fsURL=$fragmentShaderURL");
    
    Shader compileShader(String shaderURL, String shaderSource, int shaderType) {
      Shader shader = gl.createShader(shaderType);
      gl.shaderSource(shader, shaderSource);
      gl.compileShader(shader);
      if (!gl.getShaderParameter(shader, RenderingContext.COMPILE_STATUS) && !gl.isContextLost()) { 
        print("compileShader: compilation FAILURE: " + shaderURL + ": " + gl.getShaderInfoLog(shader));
        return null;
      }
      
      shaderCache[shaderURL] = shader;
      print("compileShader: " + shaderURL + ": compiled and cached");
      
      return shader;
    }

    Shader vertexShader, fragmentShader;

    void tryLink() {
      if (vertexShader == null || fragmentShader == null) {
        // not ready
        return;
      }
      
      Program p = gl.createProgram();
      gl.attachShader(p, vertexShader);
      gl.attachShader(p, fragmentShader);
      gl.linkProgram(p);
      if (!gl.getProgramParameter(p, RenderingContext.LINK_STATUS) && !gl.isContextLost()) { 
        print(gl.getProgramInfoLog(p));
      }
      
      this.a_Position = gl.getAttribLocation(p, "a_Position");
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
        vertexShader = compileShader(vertexShaderURL, response, RenderingContext.VERTEX_SHADER);
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
        fragmentShader = compileShader(fragmentShaderURL, response, RenderingContext.FRAGMENT_SHADER);
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
   
  void addModel(Model m) {
    this.modelList.add(m);
  }
  
  void drawModels(mat4 pMatrix) {
    
    gl.useProgram(program);
    gl.enableVertexAttribArray(a_Position);

    modelList.forEach((Model m) => m.drawInstances());

    // clean up
    gl.bindBuffer(RenderingContext.ARRAY_BUFFER, null);
    gl.bindBuffer(RenderingContext.ELEMENT_ARRAY_BUFFER, null);
    
    //gl.disableVertexAttribArray(a_Position); // needed ??
  }
}



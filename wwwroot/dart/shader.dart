library shader;

import 'dart:html';

class Program {
  
  WebGLProgram program;
  int aVertexPosition;
  
  Program._load(WebGLRenderingContext gl, String vertexShaderURL, fragmentShaderURL) {
    print("Program._load: vsUrl=$vertexShaderURL fsURL=$fragmentShaderURL");
        
    var requestVert = new HttpRequest();
    requestVert.open("GET", vertexShaderURL);
    requestVert.onLoad.listen((ProgressEvent e) {
      String response = requestVert.responseText;
      if (requestVert.status != 200) {
        print("vertexShader: error: [$response]");
        return;
      }
      print("vertexShader: loaded: [$response]");
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
    });
    requestFrag.onError.listen((e) {
      print("fragmentShader: error: [$e]");
    });
    requestFrag.send();

    //aVertexPosition = gl.getAttribLocation(program, "aVertexPosition");    
  }
  
  factory Program(WebGLRenderingContext gl, String vertexShader, fragmentShader) {
    return new Program._load(gl, vertexShader, fragmentShader);
  }
}


library shader;

import 'dart:html';
import 'dart:async';
import 'dart:json';
import 'dart:math' as math;
import 'dart:web_gl';
import 'dart:typed_data';

import 'package:vector_math/vector_math.dart';
import 'package:game_loop/game_loop_html.dart';

import 'camera.dart';
import 'texture.dart';
import 'obj.dart';
import 'asset.dart';

part 'buffer.dart';
part 'shader_tex.dart';

class ShaderProgram {
  
  Program program;
  RenderingContext gl;
  int a_Position;
  UniformLocation u_MV;
  UniformLocation u_P;
  bool shaderReady = false;
  
  List<Model> modelList = new List<Model>();  
 
  ShaderProgram(RenderingContext this.gl);
  
  void initContext(RenderingContext gl, Map<String,Texture> textureTable) {
  }

  void getLocations() {
    a_Position = gl.getAttribLocation(program, "a_Position");
    u_MV       = gl.getUniformLocation(program, "u_MV");
    u_P        = gl.getUniformLocation(program, "u_P");      

    print("ShaderProgram: locations ready");      
  }
  
  void fetch(Map<String,Shader> shaderCache, String vertexShaderURL, fragmentShaderURL) {
    print("Program.fetch: vsUrl=$vertexShaderURL fsURL=$fragmentShaderURL");
    
    Shader compileShader(String shaderURL, String shaderSource, int shaderType) {
      Shader shader = gl.createShader(shaderType);
      gl.shaderSource(shader, shaderSource);
      gl.compileShader(shader);
      var parameter = gl.getShaderParameter(shader, RenderingContext.COMPILE_STATUS);
      if (parameter == null) {
        String infoLog = gl.getShaderInfoLog(shader);
        print("compileShader: compilation FAILURE: $shaderURL: $infoLog");
        if (gl.isContextLost()) {
          print("compileShader: compilation FAILURE: $shaderURL: $infoLog: context is lost");
        }
        return null;
      }
      
      shaderCache[shaderURL] = shader;
      print("compileShader: " + shaderURL + ": compiled and cached");
      
      return shader;
    }

    Shader vertexShader, fragmentShader;
    
    void tryLink() {
      if (vertexShader == null || fragmentShader == null) {
        // not ready to link
        return;
      }
      
      Program p = gl.createProgram();
      gl.attachShader(p, vertexShader);
      gl.attachShader(p, fragmentShader);
      gl.linkProgram(p);
      var parameter = gl.getProgramParameter(p, RenderingContext.LINK_STATUS);
      if (parameter == null) {
        String infoLog = gl.getProgramInfoLog(p);
        print("tryLink: shader program link FAILURE: $infoLog");
        if (gl.isContextLost()) {
          print("tryLink: shader program link FAILURE: $infoLog: context is lost");          
        }
        return;
      }

      print("ShaderProgram: program linked");      

      this.program = p;
      
      getLocations();
      
      shaderReady = true;
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
        //print("vertexShader: loaded: [$response]");
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
        //print("fragmentShader: loaded: [$response]");
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
  
  void drawModels(GameLoopHtml gameLoop, Camera cam, Matrix4 pMatrix) {
    
    if (!shaderReady) {
      return;
    }
    
    gl.useProgram(program);
    gl.enableVertexAttribArray(a_Position);

    // send perspective projection matrix uniform
    gl.uniformMatrix4fv(this.u_P, false, pMatrix.storage);

    modelList.forEach((Model m) => m.drawInstances(gameLoop, cam));

    // clean up
    gl.bindBuffer(RenderingContext.ARRAY_BUFFER, null);
    gl.bindBuffer(RenderingContext.ELEMENT_ARRAY_BUFFER, null);
    
    //gl.disableVertexAttribArray(a_Position); // needed ??
  }
  
  void update(GameLoopHtml gameLoop) {
    modelList.forEach((Model m) => m.update(gameLoop));    
  }
}



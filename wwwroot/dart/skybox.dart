library skybox;

import 'dart:html';
import 'dart:web_gl';
import 'dart:math' as math;

import 'package:vector_math/vector_math.dart';
import 'package:game_loop/game_loop_html.dart';

import 'shader.dart';
import 'buffer.dart';
import 'camera.dart';

class SkyboxProgram extends ShaderProgram {
  
  //UniformLocation u_P;
  UniformLocation u_Skybox;
  
  SkyboxProgram(RenderingContext gl) : super(gl) {
  }

  void fetch(Map<String,Shader> shaderCache, String vertexShaderURL, fragmentShaderURL) {
    super.fetch(shaderCache, vertexShaderURL, fragmentShaderURL);
    
    //this.u_P      = gl.getUniformLocation(this.program, "u_P");
    this.u_Skybox = gl.getUniformLocation(this.program, "u_Skybox");
  }

  void drawModels(GameLoopHtml gameLoop, Camera cam, mat4 pMatrix) {
    
    gl.useProgram(program);
    gl.enableVertexAttribArray(a_Position);
    
    int unit = 0;
    gl.activeTexture(RenderingContext.TEXTURE0 + unit);
    gl.uniform1i(this.u_Skybox, unit);

    // send perspective projection matrix uniform
    /*
    List<num> pTmp = new List<num>(16); 
    pMatrix.copyIntoArray(pTmp);   
    gl.uniformMatrix4fv(this.u_P, false, pTmp);
    */
    gl.uniformMatrix4fv(this.u_P, false, pMatrix.storage);

    modelList.forEach((Model m) => m.drawInstances(gameLoop, cam));

    // clean up
    gl.bindBuffer(RenderingContext.ARRAY_BUFFER, null);
    gl.bindBuffer(RenderingContext.ELEMENT_ARRAY_BUFFER, null);
    
    //gl.disableVertexAttribArray(a_Position); // needed ??
  }

}

class SkyboxModel extends Model {
  
  Texture cubemapTexture;
  
  SkyboxModel.fromJson(RenderingContext gl, SkyboxProgram prog, String URL, bool reverse, num rescale) : super.fromJson(gl, prog, URL) {
    cubemapTexture = gl.createTexture();
    
    /*
    gl.bindTexture(RenderingContext.TEXTURE_CUBE_MAP, cubemapTexture);
    gl.texParameteri(RenderingContext.TEXTURE_CUBE_MAP, RenderingContext.TEXTURE_MAG_FILTER, RenderingContext.LINEAR);
    gl.texParameteri(RenderingContext.TEXTURE_CUBE_MAP, RenderingContext.TEXTURE_MIN_FILTER, RenderingContext.LINEAR);
    */
  }
    
  void addCubemapFace(int face, String URL) {

    ImageElement image = new ImageElement();

    void handleDone(Event e) {
      print("addCubemapFace: handleDone: loaded image from URL: $URL");
      
      RenderingContext gl = this.program.gl;
      
      gl.bindTexture(RenderingContext.TEXTURE_CUBE_MAP, cubemapTexture);      
      gl.texParameteri(RenderingContext.TEXTURE_CUBE_MAP, RenderingContext.TEXTURE_MAG_FILTER, RenderingContext.NEAREST);
      gl.texParameteri(RenderingContext.TEXTURE_CUBE_MAP, RenderingContext.TEXTURE_MIN_FILTER, RenderingContext.NEAREST);
      
      gl.texImage2D(face, 0, RenderingContext.RGBA, RenderingContext.RGBA, RenderingContext.UNSIGNED_BYTE, image);

      gl.texParameteri(RenderingContext.TEXTURE_CUBE_MAP, RenderingContext.TEXTURE_WRAP_S, RenderingContext.CLAMP_TO_EDGE);
      gl.texParameteri(RenderingContext.TEXTURE_CUBE_MAP, RenderingContext.TEXTURE_WRAP_T, RenderingContext.CLAMP_TO_EDGE);
      
      gl.bindTexture(RenderingContext.TEXTURE_CUBE_MAP, null);
    }

    void handleError(Event e) {
      print("addCubemapFace: handleError: failure loading image from URL: $URL: $e");
    }

    image
      ..onLoad.listen(handleDone)
      ..onError.listen(handleError)
      ..src = URL;
  }
  
  void drawInstances(GameLoopHtml gameLoop, Camera cam) {
    
    RenderingContext gl = program.gl;
    
    gl.bindTexture(RenderingContext.TEXTURE_CUBE_MAP, cubemapTexture);

    this.instanceList.forEach((Instance i) => i.draw(gameLoop, cam));
    
    gl.bindTexture(RenderingContext.TEXTURE_CUBE_MAP, null);
  }  
}

class SkyboxInstance extends Instance {
  
  SkyboxInstance(Model model, vec3 center, double scale) : super(model, center, scale);
    
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

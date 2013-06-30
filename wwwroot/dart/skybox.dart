library skybox;

import 'dart:html';
import 'dart:web_gl';
import 'dart:math' as math;

import 'package:vector_math/vector_math.dart';
import 'package:game_loop/game_loop_html.dart';

import 'shader.dart';
import 'camera.dart';

class SkyboxProgram extends ShaderProgram {
  
  UniformLocation u_Skybox;
  
  SkyboxProgram(RenderingContext gl) : super(gl) {
  }

  void getLocations() {
    super.getLocations();
    
    u_Skybox = gl.getUniformLocation(program, "u_Skybox");
    
    print("SkyboxProgram: locations ready");      
  }

  void fetch(Map<String,Shader> shaderCache, String vertexShaderURL, fragmentShaderURL) {
    super.fetch(shaderCache, vertexShaderURL, fragmentShaderURL);
  }

  void drawModels(GameLoopHtml gameLoop, Camera cam, Matrix4 pMatrix) {

    if (!shaderReady) {
      return;
    }
    
    gl.useProgram(program);
    gl.enableVertexAttribArray(a_Position);
    
    int unit = 0;
    gl.activeTexture(RenderingContext.TEXTURE0 + unit);
    gl.uniform1i(this.u_Skybox, unit);

    // send perspective projection matrix uniform
    gl.uniformMatrix4fv(this.u_P, false, pMatrix.storage);

    modelList.forEach((Model m) => m.drawInstances(gameLoop, this, cam));

    // clean up
    gl.bindBuffer(RenderingContext.ARRAY_BUFFER, null);
    gl.bindBuffer(RenderingContext.ELEMENT_ARRAY_BUFFER, null);
    
    //gl.disableVertexAttribArray(a_Position); // needed ??
  }

}

class SkyboxModel extends Model {
  
  Texture cubemapTexture;
  
  SkyboxModel.fromJson(RenderingContext gl, String URL, bool reverse, num rescale):
    super.fromJson(gl, URL) {
    cubemapTexture = gl.createTexture();
  }
    
  void addCubemapFace(RenderingContext gl, int face, String URL) {

    ImageElement image = new ImageElement();

    void handleDone(Event e) {
      print("addCubemapFace: handleDone: loaded image from URL: $URL");
      
      //RenderingContext gl = program.gl;
      
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
  
  void drawInstances(GameLoopHtml gameLoop, ShaderProgram program, Camera cam) {

    RenderingContext gl = program.gl;
    
    gl.bindTexture(RenderingContext.TEXTURE_CUBE_MAP, cubemapTexture);

    instanceList.forEach((Instance i) => i.draw(gameLoop, program, cam));
    
    gl.bindTexture(RenderingContext.TEXTURE_CUBE_MAP, null);
  }  
}

class SkyboxInstance extends Instance {
  
  SkyboxInstance(Model model, Vector3 center, double scale) : super(model, center, scale);
  
  void draw(GameLoopHtml gameLoop, ShaderProgram prog, Camera cam) {
    
    double r = cam.getRad(gameLoop.renderInterpolationFactor);

    double size = 15 * math.sin(r).abs() + 1;

    setViewMatrix(MV, cam.eye, cam.center, cam.up);
    
    MV.translate(center[0], center[1], center[2]);
    
    double s = scale * size;
    MV.scale(s, s, s);

    //ShaderProgram prog = model.program;
    RenderingContext gl = prog.gl;

    gl.uniformMatrix4fv(prog.u_MV, false, MV.storage);

    gl.bindBuffer(RenderingContext.ARRAY_BUFFER, model.vertexPositionBuffer);
    gl.vertexAttribPointer(prog.a_Position, model.vertexPositionBufferItemSize, RenderingContext.FLOAT, false, 0, 0);
  
    gl.bindBuffer(RenderingContext.ELEMENT_ARRAY_BUFFER, model.vertexIndexBuffer);
    
    model.pieceList.forEach((Piece piece) {
      gl.drawElements(RenderingContext.TRIANGLES, piece.vertexIndexLength, RenderingContext.UNSIGNED_SHORT,
          piece.vertexIndexOffset * model.vertexIndexBufferItemSize);
    });
    
  }
}

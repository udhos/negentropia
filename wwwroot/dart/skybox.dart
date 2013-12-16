library skybox;

import 'dart:html';
import 'dart:web_gl';
import 'dart:math' as math;

import 'package:vector_math/vector_math.dart';
import 'package:game_loop/game_loop_html.dart';

import 'shader.dart';
import 'camera.dart';
import 'interpolate.dart';

class SkyboxProgram extends ShaderProgram {
  
  UniformLocation u_Skybox;
  
  SkyboxProgram(RenderingContext gl) : super(gl, "skyboxShader");

  void getLocations() {
    super.getLocations();
    
    u_Skybox = gl.getUniformLocation(program, "u_Skybox");
  }

  /*
  void fetch(Map<String,Shader> shaderCache, String vertexShaderURL, fragmentShaderURL) {
    super.fetch(shaderCache, vertexShaderURL, fragmentShaderURL);
  }
  */

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

    modelList.forEach((m) => m.drawInstances(gameLoop, this, cam));

    // clean up
    gl.bindBuffer(RenderingContext.ARRAY_BUFFER, null);
    gl.bindBuffer(RenderingContext.ELEMENT_ARRAY_BUFFER, null);
    
    //gl.disableVertexAttribArray(a_Position); // needed ??
  }

}

class SkyboxModel extends Model {
  
  Texture cubemapTexture;
  bool cubemapReady = false;
  
  SkyboxModel.fromJson(RenderingContext gl, String URL, bool reverse, num rescale):
    super.fromJson(gl, URL, reverse) {
    cubemapTexture = gl.createTexture();
  }
    
  void addCubemapFace(RenderingContext gl, int face, String URL) {

    ImageElement image = new ImageElement();

    void handleDone(Event e) {
      gl.bindTexture(RenderingContext.TEXTURE_CUBE_MAP, cubemapTexture);      
      gl.texParameteri(RenderingContext.TEXTURE_CUBE_MAP, RenderingContext.TEXTURE_MAG_FILTER, RenderingContext.NEAREST);
      gl.texParameteri(RenderingContext.TEXTURE_CUBE_MAP, RenderingContext.TEXTURE_MIN_FILTER, RenderingContext.NEAREST);
      
      gl.texImage2DImage(face, 0, RenderingContext.RGBA, RenderingContext.RGBA, RenderingContext.UNSIGNED_BYTE, image);

      gl.texParameteri(RenderingContext.TEXTURE_CUBE_MAP, RenderingContext.TEXTURE_WRAP_S, RenderingContext.CLAMP_TO_EDGE);
      gl.texParameteri(RenderingContext.TEXTURE_CUBE_MAP, RenderingContext.TEXTURE_WRAP_T, RenderingContext.CLAMP_TO_EDGE);
      
      //anisotropic_filtering_enable(gl);
      
      gl.bindTexture(RenderingContext.TEXTURE_CUBE_MAP, null);
      
      cubemapReady = true;
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
    if (!modelReady || !cubemapReady) {
      return;
    }

    RenderingContext gl = program.gl;
    
    gl.bindTexture(RenderingContext.TEXTURE_CUBE_MAP, cubemapTexture);

    instanceList.forEach((Instance i) => i.draw(gameLoop, program, cam));
    
    gl.bindTexture(RenderingContext.TEXTURE_CUBE_MAP, null);
  }  
}

class SkyboxInstance extends Instance {

  // demo animate fields:
  final double degreesPerSec = 20.0;
  bool demoAnimate;
  double _oldAngle = 0.0;
  double _angle = 0.0;
  
  SkyboxInstance(Model model, Vector3 center, double scale, this.demoAnimate) : super(model, center, scale);

  void update(GameLoopHtml gameLoop) {
    
    if (demoAnimate) {
      _oldAngle = _angle;
      _angle = gameLoop.gameTime * this.degreesPerSec % 360.0;
    }
  }
  
  void draw(GameLoopHtml gameLoop, ShaderProgram prog, Camera cam) {

    double rescale;
    
    if (demoAnimate) {
      double deg = interpolateDegree(_angle, _oldAngle, gameLoop.renderInterpolationFactor);
      double rad = deg * math.PI / 180.0;
      double size = 15 * math.sin(rad).abs() + 1;
      rescale = scale * size;
    }
    else {
      rescale = scale;
    }      
    
    RenderingContext gl = prog.gl;

    modelView(gl, prog.u_MV, cam, rescale); // set up MV matrix    

    gl.bindBuffer(RenderingContext.ARRAY_BUFFER, model.vertexPositionBuffer);
    gl.vertexAttribPointer(prog.a_Position, model.vertexPositionBufferItemSize, RenderingContext.FLOAT, false, 0, 0);
  
    gl.bindBuffer(RenderingContext.ELEMENT_ARRAY_BUFFER, model.vertexIndexBuffer);
    
    model.pieceList.forEach((piece) {
      gl.drawElements(RenderingContext.TRIANGLES, piece.vertexIndexLength, RenderingContext.UNSIGNED_SHORT,
          piece.vertexIndexOffset * model.vertexIndexBufferItemSize);
    });
    
  }
}

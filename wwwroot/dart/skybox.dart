library skybox;

import 'dart:html';
import 'dart:web_gl';

import 'package:vector_math/vector_math.dart';

import 'shader.dart';
import 'buffer.dart';
import 'camera.dart';

class SkyboxProgram extends ShaderProgram {
  
  UniformLocation u_P;
  UniformLocation u_Skybox;
  
  SkyboxProgram(RenderingContext gl) : super(gl) {
  }

  void fetch(Map<String,Shader> shaderCache, String vertexShaderURL, fragmentShaderURL) {
    super.fetch(shaderCache, vertexShaderURL, fragmentShaderURL);
    
    this.u_P      = gl.getUniformLocation(this.program, "u_P");
    this.u_Skybox = gl.getUniformLocation(this.program, "u_Skybox");
  }

  void drawModels(Camera cam, mat4 pMatrix) {
    
    gl.useProgram(program);
    gl.enableVertexAttribArray(a_Position);
    
    int unit = 0;
    gl.activeTexture(RenderingContext.TEXTURE0 + unit);
    gl.uniform1i(this.u_Skybox, unit);

    // send perspective projection matrix uniform
    gl.uniformMatrix4fv(this.u_P, false, pMatrix);

    modelList.forEach((Model m) => m.drawInstances(cam));

    // clean up
    gl.bindBuffer(RenderingContext.ARRAY_BUFFER, null);
    gl.bindBuffer(RenderingContext.ELEMENT_ARRAY_BUFFER, null);
    
    //gl.disableVertexAttribArray(a_Position); // needed ??
  }

}

class SkyboxModel extends Model {
  
  Texture cubemapTexture;
  
  SkyboxModel.fromURL(RenderingContext gl, SkyboxProgram prog, String URL, bool reverse, num rescale) : super.fromURL(gl, prog, URL) {
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
  
  void drawInstances(Camera cam) {
    gl.bindTexture(RenderingContext.TEXTURE_CUBE_MAP, cubemapTexture);

    this.instanceList.forEach((Instance i) => i.draw(cam));
    
    gl.bindTexture(RenderingContext.TEXTURE_CUBE_MAP, null);
  }  
}
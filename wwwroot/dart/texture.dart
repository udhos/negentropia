library texture;

import 'dart:html';
import 'dart:async';
import 'dart:web_gl';
import 'dart:typed_data';

class TextureInfo {
  
  String textureName;
  Texture texture;
  List<int> temporaryColor;

  void loadSolidColor(RenderingContext gl) {
    gl.bindTexture(RenderingContext.TEXTURE_2D, texture);
    gl.texImage2DTyped(RenderingContext.TEXTURE_2D, 0, RenderingContext.RGBA, 1, 1, 0, RenderingContext.RGBA, RenderingContext.UNSIGNED_BYTE, new Uint8List.fromList(temporaryColor));
    gl.texParameteri(RenderingContext.TEXTURE_2D, RenderingContext.TEXTURE_MAG_FILTER, RenderingContext.NEAREST);
    gl.texParameteri(RenderingContext.TEXTURE_2D, RenderingContext.TEXTURE_MIN_FILTER, RenderingContext.NEAREST);
    gl.texParameteri(RenderingContext.TEXTURE_2D, RenderingContext.TEXTURE_WRAP_S, RenderingContext.CLAMP_TO_EDGE);
    gl.texParameteri(RenderingContext.TEXTURE_2D, RenderingContext.TEXTURE_WRAP_T, RenderingContext.CLAMP_TO_EDGE);  
    gl.bindTexture(RenderingContext.TEXTURE_2D, null);    
  }
  
  void _loadTexture2D(RenderingContext gl, Map<String,Texture> textureTable) {
    
    assert(texture != null);
        
    ImageElement image = new ImageElement();

    if (textureName != null) {
      textureTable[textureName] = texture;
    }
    
    // temporary solid color texture
    loadSolidColor(gl);
    
    void onDone(Event e) {
      
      gl.bindTexture(RenderingContext.TEXTURE_2D, texture);
      //gl.pixelStorei(RenderingContext.UNPACK_FLIP_Y_WEBGL, true);
      gl.pixelStorei(RenderingContext.UNPACK_FLIP_Y_WEBGL, 1);
      
      gl.texImage2DImage(RenderingContext.TEXTURE_2D, 0, RenderingContext.RGBA, RenderingContext.RGBA, RenderingContext.UNSIGNED_BYTE, image);
      
      gl.texParameteri(RenderingContext.TEXTURE_2D, RenderingContext.TEXTURE_MAG_FILTER, RenderingContext.NEAREST);
      gl.texParameteri(RenderingContext.TEXTURE_2D, RenderingContext.TEXTURE_MIN_FILTER, RenderingContext.NEAREST);
      gl.texParameteri(RenderingContext.TEXTURE_2D, RenderingContext.TEXTURE_WRAP_S, RenderingContext.CLAMP_TO_EDGE);
      gl.texParameteri(RenderingContext.TEXTURE_2D, RenderingContext.TEXTURE_WRAP_T, RenderingContext.CLAMP_TO_EDGE);   
      gl.bindTexture(RenderingContext.TEXTURE_2D, null);
    }

    void onError(Event e) {
      print("TextureInfo: handleError: failure loading image from URL: $textureName: $e");
    }

    // fetch definitive texture
    image
    ..onLoad.listen(onDone)
    ..onError.listen(onError)
    ..src = textureName;
  }
  
  void createTexture(RenderingContext gl, Map<String,Texture> textureTable) {

    texture = gl.createTexture();
    if (texture == null) {
      print("TextureInfo: could not create texture for: $textureName");
      return;    
    }

    if (textureName == null) {
      loadSolidColor(gl);
      return;
    }
    
    _loadTexture2D(gl, textureTable);    
  }

  void loadTexture(RenderingContext gl, Map<String,Texture> textureTable) {
    
    if (textureName != null) {
      texture = textureTable[textureName];
      if (texture != null) {
        return;
      }
    }

    createTexture(gl, textureTable);    
  }
  
  TextureInfo(RenderingContext gl, Map<String,Texture> textureTable, this.textureName,
              List<int> this.temporaryColor) {    
    loadTexture(gl, textureTable);    
  }
}
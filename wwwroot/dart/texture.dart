library texture;

import 'dart:html';
import 'dart:async';
import 'dart:web_gl';
import 'dart:typed_data';

class TextureInfo {
  
  //int indexOffset;
  //int indexNumber;
  String textureName;
  Texture texture;
  List<int> temporaryColor;

  void _loadTexture2D(RenderingContext gl, Map<String,Texture> textureTable, String textureName, List<int> temporaryColor,
                      void handleDone(Event e), void handleError(Event e)) {
    
    assert(texture != null);
        
    ImageElement image = new ImageElement();

    void onDone(Event e) {
      
      gl.bindTexture(RenderingContext.TEXTURE_2D, texture);
      //gl.pixelStorei(RenderingContext.UNPACK_FLIP_Y_WEBGL, true);
      gl.pixelStorei(RenderingContext.UNPACK_FLIP_Y_WEBGL, 1);
      gl.texImage2D(RenderingContext.TEXTURE_2D, 0, RenderingContext.RGBA, RenderingContext.RGBA, RenderingContext.UNSIGNED_BYTE, image);
      gl.texParameteri(RenderingContext.TEXTURE_2D, RenderingContext.TEXTURE_MAG_FILTER, RenderingContext.NEAREST);
      gl.texParameteri(RenderingContext.TEXTURE_2D, RenderingContext.TEXTURE_MIN_FILTER, RenderingContext.NEAREST);
      gl.texParameteri(RenderingContext.TEXTURE_2D, RenderingContext.TEXTURE_WRAP_S, RenderingContext.CLAMP_TO_EDGE);
      gl.texParameteri(RenderingContext.TEXTURE_2D, RenderingContext.TEXTURE_WRAP_T, RenderingContext.CLAMP_TO_EDGE);   
      gl.bindTexture(RenderingContext.TEXTURE_2D, null);
            
      handleDone(e);
    }

    textureTable[textureName] = texture;
    
    // temporary solid color texture
    loadSolidColor(gl);

    // fetch definitive texture
    image
    ..onLoad.listen(onDone)
    ..onError.listen(handleError)
    ..src = textureName;
  }
  
  void loadSolidColor(RenderingContext gl) {
    gl.bindTexture(RenderingContext.TEXTURE_2D, texture);
    gl.texImage2D(RenderingContext.TEXTURE_2D, 0, RenderingContext.RGBA, 1, 1, 0, RenderingContext.RGBA, RenderingContext.UNSIGNED_BYTE, new Uint8List.fromList(temporaryColor));
    gl.texParameteri(RenderingContext.TEXTURE_2D, RenderingContext.TEXTURE_MAG_FILTER, RenderingContext.NEAREST);
    gl.texParameteri(RenderingContext.TEXTURE_2D, RenderingContext.TEXTURE_MIN_FILTER, RenderingContext.NEAREST);
    gl.texParameteri(RenderingContext.TEXTURE_2D, RenderingContext.TEXTURE_WRAP_S, RenderingContext.CLAMP_TO_EDGE);
    gl.texParameteri(RenderingContext.TEXTURE_2D, RenderingContext.TEXTURE_WRAP_T, RenderingContext.CLAMP_TO_EDGE);  
    gl.bindTexture(RenderingContext.TEXTURE_2D, null);    
  }

  void forceCreateTexture(RenderingContext gl, Map<String,Texture> textureTable) {

    texture = gl.createTexture();
    if (texture == null) {
      String fail = "could not create texture for: $textureName"; 
      print("TextureInfo: $fail");
      return;    
    }

    if (textureName == null) {
      loadSolidColor(gl);
      return;
    }
    
    void handleDone(Event e) {
      print("TextureInfo: handleDone: loaded image from URL: $textureName");
    }

    void handleError(Event e) {
      print("TextureInfo: handleError: failure loading image from URL: $textureName: $e");
    }
    
    _loadTexture2D(gl, textureTable, textureName, temporaryColor, handleDone, handleError);    
  }
  
  TextureInfo(RenderingContext gl, Map<String,Texture> textureTable, this.textureName,
              List<int> this.temporaryColor) {
    
    //print("TextureInfo: indexOffset=$indexOffset indexNumber=$indexNumber");
        
    texture = textureTable[textureName];
    if (texture != null) {
      print("TextureInfo: texture table HIT: $textureName");
      return;
    }

    forceCreateTexture(gl, textureTable);    
  }
}
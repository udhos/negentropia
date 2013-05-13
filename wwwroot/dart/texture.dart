library texture;

import 'dart:html';
import 'dart:async';
import 'dart:web_gl';

void loadTexture2D(RenderingContext gl, Map<String,Texture> textureTable, String textureName, void handleDone(Event e), void handleError(Event e)) {
  
  ImageElement image = new ImageElement();
  
  void onDone(Event e) {
        
    Texture tex = gl.createTexture();
    
    gl.bindTexture(RenderingContext.TEXTURE_2D, tex);
    //gl.pixelStorei(RenderingContext.UNPACK_FLIP_Y_WEBGL, true);
    gl.pixelStorei(RenderingContext.UNPACK_FLIP_Y_WEBGL, 1);

    gl.texParameteri(RenderingContext.TEXTURE_2D, RenderingContext.TEXTURE_MAG_FILTER, RenderingContext.NEAREST);
    gl.texParameteri(RenderingContext.TEXTURE_2D, RenderingContext.TEXTURE_MIN_FILTER, RenderingContext.NEAREST);

    gl.texImage2D(RenderingContext.TEXTURE_2D, 0, RenderingContext.RGBA, RenderingContext.RGBA, RenderingContext.UNSIGNED_BYTE, image);

    gl.texParameteri(RenderingContext.TEXTURE_2D, RenderingContext.TEXTURE_WRAP_S, RenderingContext.CLAMP_TO_EDGE);
    gl.texParameteri(RenderingContext.TEXTURE_2D, RenderingContext.TEXTURE_WRAP_T, RenderingContext.CLAMP_TO_EDGE);  
    
    gl.bindTexture(RenderingContext.TEXTURE_2D, null);
    
    textureTable[textureName] = tex;
    
    print("loadTexture2D: fetched: $textureName");
    
    handleDone(e);
  }

  image
    ..onLoad.listen(onDone)
    ..onError.listen(handleError)
    ..src = textureName;
}

class TextureInfo {
  
  int indexOffset;
  int indexNumber;
  String textureName;
  Texture texture;
  
  void forceCreateTexture(RenderingContext gl, Map<String,Texture> textureTable) {
    
    void handleDone(Event e) {
      print("TextureInfo: handleDone: loaded image from URL: $textureName");
      
      texture = textureTable[textureName];
    }

    void handleError(Event e) {
      print("TextureInfo: handleError: failure loading image from URL: $textureName: $e");
    }
    
    loadTexture2D(gl, textureTable, textureName, handleDone, handleError);    
  }
  
  TextureInfo(RenderingContext gl, Map<String,Texture> textureTable, this.indexOffset, this.indexNumber, this.textureName) {
    
    texture = textureTable[textureName];
    if (texture != null) {
      print("TextureInfo: texture table HIT: $textureName");
      return;
    }

    forceCreateTexture(gl, textureTable);    
  }
}
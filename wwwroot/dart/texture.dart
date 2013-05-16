library texture;

import 'dart:html';
import 'dart:async';
import 'dart:web_gl';
import 'dart:typed_data';

class TextureInfo {
  
  int indexOffset;
  int indexNumber;
  String textureName;
  Texture texture;

  void loadTexture2D(RenderingContext gl, Map<String,Texture> textureTable, String textureName, List<int> temporaryColor,
                     void handleDone(Event e), void handleError(Event e)) {
    
    ImageElement image = new ImageElement();
    
    texture = gl.createTexture();
    
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

    if (texture == null) {
      String fail = "could not create texture for: $textureName"; 
      print("loadTexture2D: $fail");
      handleError(new Event(fail));
      return;    
    }

    textureTable[textureName] = texture;

    
    // temporary solid color texture
    gl.bindTexture(RenderingContext.TEXTURE_2D, texture);
    gl.texImage2D(RenderingContext.TEXTURE_2D, 0, RenderingContext.RGBA, 1, 1, 0, RenderingContext.RGBA, RenderingContext.UNSIGNED_BYTE, new Uint8List.fromList(temporaryColor));
    gl.texParameteri(RenderingContext.TEXTURE_2D, RenderingContext.TEXTURE_MAG_FILTER, RenderingContext.NEAREST);
    gl.texParameteri(RenderingContext.TEXTURE_2D, RenderingContext.TEXTURE_MIN_FILTER, RenderingContext.NEAREST);
    gl.texParameteri(RenderingContext.TEXTURE_2D, RenderingContext.TEXTURE_WRAP_S, RenderingContext.CLAMP_TO_EDGE);
    gl.texParameteri(RenderingContext.TEXTURE_2D, RenderingContext.TEXTURE_WRAP_T, RenderingContext.CLAMP_TO_EDGE);  
    gl.bindTexture(RenderingContext.TEXTURE_2D, null);
    

    // load definitive texture
    image
    ..onLoad.listen(onDone)
    ..onError.listen(handleError)
    ..src = textureName;
  }

  void forceCreateTexture(RenderingContext gl, Map<String,Texture> textureTable, List<int> temporaryColor) {
    
    void handleDone(Event e) {
      print("TextureInfo: handleDone: loaded image from URL: $textureName");
    }

    void handleError(Event e) {
      print("TextureInfo: handleError: failure loading image from URL: $textureName: $e");
    }
    
    loadTexture2D(gl, textureTable, textureName, temporaryColor, handleDone, handleError);    
  }
  
  TextureInfo(RenderingContext gl, Map<String,Texture> textureTable, this.indexOffset, this.indexNumber, this.textureName,
              List<int> temporaryColor) {
    
    //print("TextureInfo: indexOffset=$indexOffset indexNumber=$indexNumber");
    
    texture = textureTable[textureName];
    if (texture != null) {
      print("TextureInfo: texture table HIT: $textureName");
      return;
    }

    forceCreateTexture(gl, textureTable, temporaryColor);    
  }
}
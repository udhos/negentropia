library texture;

import 'dart:html';
import 'dart:web_gl';
import 'dart:typed_data';

import 'anisotropic.dart';
import 'logg.dart';

int defaultTextureUnit = 0;

class TextureInfo {
  String textureName;
  Texture texture;
  List<int> temporaryColor;
  int textureUnit;
  int wrap;

  void loadSolidColor(RenderingContext gl) {

    // bind textureUnit to texture
    //gl.activeTexture(RenderingContext.TEXTURE0 + textureUnit);
    gl.bindTexture(RenderingContext.TEXTURE_2D, texture);

    gl.texImage2DTyped(RenderingContext.TEXTURE_2D, 0, RenderingContext.RGBA, 1,
        1, 0, RenderingContext.RGBA, RenderingContext.UNSIGNED_BYTE,
        new Uint8List.fromList(temporaryColor));
    gl.texParameteri(RenderingContext.TEXTURE_2D,
        RenderingContext.TEXTURE_MAG_FILTER, RenderingContext.NEAREST);
    gl.texParameteri(RenderingContext.TEXTURE_2D,
        RenderingContext.TEXTURE_MIN_FILTER, RenderingContext.NEAREST);
    gl.texParameteri(RenderingContext.TEXTURE_2D,
        RenderingContext.TEXTURE_WRAP_S, RenderingContext.CLAMP_TO_EDGE);
    gl.texParameteri(RenderingContext.TEXTURE_2D,
        RenderingContext.TEXTURE_WRAP_T, RenderingContext.CLAMP_TO_EDGE);

    gl.bindTexture(RenderingContext.TEXTURE_2D, null);
  }

  bool isPowerOfTwo(int v) {
    return v != 0 && (v & (v - 1)) == 0;
  }

  void _loadTexture2D(RenderingContext gl, Map<String, Texture> textureTable) {
    assert(texture != null);

    ImageElement image = new ImageElement();

    if (textureName != null) {
      textureTable[textureName] = texture;
    }

    // temporary solid color texture
    loadSolidColor(gl);

    void onDone(Event e) {

      // bind textureUnit to texture
      //gl.activeTexture(RenderingContext.TEXTURE0 + textureUnit);
      gl.bindTexture(RenderingContext.TEXTURE_2D, texture);

      //gl.pixelStorei(RenderingContext.UNPACK_FLIP_Y_WEBGL, true);
      gl.pixelStorei(RenderingContext.UNPACK_FLIP_Y_WEBGL, 1);

      gl.texImage2DImage(RenderingContext.TEXTURE_2D, 0, RenderingContext.RGBA,
          RenderingContext.RGBA, RenderingContext.UNSIGNED_BYTE, image);

      // undo flip Y otherwise it could affect other texImage calls
      gl.pixelStorei(RenderingContext.UNPACK_FLIP_Y_WEBGL, 0);

      bool mipmap = isPowerOfTwo(image.width) && isPowerOfTwo(image.height);

      if (mipmap) {
        gl.texParameteri(RenderingContext.TEXTURE_2D,
            RenderingContext.TEXTURE_MAG_FILTER, RenderingContext.LINEAR);
        gl.texParameteri(RenderingContext.TEXTURE_2D,
            RenderingContext.TEXTURE_MIN_FILTER,
            RenderingContext.LINEAR_MIPMAP_NEAREST);
        gl.generateMipmap(RenderingContext.TEXTURE_2D);
      } else {
        log("can't enable MIPMAP for NPOT texture: $textureName");
        gl.texParameteri(RenderingContext.TEXTURE_2D,
            RenderingContext.TEXTURE_MAG_FILTER, RenderingContext.NEAREST);
        gl.texParameteri(RenderingContext.TEXTURE_2D,
            RenderingContext.TEXTURE_MIN_FILTER, RenderingContext.NEAREST);
      }

      gl.texParameteri(
          RenderingContext.TEXTURE_2D, RenderingContext.TEXTURE_WRAP_S, wrap);
      gl.texParameteri(
          RenderingContext.TEXTURE_2D, RenderingContext.TEXTURE_WRAP_T, wrap);

      int anisotropy = anisotropic_filtering_enable(gl, textureName);

      gl.bindTexture(RenderingContext.TEXTURE_2D, null);

      log("texture loaded: url=$textureName size=${image.width}x${image.height} mipmap=$mipmap anisotropy=$anisotropy");
    }

    void onError(Event e) {
      err("TextureInfo: handleError: failure loading image from URL: $textureName: $e");
    }

    // fetch definitive texture
    image
      ..onLoad.listen(onDone)
      ..onError.listen(onError)
      ..src = textureName;
  }

  void createTexture(RenderingContext gl, Map<String, Texture> textureTable) {
    texture = gl.createTexture();
    if (texture == null) {
      err("TextureInfo: could not create texture for: $textureName");
      return;
    }

    if (textureName == null) {
      loadSolidColor(gl);
      return;
    }

    _loadTexture2D(gl, textureTable);
  }

  void loadTexture(RenderingContext gl, Map<String, Texture> textureTable) {
    if (textureName != null) {
      texture = textureTable[textureName];
      if (texture != null) {
        return;
      }
    }

    createTexture(gl, textureTable);
  }

  TextureInfo(RenderingContext gl, Map<String, Texture> textureTable,
      this.textureName, List<int> this.temporaryColor, this.textureUnit,
      this.wrap) {
    loadTexture(gl, textureTable);
  }
}

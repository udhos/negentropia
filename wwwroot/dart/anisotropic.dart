library anisotropic;

import 'dart:web_gl';

import 'logg.dart';

final List<String> _names = [
  "EXT_texture_filter_anisotropic",
  "MOZ_EXT_texture_filter_anisotropic",
  "WEBKIT_EXT_texture_filter_anisotropic"
];
ExtTextureFilterAnisotropic _extAnisotropic = null;
int _anisotropy = 16;

void anisotropic_filtering_detect(RenderingContext gl) {
  void enable(ExtTextureFilterAnisotropic ext, String name) {
    int max_anisotropy = gl.getParameter(
        ExtTextureFilterAnisotropic.MAX_TEXTURE_MAX_ANISOTROPY_EXT);

    if (_anisotropy > max_anisotropy) {
      _anisotropy = max_anisotropy;
    }

    log("using anisotropy=$_anisotropy");

    _extAnisotropic = ext;
  }

  for (int i = 0; i < _names.length; ++i) {
    String extName = _names[i];
    ExtTextureFilterAnisotropic ext;
    try {
      ext = gl.getExtension(extName);
    } catch (exc) {
      warn("gl.getExtension('$extName') exception: $exc");
    }
    if (ext != null) {
      enable(ext, extName);
      return;
    }
  }

  _anisotropy = 0; // disabled

  warn("anisotropic filtering: NOT SUPPORTED");
}

int anisotropic_filtering_enable(RenderingContext gl, String url) {
  if (_anisotropy < 2 || _extAnisotropic == null) {
    // not supported
    return 0;
  }

  gl.texParameterf(RenderingContext.TEXTURE_2D,
      ExtTextureFilterAnisotropic.TEXTURE_MAX_ANISOTROPY_EXT,
      _anisotropy.toDouble());

  double result = gl.getTexParameter(RenderingContext.TEXTURE_2D,
      ExtTextureFilterAnisotropic.TEXTURE_MAX_ANISOTROPY_EXT);

  //log("texture url=$url anisotropy=$result");

  if (result != _anisotropy.toDouble()) {
    warn(
        "anisotropic_filtering_enable: anisotropy set to texParameterf=${_anisotropy.toDouble()} but got getTexParameter=$result");
  }

  return result.toInt();
}

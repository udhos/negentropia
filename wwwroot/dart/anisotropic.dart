library anisotropic;

import 'dart:web_gl';

import 'logg.dart';

final List<String> _names = ["EXT_texture_filter_anisotropic",
                             "MOZ_EXT_texture_filter_anisotropic",
                             "WEBKIT_EXT_texture_filter_anisotropic"];
ExtTextureFilterAnisotropic _extAnisotropic = null;
int _anisotropy = 16;

void anisotropic_filtering_detect(RenderingContext gl) {
  
  void enable(ExtTextureFilterAnisotropic ext, String name) {
    //print("anisotropic extension: name=$name class=$ext");
        
    int max_anisotropy = gl.getParameter(ExtTextureFilterAnisotropic.MAX_TEXTURE_MAX_ANISOTROPY_EXT);
    //print("max anisotropy: $max_anisotropy");
    
    if (_anisotropy > max_anisotropy) {
      _anisotropy = max_anisotropy;
    }

    debug("using anisotropy=$_anisotropy");
    
    _extAnisotropic = ext;
  }

  for (int i = 0; i < _names.length; ++i) {
    String extName = _names[i];
    ExtTextureFilterAnisotropic ext;
    try {
      ext = gl.getExtension(extName);
    }
    catch (exc) {
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

void anisotropic_filtering_enable(RenderingContext gl) {

  if (_anisotropy < 2 || _extAnisotropic == null) {
    // not supported
    return;
  }
  
  //print("enabling anisotropy=$_anisotropy on texture");
  
  gl.texParameterf(RenderingContext.TEXTURE_2D, ExtTextureFilterAnisotropic.TEXTURE_MAX_ANISOTROPY_EXT, _anisotropy.toDouble());

  int result = gl.getTexParameter(RenderingContext.TEXTURE_2D, ExtTextureFilterAnisotropic.TEXTURE_MAX_ANISOTROPY_EXT);
  
  debug("texture anisotropy=$result");        
  
  if (result != _anisotropy) {
    warn("anisotropic_filtering_enable: anisotropy set to texParameterf=$_anisotropy but got getTexParameter=$result");
  }
}

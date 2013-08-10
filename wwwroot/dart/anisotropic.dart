library anisotropic;

import 'dart:web_gl';

ExtTextureFilterAnisotropic _extAnisotropic = null;
double _anisotropy = 16.0;

void anisotropic_filtering_detect(RenderingContext gl) {
  
  void enable(ExtTextureFilterAnisotropic ext, String name) {
    print("anisotropic extension: name=$name class=$ext");
        
    double max_anisotropy;
    try {
      max_anisotropy = gl.getParameter(ext.MAX_TEXTURE_MAX_ANISOTROPY_EXT);
    }
    catch (e) {
      print("FIXME gl.getParameter: ext.MAX_TEXTURE_MAX_ANISOTROPY_EXT: exception: $e");
      int max = gl.getParameter(0x84FF);
      max_anisotropy = max.toDouble();
    }
    print("max anisotropy: $max_anisotropy");
    
    if (_anisotropy > max_anisotropy) {
      _anisotropy = max_anisotropy;
    }

    print("using anisotropy=$_anisotropy");
    
    _extAnisotropic = ext;
  }
  
  String e = "EXT_texture_filter_anisotropic";
  ExtTextureFilterAnisotropic ext = gl.getExtension(e);
  if (ext != null) {
    enable(ext, e);
    return;
  }

  e = "MOZ_EXT_texture_filter_anisotropic";
  ext = gl.getExtension(e);
  if (ext != null) {
    enable(ext, e);
    return;
  }
  
  e = "WEBKIT_EXT_texture_filter_anisotropic";
  ext = gl.getExtension(e);
  if (ext != null) {
    enable(ext, e);
    return;
  }
  
  _anisotropy = 0.0; // disabled  

  print("anisotropic filtering: NOT SUPPORTED");
}

void anisotropic_filtering_enable(RenderingContext gl) {

  if (_anisotropy < 1.0 || _extAnisotropic == null) {
    return;
  }
  
  print("enabling anisotropy=$_anisotropy on texture");
  
  try {
    gl.texParameterf(RenderingContext.TEXTURE_2D, _extAnisotropic.TEXTURE_MAX_ANISOTROPY_EXT, _anisotropy);
  }
  catch (e) {
    print("FIXME gl.texParameterf: ext.TEXTURE_MAX_ANISOTROPY_EXT: exception: $e");
    gl.texParameterf(RenderingContext.TEXTURE_2D, 0x84FE, _anisotropy);
  }

  double result;
  try {
    result = gl.getTexParameter(RenderingContext.TEXTURE_2D, _extAnisotropic.TEXTURE_MAX_ANISOTROPY_EXT);
  }
  catch (e) {
    print("FIXME gl.getTexParameter: ext.TEXTURE_MAX_ANISOTROPY_EXT: exception: $e");
    int r = gl.getTexParameter(RenderingContext.TEXTURE_2D, 0x84FE);
    result = r.toDouble();
  }
  
  print("texture anisotropy=$result");        
}

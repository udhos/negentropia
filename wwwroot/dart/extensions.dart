library extensions;

import 'dart:web_gl';

import 'logg.dart';

bool _ext_element_uint = false;
int _ext_element_type = RenderingContext.UNSIGNED_SHORT;
int _ext_element_size = 2;

bool get ext_element_uint => _ext_element_uint;

int get ext_get_element_type => _ext_element_type;

int get ext_get_element_size => _ext_element_size;

void enable_element_uint(RenderingContext gl) {
  String extName = "OES_element_index_uint";
  OesElementIndexUint ext;

  try {
    ext = gl.getExtension(extName);
  } catch (exc) {
    warn("gl.getExtension('$extName') exception: $exc");
  }

  _ext_element_uint = ext != null;
  
  if (_ext_element_uint) {
    _ext_element_type = RenderingContext.UNSIGNED_INT;
    _ext_element_size = 4;
  }
  else {
    _ext_element_type = RenderingContext.UNSIGNED_SHORT;
    _ext_element_size = 2;
  }

  log("gl.getExtension('$extName'): available = $ext_element_uint");
}

void enable_extensions(RenderingContext gl) {
  enable_element_uint(gl);
}

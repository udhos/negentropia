library extensions;

import 'dart:web_gl';

import 'logg.dart';

bool _ext_element_uint = false;

bool get ext_element_uint => _ext_element_uint;

int get ext_get_element_type {
  if (ext_element_uint) {
    return RenderingContext.UNSIGNED_INT;
  }

  return RenderingContext.UNSIGNED_SHORT;
}

void enable_element_uint(RenderingContext gl) {
  String extName = "OES_element_index_uint";
  OesElementIndexUint ext;

  try {
    ext = gl.getExtension(extName);
  } catch (exc) {
    warn("gl.getExtension('$extName') exception: $exc");
  }

  _ext_element_uint = ext != null;

  log("gl.getExtension('$extName'): available = $ext_element_uint");
}

void enable_extensions(RenderingContext gl) {
  enable_element_uint(gl);
}

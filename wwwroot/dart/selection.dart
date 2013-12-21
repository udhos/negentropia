library selection;

//import 'dart:html';
import 'dart:web_gl';
import 'dart:typed_data';
import 'dart:collection';

//import 'package:game_loop/game_loop_html.dart';

import 'shader.dart';
import 'logg.dart';

Set<PickerInstance> _selection = new HashSet<PickerInstance>();

PickerInstance colorHit(Iterable<Instance> list, int r,g,b) {

  /*
  bool matchColor(Float32List f, int r,g,b) {
    
    double d0 = (255.0*f[0] - r.toDouble()).abs();
    double d1 = (255.0*f[1] - g.toDouble()).abs();
    double d2 = (255.0*f[2] - b.toDouble()).abs();
    
    return d0 < 1.0 && d1 < 1.0 && d2 < 1.0;
  }
  */

  bool match(Instance i) {
    Float32List f = i.pickColor;
    return (255.0*f[0] - r.toDouble()).abs() < 1.0 &&
        (255.0*f[1] - g.toDouble()).abs() < 1.0 &&
        (255.0*f[2] - b.toDouble()).abs() < 1.0;
  }
  
  Instance pi;
    
  try {
    pi = list.firstWhere(match);
  } catch (e) {
    return null;
  }
  
  return pi as PickerInstance;
}

void handleSelection(PickerInstance pi, bool shift) {
  
  assert(shift != null);
  
  if (pi == null) {
    // didn't hit anything
    if (!shift) {
      // shift is released
      _selection.clear();
    }
    return;
  }
  
  assert(pi != null);
  
  if (shift) {
    if (_selection.contains(pi)) {
      _selection.remove(pi);
    }
    else {
      _selection.add(pi);      
    }
    return;
  }
  
  _selection.clear();
  _selection.add(pi);
}

void mouseSelection(PickerInstance pi, bool shift) {
  handleSelection(pi, shift);
  debug("mouseSelection: $_selection");
}

Uint8List _color = new Uint8List(4);

void bandSelection(int x, y, width, height, PickerShader picker, RenderingContext gl, bool shift) {
  //debug("bandSelection: x=$x y=$y width=$width height=$height");
  
  if (picker == null) {
    err("bandSelection: picker not available");
    return;
  }

  int size = 4 * width * height;
  
  debug("bandSelection: color buffer current=${_color.length} needed=$size");
  
  if (size > _color.length) {
    _color = new Uint8List(size);
  }
  
  gl.bindFramebuffer(RenderingContext.FRAMEBUFFER, picker.framebuffer);
  gl.readPixels(x, y, width, height, RenderingContext.RGBA, RenderingContext.UNSIGNED_BYTE, _color);
  
  if (!shift) {
    // shift is released
    _selection.clear();
  }
  
  for (int i = 0; i < size; i += 4) {
    
    if (_selection.length >= picker.numberOfInstances) {
      // selected all available objects, no need to keep searching
      break;
    }
    
    PickerInstance pi = picker.findInstanceByColor(_color[i], _color[i+1], _color[i+2]);
    if (pi == null) {
      continue;
    }
    
    _selection.add(pi);
  }
  
  debug("bandSelection: $_selection");  
}

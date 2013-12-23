library selection;

import 'dart:web_gl';
import 'dart:typed_data';
import 'dart:collection';

import 'shader.dart';
import 'logg.dart';

Set<PickerInstance> _selection = new HashSet<PickerInstance>();

PickerInstance colorHit(Iterable<Instance> list, int r,g,b) {

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

double _bgColorR;
double _bgColorG;
double _bgColorB;

void pickerClearColor(double r, g, b) {
  _bgColorR = r;
  _bgColorG = g;
  _bgColorB = b;  
}

bool backgroundColor(int r, g, b) {
  return (r.toDouble() - 255.0 * _bgColorR).abs() < 1.0
      && (g.toDouble() - 255.0 * _bgColorG).abs() < 1.0
      && (b.toDouble() - 255.0 * _bgColorB).abs() < 1.0;
}

Uint8List _color = new Uint8List(4);

void bandSelection(int x, y, width, height, PickerShader picker, RenderingContext gl, bool shift) {
  //debug("bandSelection: x=$x y=$y width=$width height=$height");
  
  if (picker == null) {
    err("bandSelection: picker not available");
    return;
  }
  
  int size = 4 * width * height;
  if (size > _color.length) {
    _color = new Uint8List(size);
  }
  
  gl.bindFramebuffer(RenderingContext.FRAMEBUFFER, picker.framebuffer);
  gl.readPixels(x, y, width, height, RenderingContext.RGBA, RenderingContext.UNSIGNED_BYTE, _color);
  
  if (!shift) {
    // shift is released
    _selection.clear();
  }
  
  Map<int,PickerInstance> cache = new Map<int,PickerInstance>(); 
  
  int pixels = 0;
  int bgHits = 0;
  int cacheHits = 0;
  
  for (int i = 0; i < size; i += 4, ++pixels) {
    
    if (_selection.length >= picker.numberOfInstances) {
      // optimization: selected all available objects, no need to keep searching
      break;
    }
    
    int r = _color[i];
    int g = _color[i + 1];
    int b = _color[i + 2];
    
    if (backgroundColor(r, g, b)) {
      // optimization: hit background clear color -- no object
      ++bgHits;
      continue;
    }
    
    int cacheKey = r.toInt() << 16 + g.toInt() << 8 + b.toInt();
    PickerInstance cacheEntry = cache[cacheKey];
    if (cacheEntry != null) {
      // optimization: cache
      ++cacheHits;
      continue;
    }
    
    PickerInstance pi = picker.findInstanceByColor(r, g, b);
    if (pi == null) {
      continue;
    }
    
    cache[cacheKey] = pi;    
    _selection.add(pi);
  }
  
  debug("bandSelection: $_selection (pixels total=$size scanned=$pixels, background hits=$bgHits, cache size=${cache.length} hits=$cacheHits)");  
}

library lost_context;

import 'dart:html';
import 'dart:web_gl';

void initDebugLostContext(RenderingContext gl, CanvasElement canvas) {

  print("FIXME: initDebugLostContext: trap webglcontextlost");
  print("FIXME: initDebugLostContext: trap webglcontextrestored");
  
  LoseContext ext = gl.getExtension('WEBGL_lose_context');
  if (ext == null) {
    print("WEBGL_lose_context: NOT AVAILABLE");
    return;
  }
  print("WEBGL_lose_context: available");

  DivElement control = query("#control");
  assert(control != null);
  
  InputElement loseContextButton = new InputElement();
  loseContextButton.type = 'button';
  loseContextButton.value = 'lose context';
  //loseContextButton.onClick.listen((Event e) { WebGLLoseContext.internal().loseContext(); });
  loseContextButton.onClick.listen((Event e) { print("lose context button: FIXME"); });
  control.append(loseContextButton);
  
  InputElement restoreContextButton = new InputElement();
  restoreContextButton.type = 'button';
  restoreContextButton.value = 'restore context';
  restoreContextButton.onClick.listen((Event e) { print("restore context button: FIXME"); });
  control.append(restoreContextButton);  
}



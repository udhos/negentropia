library lost_context;

import 'dart:html';

void initDebugLostContext(CanvasElement canvas) {
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
  
  print("FIXME: initDebugLostContext: trap webglcontextlost");
  print("FIXME: initDebugLostContext: trap webglcontextrestored"); 
}



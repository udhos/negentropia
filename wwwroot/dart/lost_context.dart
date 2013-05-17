library lost_context;

import 'dart:html';
import 'dart:web_gl';

import 'package:game_loop/game_loop_html.dart';

void initDebugLostContext(RenderingContext gl, CanvasElement canvas, GameLoopHtml gameLoop,
                          void initContextCall(RenderingContext gl, GameLoopHtml gameLoop)) {
  
  void onContextLost(Event e) {
    print("webgl context: lost");
    e.preventDefault();
    gameLoop.stop();
  }

  void onContextRestored(Event e) {
    print("webgl context: restored");
    initContextCall(gl, gameLoop); // recreate resources and restart gameLoop
  }

  canvas.on['webglcontextlost'].listen((Event e) => onContextLost(e));
  canvas.on['webglcontextrestored'].listen((Event e) => onContextRestored(e));

  print("initDebugLostContext: webglcontextlost trapped");
  print("initDebugLostContext: webglcontextrestored trapped");

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
  loseContextButton.onClick.listen((Event e) { ext.loseContext(); });
  control.append(loseContextButton);
  
  InputElement restoreContextButton = new InputElement();
  restoreContextButton.type = 'button';
  restoreContextButton.value = 'restore context';
  restoreContextButton.onClick.listen((Event e) { ext.restoreContext(); });
  control.append(restoreContextButton);  
}

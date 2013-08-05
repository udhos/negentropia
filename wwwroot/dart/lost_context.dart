library lost_context;

import 'dart:html';
import 'dart:web_gl';

import 'package:game_loop/game_loop_html.dart';

void initDebugLostContext(RenderingContext gl, CanvasElement canvas, GameLoopHtml gameLoop,
                          void initContextCall(RenderingContext gl, GameLoopHtml gameLoop)) {
  
  void onContextLost(Event e) {
    e.preventDefault();
    gameLoop.stop();
    print("webgl context: lost");
  }

  void onContextRestored(Event e) {
    initContextCall(gl, gameLoop); // recreate resources and restart gameLoop
    print("webgl context: restored");
  }

  //canvas.on['webglcontextlost'].listen((Event e) => onContextLost(e));
  //canvas.on['webglcontextrestored'].listen((Event e) => onContextRestored(e));
  canvas.onWebGlContextLost.listen(onContextLost);
  canvas.onWebGlContextRestored.listen(onContextRestored);
  
  /*
  LoseContext ext;
  print("initDebugLostContext: FIXME: work-around for 'dart2js -c' bug affecting Firefox 22");
  try {
    ext = gl.getExtension('WEBGL_lose_context');
  }
  catch (e) {
    print("getExtension('WEBGL_lose_context'): exception: $e");
  }
  if (ext == null) {
    print("WEBGL_lose_context: NOT AVAILABLE");
    return;
  }
  */
  LoseContext ext = gl.getExtension('WEBGL_lose_context');
  if (ext == null) {
    print("WEBGL_lose_context: NOT AVAILABLE");
    return;
  }

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

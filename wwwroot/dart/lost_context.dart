library lost_context;

import 'dart:html';
import 'dart:web_gl';

import 'package:game_loop/game_loop_html.dart';

import 'visibility.dart';

bool _lost;

bool contextIsLost() {
  assert(_lost != null);
  return _lost;
}

void initHandleLostContext(RenderingContext gl, CanvasElement canvas, GameLoopHtml gameLoop,
                           void initContextCall(RenderingContext, GameLoopHtml)) {
  
  _lost = gl.isContextLost();
  assert(_lost != null);
  
  void onContextLost(Event e) {
    e.preventDefault();
    _lost = true;
    print("webgl context: lost");
    updateGameLoop(gameLoop, contextIsLost(), pageHidden()); // gameLoop.stop();
  }

  void onContextRestored(Event e) {
    _lost = false;
    print("webgl context: restored");
    initContextCall(gl, gameLoop); // recreate resources and // gameLoop.start();
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

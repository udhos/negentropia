library lost_context;

import 'dart:html';
import 'dart:web_gl';

import 'package:game_loop/game_loop_html.dart';

import 'visibility.dart';
import 'logg.dart';

bool _lost;

bool contextIsLost() {
  assert(_lost != null);
  return _lost;
}

void initHandleLostContext(RenderingContext gl, CanvasElement
    canvas, GameLoopHtml gameLoop, void
    initContextCall(RenderingContext, GameLoopHtml)) {

  _lost = gl.isContextLost();
  assert(_lost != null);

  void onContextLost(Event e) {
    e.preventDefault();
    _lost = true;
    debug("webgl context: lost");
    updateGameLoop(gameLoop, contextIsLost(), pageHidden()); // gameLoop.stop();
  }

  void onContextRestored(Event e) {
    _lost = false;
    debug("webgl context: restored");
    initContextCall(gl, gameLoop);
    // recreate resources and // gameLoop.start();
  }

  canvas.onWebGlContextLost.listen(onContextLost);
  canvas.onWebGlContextRestored.listen(onContextRestored);

  LoseContext ext = gl.getExtension('WEBGL_lose_context');
  if (ext == null) {
    warn("WEBGL_lose_context: NOT AVAILABLE");
    return;
  }

  DivElement control = querySelector("#control");
  assert(control != null);

  InputElement loseContextButton = new InputElement();
  loseContextButton.type = 'button';
  loseContextButton.value = 'lose context';
  loseContextButton.onClick.listen((Event e) {
    ext.loseContext();
  });
  control.append(loseContextButton);

  InputElement restoreContextButton = new InputElement();
  restoreContextButton.type = 'button';
  restoreContextButton.value = 'restore context';
  restoreContextButton.onClick.listen((Event e) {
    ext.restoreContext();
  });
  control.append(restoreContextButton);
}

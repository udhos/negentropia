import 'dart:html';

import 'package:game_loop/game_loop_html.dart';

DivElement _log = new DivElement();
int _line = 0;

void log(String msg) {
  ++_line;
  msg = "$_line $msg";
  print(msg);
  DivElement m = new DivElement();
  m.text = msg;
  _log.append(m);

  List<Element> children = _log.children;
  if (children.length > 10) {
    children.removeAt(0);
  }
}

void main() {

  CanvasElement canvas = new CanvasElement();
  canvas.id = 'webgl_canvas';
  canvas.width = 100;
  canvas.height = 100;
  canvas.style.border = '2px solid black';
  document.body.append(canvas);

  document.body.append(_log);

  log(
      "canvas '${canvas.id}' created: width=${canvas.width} height=${canvas.height}");

  GameLoopHtml gameLoop = new GameLoopHtml(canvas);

  gameLoop.pointerLock.lockOnClick = false; // disable pointer lock

  gameLoop.onUpdate = ((gameLoop) {
    Mouse m = gameLoop.mouse;

    if (m.pressed(Mouse.LEFT)) {
      log("mouse left button pressed");
    }

    if (m.released(Mouse.LEFT)) {
      log("mouse left button released");
    }

    if (m.dx != 0 || m.dy != 0) {
      log("mouse moved dx=${m.dx} dy=${m.dy}");
    }

    if (m.wheelDx != 0 || m.wheelDy != 0) {
      log("mouse wheel dx=${m.wheelDx} dy=${m.wheelDy}");
    }
  });

  gameLoop.onRender = ((gameLoop) {
  });

  gameLoop.start();

  void timer(GameLoopTimer t) {
    gameLoop.addTimer(timer, 3.0);
    log("timer fired");
  }

  GameLoopTimer t = gameLoop.addTimer(timer, 3.0);
}


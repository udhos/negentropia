library visibility;

import 'dart:html';

import 'package:game_loop/game_loop_html.dart';

import 'lost_context.dart';

bool pageHidden() {
  bool hidden = document.hidden;
  if (hidden == null) { 
    print("FIXME work-around for https://code.google.com/p/dart/issues/detail?id=13411");
    hidden = false;
  }
  assert(hidden != null);
  assert(hidden is bool);
  return hidden;
}

void initPageVisibility(GameLoopHtml gameLoop) {

  document.onVisibilityChange.listen((e) {
    bool hidden = pageHidden();
    print("onVisibilityChange: visibility changed to hidden=$hidden");
    updateGameLoop(gameLoop, contextIsLost(), hidden);    
  });

}

void updateGameLoop(GameLoopHtml gameLoop, bool contextLost, bool pageHidden) {
  
  assert(contextLost != null);
  assert(pageHidden != null);
  assert(contextLost is bool);
  assert(pageHidden is bool);
  
  if (contextLost || pageHidden) {
    gameLoop.stop();
    print("updateGameLoop: game loop stopped");
  }
  else {
    gameLoop.start();
    print("updateGameLoop: game loop started");
  }
}

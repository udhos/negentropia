library visibility;

import 'dart:html';

import 'package:game_loop/game_loop_html.dart';

import 'lost_context.dart';

bool _pageHidden; 

bool pageHidden() {
  assert(pageHidden != null);
  return _pageHidden;  
}

void _updateHidden() {
  _pageHidden = document.hidden;
  if (_pageHidden == null) {
    _pageHidden = false;
  }
}

void initPageVisibility(GameLoopHtml gameLoop) {

  _updateHidden();
  
  void register(String eventName) {
    
    print("registering page visibility event: $eventName");

    document.on[eventName].listen((e) {
      _updateHidden(); 
      print("$eventName visibility changed to hidden=${pageHidden()}");
      updateGameLoop(gameLoop, contextIsLost(), pageHidden());
    });
    
  }
  
  register('visibilitychange'); 
  register('webkitvisibilitychange'); 
  register('mozvisibilitychange');   
  register('msvisibilitychange'); 
}

void updateGameLoop(GameLoopHtml gameLoop, bool contextLost, bool pageHidden) {
  if (contextLost || pageHidden) {
    gameLoop.stop();
    print("updateGameLoop: game loop stopped");
  }
  else {
    gameLoop.start();
    print("updateGameLoop: game loop started");
  }
}

library fullscreen;

import 'dart:html';

import 'logg.dart';

void toggleFullscreen(CanvasElement c) {
  if (document.fullscreenElement != null) {
    log("exiting fullscreen");
    document.exitFullscreen();
  } else {
    log("requesting fullscreen");
    c.requestFullscreen();
  }
}

library fullscreen;

import 'dart:html';

import 'logg.dart';

void trapFullscreenError() {
  document.onFullscreenError.listen((e) {
    log("fullscreenerror: $e");
  });
}

void toggleFullscreen(CanvasElement c) {
  log(
      "fullscreenSupport=${document.fullscreenEnabled} fullscreenElement=${document.fullscreenElement}"
      );

  if (document.fullscreenElement != null) {
    log("exiting fullscreen");
    document.exitFullscreen();
  } else {
    log("requesting fullscreen");
    c.requestFullscreen();
  }
}

library fullscreen;

import 'dart:html';

import 'logg.dart';

void trapFullscreenError() {
  document.onFullscreenError.listen((e) {
    err("fullscreenerror: $e");
  });
}

void toggleFullscreen(CanvasElement c) {
  debug(
      "fullscreenSupport=${document.fullscreenEnabled} fullscreenElement=${document.fullscreenElement}"
      );

  if (document.fullscreenElement != null) {
    debug("exiting fullscreen");
    document.exitFullscreen();
  } else {
    debug("requesting fullscreen");
    c.requestFullscreen();
  }
}

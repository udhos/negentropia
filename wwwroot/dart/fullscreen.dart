library fullscreen;

import 'dart:html';
import 'dart:web_gl';

import 'package:game_loop/game_loop_html.dart';

import 'logg.dart';
import 'message.dart';

double canvasAspect;
const int CANVAS_WIDTH = 780;
const int CANVAS_HEIGHT = 500;

void setViewport(CanvasElement c, RenderingContext gl, int w, int h) {

  /*
    canvas.width, canvas.height = size you requested the canvas's drawingBuffer to be
    gl.drawingBufferWidth, gl.drawingBufferHeight = size you actually got.
    canvas.clientWidth, canvas.clientHeight = size the browser is displaying your canvas.
   */
  c.width = w;
  c.height = h;
  c.style.width = "${w}px";
  c.style.height = "${h}px";

  gl.bindFramebuffer(RenderingContext.FRAMEBUFFER, null);
  // viewport for default on-screen canvas
  debug(
      "viewport: canvas=${c.width}x${c.height} drawingBuffer=${gl.drawingBufferWidth}x${gl.drawingBufferHeight}"
      );
  gl.viewport(0, 0, gl.drawingBufferWidth, gl.drawingBufferHeight);

  debug(
      "canvas size=${c.width}x${c.height} clientSize=${c.clientWidth}x${c.clientHeight}"
      );
  canvasAspect = c.clientWidth.toDouble() / c.clientHeight.toDouble();
  // save aspect for render loop setPerspectiveMatrix()
  debug("canvas aspect ratio: $canvasAspect");

  repositionMessagebox(c);
}

void trapFullscreen(CanvasElement c, RenderingContext gl, GameLoopHtml gameLoop)
    {
  document.onFullscreenError.listen((e) {
    err("fullscreenerror: $e");
  });

  document.onFullscreenChange.listen((e) {
    if (gameLoop.isFullscreen) {
      int w = window.screen.width;
      int h = window.screen.height;
      debug("fullscreen canvas: $w x $h");
      setViewport(c, gl, w, h);
      return;
    }

    setViewport(c, gl, CANVAS_WIDTH, CANVAS_HEIGHT);
  });
}

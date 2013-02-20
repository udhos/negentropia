import 'dart:html';

import 'cookies/cookies.dart';
import 'ws.dart';

void main() {
  
  CanvasElement canvas = new CanvasElement();
  assert(canvas != null);
  canvas.id = "main_canvas";
  query("#canvasbox").append(canvas);  
  print("canvas '${canvas.id}' created");
  
  var sid = Cookie.getCookie("sid");
  assert(sid != null);
  print("session id sid=${sid}");
  
  var wsUri = query("#wsUri").text;
  assert(wsUri != null);
  
  var statusElem = query("#ws_status");
  assert(statusElem != null);
  
  print("WebGL: initializing");  
  
  WebGLRenderingContext gl = canvas.getContext("experimental-webgl");
  if (gl == null) {
    print("WebGL: initialization failure: experimental-webgl");
    gl = canvas.getContext("webgl");
    if (gl == null) {
      print("WebGL: initialization failure: webgl");
      return;
    }
  }

  print("WebGL: initialized");  

  initWebSocket(wsUri, sid, 1, statusElem);
}

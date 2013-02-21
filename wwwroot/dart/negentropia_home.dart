import 'dart:html';

import 'cookies/cookies.dart';
import 'ws.dart';

WebGLRenderingContext initGL(CanvasElement canvas) {
  print("WebGL: initializing");

  WebGLRenderingContext gl;

  gl = canvas.getContext3d();
  if (gl != null) {
    print("WebGL: initialized");
    return gl;
  }

  var names = ["webgl", "experimental-webgl", "webkit-3d", "moz-webgl"];  
  
  for (var n in names) {
    gl = canvas.getContext(n);
    print("WebGL: trying context: $n");
    if (gl != null) {
      print("WebGL: initialized context: $n");
      return gl;
    }
  }

  print("WebGL: initialization failure");
  
  return null;
}

void main() {
  
  CanvasElement canvas = new CanvasElement();
  assert(canvas != null);
  canvas.id = "main_canvas";
  var canvasbox = query("#canvasbox");
  assert(canvasbox != null);  
  canvasbox.append(canvas);  
  print("canvas '${canvas.id}' created");
  
  WebGLRenderingContext gl = initGL(canvas);
  if (gl == null) {
    canvas.remove();
    var p = new ParagraphElement();
    p.text = 'WebGL is not supported by this browser.';
    canvasbox.append(p);
    var a = new AnchorElement();
    a.href = 'http://get.webgl.org/';
    a.text = 'Get more information';
    canvasbox.append(a);
    canvasbox.style.backgroundColor = 'lightblue';    
    return;
  }
  
  var sid = Cookie.getCookie("sid");
  assert(sid != null);
  print("session id sid=${sid}");
  
  var wsUri = query("#wsUri").text;
  assert(wsUri != null);
  
  var statusElem = query("#ws_status");
  assert(statusElem != null);  

  initWebSocket(wsUri, sid, 1, statusElem);
}

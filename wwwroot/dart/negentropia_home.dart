import 'dart:html';

import 'cookies/cookies.dart';
import 'ws.dart';
import 'shader.dart';
import 'buffer.dart';

int requestId;
CanvasElement canvas;
Program shaderProgram;
Model squareModel;
bool drawOnce = false;

WebGLRenderingContext initGL(CanvasElement canvas) {
  print("WebGL: initializing");

  WebGLRenderingContext gl;

  gl = canvas.getContext3d();
  if (gl != null) {
    print("WebGL: initialized");
    return gl;
  }

  print("WebGL: initialization failure");
  
  return null;
}

WebGLRenderingContext boot() {
  canvas = new CanvasElement();
  assert(canvas != null);
  canvas.id = "main_canvas";
  canvas.width = 780;
  canvas.height = 500;
  var canvasbox = query("#canvasbox");
  assert(canvasbox != null);  
  canvasbox.append(canvas);  
  print("canvas '${canvas.id}' created: width=${canvas.width} height=${canvas.height}");
  
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
    return null;
  }
  
  var sid = Cookie.getCookie("sid");
  assert(sid != null);
  print("session id sid=${sid}");
  
  var wsUri = query("#wsUri").text;
  assert(wsUri != null);
  
  var statusElem = query("#ws_status");
  assert(statusElem != null);  

  initWebSocket(wsUri, sid, 1, statusElem);
  
  return gl;
}

void initBuffers(WebGLRenderingContext gl) {
  print("initBuffers: square model: fetching");
  fetchSquare(gl, "/mesh/square.json", (Model square) {
    squareModel = square;
    print("initBuffers: square model: done");
  });
}

void initContext(WebGLRenderingContext gl) {
  // load shaders
  shaderProgram = new Program(gl, "/shader/min_vs.txt", "/shader/min_fs.txt");
  assert(shaderProgram != null);
  
  // init buffers
  initBuffers(gl);
    
  gl.clearColor(0.5, 0.5, 0.5, 1.0);            // clear color
  gl.enable(WebGLRenderingContext.DEPTH_TEST);  // enable depth testing
  gl.depthFunc(WebGLRenderingContext.LESS);     // gl.LESS is default depth test

  // enable backface culling
  gl.frontFace(WebGLRenderingContext.CCW);
  gl.cullFace(WebGLRenderingContext.BACK);
  gl.enable(WebGLRenderingContext.CULL_FACE);    
}

void animate() {
}

void render(WebGLRenderingContext gl) {
  gl.viewport(0, 0, canvas.width, canvas.height); // define viewport size
  gl.depthRange(0.0, 1.0); // default
  
  // http://www.opengl.org/sdk/docs/man/xhtml/glClear.xml
  gl.clear(WebGLRenderingContext.COLOR_BUFFER_BIT | WebGLRenderingContext.DEPTH_BUFFER_BIT);    // clear color buffer and depth buffer
  
  // set perspective matrix
  // field of view y: 45 degrees
  // width to height ratio
  // view from 1.0 to 1000.0 distance units
  //
  // tan(45/2) = (h/2) / near
  // h = 2 * tan(45/2) * near
  // h = 2 * 0.414 * 1.0
  // h = 0.828
  //
  //mat4.perspective(neg.fieldOfViewY, neg.canvas.width / neg.canvas.height, 1.0, 1000.0, neg.pMatrix);
  
  drawSquare(gl);
}

void drawSquare(WebGLRenderingContext gl) {

  if (!shaderProgram.ready) {
    // shader program is not loaded yet
    return;
  }

  if (squareModel == null) {
    // square buffers are not loaded yet
    return;
  }

  int aVertexPosition = shaderProgram.aVertexPosition;

  gl.bindBuffer(WebGLRenderingContext.ARRAY_BUFFER, squareModel.vertexPositionBuffer);
  gl.vertexAttribPointer(aVertexPosition, squareModel.vertexPositionBufferItemSize, WebGLRenderingContext.FLOAT, false, 0, 0);
  gl.enableVertexAttribArray(aVertexPosition);
  
  gl.bindBuffer(WebGLRenderingContext.ELEMENT_ARRAY_BUFFER, squareModel.vertexIndexBuffer);

  gl.drawElements(WebGLRenderingContext.TRIANGLES, squareModel.vertexIndexLength, WebGLRenderingContext.UNSIGNED_SHORT, 0 * squareModel.vertexIndexBufferItemSize);

  // clean up
  gl.bindBuffer(WebGLRenderingContext.ARRAY_BUFFER, null);
  gl.bindBuffer(WebGLRenderingContext.ELEMENT_ARRAY_BUFFER, null);       
}

void loop(WebGLRenderingContext gl) {
  // FIXME update framerate statistics

  if (drawOnce) {
    print("loop: drawOnce ON: will render only one frame");
  }
  else {
    requestId = window.requestAnimationFrame((num time) { loop(gl); });
    if (requestId == 0) {
      print("loop: could not obtain requestId from requestAnimationFrame");
    }
  }
  
  animate();    // update state
  render(gl);   // draw
}

void main() {
  WebGLRenderingContext gl = boot();
  
  if (gl == null) {
    print("WebGL: not available");
    return;
  }
  
  initContext(gl);
  
  loop(gl);
}

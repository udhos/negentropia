import 'dart:html';
import 'dart:async';
import 'dart:web_gl';

import 'package:stats/stats.dart';

import 'cookies/cookies.dart';
import 'ws.dart';
import 'shader.dart';
import 'buffer.dart';
import 'lost_context.dart';

int requestId;
CanvasElement canvas;
num canvasAspect;
ShaderProgram shaderProgram;
bool debugLostContext = true;
List<ShaderProgram> programList = new List<ShaderProgram>();
Map<String,Shader> shaderCache = new Map<String,Shader>();

// >0  : render at max rate then stop
// <=0 : periodic rendering
//int fullRateFrames = 30000;
int fullRateFrames = 0; // periodic rendering

RenderingContext initGL(CanvasElement canvas) {
  print("WebGL: initializing");

  RenderingContext gl;

  gl = canvas.getContext3d();
  if (gl != null) {
    print("WebGL: initialized");
    return gl;
  }

  print("WebGL: initialization failure");
  
  return null;
}

/*
	https://github.com/jwill/stats.dart
	https://github.com/toji/dart-render-stats
	https://github.com/Dartist/stats.dart
*/

Stats stats = null;

void initStats() {
  DivElement div = query("#framerate");
  assert(div != null);
  stats = new Stats();
  div.children.add(stats.container);
}

RenderingContext boot() {
  canvas = new CanvasElement();
  assert(canvas != null);
  canvas.id = "main_canvas";
  canvas.width = 780;
  canvas.height = 500;
  var canvasbox = query("#canvasbox");
  assert(canvasbox != null);  
  canvasbox.append(canvas);  
  print("canvas '${canvas.id}' created: width=${canvas.width} height=${canvas.height}");
  
  RenderingContext gl = initGL(canvas);
  if (gl == null) {
    canvas.remove();
    var p = new ParagraphElement();
    p.text = 'WebGL is currently not available on this system.';
    canvasbox.append(p);
    var a = new AnchorElement();
    a.href = 'http://get.webgl.org/';
    a.text = 'Get more information';
    canvasbox.append(a);
    canvasbox.style.backgroundColor = 'lightblue';    
    return null;
  }
  
  if (debugLostContext) {
    //initDebugLostContext(canvas, cfg);
    initDebugLostContext(canvas);
  }
  
  var sid = Cookie.getCookie("sid");
  assert(sid != null);
  print("session id sid=${sid}");
  
  var wsUri = query("#wsUri").text;
  assert(wsUri != null);
  
  var statusElem = query("#ws_status");
  assert(statusElem != null);  

  initWebSocket(wsUri, sid, 1, statusElem);
  
  initStats();
  
  return gl;
}

void initContext(RenderingContext gl) {

  programList = new List<ShaderProgram>();           // drop existing programs 
  shaderCache = new Map<String,Shader>(); // drop existing compile shader cache

  ShaderProgram squareProgram = new ShaderProgram(gl, shaderCache, "/shader/min_vs.txt", "/shader/min_fs.txt");
  programList.add(squareProgram);
  Model squareModel = new Model.fromURL(gl, squareProgram, "/mesh/square.json");
  squareProgram.addModel(squareModel);
  Instance squareInstance = new Instance(squareModel);
  squareModel.addInstance(squareInstance);

  // execute after 2 secs, giving time to first program populate shadeCache
  new Timer(new Duration(seconds:2), () {
    ShaderProgram squareProgram2 = new ShaderProgram(gl, shaderCache, "/shader/min_vs.txt", "/shader/min2_fs.txt");
    programList.add(squareProgram2);
    Model squareModel2 = new Model.fromURL(gl, squareProgram2, "/mesh/square2.json");
    squareProgram2.addModel(squareModel2);
    Instance squareInstance2 = new Instance(squareModel2);
    squareModel2.addInstance(squareInstance2);
  });
  
  gl.clearColor(0.5, 0.5, 0.5, 1.0);            // clear color
  gl.enable(RenderingContext.DEPTH_TEST);  // enable depth testing
  gl.depthFunc(RenderingContext.LESS);     // gl.LESS is default depth test
  gl.depthRange(0.0, 1.0);                      // default
  
  // define viewport size
  gl.viewport(0, 0, canvas.width, canvas.height);
  canvasAspect = canvas.width / canvas.height; // save aspect for render loop mat4.perspective

  // enable backface culling
  gl.frontFace(RenderingContext.CCW);
  gl.cullFace(RenderingContext.BACK);
  gl.enable(RenderingContext.CULL_FACE);
  
  loop(gl); // render loop
}

void animate() {
    // TODO: FIXME: WRITEME: update state
}

void render(RenderingContext gl) {
  
  // http://www.opengl.org/sdk/docs/man/xhtml/glClear.xml
  gl.clear(RenderingContext.COLOR_BUFFER_BIT | RenderingContext.DEPTH_BUFFER_BIT);    // clear color buffer and depth buffer
  
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
  // aspect = canvas.width / canvas.height
  //mat4.perspective(neg.fieldOfViewY, canvasAspect, 1.0, 1000.0, neg.pMatrix);
    
  programList.forEach((ShaderProgram p) => p.drawModels());
}

void loop(RenderingContext gl) {
  
  if (fullRateFrames > 0) {
    
    print("loop: firing $fullRateFrames frames at full rate");
    
    var before = new DateTime.now();
        
    for (int i = 0; i < fullRateFrames; ++i) {
      stats.begin();      
      animate();    // update state
      render(gl);   // draw
      stats.end();      
    };

    var after = new DateTime.now();
    var duration = after.difference(before);
    var rate = fullRateFrames / duration.inSeconds;
    
    print("loop: duration = $duration framerate = $rate fps");

    return;
  }

  stats.begin();      

  requestId = window.requestAnimationFrame((num time) { loop(gl); });
  if (requestId == 0) {
    print("loop: could not obtain requestId from requestAnimationFrame");
  }

  animate();    // update state
  render(gl);   // draw
  
  stats.end();
}

void main() {
  RenderingContext gl = boot();
  
  if (gl == null) {
    print("WebGL: not available");
    return;
  }
  
  initContext(gl); // calls loop(gl)
}

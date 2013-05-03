import 'dart:html';
import 'dart:async';
import 'dart:web_gl';
import 'dart:math' as math;

import 'package:stats/stats.dart';
import 'package:vector_math/vector_math.dart';
import 'package:game_loop/game_loop_html.dart';

import 'cookies/cookies.dart';
import 'ws.dart';
import 'shader.dart';
import 'skybox.dart';
import 'buffer.dart';
import 'lost_context.dart';
import 'camera.dart';

int requestId;
CanvasElement canvas;
num canvasAspect;
ShaderProgram shaderProgram;
bool debugLostContext = true;
List<ShaderProgram> programList = new List<ShaderProgram>();
Map<String,Shader> shaderCache = new Map<String,Shader>();
mat4 pMatrix = new mat4.zero();
num fieldOfViewYRadians = 45 * math.PI / 180;
Camera cam = new Camera(new vec3(0.0,0.0,10.0), new vec3(0.0,0.0,-1.0), new vec3(0.0,1.0,0.0));
bool backfaceCulling = false;

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

void initSquares(RenderingContext gl) {
  ShaderProgram squareProgram = new ShaderProgram(gl);
  programList.add(squareProgram);
  squareProgram.fetch(shaderCache, "/shader/clip_vs.txt", "/shader/clip_fs.txt");
  Model squareModel = new Model.fromURL(gl, squareProgram, "/mesh/square.json");
  squareProgram.addModel(squareModel);
  Instance squareInstance = new Instance(squareModel, new vec3(0.0, 0.0, 0.0), 1.0);
  squareModel.addInstance(squareInstance);

  ShaderProgram squareProgram2 = new ShaderProgram(gl);
  programList.add(squareProgram2);
  // execute after 2 secs, giving time to first program populate shadeCache
  new Timer(new Duration(seconds:2), () {
    squareProgram2.fetch(shaderCache, "/shader/clip_vs.txt", "/shader/clip2_fs.txt");
  });
  Model squareModel2 = new Model.fromURL(gl, squareProgram2, "/mesh/square2.json");
  squareProgram2.addModel(squareModel2);
  Instance squareInstance2 = new Instance(squareModel2, new vec3(0.0, 0.0, 0.0), 1.0);
  squareModel2.addInstance(squareInstance2);
  
  ShaderProgram squareProgram3 = new ShaderProgram(gl);
  programList.add(squareProgram3);
  squareProgram3.fetch(shaderCache, "/shader/clip_vs.txt", "/shader/clip3_fs.txt");
  Model squareModel3 = new Model.fromURL(gl, squareProgram3, "/mesh/square3.json");
  squareProgram3.addModel(squareModel3);
  Instance squareInstance3 = new Instance(squareModel3, new vec3(0.0, 0.0, 0.0), 1.0);
  squareModel3.addInstance(squareInstance3);  
}

void initSkybox(RenderingContext gl) {
  SkyboxProgram skyboxProgram = new SkyboxProgram(gl);
  programList.add(skyboxProgram);
  skyboxProgram.fetch(shaderCache, "/shader/skybox_vs.txt", "/shader/skybox_fs.txt");
  SkyboxModel skyboxModel = new SkyboxModel.fromURL(gl, skyboxProgram, "/mesh/cube.json", true, 0);
  skyboxModel.addCubemapFace(RenderingContext.TEXTURE_CUBE_MAP_POSITIVE_X, '/texture/space_rt.jpg');
  skyboxModel.addCubemapFace(RenderingContext.TEXTURE_CUBE_MAP_NEGATIVE_X, '/texture/space_lf.jpg');
  skyboxModel.addCubemapFace(RenderingContext.TEXTURE_CUBE_MAP_POSITIVE_Y, '/texture/space_up.jpg');
  skyboxModel.addCubemapFace(RenderingContext.TEXTURE_CUBE_MAP_NEGATIVE_Y, '/texture/space_dn.jpg');
  skyboxModel.addCubemapFace(RenderingContext.TEXTURE_CUBE_MAP_POSITIVE_Z, '/texture/space_fr.jpg');
  skyboxModel.addCubemapFace(RenderingContext.TEXTURE_CUBE_MAP_NEGATIVE_Z, '/texture/space_bk.jpg');  
  skyboxProgram.addModel(skyboxModel);
  Instance skyboxInstance = new Instance(skyboxModel, new vec3(0.0, 0.0, 0.0), 1.0);
  skyboxModel.addInstance(skyboxInstance);
}

void initContext(RenderingContext gl, GameLoopHtml gameLoop) {

  programList = new List<ShaderProgram>(); // drop existing programs 
  shaderCache = new Map<String,Shader>();  // drop existing compile shader cache

  initSquares(gl);
  initSkybox(gl);
  
  gl.clearColor(0.5, 0.5, 0.5, 1.0);            // clear color
  gl.enable(RenderingContext.DEPTH_TEST);  // enable depth testing
  gl.depthFunc(RenderingContext.LESS);     // gl.LESS is default depth test
  gl.depthRange(0.0, 1.0);                      // default
  
  // define viewport size
  gl.viewport(0, 0, canvas.width, canvas.height);
  canvasAspect = canvas.width / canvas.height; // save aspect for render loop mat4.perspective

  // enable backface culling
  if (backfaceCulling) {
    gl.frontFace(RenderingContext.CCW);
    gl.cullFace(RenderingContext.BACK);
    gl.enable(RenderingContext.CULL_FACE);
  }
  
  if (fullRateFrames > 0) {
    print("firing $fullRateFrames frames at full rate");
    
    var before = new DateTime.now();
        
    for (int i = 0; i < fullRateFrames; ++i) {
      stats.begin();      
      draw(gl, gameLoop);
      stats.end();      
    };

    var after = new DateTime.now();
    var duration = after.difference(before);
    var rate = fullRateFrames / duration.inSeconds;
    
    print("duration = $duration framerate = $rate fps");
  }

  //render(gl); // render loop
  gameLoop.start();
}

void draw(RenderingContext gl, GameLoopHtml gameLoop) {
  
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
  setPerspectiveMatrix(pMatrix, fieldOfViewYRadians, canvasAspect, 1.0, 1000.0);

  cam.render(gameLoop);
  
  programList.forEach((ShaderProgram p) => p.drawModels(gameLoop, cam, pMatrix));
}

/*
void render(RenderingContext gl) {
  
  if (fullRateFrames > 0) {
    
    print("render: firing $fullRateFrames frames at full rate");
    
    var before = new DateTime.now();
        
    for (int i = 0; i < fullRateFrames; ++i) {
      stats.begin();      
      update();    // update state
      draw(gl);   // draw
      stats.end();      
    };

    var after = new DateTime.now();
    var duration = after.difference(before);
    var rate = fullRateFrames / duration.inSeconds;
    
    print("render: duration = $duration framerate = $rate fps");

    return;
  }

  stats.begin();      

  requestId = window.requestAnimationFrame((num time) { render(gl); });
  if (requestId == 0) {
    print("render: could not obtain requestId from requestAnimationFrame");
  }

  update();    // update state
  draw(gl);   // draw
  
  stats.end();
}
*/

void render(RenderingContext gl, GameLoopHtml gameLoop) {
  stats.begin();
  draw(gl, gameLoop);
  stats.end();
}

void update(GameLoopHtml gameLoop) {
  //print('${gameLoop.frame}: ${gameLoop.frameTime} [dt = ${gameLoop.dt}].');

  cam.update(gameLoop);
    
  programList.forEach((ShaderProgram p) => p.update(gameLoop));
}

void main() {
  RenderingContext gl = boot();
  
  if (gl == null) {
    print("WebGL: not available");
    return;
  }

  GameLoopHtml gameLoop = new GameLoopHtml(canvas);
  gameLoop.onUpdate = ((GameLoopHtml gameLoop) { 
    update(gameLoop);
  });
  gameLoop.onRender = ((GameLoopHtml gameLoop) {
    render(gl, gameLoop);
  });

  initContext(gl, gameLoop);
}


import 'dart:html';
import 'dart:async';
import 'dart:web_gl';
import 'dart:math' as math;
import 'dart:typed_data';

import 'package:stats/stats.dart';
import 'package:vector_math/vector_math.dart';
import 'package:game_loop/game_loop_html.dart';

import 'cookies/cookies.dart';
import 'ws.dart';
import 'shader.dart';
import 'skybox.dart';
import 'lost_context.dart';
import 'camera.dart';
import 'texture.dart';
import 'obj.dart';
import 'asset.dart';

CanvasElement canvas;
double canvasAspect;
ShaderProgram shaderProgram;
bool debugLostContext = true;
List<ShaderProgram> programList;
Map<String,Shader> shaderCache;
Map<String,Texture> textureTable;
Matrix4 pMatrix = new Matrix4.zero();
double fieldOfViewYRadians = 45 * math.PI / 180;
Camera cam = new Camera(new Vector3(0.0,0.0,15.0), new Vector3(0.0,0.0,-1.0), new Vector3(0.0,1.0,0.0));
bool backfaceCulling = false;
bool showPicking = false;
Asset asset = new Asset("/");
PickerShader picker;

// >0  : render at max rate then stop
// <=0 : periodic rendering
//int fullRateFrames = 30000;
int fullRateFrames = 0; // periodic rendering

RenderingContext initGL(CanvasElement canvas) {
  print("WebGL: initializing");

  RenderingContext gl;

  // FIXME: ERASEME: preserveDrawingBuffer: true
  gl = canvas.getContext3d(preserveDrawingBuffer: true);
  //gl = canvas.getContext3d();
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

void initShowPicking() {
  DivElement control = query("#control");
  assert(control != null);

  InputElement showPickingCheck = new InputElement();
  showPickingCheck.type = 'checkbox';
  showPickingCheck.id = 'show_picking';
  showPickingCheck.checked = showPicking;
  showPickingCheck.onClick.listen((Event e) { showPicking = showPickingCheck.checked; });
  control.append(showPickingCheck);

  LabelElement label = new LabelElement();
  label.htmlFor = showPickingCheck.id;
  label.appendText('show picking');
  control.append(label);
}

RenderingContext boot() {
  canvas = new CanvasElement();
  assert(canvas != null);
  canvas.id = "main_canvas";
  canvas.width = 780;
  canvas.height = 500;
  DivElement canvasbox = query("#canvasbox");
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
  
  initShowPicking();
    
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
  squareProgram.fetch(shaderCache, "${asset.shader}/clip_vs.txt", "${asset.shader}/clip_fs.txt");
  Model squareModel = new Model.fromJson(gl, "${asset.mesh}/square.json");
  squareProgram.addModel(squareModel);
  Instance squareInstance = new Instance(squareModel, new Vector3(0.0, 0.0, 0.0), 1.0);
  squareModel.addInstance(squareInstance);

  ShaderProgram squareProgram2 = new ShaderProgram(gl);
  programList.add(squareProgram2);
  // execute after 2 secs, giving time to first program populate shadeCache
  new Timer(new Duration(seconds:2), () {
    squareProgram2.fetch(shaderCache, "${asset.shader}/clip_vs.txt", "${asset.shader}/clip2_fs.txt");
  });
  Model squareModel2 = new Model.fromJson(gl, "${asset.mesh}/square2.json");
  squareProgram2.addModel(squareModel2);
  Instance squareInstance2 = new Instance(squareModel2, new Vector3(0.0, 0.0, 0.0), 1.0);
  squareModel2.addInstance(squareInstance2);
  
  ShaderProgram squareProgram3 = new ShaderProgram(gl);
  programList.add(squareProgram3);
  squareProgram3.fetch(shaderCache, "${asset.shader}/clip_vs.txt", "${asset.shader}/clip3_fs.txt");
  Model squareModel3 = new Model.fromJson(gl, "${asset.mesh}/square3.json");
  squareProgram3.addModel(squareModel3);
  Instance squareInstance3 = new Instance(squareModel3, new Vector3(0.0, 0.0, 0.0), 1.0);
  squareModel3.addInstance(squareInstance3);  
}

void initSkybox(RenderingContext gl) {
  SkyboxProgram skyboxProgram = new SkyboxProgram(gl);
  programList.add(skyboxProgram);
  skyboxProgram.fetch(shaderCache, "${asset.shader}/skybox_vs.txt", "${asset.shader}/skybox_fs.txt");
  SkyboxModel skyboxModel = new SkyboxModel.fromJson(gl, "/mesh/cube.json", true, 0);
  skyboxModel.addCubemapFace(gl, RenderingContext.TEXTURE_CUBE_MAP_POSITIVE_X, '/texture/space_rt.jpg');
  skyboxModel.addCubemapFace(gl, RenderingContext.TEXTURE_CUBE_MAP_NEGATIVE_X, '/texture/space_lf.jpg');
  skyboxModel.addCubemapFace(gl, RenderingContext.TEXTURE_CUBE_MAP_POSITIVE_Y, '/texture/space_up.jpg');
  skyboxModel.addCubemapFace(gl, RenderingContext.TEXTURE_CUBE_MAP_NEGATIVE_Y, '/texture/space_dn.jpg');
  skyboxModel.addCubemapFace(gl, RenderingContext.TEXTURE_CUBE_MAP_POSITIVE_Z, '/texture/space_fr.jpg');
  skyboxModel.addCubemapFace(gl, RenderingContext.TEXTURE_CUBE_MAP_NEGATIVE_Z, '/texture/space_bk.jpg');  
  skyboxProgram.addModel(skyboxModel);
  SkyboxInstance skyboxInstance = new SkyboxInstance(skyboxModel, new Vector3(0.0, 0.0, 0.0), 1.0);
  skyboxModel.addInstance(skyboxInstance);
}

void initAirship(RenderingContext gl) {
  ShaderProgram prog = new ShaderProgram(gl);
  programList.add(prog);
  prog.fetch(shaderCache, "${asset.shader}/simple_vs.txt", "${asset.shader}/simple_fs.txt");
  Model airshipModel = new Model.fromOBJ(gl, "${asset.obj}/airship.obj");
  prog.addModel(airshipModel);
  Instance airshipInstance = new Instance(airshipModel, new Vector3(-8.0, 0.0, 0.0), 1.0, generatePickColor());
  airshipModel.addInstance(airshipInstance);  
}

void initAirshipTex(RenderingContext gl) {
  TexShaderProgram prog = new TexShaderProgram(gl);
  programList.add(prog);
  prog.fetch(shaderCache, "${asset.shader}/simpleTex_vs.txt", "${asset.shader}/simpleTex_fs.txt");
  
  String objURL = "${asset.obj}/airship.obj"; 

  TexModel airshipModel = new TexModel.fromOBJ(gl, objURL, textureTable, asset);
  prog.addModel(airshipModel);
  TexInstance airshipInstance = new TexInstance(airshipModel, new Vector3(0.0, 0.0, 0.0), 1.0, generatePickColor());
  airshipModel.addInstance(airshipInstance);

  TexModel airshipModel2 = new TexModel.fromOBJ(gl, objURL, textureTable, asset);
  prog.addModel(airshipModel2);
  TexInstance airshipInstance2 = new TexInstance(airshipModel2, new Vector3(8.0, 0.0, 0.0), 1.0, generatePickColor());
  airshipModel2.addInstance(airshipInstance2);
  
  String colonyShipURL = "${asset.obj}/Colony Ship Ogame Fleet.obj";  
  TexModel colonyShipModel = new TexModel.fromOBJ(gl, colonyShipURL, textureTable, asset);
  prog.addModel(colonyShipModel);
  TexInstance colonyShipInstance = new TexInstance(colonyShipModel, new Vector3(0.0, -5.0, -50.0), 1.0, generatePickColor());
  colonyShipModel.addInstance(colonyShipInstance);
    
  String coneURL = "${asset.obj}/cone.obj";  
  TexModel coneModel = new TexModel.fromOBJ(gl, coneURL, textureTable, asset);
  prog.addModel(coneModel);
  TexInstance coneInstance = new TexInstance(coneModel, new Vector3(0.0, 2.0, -10.0), 1.0, generatePickColor());
  coneModel.addInstance(coneInstance);
}

void initShips(RenderingContext gl) {
  initAirship(gl);
  initAirshipTex(gl);
}

void initPicker(RenderingContext gl) {
  picker = new PickerShader(gl, programList, canvas.width, canvas.height);
  programList.add(picker);
  picker.fetch(shaderCache, "${asset.shader}/picker_vs.txt", "${asset.shader}/picker_fs.txt");  
}

void initContext(RenderingContext gl, GameLoopHtml gameLoop) {
    
  programList = new List<ShaderProgram>();  // drop existing programs 
  shaderCache = new Map<String,Shader>();   // drop existing compile shader cache
  textureTable = new Map<String,Texture>(); // drop existing texture table

  programList.forEach((ShaderProgram p) => p.initContext(gl, textureTable));
  
  initSquares(gl);
  initShips(gl);
  initSkybox(gl);
  initPicker(gl);
  
  gl.clearColor(0.5, 0.5, 0.5, 1.0);       // clear color
  gl.enable(RenderingContext.DEPTH_TEST);  // enable depth testing
  gl.depthFunc(RenderingContext.LESS);     // gl.LESS is default depth test
  gl.depthRange(0.0, 1.0);                 // default
  
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

  // clear offscreen framebuffer
  gl.bindFramebuffer(RenderingContext.FRAMEBUFFER, picker.framebuffer);
  gl.clear(RenderingContext.COLOR_BUFFER_BIT | RenderingContext.DEPTH_BUFFER_BIT);

  // clear canvas framebuffer
  gl.bindFramebuffer(RenderingContext.FRAMEBUFFER, null);
  gl.clear(RenderingContext.COLOR_BUFFER_BIT | RenderingContext.DEPTH_BUFFER_BIT);  
  
  if (showPicking) {
    // draw picking on both framebuffers
    
    picker.offscreen = true; // offscreen framebuffer
    picker.drawModels(gameLoop, cam, pMatrix);
    
    picker.offscreen = false; // canvas framebuffer
    picker.drawModels(gameLoop, cam, pMatrix);
  }
  else {
    // draw picking on offscreen framebuffer
    // and actual object colors on canvas framebuffer
    picker.offscreen = true; 
    programList.forEach((p) => p.drawModels(gameLoop, cam, pMatrix));     
  }
}

void render(RenderingContext gl, GameLoopHtml gameLoop) {
  stats.begin();
  draw(gl, gameLoop);
  stats.end();
}

void readColor(String label, RenderingContext gl, int x, int y, Framebuffer framebuffer) {
  Uint8List color = new Uint8List(4);
  gl.bindFramebuffer(RenderingContext.FRAMEBUFFER, framebuffer);
  gl.readPixels(x, y, 1, 1, RenderingContext.RGBA, RenderingContext.UNSIGNED_BYTE, color);
  print("$label: readPixels: x=$x y=$y color=$color");     
}

void update(RenderingContext gl, GameLoopHtml gameLoop) {
  //print('${gameLoop.frame}: ${gameLoop.frameTime} [dt = ${gameLoop.dt}].');

  Mouse m = gameLoop.mouse;
  if (m.pressed(Mouse.LEFT)) {

    print("Mouse.LEFT pressed: withinCanvas=${m.withinCanvas}");
    
    int y = canvas.height - m.y;

    readColor("canvas-framebuffer", gl, m.x, y, null);
    readColor("offscreen-framebuffer", gl, m.x, y, picker.framebuffer);    
  }  
  
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
  
  gameLoop.pointerLock.lockOnClick = false; // disable pointer lock

  print("gameLoop lockOnClick = ${gameLoop.pointerLock.lockOnClick}");

  print("gameLoop updateStep = ${gameLoop.updateTimeStep} seconds");

  if (debugLostContext) {
    initDebugLostContext(gl, canvas, gameLoop, initContext);
  }
  
  gameLoop.onUpdate = ((GameLoopHtml gameLoop) { 
    update(gl, gameLoop);
  });
  gameLoop.onRender = ((GameLoopHtml gameLoop) {
    render(gl, gameLoop);
  });

  initContext(gl, gameLoop);
}

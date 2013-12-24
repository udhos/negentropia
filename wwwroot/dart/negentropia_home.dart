
import 'dart:html';
import 'dart:async';
import 'dart:web_gl';
import 'dart:math' as math;
import 'dart:typed_data';
import 'dart:convert';
import 'dart:collection';

import 'package:stats/stats.dart';
import 'package:vector_math/vector_math.dart';
import 'package:game_loop/game_loop_html.dart';

import 'cookies/cookies.dart';
import 'ws.dart';
import 'shader.dart';
import 'skybox.dart';
//import 'solid.dart';
import 'lost_context.dart';
import 'camera.dart';
import 'asset.dart';
import 'anisotropic.dart';
import 'visibility.dart';
import 'logg.dart';
import 'interpolate.dart';
import 'vec.dart';
import 'selection.dart';

CanvasElement canvas;
DivElement messagebox;
double canvasAspect;
ShaderProgram shaderProgram;
bool debugLostContext = true;
List<ShaderProgram> programList;
Map<String,Shader> shaderCache;
Map<String,Texture> textureTable;
Matrix4 pMatrix = new Matrix4.zero();
double fieldOfViewYRadians = 45 * math.PI / 180;
Camera cam = new Camera(new Vector3(0.0, 0.0, 15.0));
bool backfaceCulling = false;
bool showPicking = false;
Asset asset = new Asset("/");
SkyboxProgram skybox;
PickerShader picker;
SolidShader solidShader;
double planeNear   = 1.0;
double planeFar    = 2000.0;
double skyboxScale = 1000.0;
int mouseDragBeginX = null;
int mouseDragBeginY = null;
int mouseDragCurrX = null;
int mouseDragCurrY = null;

// >0  : render at max rate then stop
// <=0 : periodic rendering
//int fullRateFrames = 30000;
int fullRateFrames = 0; // periodic rendering

RenderingContext initGL(CanvasElement canvas) {
  //print("WebGL: initializing");

  RenderingContext gl;

  //print("initGL: FIXME: ERASEME: preserveDrawingBuffer: true");
  //gl = canvas.getContext3d(preserveDrawingBuffer: true);
  gl = canvas.getContext3d();
  if (gl != null) {
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
  DivElement div = querySelector("#framerate");
  assert(div != null);
  stats = new Stats();
  div.children.add(stats.container);
}

void initShowPicking() {
  DivElement control = querySelector("#control");
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

void loadDemo(RenderingContext gl) {
  demoInitSquares(gl);
  demoInitShips(gl);
  demoInitSkybox(gl);
  demoInitPicker(gl);
}

TexShaderProgram findTexShader(String programName) {
  TexShaderProgram prog;
  try {
    prog = programList.firstWhere((p) { return p.programName == programName; });
  }
  on StateError {
    // not found
  }
  return prog;
}

int maxList = 10;
ListQueue<String> msgList = new ListQueue<String>(maxList);

void messageUser(String m) {
  
  msgList.add(m);
  
  while (msgList.length > maxList) {
    msgList.removeFirst();
  }
  
  messagebox.children.clear();
  
  msgList.forEach((m) {
    DivElement d = new DivElement();
    d.text = m;
    messagebox.children.add(d);
  });
}

void dispatcher(RenderingContext gl, int code, String data, Map<String,String> tab) {
  
  switch (code) {
    case CM_CODE_INFO:
      
      print("dispatcher: server sent info: $data");
      
      if (data.startsWith("welcome")) {
        // test echo loop thru server
        /*
        var m = new Map();
        m["Code"] = CM_CODE_ECHO;
        m["Data"] = "hi there";
        */
        Map m = {
          "Code": CM_CODE_ECHO,
          "Data": "hi there"
        };
        wsSend(JSON.encode(m));        
      }
      break;
      
    case CM_CODE_ZONE:

      if (tab != null) {
        
        String culling = tab['backfaceCulling'];
        if (culling != null) {
          backfaceCulling = culling.toLowerCase().startsWith("t");
          //print("dispatcher: backfaceCulling=$backfaceCulling");
          updateCulling(gl);
        }

        String camCoord = tab['cameraCoord'];
        /*
        if (camCoord != null) {
          List<double> c = new List<double>();
          camCoord.split(',').forEach((i) => c.add(double.parse(i)));
          Vector3 coord = new Vector3.array(c);
          debug("cameraCoord: $camCoord => $coord");
          if (c.length == 3) {
            cam.moveTo(coord);
          }
          else {
            err("cameraCoord: $camCoord => $coord: bad length=${c.length}");
          }
        }
        */
        Vector3 coord = parseVector3(camCoord);
        if (coord != null) {
          cam.moveTo(coord);
        }
        else {
          err("cameraCoord: parsing failure: camCoord=$camCoord");
        }
        
      }
      
      resetZone(gl);
            
      if (data == "demo") {
        loadDemo(gl);
      }
      break;
      
    case CM_CODE_SKYBOX:

      String skyboxURL = tab['skyboxURL'];

      void handleResponse(String response) {
        Map<String,String> skybox = JSON.decode(response);
        addSkybox(gl, skybox);
      }
      
      HttpRequest.getString(skyboxURL)
        .then(handleResponse)
        .catchError((e) { print("dispatcher: failure fetching skyboxURL=$skyboxURL: $e"); });

      break;
      
    case CM_CODE_PROGRAM:
      
      String programName = tab['programName'];
      TexShaderProgram prog = findTexShader(programName);
      if (prog != null) {
        print("dispatcher: failure redefining program programName=$programName");
      }
      else {
        prog = new TexShaderProgram(gl, programName);
        programList.add(prog);
        prog.fetch(shaderCache, tab['vertexShader'], tab['fragmentShader']);        
      }
      
      break;

    case CM_CODE_INSTANCE:
      
      String objURL      = tab['objURL'];
      String programName = tab['programName'];
      String front       = tab['directionFront'];
      String up          = tab['directionUp'];
      String coord       = tab['coord'];
      String scale       = tab['scale'];
      
      debug("dispatcher: instance: obj='$objURL'");

      Vector3 f = parseVector3(front);
      if (f == null) {
        err("obj=$objURL: parsing failure: front=$front");        
        return;
      }

      Vector3 u = parseVector3(up);
      if (u == null) {
        err("obj=$objURL: parsing failure: up=$up");        
        return;
      }
      
      /*
      List<String> coordList = coord.split(',');
      Vector3 vec3 = new Vector3(double.parse(coordList[0]), double.parse(coordList[1]), double.parse(coordList[2]));
      */
      Vector3 c = parseVector3(coord);
      if (c == null) {
        err("obj=$objURL: parsing failure: coord=$coord");        
        return;
      }
      
      double sc = double.parse(scale);
      
      TexShaderProgram prog = findTexShader(programName);
      if (prog == null) {
        print("dispatcher: instance: could not find programName=$programName");
        return;
      }

      TexModel model = prog.findModel(objURL);
      if (model == null) {
        model = new TexModel.fromOBJ(gl, objURL, f, u, textureTable, asset);
        prog.addModel(model);
      }
      
      TexInstance instance = new TexInstance(model, c, sc, generatePickColor());
      model.addInstance(instance);
      
      fixme("dispatcher: update picker incrementally instead of fully rebuilding it for each instance");
      addPicker(gl);
      
      fixme("dispatcher: update axis shader incrementally instead of fully rebuilding it for each instance");
      addSolidShader(gl);     
      
      break;

    case CM_CODE_MESSAGE:
      messageUser(data);
      break;      

    default:
      print("dispatcher: unknown code=$code");
  }  
}

DivElement createMessagebox(String id, CanvasElement c) {
  
  DivElement mbox = new DivElement();
  mbox.id = id;
  
  int left = 10 + c.offsetLeft;
  int top  = 28 + c.offsetTop;

  mbox.style.border = '2px solid #FFF';
  mbox.style.zIndex = "1";
  mbox.style.position = "absolute";
  mbox.style.width = "300px";
  mbox.style.color = "lightgreen";
  mbox.style.background = "rgba(50,50,50,0.7)";
  mbox.style.textAlign = "left";
  mbox.style.padding = "2px";
  mbox.style.fontSize = 'x-small';
  
  void repositionBox(Event e) {
    int left = 10 + c.offsetLeft;
    int top  = 28 + c.offsetTop;
    
    mbox.style.left = "${left}px";
    mbox.style.top = "${top}px";
    
    print("repositionBox: event=$e: left=${mbox.style.left} top=${mbox.style.top}");
  }
  
  repositionBox(null);
  
  c.onChange.listen(repositionBox);
  
  return mbox;
}

DivElement canvasbox;

RenderingContext boot() {
  canvas = new CanvasElement();
  assert(canvas != null);
  canvas.id = "main_canvas";
  canvas.width = 780;
  canvas.height = 500;
  canvasbox = querySelector("#canvasbox");
  assert(canvasbox != null);  
  canvasbox.append(canvas);
    
  RenderingContext gl = initGL(canvas);
  if (gl == null) {
    canvas.remove();
    ParagraphElement p = new ParagraphElement();
    p.text = 'WebGL is currently not available on this system.';
    canvasbox.append(p);
    AnchorElement a = new AnchorElement();
    a.href = 'http://get.webgl.org/';
    a.text = 'Get more information';
    canvasbox.append(a);
    canvasbox.style.backgroundColor = 'lightblue';    
    return null;
  }
  
  messagebox = createMessagebox('messagebox', canvas);
  canvasbox.append(messagebox);
  
  initShowPicking();
    
  String sid = Cookie.getCookie("sid");
  assert(sid != null);
  assert(sid is String);
  
  String wsUri = querySelector("#wsUri").text;
  assert(wsUri != null);
  assert(wsUri is String);
  
  Element statusElem = querySelector("#ws_status");
  assert(statusElem != null);
  assert(statusElem is Element);
  
  void dispatch(int code, String data, Map<String,String> tab) {
    dispatcher(gl, code, data, tab);
  }

  initWebSocket(wsUri, sid, 1, statusElem, dispatch);
  
  initStats();
  
  return gl;
}

void demoInitSquares(RenderingContext gl) {
  ShaderProgram squareProgram = new ShaderProgram(gl, "clip");
  programList.add(squareProgram);
  squareProgram.fetch(shaderCache, "${asset.shader}/clip_vs.txt", "${asset.shader}/clip_fs.txt");
  Model squareModel = new Model.fromJson(gl, "${asset.mesh}/square.json", false);
  squareProgram.addModel(squareModel);
  Instance squareInstance = new Instance(squareModel, new Vector3(0.0, 0.0, 0.0), 1.0);
  squareModel.addInstance(squareInstance);

  ShaderProgram squareProgram2 = new ShaderProgram(gl, "clip2");
  programList.add(squareProgram2);
  // execute after 2 secs, giving time to first program populate shadeCache
  new Timer(new Duration(seconds:2), () {
    squareProgram2.fetch(shaderCache, "${asset.shader}/clip_vs.txt", "${asset.shader}/clip2_fs.txt");
  });
  Model squareModel2 = new Model.fromJson(gl, "${asset.mesh}/square2.json", false);
  squareProgram2.addModel(squareModel2);
  Instance squareInstance2 = new Instance(squareModel2, new Vector3(0.0, 0.0, 0.0), 1.0);
  squareModel2.addInstance(squareInstance2);
  
  ShaderProgram squareProgram3 = new ShaderProgram(gl, "clip3");
  programList.add(squareProgram3);
  squareProgram3.fetch(shaderCache, "${asset.shader}/clip_vs.txt", "${asset.shader}/clip3_fs.txt");
  Model squareModel3 = new Model.fromJson(gl, "${asset.mesh}/square3.json", false);
  squareProgram3.addModel(squareModel3);
  Instance squareInstance3 = new Instance(squareModel3, new Vector3(0.0, 0.0, 0.0), 1.0);
  squareModel3.addInstance(squareInstance3);  
}

void addSkybox(RenderingContext gl, Map<String,String> s) {
  skybox = new SkyboxProgram(gl);
  skybox.fetch(shaderCache, s['vertexShader'], s['fragmentShader']);
  SkyboxModel skyboxModel = new SkyboxModel.fromJson(gl, s['cube'], true, 0);
  skyboxModel.addCubemapFace(gl, RenderingContext.TEXTURE_CUBE_MAP_POSITIVE_X, s['faceRight']);
  skyboxModel.addCubemapFace(gl, RenderingContext.TEXTURE_CUBE_MAP_NEGATIVE_X, s['faceLeft']);
  skyboxModel.addCubemapFace(gl, RenderingContext.TEXTURE_CUBE_MAP_POSITIVE_Y, s['faceUp']);
  skyboxModel.addCubemapFace(gl, RenderingContext.TEXTURE_CUBE_MAP_NEGATIVE_Y, s['faceDown']);
  skyboxModel.addCubemapFace(gl, RenderingContext.TEXTURE_CUBE_MAP_POSITIVE_Z, s['faceFront']);
  skyboxModel.addCubemapFace(gl, RenderingContext.TEXTURE_CUBE_MAP_NEGATIVE_Z, s['faceBack']);  
  skybox.addModel(skyboxModel);
  SkyboxInstance skyboxInstance = new SkyboxInstance(skyboxModel, new Vector3(0.0, 0.0, 0.0), skyboxScale, false);
  skyboxModel.addInstance(skyboxInstance);
}

void demoInitSkybox(RenderingContext gl) {
  skybox = new SkyboxProgram(gl);
  skybox.fetch(shaderCache, "${asset.shader}/skybox_vs.txt", "${asset.shader}/skybox_fs.txt");
  SkyboxModel skyboxModel = new SkyboxModel.fromJson(gl, "${asset.mesh}/cube.json", true, 0);
  skyboxModel.addCubemapFace(gl, RenderingContext.TEXTURE_CUBE_MAP_POSITIVE_X, '${asset.texture}/space_rt.jpg');
  skyboxModel.addCubemapFace(gl, RenderingContext.TEXTURE_CUBE_MAP_NEGATIVE_X, '${asset.texture}/space_lf.jpg');
  skyboxModel.addCubemapFace(gl, RenderingContext.TEXTURE_CUBE_MAP_POSITIVE_Y, '${asset.texture}/space_up.jpg');
  skyboxModel.addCubemapFace(gl, RenderingContext.TEXTURE_CUBE_MAP_NEGATIVE_Y, '${asset.texture}/space_dn.jpg');
  skyboxModel.addCubemapFace(gl, RenderingContext.TEXTURE_CUBE_MAP_POSITIVE_Z, '${asset.texture}/space_fr.jpg');
  skyboxModel.addCubemapFace(gl, RenderingContext.TEXTURE_CUBE_MAP_NEGATIVE_Z, '${asset.texture}/space_bk.jpg');  
  skybox.addModel(skyboxModel);
  SkyboxInstance skyboxInstance = new SkyboxInstance(skyboxModel, new Vector3(0.0, 0.0, 0.0), 1.0, true);
  skyboxModel.addInstance(skyboxInstance);
}

void demoInitAirship(RenderingContext gl) {
  ShaderProgram prog = new ShaderProgram(gl, "simple");
  programList.add(prog);
  prog.fetch(shaderCache, "${asset.shader}/simple_vs.txt", "${asset.shader}/simple_fs.txt");
  Model airshipModel = new Model.fromOBJ(gl, "${asset.obj}/airship.obj", new Vector3.zero(), new Vector3.zero());
  prog.addModel(airshipModel);
  Instance airshipInstance = new Instance(airshipModel, new Vector3(-8.0, 0.0, 0.0), 1.0, generatePickColor());
  airshipModel.addInstance(airshipInstance);  
}

void demoInitAirshipTex(RenderingContext gl) {
  TexShaderProgram prog = new TexShaderProgram(gl, "simpleTexturizer");
  programList.add(prog);
  prog.fetch(shaderCache, "${asset.shader}/simpleTex_vs.txt", "${asset.shader}/simpleTex_fs.txt");
  
  String objURL = "${asset.obj}/airship.obj"; 

  TexModel airshipModel = new TexModel.fromOBJ(gl, objURL, new Vector3.zero(), new Vector3.zero(), textureTable, asset);
  prog.addModel(airshipModel);
  TexInstance airshipInstance = new TexInstance(airshipModel, new Vector3(0.0, 0.0, 0.0), 1.0, generatePickColor());
  airshipModel.addInstance(airshipInstance);

  TexModel airshipModel2 = new TexModel.fromOBJ(gl, objURL, new Vector3.zero(), new Vector3.zero(), textureTable, asset);
  prog.addModel(airshipModel2);
  TexInstance airshipInstance2 = new TexInstance(airshipModel2, new Vector3(8.0, 0.0, 0.0), 1.0, generatePickColor());
  airshipModel2.addInstance(airshipInstance2);
  
  String colonyShipURL = "${asset.obj}/Colony Ship Ogame Fleet.obj";  
  TexModel colonyShipModel = new TexModel.fromOBJ(gl, colonyShipURL, new Vector3.zero(), new Vector3.zero(), textureTable, asset);
  prog.addModel(colonyShipModel);
  TexInstance colonyShipInstance = new TexInstance(colonyShipModel, new Vector3(0.0, -5.0, -50.0), 1.0, generatePickColor());
  colonyShipModel.addInstance(colonyShipInstance);
    
  String coneURL = "${asset.obj}/cone.obj";  
  TexModel coneModel = new TexModel.fromOBJ(gl, coneURL, new Vector3.zero(), new Vector3.zero(), textureTable, asset);
  prog.addModel(coneModel);
  TexInstance coneInstance = new TexInstance(coneModel, new Vector3(0.0, 2.0, -10.0), 1.0, generatePickColor());
  coneModel.addInstance(coneInstance);
}

void demoInitShips(RenderingContext gl) {
  demoInitAirship(gl);
  demoInitAirshipTex(gl);
}

void addPicker(RenderingContext gl) {
  picker = new PickerShader(gl, programList, canvas.width, canvas.height);
  picker.fetch(shaderCache, "${asset.shader}/picker_vs.txt", "${asset.shader}/picker_fs.txt");  
}

void demoInitPicker(RenderingContext gl) {
  picker = new PickerShader(gl, programList, canvas.width, canvas.height);
  picker.fetch(shaderCache, "${asset.shader}/picker_vs.txt", "${asset.shader}/picker_fs.txt");  
}

void addSolidShader(RenderingContext gl) {
  solidShader = new SolidShader(gl, programList);
  solidShader.fetch(shaderCache, "${asset.shader}/uniformColor_vs.txt", "${asset.shader}/uniformColor_fs.txt");  
}

void resetZone(RenderingContext gl) {
  programList = new List<ShaderProgram>();  // drop existing programs 
  shaderCache = new Map<String,Shader>();   // drop existing compile shader cache
  textureTable = new Map<String,Texture>(); // drop existing texture table  
}

void updateCulling(RenderingContext gl) {
  if (backfaceCulling) {
    gl.frontFace(RenderingContext.CCW);
    gl.cullFace(RenderingContext.BACK);
    gl.enable(RenderingContext.CULL_FACE);
    return;
  }  
  
  gl.disable(RenderingContext.CULL_FACE);
}

void clearColor(RenderingContext gl, double r, g, b, a) {
  pickerClearColor(r, g, b); // save clear color for picking
  gl.clearColor(r, g, b, a);
}

void initContext(RenderingContext gl, GameLoopHtml gameLoop) {
    
  requestZone();  

  clearColor(gl, 0.5, 0.5, 0.5, 1.0);
  gl.enable(RenderingContext.DEPTH_TEST);  // enable depth testing
  gl.depthFunc(RenderingContext.LESS);     // gl.LESS is default depth test
  gl.depthRange(0.0, 1.0);                 // default
  
  // define viewport size
  gl.viewport(0, 0, canvas.width, canvas.height);
  canvasAspect = canvas.width.toDouble() / canvas.height.toDouble(); // save aspect for render loop mat4.perspective
  debug("canvas aspect ratio: $canvasAspect");

  updateCulling(gl);
  
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
  
  updateGameLoop(gameLoop, contextIsLost(), pageHidden());
}

void regularDraw(RenderingContext gl, GameLoopHtml gameLoop) {
  if (solidShader != null) {
    solidShader.drawModels(gameLoop, cam, pMatrix);
  }
  if (programList != null) {
    programList.where((p) => !p.modelList.isEmpty).forEach((p) => p.drawModels(gameLoop, cam, pMatrix));
  }
  if (skybox != null) {
    skybox.drawModels(gameLoop, cam, pMatrix);
  }
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
  setPerspectiveMatrix(pMatrix, fieldOfViewYRadians, canvasAspect, planeNear, planeFar);

  cam.render(gameLoop.renderInterpolationFactor);

  // clear canvas framebuffer
  gl.bindFramebuffer(RenderingContext.FRAMEBUFFER, null);
  gl.clear(RenderingContext.COLOR_BUFFER_BIT | RenderingContext.DEPTH_BUFFER_BIT);
  
  if (picker == null) {
    // only regular draw -- skip picking drawing
    regularDraw(gl, gameLoop);
    return;
  }
  
  // clear offscreen framebuffer
  gl.bindFramebuffer(RenderingContext.FRAMEBUFFER, picker.framebuffer);
  gl.clear(RenderingContext.COLOR_BUFFER_BIT | RenderingContext.DEPTH_BUFFER_BIT);
  
  // restore drawing to default canvas framebuffer
  gl.bindFramebuffer(RenderingContext.FRAMEBUFFER, null);
  
  if (showPicking) {
    // draw only picking -- draw picking on both framebuffers
    
    picker.offscreen = true; // offscreen framebuffer
    picker.drawModels(gameLoop, cam, pMatrix);
    
    picker.offscreen = false; // canvas framebuffer
    picker.drawModels(gameLoop, cam, pMatrix);
    
    return;
  }

  // draw picking on offscreen framebuffer
  picker.offscreen = true;
  picker.drawModels(gameLoop, cam, pMatrix);
  
  regularDraw(gl, gameLoop);
}

void render(RenderingContext gl, GameLoopHtml gameLoop) {
  stats.begin();
  draw(gl, gameLoop);
  stats.end();
}

void readColor(String label, RenderingContext gl, int x, int y, Framebuffer framebuffer, Uint8List color) {
  gl.bindFramebuffer(RenderingContext.FRAMEBUFFER, framebuffer);
  gl.readPixels(x, y, 1, 1, RenderingContext.RGBA, RenderingContext.UNSIGNED_BYTE, color);
  //print("$label: readPixels: x=$x y=$y color=$color");     
}

DivElement dragBox;

void deleteBandSelectionBox(RenderingContext gl, CanvasElement c, bool shift) {
  if (dragBox == null) {
    return;
  }

  int minX = math.min(mouseDragBeginX, mouseDragCurrX);
  int minY = c.height - math.max(mouseDragBeginY, mouseDragCurrY);
  int w = 1 + (mouseDragCurrX - mouseDragBeginX).abs();
  int h = 1 + (mouseDragCurrY - mouseDragBeginY).abs(); 
    
  bandSelection(minX, minY, w, h, picker, gl, shift);  

  dragBox.remove();
  dragBox = null;
}

void createBandSelectionBox(RenderingContext gl, CanvasElement c, bool shift) {
  
  assert(canvasbox != null);

  if (dragBox == null) {
    dragBox = new DivElement();
    
    dragBox.style.border = '1px solid #FFF';
    dragBox.style.zIndex = "2";
    dragBox.style.position = "absolute";
    //dragBo.style.color = "lightgreen";
    //dragBo.style.background = "rgba(50,50,50,0.7)";
    //dragBo.style.textAlign = "left";
    //dragBo.style.padding = "2px";
    //dragBox.style.fontSize = 'x-small';
    
    // Pass through pointer events
    // http://stackoverflow.com/questions/1009753/pass-mouse-events-through-absolutely-positioned-element
    // https://developer.mozilla.org/en/css/pointer-events
    dragBox.style.pointerEvents = "none";
    
    canvasbox.append(dragBox);
  }
  
/*
// show drag box coordinates
dragBox.children.clear();
DivElement d = new DivElement();
d.text = "($mouseDragBeginX,$mouseDragBeginY) - ($mouseDragCurrX,$mouseDragCurrY)";
dragBox.children.add(d);
*/
  
  int minX = math.min(mouseDragBeginX, mouseDragCurrX);
  int minY = c.height - math.max(mouseDragBeginY, mouseDragCurrY);  
  
  int left = minX + c.offsetLeft;
  int top  = math.min(mouseDragBeginY, mouseDragCurrY) + c.offsetTop;  
  int w = 1 + (mouseDragCurrX - mouseDragBeginX).abs();
  int h = 1 + (mouseDragCurrY - mouseDragBeginY).abs();
    
  dragBox.style.left = "${left}px";
  dragBox.style.top = "${top}px";
  dragBox.style.width = "${w}px";
  dragBox.style.height = "${h}px";
}

PickerInstance mouseLeftClick(RenderingContext gl, Mouse m) {
  
  if (picker == null) {
    err("mouseLeftClick: picker not available");
    return null;
  }
  
  int y = canvas.height - m.y;
  
  debug("mouseLeftClick: mouse=${m.x},${m.y} webgl=${m.x},${y}");
  
  Uint8List color = new Uint8List(4);
  //readColor("canvas-framebuffer", gl, m.x, y, null, color);
  readColor("offscreen-framebuffer", gl, m.x, y, picker.framebuffer, color);
  
  //PickerInstance pi = mouseClickHit(picker.instanceList, color);
  PickerInstance pi = picker.findInstanceByColor(color[0], color[1], color[2]);
  
  return pi;
}

void update(RenderingContext gl, GameLoopHtml gameLoop) {
  //print('${gameLoop.frame}: ${gameLoop.frameTime} [dt = ${gameLoop.dt}].');
  
  //
  // handle input
  //

  Mouse m = gameLoop.mouse;
  Keyboard k = gameLoop.keyboard;
  
  bool mouseLeftPressed = m.pressed(Mouse.LEFT);
  bool shiftDown = k.isDown(Keyboard.SHIFT);
  bool ctrlPressed = k.pressed(Keyboard.CTRL);
  bool ctrlReleased = k.released(Keyboard.CTRL);
    
  if (ctrlReleased) {
    deleteBandSelectionBox(gl, canvas, shiftDown);
    mouseDragBeginX = null;
    mouseDragBeginY = null;    
    mouseDragCurrX = null;
    mouseDragCurrY = null;
  }
  
  if (mouseLeftPressed) {
    PickerInstance pi = mouseLeftClick(gl, m);
    mouseSelection(pi, shiftDown);
  }
  if (ctrlPressed) {
    mouseDragBeginX = m.x;
    mouseDragBeginY = m.y;    
    mouseDragCurrX = null;
    mouseDragCurrY = null;
  }
  if (mouseDragBeginX != null) {
    if ((mouseDragCurrX != m.x) || (mouseDragCurrY != m.y)) {
      // mouse moved
      mouseDragCurrX = m.x;
      mouseDragCurrY = m.y;
      createBandSelectionBox(gl, canvas, shiftDown);
    }
  }

  pauseKey(k.isDown(Keyboard.P));
  
  if (paused()) {
    return; // skip all updates below
  }

  //
  // handle non-input updates
  //
  
  cam.update(gameLoop.gameTime);
    
  if (programList != null) {
    programList.forEach((p) => p.update(gameLoop));
  }
  
  if (skybox != null) {
    skybox.update(gameLoop);
  }
}

void checkAntialias(RenderingContext gl) {
  ContextAttributes attr = gl.getContextAttributes();
  
  if (attr == null) {
    print("ugh: gl.getContextAttributes() returned null -- gl.isContextLost() is ${gl.isContextLost()}");
    print("antialias: UNKNOWN");
  }
  else if (attr is! ContextAttributes) {
    print("ugh: gl.getContextAttributes() returned non-ContextAttributes: $attr");    
    print("antialias: UNKNOWN");
  }
  else if (attr.antialias == null) {
    print("ugh: attr.antialias == null");    
    print("antialias: UNKNOWN");    
  }
  else if (attr.antialias is! bool) {
    print("ugh: attr.antialias is! bool");    
    print("antialias: UNKNOWN");    
  }
  else {
    bool antialias = attr.antialias;
    print("antialias: $antialias");
  }
  
  int size = gl.getParameter(RenderingContext.SAMPLES);
  print("antialias MSSA size: $size");  
}

void main() {
  print("--");
  print("main: negentropia dart client starting");
  
  RenderingContext gl = boot();
  if (gl == null) {
    print("WebGL: not available");
    return;
  }
  
  checkAntialias(gl);
  
  anisotropic_filtering_detect(gl);
  
  GameLoopHtml gameLoop = new GameLoopHtml(canvas);
  
  gameLoop.pointerLock.lockOnClick = false; // disable pointer lock

  //print("gameLoop lockOnClick = ${gameLoop.pointerLock.lockOnClick}");

  //print("gameLoop updateStep = ${gameLoop.updateTimeStep} seconds");

  if (debugLostContext) {
    initHandleLostContext(gl, canvas, gameLoop, initContext);
  }
  
  initPageVisibility(gameLoop);
    
  gameLoop.onUpdate = ((gLoop) { 
    update(gl, gLoop);
  });
  gameLoop.onRender = ((gLoop) {
    render(gl, gLoop);
  });

  initContext(gl, gameLoop);
  
  print("main: negentropia dart client ready");
}

import 'dart:html';
import 'dart:web_gl';
import 'dart:typed_data';
import 'dart:math' as math;

import 'package:vector_math/vector_math.dart';
import 'package:vector_math/vector_math_geometry.dart';
//import 'package:vector_math/vector_math_lists.dart';
import 'package:stats/stats.dart';

final String skyboxVertexShaderSource =
"""
attribute vec3 a_Position;
uniform mat4 u_MV;         // model-view
uniform mat4 u_P;          // projection
varying vec4 v_objCoord;   // object-space

void main() {
  v_objCoord = vec4(a_Position, 1.0); // send obj-space coord to fragment shader
  gl_Position = u_P * u_MV * v_objCoord;
}
""";

final String skyboxFragmentShaderSource =
"""
precision mediump float; // required

uniform samplerCube u_Skybox;
varying vec4 v_objCoord;   // object-space

void main() {
  gl_FragColor = textureCube(u_Skybox, v_objCoord.xyz / v_objCoord.w);
}
""";

double canvasAspect = 1.0;
Matrix4 pMatrix = new Matrix4.zero();
double fieldOfViewYRadians = 45 * math.PI / 180; 
double planeNear = 1.0;
double planeFar = 1000.0;

void log(String msg) {
  print(msg);
  DivElement div = new DivElement();
  div.text = msg;
  document.body.append(div);
}

void updateCulling(RenderingContext gl, bool backfaceCulling) {
  if (backfaceCulling) {
    gl.frontFace(RenderingContext.CCW);
    gl.cullFace(RenderingContext.BACK);
    gl.enable(RenderingContext.CULL_FACE);
    return;
  }

  gl.disable(RenderingContext.CULL_FACE);
}

void updatePerspectiveMatrix(RenderingContext gl, UniformLocation u_P) {
  setPerspectiveMatrix(pMatrix, fieldOfViewYRadians, canvasAspect, planeNear, planeFar);
  gl.uniformMatrix4fv(u_P, false, pMatrix.storage);  
}

void activateTextureUnit(RenderingContext gl, UniformLocation u_Skybox) {
  int unit = 0;
  gl.activeTexture(RenderingContext.TEXTURE0 + unit);
  gl.uniform1i(u_Skybox, unit); 
}

Matrix4 MV = new Matrix4.zero();
double radius = 10.0;
Vector3 camPosition = new Vector3(0.0, 0.0, radius);
Vector3 camFocusPosition = new Vector3(0.0, 0.0, 0.0);
Vector3 camUpDirection = new Vector3(0.0, 1.0, 0.0);
double scale = 1.0;

void updateModelView(RenderingContext gl, UniformLocation u_MV) {
  setViewMatrix(MV, camPosition, camFocusPosition, camUpDirection);
  
  // 1. obj scale
  MV.scale(scale, scale, scale);

  gl.uniformMatrix4fv(u_MV, false, MV.storage);  
}

Texture cubemapTexture;

void addCubemapFace(RenderingContext gl, int face, String URL) {

  ImageElement image = new ImageElement();

  void handleDone(Event e) {
    gl.bindTexture(RenderingContext.TEXTURE_CUBE_MAP, cubemapTexture);
    gl.texParameteri(RenderingContext.TEXTURE_CUBE_MAP,
        RenderingContext.TEXTURE_MAG_FILTER, RenderingContext.NEAREST);
    gl.texParameteri(RenderingContext.TEXTURE_CUBE_MAP,
        RenderingContext.TEXTURE_MIN_FILTER, RenderingContext.NEAREST);

    gl.texImage2DImage(face, 0, RenderingContext.RGBA, RenderingContext.RGBA,
        RenderingContext.UNSIGNED_BYTE, image);

    gl.texParameteri(RenderingContext.TEXTURE_CUBE_MAP,
        RenderingContext.TEXTURE_WRAP_S, RenderingContext.CLAMP_TO_EDGE);
    gl.texParameteri(RenderingContext.TEXTURE_CUBE_MAP,
        RenderingContext.TEXTURE_WRAP_T, RenderingContext.CLAMP_TO_EDGE);
    
    log("loaded cubemap texture: face=$face url=$URL");

    //anisotropic_filtering_enable(gl);

    gl.bindTexture(RenderingContext.TEXTURE_CUBE_MAP, null);
  }

  void handleError(Event e) {
    log("addCubemapFace: handleError: failure loading image from URL: $URL: $e");
  }

  image
      ..onLoad.listen(handleDone)
      ..onError.listen(handleError)
      ..src = URL;
}

void initCubemap(RenderingContext gl) {
  cubemapTexture = gl.createTexture();
  
  addCubemapFace(gl, RenderingContext.TEXTURE_CUBE_MAP_POSITIVE_X, 'space_rt.jpg');
  addCubemapFace(gl, RenderingContext.TEXTURE_CUBE_MAP_NEGATIVE_X, 'space_lf.jpg');
  addCubemapFace(gl, RenderingContext.TEXTURE_CUBE_MAP_POSITIVE_Y, 'space_up.jpg');
  addCubemapFace(gl, RenderingContext.TEXTURE_CUBE_MAP_NEGATIVE_Y, 'space_dn.jpg');
  addCubemapFace(gl, RenderingContext.TEXTURE_CUBE_MAP_POSITIVE_Z, 'space_fr.jpg');
  addCubemapFace(gl, RenderingContext.TEXTURE_CUBE_MAP_NEGATIVE_Z, 'space_bk.jpg');  
}

Buffer bufferVertexPosition;
Buffer bufferVertexIndex;
const int bufferVertexPositionItemSize = 3; // coord x,y,z
//const int bufferVertexIndexItemSize = 2; // size of Uint16Array
int cubeIndexLength;

void createBuffers(RenderingContext gl) {

  /*
  final List<double> cubeVertCoord = [1.0, 1.0, 1.0, -1.0, 1.0, 1.0, -1.0, -1.0, 1.0,
      1.0, -1.0, 1.0, 1.0, 1.0, -1.0, -1.0, 1.0, -1.0, -1.0, -1.0, -1.0, 1.0, -1.0,
      -1.0, -1.0, 1.0, 1.0, -1.0, 1.0, -1.0, -1.0, -1.0, -1.0, -1.0, -1.0, 1.0, 1.0,
      1.0, 1.0, 1.0, -1.0, 1.0, 1.0, -1.0, -1.0, 1.0, 1.0, -1.0, 1.0, 1.0, 1.0, 1.0,
      1.0, -1.0, -1.0, 1.0, -1.0, -1.0, 1.0, 1.0, 1.0, -1.0, 1.0, 1.0, -1.0, -1.0,
      -1.0, -1.0, -1.0, -1.0, -1.0, 1.0];

  final List<int> cubeInd = [0, 1, 2, 0, 2, 3, 4, 6, 5, 4, 7, 6, 8, 9, 10, 8, 10, 11,
      12, 13, 14, 12, 14, 15, 16, 17, 18, 16, 18, 19, 20, 22, 21, 20, 23, 22];
      */
  
  CubeGenerator gen = new CubeGenerator();
  double edge = 1.0;
  MeshGeometry geo = gen.createCube(edge, edge, edge, flags: new GeometryGeneratorFlags(texCoords:false, normals:false, tangents:false));

  Float32List cubeVertCoord = geo.buffer;
  Uint16List cubeInd = geo.indices;  
  
  cubeIndexLength = cubeInd.length;
  
  log("cube indices = ${cubeInd.length} $cubeInd");  
  log("cube vertices = ${cubeVertCoord.length} $cubeVertCoord");
  
  bufferVertexPosition = gl.createBuffer();
  gl.bindBuffer(RenderingContext.ARRAY_BUFFER, bufferVertexPosition);
  gl.bufferDataTyped(RenderingContext.ARRAY_BUFFER, cubeVertCoord, RenderingContext.STATIC_DRAW);

  bufferVertexIndex = gl.createBuffer();
  gl.bindBuffer(RenderingContext.ELEMENT_ARRAY_BUFFER, bufferVertexIndex);
  //gl.bufferDataTyped(RenderingContext.ELEMENT_ARRAY_BUFFER, new Uint16List.fromList(cubeInd.reversed.toList()), RenderingContext.STATIC_DRAW);
  gl.bufferDataTyped(RenderingContext.ELEMENT_ARRAY_BUFFER, cubeInd, RenderingContext.STATIC_DRAW);
}

Stats stats;

void initStats() {
  stats = new Stats();
  document.body.append(stats.container);
}

void main() {
  CanvasElement canvas = new CanvasElement();
  canvas.id = 'webgl_canvas';
  canvas.width = 600;
  canvas.height = 400;
  canvas.style.border = '2px solid black';
  document.body.append(canvas);
  log(
      "canvas '${canvas.id}' created: width=${canvas.width} height=${canvas.height}");

  RenderingContext gl = canvas.getContext3d(preserveDrawingBuffer: false);
  if (gl == null) {
    log("WebGL: initialization failure");
    return;
  }
  log("WebGL: initialized");

  Shader vertShader = gl.createShader(RenderingContext.VERTEX_SHADER);
  gl.shaderSource(vertShader, skyboxVertexShaderSource);
  gl.compileShader(vertShader);
  bool vparameter = gl.getShaderParameter(vertShader, RenderingContext.COMPILE_STATUS);
  if (!vparameter) {
    String infoLog = gl.getShaderInfoLog(vertShader);
    log("compileShader: FAILURE: infoLog=$infoLog vertShader=$skyboxVertexShaderSource");
    return;
  }

  Shader fragShader = gl.createShader(RenderingContext.FRAGMENT_SHADER);
  gl.shaderSource(fragShader, skyboxFragmentShaderSource);
  gl.compileShader(fragShader);
  bool fparameter = gl.getShaderParameter(fragShader, RenderingContext.COMPILE_STATUS);
  if (!fparameter) {
    String infoLog = gl.getShaderInfoLog(fragShader);
    log("compileShader: FAILURE: infoLog=$infoLog fragShader=$skyboxFragmentShaderSource");
    return;
  }

  Program p = gl.createProgram();
  gl.attachShader(p, vertShader);
  gl.attachShader(p, fragShader);
  gl.linkProgram(p);
  bool lparameter = gl.getProgramParameter(p, RenderingContext.LINK_STATUS);
  if (!lparameter) {
    String infoLog = gl.getProgramInfoLog(p);
    log("linkProgram: FAILURE: $infoLog");
    return;
  }
  
  int a_Position = gl.getAttribLocation(p, "a_Position");
  UniformLocation u_MV = gl.getUniformLocation(p, "u_MV");
  UniformLocation u_P = gl.getUniformLocation(p, "u_P");
  UniformLocation u_Skybox = gl.getUniformLocation(p, "u_Skybox");
  
  gl.useProgram(p);
  gl.enableVertexAttribArray(a_Position);

  gl.clearColor(0.5, 0.5, 0.5, 1.0); // clear color
  gl.enable(RenderingContext.DEPTH_TEST); // enable depth testing
  gl.depthFunc(RenderingContext.LESS); // gl.LESS is default depth test
  gl.depthRange(0.0, 1.0); // default
  gl.viewport(0, 0, canvas.width, canvas.height);
  canvasAspect = canvas.width.toDouble() / canvas.height.toDouble();
  
  initStats();
  
  initCubemap(gl);
  
  createBuffers(gl);
  
  updateCulling(gl, false);

  updatePerspectiveMatrix(gl, u_P);
  
  activateTextureUnit(gl, u_Skybox); 
  
  window.animationFrame.then((delta) => render(delta, gl, u_MV, a_Position));
}

void render(num delta, RenderingContext gl, UniformLocation u_MV, int a_Position) {
  stats.begin();
  
  double degreesPerSec = 30.0;
  double radiansPerSec = degreesPerSec * math.PI / 180.0;
  double radians = radiansPerSec * delta.toDouble() / 1000.0;
  
  double sin = math.sin(radians);
  double cos = math.cos(radians);
  
  scale = 1.0 + 20.0 * sin.abs();
  
  camPosition[0] = radius * sin;
  camPosition[2] = radius * cos;
  
  updateModelView(gl, u_MV);
  
  gl.clear(RenderingContext.COLOR_BUFFER_BIT | RenderingContext.DEPTH_BUFFER_BIT);
  
  gl.bindTexture(RenderingContext.TEXTURE_CUBE_MAP, cubemapTexture);
  
  gl.bindBuffer(RenderingContext.ARRAY_BUFFER, bufferVertexPosition);
  gl.vertexAttribPointer(a_Position, bufferVertexPositionItemSize, RenderingContext.FLOAT, false, 0, 0);

  gl.bindBuffer(RenderingContext.ELEMENT_ARRAY_BUFFER, bufferVertexIndex);
  gl.drawElements(RenderingContext.TRIANGLES, cubeIndexLength, RenderingContext.UNSIGNED_SHORT, 0);
  
  gl.bindTexture(RenderingContext.TEXTURE_CUBE_MAP, null);
  gl.bindBuffer(RenderingContext.ARRAY_BUFFER, null);
  gl.bindBuffer(RenderingContext.ELEMENT_ARRAY_BUFFER, null);
  
  window.animationFrame.then((d) => render(d, gl, u_MV, a_Position)); // reschedule
  
  stats.end();  
}

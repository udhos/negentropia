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

void main() {
  gl_Position = u_P * u_MV * vec4(a_Position, 1.0);
}
""";

final String skyboxFragmentShaderSource =
"""
precision mediump float;
uniform vec4 u_Color;
void main(void) {
  gl_FragColor = u_Color;
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

Matrix4 MV = new Matrix4.zero();
/*
Vector3 camPosition = new Vector3(0.0, 0.0, 10.0);
Vector3 camFocusPosition = new Vector3(0.0, 0.0, 0.0);
Vector3 camUpDirection = new Vector3(0.0, 1.0, 0.0);
*/
Vector3 camPosition = new Vector3(0.0, 10.0, 0.0);
Vector3 camFocusPosition = new Vector3(0.0, 0.0, 0.0);
Vector3 camUpDirection = new Vector3(0.0, 0.0, -1.0);
double scale = 3.0;

void updateModelView(RenderingContext gl, UniformLocation u_MV) {
  setViewMatrix(MV, camPosition, camFocusPosition, camUpDirection);
  
  // 1. obj scale
  MV.scale(scale, scale, scale);

  gl.uniformMatrix4fv(u_MV, false, MV.storage);  
}

Buffer bufferVertexPosition;
Buffer bufferVertexIndex;
const int bufferVertexPositionItemSize = 3; // coord x,y,z
//const int bufferVertexIndexItemSize = 2; // size of Uint16Array
int indexLength;

void createBuffers(RenderingContext gl) {
 
  CircleGenerator gen = new CircleGenerator();
  double radius = 1.0;
  MeshGeometry geo = gen.createCircle(radius, flags: new GeometryGeneratorFlags(texCoords:false, normals:false, tangents:false));

  Float32List vertCoord = geo.buffer;
  Uint16List ind = geo.indices;  
  
  indexLength = ind.length;
  
  log("indices = ${ind.length} $ind");  
  log("vertices = ${vertCoord.length} $vertCoord");
  
  bufferVertexPosition = gl.createBuffer();
  gl.bindBuffer(RenderingContext.ARRAY_BUFFER, bufferVertexPosition);
  gl.bufferDataTyped(RenderingContext.ARRAY_BUFFER, vertCoord, RenderingContext.STATIC_DRAW);

  bufferVertexIndex = gl.createBuffer();
  gl.bindBuffer(RenderingContext.ELEMENT_ARRAY_BUFFER, bufferVertexIndex);
  gl.bufferDataTyped(RenderingContext.ELEMENT_ARRAY_BUFFER, ind, RenderingContext.STATIC_DRAW);
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
  UniformLocation u_Color = gl.getUniformLocation(p, "u_Color");
  
  gl.useProgram(p);
  gl.enableVertexAttribArray(a_Position);

  gl.clearColor(0.5, 0.5, 0.5, 1.0); // clear color
  gl.enable(RenderingContext.DEPTH_TEST); // enable depth testing
  gl.depthFunc(RenderingContext.LESS); // gl.LESS is default depth test
  gl.depthRange(0.0, 1.0); // default
  gl.viewport(0, 0, canvas.width, canvas.height);
  
  void setAspect(double aspect) {
    canvasAspect = aspect;
    log("aspect: $canvasAspect");
  }
  setAspect(canvas.width.toDouble() / canvas.height.toDouble());
  
  void keyPress(KeyboardEvent e) {
    log("keyCode=${e.keyCode}");
    switch (e.keyCode) {
      case 43: // +
        setAspect(canvasAspect + 0.1);
        break;
      case 45: // -
        setAspect(canvasAspect - 0.1);
        break;
      case 61: // =
        setAspect(canvas.width.toDouble() / canvas.height.toDouble());
        break;
    }
  }
  document.onKeyPress.listen(keyPress);  
  
  initStats();
  
  createBuffers(gl);
  
  updateCulling(gl, false);
  
  gl.uniform4f(u_Color, 1.0, 0.0, 0.0, 1.0);
  
  window.animationFrame.then((delta) => render(delta, gl, u_MV, u_P, a_Position));
}

void render(num delta, RenderingContext gl, UniformLocation u_MV, UniformLocation u_P, int a_Position) {
  stats.begin();

  updatePerspectiveMatrix(gl, u_P);

  updateModelView(gl, u_MV);
  
  gl.clear(RenderingContext.COLOR_BUFFER_BIT | RenderingContext.DEPTH_BUFFER_BIT);
   
  gl.bindBuffer(RenderingContext.ARRAY_BUFFER, bufferVertexPosition);
  gl.vertexAttribPointer(a_Position, bufferVertexPositionItemSize, RenderingContext.FLOAT, false, 0, 0);

  gl.bindBuffer(RenderingContext.ELEMENT_ARRAY_BUFFER, bufferVertexIndex);
  gl.drawElements(RenderingContext.TRIANGLES, indexLength, RenderingContext.UNSIGNED_SHORT, 0);
  
  gl.bindTexture(RenderingContext.TEXTURE_CUBE_MAP, null);
  gl.bindBuffer(RenderingContext.ARRAY_BUFFER, null);
  gl.bindBuffer(RenderingContext.ELEMENT_ARRAY_BUFFER, null);
  
  window.animationFrame.then((d) => render(d, gl, u_MV, u_P, a_Position)); // reschedule
  
  stats.end();  
}

import 'dart:html';
import 'dart:web_gl';
import 'dart:typed_data';

final String vertexShaderSource = """
attribute vec3 a_Position;
void main(void) {
  gl_Position = vec4(a_Position, 1.0);
}
""";

final String fragmentShaderSource = """
precision mediump float;
uniform vec4 u_Color;
void main(void) {
  gl_FragColor = u_Color;
}
""";

Float32List vertCoord = new Float32List.fromList(
    [0.5, 0.0, -1.0,
     -0.5, 0.0, -1.0,
     -0.5, -0.5, -1.0,
     0.5, -0.5, -1.0]);

Uint16List indices = new Uint16List.fromList(
    [0, 1, 2,
     0, 2, 3]);

void log(String msg) {
  print(msg);
  DivElement m = new DivElement();
  m.text = msg;
  document.body.append(m); 
}

Framebuffer createRenderbuffer(RenderingContext gl, int width, int height) {
  // 1. Init Picking Texture
  Texture texture = gl.createTexture();
  gl.bindTexture(RenderingContext.TEXTURE_2D, texture);
  try {
    gl.texImage2DTyped(RenderingContext.TEXTURE_2D, 0, RenderingContext.RGBA, width, height, 0, RenderingContext.RGBA, RenderingContext.UNSIGNED_BYTE, null);
  }
  catch (e) {
    // https://code.google.com/p/dart/issues/detail?id=11498
    log("gl.texImage2D: exception: $e"); 
  }
  
  // 2. Init Render Buffer
  Renderbuffer renderbuffer = gl.createRenderbuffer();
  gl.bindRenderbuffer(RenderingContext.RENDERBUFFER, renderbuffer);
  gl.renderbufferStorage(RenderingContext.RENDERBUFFER, RenderingContext.DEPTH_COMPONENT16, width, height); 
  
  // 3. Init Frame Buffer
  Framebuffer framebuffer = gl.createFramebuffer();
  gl.bindFramebuffer(RenderingContext.FRAMEBUFFER, framebuffer);
  gl.framebufferTexture2D(RenderingContext.FRAMEBUFFER, RenderingContext.COLOR_ATTACHMENT0, RenderingContext.TEXTURE_2D, texture, 0);
  gl.framebufferRenderbuffer(RenderingContext.FRAMEBUFFER, RenderingContext.DEPTH_ATTACHMENT, RenderingContext.RENDERBUFFER, renderbuffer);

  // 4. Clean up
  gl.bindTexture(RenderingContext.TEXTURE_2D, null);
  gl.bindRenderbuffer(RenderingContext.RENDERBUFFER, null);
  gl.bindFramebuffer(RenderingContext.FRAMEBUFFER, null);
  
  return framebuffer;
}

void render(RenderingContext gl, Framebuffer framebuffer, UniformLocation u_Color, int a_Position, Buffer vertexPositionBuffer, int vertexPositionBufferItemSize, Buffer vertexIndexBuffer, int indicesLength) {
  gl.bindFramebuffer(RenderingContext.FRAMEBUFFER, framebuffer);
  gl.clear(RenderingContext.COLOR_BUFFER_BIT | RenderingContext.DEPTH_BUFFER_BIT);
  
  Float32List red = new Float32List.fromList([1.0, 0.0, 0.0, 1.0]);
  
  gl.uniform4fv(u_Color, red);
  
  gl.bindBuffer(RenderingContext.ARRAY_BUFFER, vertexPositionBuffer);
  gl.vertexAttribPointer(a_Position, vertexPositionBufferItemSize, RenderingContext.FLOAT, false, 0, 0);
  
  gl.bindBuffer(RenderingContext.ELEMENT_ARRAY_BUFFER, vertexIndexBuffer);
  gl.drawElements(RenderingContext.TRIANGLES, indices.length, RenderingContext.UNSIGNED_SHORT, 0);    
}

void readColor(String label, RenderingContext gl, int x, int y, Framebuffer framebuffer) {
  Uint8List color = new Uint8List(4);
  gl.bindFramebuffer(RenderingContext.FRAMEBUFFER, framebuffer);
  gl.readPixels(x, y, 1, 1, RenderingContext.RGBA, RenderingContext.UNSIGNED_BYTE, color);
  log("$label: readPixels: x=$x y=$y color=$color");     
}

void main() {
  ParagraphElement title = new ParagraphElement();
  title.text = 'Click on canvas';
  document.body.append(title);
  
  CanvasElement canvas = new CanvasElement();
  canvas.id = 'webgl_canvas';
  canvas.width = 600;
  canvas.height = 400;
  canvas.style.border = '2px solid black';
  document.body.append(canvas);
  log("canvas '${canvas.id}' created: width=${canvas.width} height=${canvas.height}");

  RenderingContext gl = canvas.getContext3d(preserveDrawingBuffer: true);
  if (gl == null) {
    log("WebGL: initialization failure");
    return;
  }
  log("WebGL: initialized");
  
  Buffer vertexPositionBuffer = gl.createBuffer();
  final int vertexPositionBufferItemSize = 3; // coord x,y,z
  gl.bindBuffer(RenderingContext.ARRAY_BUFFER, vertexPositionBuffer);
  gl.bufferDataTyped(RenderingContext.ARRAY_BUFFER, vertCoord, RenderingContext.STATIC_DRAW);
  
  Buffer vertexIndexBuffer = gl.createBuffer();
  gl.bindBuffer(RenderingContext.ELEMENT_ARRAY_BUFFER, vertexIndexBuffer);
  gl.bufferDataTyped(RenderingContext.ELEMENT_ARRAY_BUFFER, indices, RenderingContext.STATIC_DRAW);
    
  Shader vertShader = gl.createShader(RenderingContext.VERTEX_SHADER);
  gl.shaderSource(vertShader, vertexShaderSource);
  gl.compileShader(vertShader);
  
  Shader fragShader = gl.createShader(RenderingContext.FRAGMENT_SHADER);
  gl.shaderSource(fragShader, fragmentShaderSource);
  gl.compileShader(fragShader); 
  
  Program p = gl.createProgram();
  gl.attachShader(p, vertShader);
  gl.attachShader(p, fragShader);
  gl.linkProgram(p);
  int a_Position = gl.getAttribLocation(p, "a_Position");
  UniformLocation u_Color = gl.getUniformLocation(p, "u_Color");
  
  gl.useProgram(p);
  gl.enableVertexAttribArray(a_Position);  
  
  gl.clearColor(0.5, 0.5, 0.5, 1.0);       // clear color
  gl.enable(RenderingContext.DEPTH_TEST);  // enable depth testing
  gl.depthFunc(RenderingContext.LESS);     // gl.LESS is default depth test
  gl.depthRange(0.0, 1.0);                 // default
  gl.viewport(0, 0, canvas.width, canvas.height);
  
  Framebuffer offscreen = createRenderbuffer(gl, canvas.width, canvas.height);
  
  render(gl, null, u_Color, a_Position, vertexPositionBuffer, vertexPositionBufferItemSize, vertexIndexBuffer, indices.length);      // render on canvas framebuffer
  render(gl, offscreen, u_Color, a_Position, vertexPositionBuffer, vertexPositionBufferItemSize, vertexIndexBuffer, indices.length); // render on offscreen framebuffer
  
  canvas.onClick.listen((MouseEvent e) {
    log("mouse offset: x=${e.offset.x} y=${e.offset.y}");
    
    int x = e.offset.x;
    int y = canvas.height - e.offset.y;
    
    readColor("screen", gl, x, y, null);         // read color from canvas framebuffer
    readColor("offscreen", gl, x, y, offscreen); // read color from offscreen framebuffer
  }); 
}

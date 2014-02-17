part of shader;

void _nextColor(List<double> color) {
  const double d = 1.0 / 16.0;

  int i = 0;

  for ( ; i < 3; ++i) {
    color[i] += d;
    if (color[i] <= 1.0) {
      break;
    }
    color[i] = 0.0;
  }

  if (i == 3) {
    err("_nextColor: overflow: $color");
  }

  //log("_nextColor: $color");
}

List<double> _currentPickColor = [0.0, 0.0, 0.0, 1.0];

Float32List generatePickColor() {

  _nextColor(_currentPickColor);

  bool bgHit = backgroundColorDouble(_currentPickColor[0], _currentPickColor[1],
      _currentPickColor[2]);
  if (bgHit) {
    _nextColor(_currentPickColor);

    bgHit = backgroundColorDouble(_currentPickColor[0], _currentPickColor[1],
        _currentPickColor[2]);
    if (bgHit) {
      err("ugh: generatePickColor: background color: $_currentPickColor");
    }
  }

  debug("generatePickColor: $_currentPickColor bgHit=$bgHit");

  return new Float32List.fromList(_currentPickColor);
}

/*
PickerInstance mouseClickHit(Iterable<Instance> list, Uint8List color) {  
  return colorHit(list, color[0], color[1], color[2]);
}
*/

class PickerInstance extends Instance {

  PickerInstance(Instance i): super(i.model, i.center, i.scale, i.pickColor);

  // the whole purpose of this class is to redefine the draw() method
  // in order to send the pickColor as a uniform to the fragment shader
  void draw(GameLoopHtml gameLoop, ShaderProgram prog, Camera cam) {
    RenderingContext gl = prog.gl;
    gl.uniform4fv((prog as PickerShader).u_Color, pickColor);
    super.draw(gameLoop, prog, cam);
  }
}

class PickerShader extends ShaderProgram {

  UniformLocation u_Color;
  List<PickerInstance> _instanceList = new List<PickerInstance>();

  int get numberOfInstances => _instanceList.length;

  PickerInstance findInstanceByColor(int r, g, b) {
    return colorHit(_instanceList, r, g, b);
  }

  Framebuffer framebuffer;
  bool offscreen;

  void _createRenderbuffer(RenderingContext gl, int width, int height) {
    // 1. Init Picking Texture
    Texture texture = gl.createTexture();
    gl.bindTexture(RenderingContext.TEXTURE_2D, texture);
    //gl.pixelStorei(RenderingContext.UNPACK_FLIP_Y_WEBGL, 1);

    gl.texImage2DTyped(RenderingContext.TEXTURE_2D, 0, RenderingContext.RGBA,
        width, height, 0, RenderingContext.RGBA, RenderingContext.UNSIGNED_BYTE, null);

    // 2. Init Render Buffer
    Renderbuffer renderbuffer = gl.createRenderbuffer();
    gl.bindRenderbuffer(RenderingContext.RENDERBUFFER, renderbuffer);
    gl.renderbufferStorage(RenderingContext.RENDERBUFFER,
        RenderingContext.DEPTH_COMPONENT16, width, height);

    // 3. Init Frame Buffer
    framebuffer = gl.createFramebuffer();
    gl.bindFramebuffer(RenderingContext.FRAMEBUFFER, framebuffer);
    gl.framebufferTexture2D(RenderingContext.FRAMEBUFFER,
        RenderingContext.COLOR_ATTACHMENT0, RenderingContext.TEXTURE_2D, texture, 0);
    gl.framebufferRenderbuffer(RenderingContext.FRAMEBUFFER,
        RenderingContext.DEPTH_ATTACHMENT, RenderingContext.RENDERBUFFER, renderbuffer);

    // 4. Check Frame Buffer status
    int status = gl.checkFramebufferStatus(RenderingContext.FRAMEBUFFER);
    switch (status) {
      case RenderingContext.FRAMEBUFFER_COMPLETE:
        break;
      case RenderingContext.FRAMEBUFFER_UNSUPPORTED:
        err("_createRenderbuffer: FRAMEBUFFER_UNSUPPORTED");
        break;
      case RenderingContext.FRAMEBUFFER_INCOMPLETE_ATTACHMENT:
        err("_createRenderbuffer: FRAMEBUFFER_INCOMPLETE_ATTACHMENT");
        break;
      case RenderingContext.FRAMEBUFFER_INCOMPLETE_DIMENSIONS:
        err("_createRenderbuffer: FRAMEBUFFER_INCOMPLETE_DIMENSIONS");
        break;
      case RenderingContext.FRAMEBUFFER_INCOMPLETE_MISSING_ATTACHMENT:
        err("_createRenderbuffer: FRAMEBUFFER_INCOMPLETE_MISSING_ATTACHMENT");
        break;
      default:
        err("_createRenderbuffer: Framebuffer unexpected status: $status");
    }

    // 5. Clean up
    gl.bindTexture(RenderingContext.TEXTURE_2D, null);
    gl.bindRenderbuffer(RenderingContext.RENDERBUFFER, null);
    gl.bindFramebuffer(RenderingContext.FRAMEBUFFER, null);
  }

  PickerShader(RenderingContext gl, List<ShaderProgram> programList, int
      width, int height): super(gl, "pickerShader") {

    // copy clickable instances
    programList.forEach((p) {
      p.modelList.forEach((m) {
        m.instanceList.where((i) => i.pickColor != null).forEach((ii) {
          PickerInstance pi = new PickerInstance(ii);
          _instanceList.add(pi);
        });
      });
    });

    _createRenderbuffer(gl, width, height);
  }

  void getLocations() {
    super.getLocations();

    u_Color = gl.getUniformLocation(program, "u_Color");
  }

  void drawModels(GameLoopHtml gameLoop, Camera cam, Matrix4 pMatrix) {

    if (!shaderReady) {
      return;
    }

    if (offscreen) {
      gl.bindFramebuffer(RenderingContext.FRAMEBUFFER, framebuffer);
    }

    gl.useProgram(program);
    gl.enableVertexAttribArray(a_Position);

    // send perspective projection matrix uniform
    gl.uniformMatrix4fv(this.u_P, false, pMatrix.storage);

    _instanceList.forEach((i) => i.draw(gameLoop, this, cam));

    // clean up
    gl.bindBuffer(RenderingContext.ARRAY_BUFFER, null);
    gl.bindBuffer(RenderingContext.ELEMENT_ARRAY_BUFFER, null);
    gl.bindFramebuffer(RenderingContext.FRAMEBUFFER, null);
    // default framebuffer

    //gl.disableVertexAttribArray(a_Position); // needed ??
  }

}

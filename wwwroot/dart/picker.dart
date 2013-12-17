part of shader;

List<double> _currentPickColor = [0.0, 0.0, 0.0, 1.0]; 

Float32List generatePickColor() {
  double d = 1.0/16.0;
  
  int i = 0;
  for (; i < 3; ++i) {
    _currentPickColor[i] += d;
    if (_currentPickColor[i] <= 1.0) {
      break;
    }
    _currentPickColor[i] = 0.0;
  } 
  if (i == 3) {
    print("generatePickColor: overflow");
  }
  
  return new Float32List.fromList(_currentPickColor);  
}

bool matchColor(Uint8List i, Float32List f) {
  
  double d0 = (255.0*f[0] - i[0].toDouble()).abs();
  double d1 = (255.0*f[1] - i[1].toDouble()).abs();
  double d2 = (255.0*f[2] - i[2].toDouble()).abs();
  
  return d0 < 1.0 && d1 < 1.0 && d2 < 1.0;
}

PickerInstance mouseClickHit(Iterable<Instance> list, Uint8List color) {
  
  bool match(Instance i) {
    return matchColor(color, i.pickColor);
  }
  
  Instance pi;
    
  try {
    pi = list.firstWhere(match);
  } catch (e) {
    return null;
  }
  
  return pi as PickerInstance;
}

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
  List<PickerInstance> instanceList = new List<PickerInstance>();
  
  Framebuffer framebuffer;
  bool offscreen;
  
  void _createRenderbuffer(RenderingContext gl, int width, int height) {
    // 1. Init Picking Texture
    Texture texture = gl.createTexture();
    gl.bindTexture(RenderingContext.TEXTURE_2D, texture);
    //gl.pixelStorei(RenderingContext.UNPACK_FLIP_Y_WEBGL, 1);
    
    gl.texImage2DTyped(RenderingContext.TEXTURE_2D, 0, RenderingContext.RGBA, width, height, 0, RenderingContext.RGBA, RenderingContext.UNSIGNED_BYTE, null);
  
    // 2. Init Render Buffer
    Renderbuffer renderbuffer = gl.createRenderbuffer();
    gl.bindRenderbuffer(RenderingContext.RENDERBUFFER, renderbuffer);
    gl.renderbufferStorage(RenderingContext.RENDERBUFFER, RenderingContext.DEPTH_COMPONENT16, width, height); 
    
    // 3. Init Frame Buffer
    framebuffer = gl.createFramebuffer();
    gl.bindFramebuffer(RenderingContext.FRAMEBUFFER, framebuffer);
    gl.framebufferTexture2D(RenderingContext.FRAMEBUFFER, RenderingContext.COLOR_ATTACHMENT0, RenderingContext.TEXTURE_2D, texture, 0);
    gl.framebufferRenderbuffer(RenderingContext.FRAMEBUFFER, RenderingContext.DEPTH_ATTACHMENT, RenderingContext.RENDERBUFFER, renderbuffer);

    // 4. Clean up
    gl.bindTexture(RenderingContext.TEXTURE_2D, null);
    gl.bindRenderbuffer(RenderingContext.RENDERBUFFER, null);
    gl.bindFramebuffer(RenderingContext.FRAMEBUFFER, null);
  }

  PickerShader(RenderingContext gl, List<ShaderProgram> programList, int width, int height) : super(gl, "pickerShader") {

    // copy clickable instances
    programList.forEach((p) {
      p.modelList.forEach((m) {
        m.instanceList.where((i) => i.pickColor != null).forEach((ii) {
          PickerInstance pi = new PickerInstance(ii);
          instanceList.add(pi);
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

    instanceList.forEach((i) => i.draw(gameLoop, this, cam));

    // clean up
    gl.bindBuffer(RenderingContext.ARRAY_BUFFER, null);
    gl.bindBuffer(RenderingContext.ELEMENT_ARRAY_BUFFER, null);
    gl.bindFramebuffer(RenderingContext.FRAMEBUFFER, null); // default framebuffer
    
    //gl.disableVertexAttribArray(a_Position); // needed ??
  }
  
}
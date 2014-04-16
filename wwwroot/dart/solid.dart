part of shader;

class AxisInstance extends Instance {

  static final Float32List red = new Float32List.fromList([1.0, 0.0, 0.0, 1.0]);
  static final Float32List green = new Float32List.fromList([0.0, 1.0, 0.0, 1.0]
      );

  AxisInstance(String id, AxisModel am, Instance i) : super(id, am, i.center,
      i.scale);

  void draw(GameLoopHtml gameLoop, ShaderProgram prog, Camera cam) {

    if (!(model as AxisModel)._axisReady) {
      return;
    }

    RenderingContext gl = prog.gl;

    modelView(gl, prog.u_MV, cam, scale); // set up MV matrix

    gl.bindBuffer(RenderingContext.ARRAY_BUFFER, model.vertexPositionBuffer);
    gl.vertexAttribPointer(prog.a_Position, model.vertexPositionBufferItemSize,
        RenderingContext.FLOAT, false, 0, 0);

    gl.bindBuffer(RenderingContext.ELEMENT_ARRAY_BUFFER, model.vertexIndexBuffer
        );

    Piece p;

    assert(model.pieceList.length == 2);

    gl.uniform4fv((prog as SolidShader).u_Color, red);
    p = model.pieceList[0];
    gl.drawElements(RenderingContext.LINES, p.vertexIndexLength,
        RenderingContext.UNSIGNED_SHORT, p.vertexIndexOffset *
        model.vertexIndexBufferItemSize);

    gl.uniform4fv((prog as SolidShader).u_Color, green);
    p = model.pieceList[1];
    gl.drawElements(RenderingContext.LINES, p.vertexIndexLength,
        RenderingContext.UNSIGNED_SHORT, p.vertexIndexOffset *
        model.vertexIndexBufferItemSize);
  }

}

class AxisModel extends Model {

  bool _axisReady = false;

  AxisModel.fromModel(RenderingContext gl, Model m) : super.init() {

    debug("AxisModel: creating from model=${m._URL}");

    List<int> indices = new List<int>();
    List<double> vertCoord = new List<double>();

    void push(List<double> d, List<int> i, double x, y, z) {
      d.add(x);
      d.add(y);
      d.add(z);

      i.add(i.length);
    }

    void _frontUpReadyCall() {
      int offset = indices.length;
      push(vertCoord, indices, 0.0, 0.0, 0.0);
      push(vertCoord, indices, m._front.x, m._front.y, m._front.z);
      addPiece(offset, indices.length - offset); // red

      offset = indices.length;
      push(vertCoord, indices, 0.0, 0.0, 0.0);
      push(vertCoord, indices, m._up.x, m._up.y, m._up.z);
      addPiece(offset, indices.length - offset); // green

      assert(vertCoord.length == 12);
      assert(indices.length == 4);
      assert(pieceList.length == 2);

      _createBuffers(gl, indices, vertCoord, null, null);

      debug("AxisModel.frontUpReady: created axis model from model=${m._URL}");

      _axisReady = true;
    }

    m.callWhenFrontUpDone(_frontUpReadyCall);
  }
}

class SolidShader extends ShaderProgram {

  UniformLocation u_Color;
  List<AxisInstance> instanceList = new List<AxisInstance>();

  // Override method to scan instanceList
  Instance findInstance(String id) {
    Instance i;
    try {
      i = instanceList.firstWhere((j) => j.id == id);
    } on StateError {
      assert(i == null); // not found
    }
    return i;
  }

  SolidShader(RenderingContext gl, List<ShaderProgram> programList) : super(gl,
      "solidShader") {

    // copy clickable instances
    programList.forEach((p) {
      p.modelList.forEach((m) {
        /*
        debug("SolidShader: model=${m._URL} front=${m._front} up=${m._up}");
        if (m._front == null || m._up == null) {
          return;
        }
        */
        AxisModel am = null;
        m.instanceList.forEach((ii) {
          if (am == null) {
            am = new AxisModel.fromModel(gl, m);
          }
          AxisInstance ai = new AxisInstance(ii.id, am, ii);
          debug(
              "SolidShader: created axis instance=${ai.id} from instance=${ii.id}");
          instanceList.add(ai);
        });
      });
    });

    debug("SolidShader: ${instanceList.length} axis instances have been created"
        );

  }

  void getLocations() {
    super.getLocations();

    u_Color = gl.getUniformLocation(program, "u_Color");

    shaderReady = true;
  }

  void drawModels(GameLoopHtml gameLoop, Camera cam, Matrix4 pMatrix) {

    if (!shaderReady) {
      return;
    }

    gl.useProgram(program);
    gl.enableVertexAttribArray(a_Position);

    gl.uniformMatrix4fv(this.u_P, false, pMatrix.storage);

    instanceList.forEach((i) => i.draw(gameLoop, this, cam));

    // clean up
    gl.bindBuffer(RenderingContext.ARRAY_BUFFER, null);
    gl.bindBuffer(RenderingContext.ELEMENT_ARRAY_BUFFER, null);
  }

}

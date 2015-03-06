part of shader;

class AxisInstance extends Instance {
  static final Float32List red = new Float32List.fromList([1.0, 0.0, 0.0, 1.0]);
  static final Float32List green =
      new Float32List.fromList([0.0, 1.0, 0.0, 1.0]);
  static final Float32List blue =
      new Float32List.fromList([0.0, 0.0, 1.0, 1.0]);

  AxisInstance(String id, AxisModel am, Vector3 center, double scale)
      : super(id, am, center, scale);

  AxisInstance.fromInstance(String id, AxisModel am, Instance i)
      : super(id, am, i._center, i.scale);

  void draw(GameLoopHtml gameLoop, ShaderProgram prog, Camera cam) {
    if (!(model as AxisModel)._axisReady) {
      return;
    }

    RenderingContext gl = prog.gl;

    uploadModelView(gl, prog.u_MV, cam, scale); // set up MV matrix

    gl.bindBuffer(RenderingContext.ARRAY_BUFFER, model.vertexPositionBuffer);
    gl.vertexAttribPointer(prog.a_Position, model.vertexPositionBufferItemSize,
        RenderingContext.FLOAT, false, 0, 0);

    gl.bindBuffer(
        RenderingContext.ELEMENT_ARRAY_BUFFER, model.vertexIndexBuffer);

    Piece p;

    assert(model.pieceList.length == 3); // red, green, blue

    // draw front/red arrow
    gl.uniform4fv((prog as SolidShader).u_Color, red);
    p = model.pieceList[0];
    gl.drawElements(RenderingContext.LINES, p.vertexIndexLength,
        ext_get_element_type,
        p.vertexIndexOffset * model.vertexIndexBufferItemSize);

    // draw up/green arrow
    gl.uniform4fv((prog as SolidShader).u_Color, green);
    p = model.pieceList[1];
    gl.drawElements(RenderingContext.LINES, p.vertexIndexLength,
        ext_get_element_type,
        p.vertexIndexOffset * model.vertexIndexBufferItemSize);

    // draw right/blue arrow
    gl.uniform4fv((prog as SolidShader).u_Color, blue);
    p = model.pieceList[2];
    gl.drawElements(RenderingContext.LINES, p.vertexIndexLength,
        ext_get_element_type,
        p.vertexIndexOffset * model.vertexIndexBufferItemSize);
  }
}

class AxisModel extends Model {
  bool _axisReady = false;

  void _push(List<double> d, List<int> i, double x, double y, double z) {
    d.add(x);
    d.add(y);
    d.add(z);

    i.add(i.length);
  }

  AxisModel(RenderingContext gl, Vector3 front, Vector3 up, Vector3 right)
      : super.init() {
    List<int> indices = new List<int>();
    List<double> vertCoord = new List<double>();

    // add two vertices for front/red arrow
    int offset = indices.length;
    _push(vertCoord, indices, 0.0, 0.0, 0.0);
    _push(vertCoord, indices, front.x, front.y, front.z);
    addPiece(offset, indices.length - offset); // red

    // add two vertices for up/green arrow
    offset = indices.length;
    _push(vertCoord, indices, 0.0, 0.0, 0.0);
    _push(vertCoord, indices, up.x, up.y, up.z);
    addPiece(offset, indices.length - offset); // green

    // add two vertices for right/blue arrow
    offset = indices.length;
    _push(vertCoord, indices, 0.0, 0.0, 0.0);
    _push(vertCoord, indices, right.x, right.y, right.z);
    addPiece(offset, indices.length - offset); // blue

    assert(vertCoord.length == 18);
    assert(indices.length == 6);
    assert(pieceList.length == 3);

    _createBuffers(gl, indices, vertCoord, null, null);

    _axisReady = true;
  }

  AxisModel.fromModel(RenderingContext gl, Model m) : super.init() {
    debug("AxisModel: creating from model=${m._objURL}");

    List<int> indices = new List<int>();
    List<double> vertCoord = new List<double>();

    void _frontUpReadyCall() {

      // add two vertices for front/red arrow
      int offset = indices.length;
      _push(vertCoord, indices, 0.0, 0.0, 0.0);
      _push(vertCoord, indices, m._front.x, m._front.y, m._front.z);
      addPiece(offset, indices.length - offset); // red

      // add two vertices for up/green arrow
      offset = indices.length;
      _push(vertCoord, indices, 0.0, 0.0, 0.0);
      _push(vertCoord, indices, m._up.x, m._up.y, m._up.z);
      addPiece(offset, indices.length - offset); // green

      // add two vertices for right/blue arrow
      Vector3 right = m.right;
      offset = indices.length;
      _push(vertCoord, indices, 0.0, 0.0, 0.0);
      _push(vertCoord, indices, right.x, right.y, right.z);
      addPiece(offset, indices.length - offset); // blue

      assert(vertCoord.length == 18);
      assert(indices.length == 6);
      assert(pieceList.length == 3);

      _createBuffers(gl, indices, vertCoord, null, null);

      debug(
          "AxisModel.frontUpReady: created axis model from model=${m._objURL}");

      _axisReady = true;
    }

    m.callWhenFrontUpDone(_frontUpReadyCall);
  }
}

class SolidShader extends ShaderProgram {
  static final Vector3 FRONT = new Vector3(1.0, 0.0, 0.0);
  static final Vector3 UP = new Vector3(0.0, 1.0, 0.0);
  static final Vector3 RIGHT = new Vector3(0.0, 0.0, 1.0);

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

  void _loadDebugOrigin(RenderingContext gl) {
    double scale = 200.0;
    Vector3 origin = new Vector3.zero();
    //log("SolidShader._loadDebugOrigin: creating {$scale}-meter xyz debug marker at origin $origin");
    AxisModel m = new AxisModel(gl, FRONT, UP, RIGHT);
    AxisInstance i = new AxisInstance("origin", m, origin, scale);
    i.setRotationFromIdentity(); // do not rotate this
    instanceList.add(i);
  }

  SolidShader(RenderingContext gl, List<ShaderProgram> programList)
      : super(gl, "solidShader") {
    _loadDebugOrigin(gl);

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
          AxisInstance ai = new AxisInstance.fromInstance(ii.id, am, ii);
          debug(
              "SolidShader: created axis instance=${ai.id} from instance=${ii.id}");
          instanceList.add(ai);
        });
      });
    });

    debug(
        "SolidShader: ${instanceList.length} axis instances have been created");
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

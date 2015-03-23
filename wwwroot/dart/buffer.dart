part of shader;

class Instance {
  String id;
  Model model;
  double scale;
  Float32List pickColor;

  int inputLock = Keyboard.ZERO;

  Vector3 _center;
  String _mission;

  String toString() {
    return "${super.toString()} id=$id";
  }

  double get boundingRadius {
    double radius;

    if (model.boundingRadius == null) {
      err("$this $id model=$model: undefined model bounding radius");
      radius = 1.0;
    } else {
      radius = model.boundingRadius;
    }

    return scale * radius;
  }

  void copyLocationInto(Vector3 result) {
    result.setFrom(_center);
  }

  void set center(Vector3 c) {
    _center.setFrom(c);
  }

  Vector3 get center {
    return _center.clone();
  }

  void set mission(String m) {
    _mission = m;
  }

  String get mission => _mission;

  Matrix4 MV = new Matrix4.identity(); // model-view matrix

  Matrix4 _undoModelRotation = new Matrix4.identity();
  Matrix4 _rotation = new Matrix4.identity();

  Vector3 _front = new Vector3(1.0, 0.0, 0.0);
  Vector3 _up = new Vector3(0.0, 1.0, 0.0);
  Vector3 _right = new Vector3(0.0, 0.0, 1.0);

  /*
   // It's dangerous to expose this storage.
  Vector3 get front => _front;
  Vector3 get up => _up;
  Vector3 get right => _right;
  */

  void copyFront(Vector3 front) {
    _front.copyInto(front);
  }
  void copyUp(Vector3 up) {
    _up.copyInto(up);
  }
  void copyRight(Vector3 right) {
    _right.copyInto(right);
  }

  // setRotationFromIdentity: do not perform any rotation
  // useful for debug
  void setRotationFromIdentity() {
    _rotation.setIdentity();
  }

  // preload on _undoModelRotation a matrix to
  // undo the model intrinsic local rotation
  void _undoModelRotationFrom(Vector3 modelFront, Vector3 modelUp) {
    Vector3 zeroPosition = new Vector3.zero();

    // rotation matrix = model matrix = inverse of view matrix
    setViewMatrix(_undoModelRotation, zeroPosition, modelFront, modelUp);
  }

  void setRotationFrom(Vector3 newFront, Vector3 newUp) {
    if (inputLock == Keyboard.R) {
      setRotationFromIdentity();
      return;
    }

    // save copy of vectors
    // because we will invert the rotation matrix
    // hence we won't be able to fetch them back from rotation matrix
    newFront.copyInto(_front);
    newUp.copyInto(_up);
    _front.crossInto(_up, _right); // right = front x up
    _right.normalize();

    // rotation matrix = model matrix = inverse of view matrix

    // compound rotation T*R*U:
    // U = first undo model intrinsic local rotation
    // R = then apply the specific rotation we want for the object
    // T = finally translate the object
    /*
    setModelMatrix(
        _model_TRU,
        _front,
        _up,
        _center[0],
        _center[1],
        _center[2]); // _model_TRU = T*R
    _model_TRU.multiply(_undoModelRotation); // _model_TRU = T*R*U
     */
    setRotationMatrix(_rotation, _front, _up); // _rotation = R
    _rotation.multiply(_undoModelRotation); // _rotation = R*U
  }

  String getOrientation() {
    return "f=$_front u=$_up r=$_right";
  }

  void debugLocation([String label = ""]) {
    log("$label$this - model: orient[${this.model.debugOrientation()}] - obj: pos[$_center] orient: ${this.getOrientation()}");
  }

  Instance(this.id, this.model, this._center, this.scale,
      [this.pickColor = null]) {
    Vector3 modelFront = this.model._front.normalized();
    Vector3 modelUp = this.model._up.normalized();
    _undoModelRotationFrom(modelFront, modelUp);
    setRotationFrom(modelFront, modelUp);
    debug(
        "new instance: $this $id model=${model.modelName} center=$_center ${this.getOrientation()}");
  }

  void update(GameLoopHtml gameLoop) {}

  /**
   * Send this object's full OpenGL view matrix into GPU.
   */
  void uploadModelView(
      RenderingContext gl, UniformLocation u_MV, Camera cam, double rescale) {

    /*
      V = View (inverse of camera matrix -- translation and rotation)
      T = Translation
      R = Rotation
      U = Undo Model Local Rotation
      S = Scaling
     */
    cam.loadViewMatrixInto(MV); // MV = V

    MV.translate(_center[0], _center[1], _center[2]); // MV = V*T

    MV.multiply(_rotation); // MV = V*T*R*U

    MV.scale(rescale, rescale, rescale); // MV = V*T*R*U*S

    gl.uniformMatrix4fv(u_MV, false, MV.storage);
  }

  void draw(GameLoopHtml gameLoop, ShaderProgram prog, Camera cam) {
    RenderingContext gl = prog.gl;

    uploadModelView(gl, prog.u_MV, cam, scale); // set up MV matrix

    gl.bindBuffer(RenderingContext.ARRAY_BUFFER, model.vertexPositionBuffer);
    gl.vertexAttribPointer(prog.a_Position, model.vertexPositionBufferItemSize,
        RenderingContext.FLOAT, false, 0, 0);

    gl.bindBuffer(
        RenderingContext.ELEMENT_ARRAY_BUFFER, model.vertexIndexBuffer);

    model.pieceList.forEach((piece) {
      gl.drawElements(RenderingContext.TRIANGLES, piece.vertexIndexLength,
          model.vertexIndexElementType,
          piece.vertexIndexOffset * model.vertexIndexElementSize);
    });
  }
}

class Piece {
  int vertexIndexOffset;
  int vertexIndexLength;

  Piece(this.vertexIndexOffset, this.vertexIndexLength);
}

typedef void frontUpCallbackFunc();

class Model {
  Buffer vertexPositionBuffer;
  Buffer vertexIndexBuffer;
  final int vertexPositionBufferItemSize = 3; // coord x,y,z

  //final int vertexIndexBufferItemSize = 2; // size of Uint16Array

  int _indexElementType = RenderingContext.UNSIGNED_SHORT;
  int _indexElementSize = 2;

  bool _needBigIndex = false;

  bool get useBigIndex => _needBigIndex && ext_element_uint;

  int get vertexIndexElementType => _indexElementType;

  int get vertexIndexElementSize => _indexElementSize;

  bool modelReady = false; // buffers
  bool piecesReady = false; // multiple OBJ pieces

  String _objURL;
  String _modelName;
  String get modelName => _modelName;

  Vector3 _front = new Vector3(1.0, 0.0, 0.0);
  Vector3 _up = new Vector3(0.0, 1.0, 0.0);
  Vector3 get right => _front.cross(_up).normalized().scaled(_front.length);

  String debugOrientation() {
    return "f=$_front u=$_up r=$right";
  }

  List<Piece> pieceList = new List<Piece>();
  List<Instance> instanceList = new List<Instance>();

  double boundingRadius;

  void createIndexBuffer(RenderingContext gl, List<int> indices) {
    vertexIndexBuffer = gl.createBuffer();
    gl.bindBuffer(RenderingContext.ELEMENT_ARRAY_BUFFER, vertexIndexBuffer);

    if (this.useBigIndex) {
      gl.bufferDataTyped(RenderingContext.ELEMENT_ARRAY_BUFFER,
          new Uint32List.fromList(indices), RenderingContext.STATIC_DRAW);
    } else {
      gl.bufferDataTyped(RenderingContext.ELEMENT_ARRAY_BUFFER,
          new Uint16List.fromList(indices), RenderingContext.STATIC_DRAW);
    }

    // clean-up
    gl.bindBuffer(RenderingContext.ELEMENT_ARRAY_BUFFER, null);
  }

  void _createBuffers(RenderingContext gl, List<int> indices,
      List<double> vertCoord, List<double> textCoord, List<double> normCoord) {
    vertexPositionBuffer = gl.createBuffer();
    gl.bindBuffer(RenderingContext.ARRAY_BUFFER, vertexPositionBuffer);

    gl.bufferDataTyped(RenderingContext.ARRAY_BUFFER,
        new Float32List.fromList(vertCoord), RenderingContext.STATIC_DRAW);

    // clean-up
    gl.bindBuffer(RenderingContext.ARRAY_BUFFER, null);

    createIndexBuffer(gl, indices);

    modelReady = true;
  }

  Piece addPiece(int offset, int length) {
    Piece pi = new Piece(offset, length);
    pieceList.add(pi);
    return pi;
  }

  Model.init();

  /*
  Model.fromVert(RenderingContext gl, List<int> indices, List<double> vertCoord) {
    addPiece(0, indices.length); // single-piece model
    assert(pieceList.length == 1);

    _createBuffers(gl, indices, vertCoord, null, null);    
  }
  */

  Model.fromJson(
      RenderingContext gl, this._modelName, String this._objURL, bool reverse) {

    /*
    // load JSON from URL
    var req = new HttpRequest();
    req.open("GET", URL);
    req.onLoad.listen((ProgressEvent e) {
      String response = req.responseText;
      if (req.status != 200) {
        err("Model.fromURL: error: [$response]");
        return;
      }
      err("Model.fromURL: loaded: json=[$response]");
      Map m = parse(response);
      List<num> vertCoord = m["vertCoord"];
      List<int> vertInd = m["vertInd"];
      _createBuffers(gl, vertCoord, vertInd);
    });
    req.onError.listen((e) {
      err("Model.fromURL: error: [$e]");
    });
    req.send();
    */

    void handleResponse(String response) {
      Map m;
      try {
        m = JSON.decode(response);
      } catch (e) {
        err("Model.fromJson: failure parsing JSON: $e");
        return;
      }

      List<int> indices = m['vertInd'];
      if (reverse) {
        indices = indices.reversed.toList();
      }
      List<double> vertCoord = m['vertCoord'];

      addPiece(0, indices.length); // single-piece model
      assert(pieceList.length == 1);

      _createBuffers(gl, indices, vertCoord, null, null);
    }

    void handleError(Object err) {
      err("Model.fromJson: failure fetching JSON from URL: $_objURL: $err");
    }

    HttpRequest.getString(_objURL).then(handleResponse).catchError(handleError);
  }

  void showObjStats(Obj o) {
    double min_x = double.INFINITY;
    double min_y = double.INFINITY;
    double min_z = double.INFINITY;

    double max_x = double.NEGATIVE_INFINITY;
    double max_y = double.NEGATIVE_INFINITY;
    double max_z = double.NEGATIVE_INFINITY;

    for (int i = 0; i < o.vertCoord.length; i += 3) {
      double x = o.vertCoord[i];
      double y = o.vertCoord[i + 1];
      double z = o.vertCoord[i + 2];

      min_x = math.min(min_x, x);
      min_y = math.min(min_y, y);
      min_z = math.min(min_z, z);

      max_x = math.max(max_x, x);
      max_y = math.max(max_y, y);
      max_z = math.max(max_z, z);
    }

    double size_x = (max_x - min_x).abs();
    double size_y = (max_y - min_y).abs();
    double size_z = (max_z - min_z).abs();

    double dx = max_x - min_x;
    double dy = max_y - min_y;
    double dz = max_z - min_z;

    boundingRadius = math.sqrt(dx * dx + dy * dy + dz * dz) / 2.0;

    debug(
        "model=$modelName indices=${o.indices.length} parts=${o.partList.length} ($min_x,$min_y,$min_z)..($max_x,$max_y,$max_z)=[$size_x,$size_y,$size_z] radius=$boundingRadius");
  }

  void loadObj(RenderingContext gl, Obj o) {
    assert(!piecesReady);

    o.partList.forEach((Part pa) {
      addPiece(pa.indexFirst, pa.indexListSize);
    });

    piecesReady = true;
  }

  void saveIndexSize(int indexSize) {}

  frontUpCallbackFunc frontUpCallback;

  void callWhenFrontUpDone(frontUpCallbackFunc callback) {
    //assert(frontUpCallback == null);
    frontUpCallback = callback;
    assert(frontUpCallback != null);
  }

  Model.fromOBJ(RenderingContext gl, this._modelName, this._objURL,
      Vector3 front, Vector3 up) {
    //log("Model.fromOBJ: model=$modelName URL=$_objURL front=$_front up=$_up");

    void handleResponse(String response) {
      //log("Model.fromOBJ: fetched OBJ from URL: $URL");

      _front = front.clone();
      _up = up.clone();

      //log("Model.fromOBJ: handleResponse: model=$modelName URL=$_objURL front=$_front up=$_up");

      if (frontUpCallback != null) {
        frontUpCallback();
      }

      Obj obj = new Obj.fromString(_objURL, response,
          defaultName: "noname", fillMissingTextCoord: true, printStats: true);

      this._needBigIndex = obj.bigIndexFound;

      if (this._needBigIndex) {
        log("OBJ URL=$_objURL requires support for big index>65535");
        if (ext_element_uint) {
          _indexElementType = RenderingContext.UNSIGNED_INT;
          _indexElementSize = 4;
        } else {
          err("OBJ URL=$_objURL requires support for big index>65535, but OES_element_index_uint is unsupported");
        }
      }
      log("OBJ URL=$_objURL useBigIndex=$useBigIndex");

      showObjStats(obj);

      loadObj(gl, obj);

      _createBuffers(
          gl, obj.indices, obj.vertCoord, obj.textCoord, obj.normCoord);
    }

    void handleError(Object err) {
      err("Model.fromOBJ: failure fetching OBJ from URL=$_objURL: $err");
    }

    HttpRequest.getString(_objURL).then(handleResponse).catchError(handleError);
  }

  Model.fromGlobe(RenderingContext gl, this._modelName, double radius,
      Vector3 front, Vector3 up) {
    _front = front.clone();
    _up = up.clone();

    //log("Model.fromGlobe: model=$modelName front=$_front up=$_up");

    if (frontUpCallback != null) {
      frontUpCallback();
    }

    // generateGeometry

    Uint16List globeIndices;
    Float32List globePosCoord;
    Float32List globeTexCoord;
    int indexSize;

    /*
    // texturized rectangle - begin
    globeIndices = new Uint16List.fromList([0, 1, 2, 0, 2, 3]);
    globePosCoord = new Float32List.fromList([
      -radius,
      -radius,
      0.0,
      radius,
      -radius,
      0.0,
      radius,
      radius,
      0.0,
      -radius,
      radius,
      0.0
    ]);
    globeTexCoord =
        new Float32List.fromList([0.0, 0.0, 1.0, 0.0, 1.0, 1.0, 0.0, 1.0]);
    indexSize = globeIndices.length;
    // texturized rectangle - end
     */

    // globe begin
    SphereGenerator gen = new SphereGenerator();
    MeshGeometry geo = gen.createSphere(radius,
        flags: new GeometryGeneratorFlags(
            texCoords: false, normals: false, tangents: false));

    globeIndices = geo.indices;
    indexSize = globeIndices.length;

    Vector3List posCoordList = new Vector3List(indexSize);
    Vector2List texCoordList = new Vector2List(indexSize);
    gen.generateVertexPositions(posCoordList, globeIndices);
    gen.generateVertexTexCoords(texCoordList, posCoordList, globeIndices);

    globePosCoord = posCoordList.buffer;
    globeTexCoord = texCoordList.buffer;
    // globe end

    saveIndexSize(indexSize);

    /*
    // DEBUG
    int vertexCount1 = globePosCoord.length.toInt() ~/ 3;
    int vertexCount2 = globeTexCoord.length.toInt() ~/ 2;
    log("globe vertexCount=${vertexCount1} vertexCount=${vertexCount2}");
    log("globe indexSize=$indexSize");
    log("globe indices: size=${globeIndices.length} $globeIndices");
    log("globe positions: size=${globePosCoord.length} (3 * $vertexCount1) $globePosCoord");
    log("globe tex coord: size=${globeTexCoord.length} (2 * $vertexCount1) $globeTexCoord");
     */

    boundingRadius = radius;

    _createBuffers(gl, globeIndices, globePosCoord, globeTexCoord, null);
  }

  void addInstance(Instance i) {
    instanceList.add(i);
  }

  Instance findInstance(String id) {
    Instance i;
    try {
      i = instanceList.firstWhere((j) => j.id == id);
    } on StateError {
      assert(i == null); // not found
    }
    return i;
  }

  void drawInstances(GameLoopHtml gameLoop, ShaderProgram program, Camera cam) {
    if (!modelReady || !piecesReady) {
      return;
    }

    instanceList.forEach((i) => i.draw(gameLoop, program, cam));
  }

  void update(GameLoopHtml gameLoop) {
    instanceList.forEach((i) => i.update(gameLoop));
  }
}

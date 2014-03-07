part of shader;

class Instance {

  String id;
  Model model;
  double scale;
  Float32List pickColor;

  Vector3 _center;
  String _mission;

  Vector3 get center => _center.clone();

  set center(Vector3 c) => _center.setFrom(c);
  set mission(String m) => _mission = m;

  Matrix4 MV = new Matrix4.identity(); // model-view matrix

  Matrix4 _rotation = new Matrix4.identity();

  Vector3 get front => new Vector3(_rotation.storage[0], _rotation.storage[1],
      _rotation.storage[2]);
  Vector3 get up => new Vector3(_rotation.storage[4], _rotation.storage[5],
      _rotation.storage[6]);
  Vector3 get right => new Vector3(_rotation.storage[8], _rotation.storage[9],
      _rotation.storage[10]);

  void setRotation(Vector3 front, up) {
    /*
    Vector3 right = front.cross(up).normalize();

    _rotation.setValues(front[0], up[0], right[0], 0.0, front[1], up[1],
        right[1], 0.0, front[2], up[2], right[2], 0.0, 0.0, 0.0, 0.0, 1.0);
        */
    setRotationMatrix(_rotation, front, up);
  }

  Instance(this.id, this.model, this._center, this.scale, [this.pickColor =
      null]) {
    setRotation(this.model._front.clone().normalize(), this.model._up.clone(
        ).normalize());
    debug(
        "new instance: $this $id model=${model.modelName} center=$center front=$front up=$up right=$right"
        );
  }

  void update(GameLoopHtml gameLoop) {
  }

  void modelView(RenderingContext gl, UniformLocation u_MV, Camera cam, double
      rescale) {

    // grand world coordinate system:
    // 1. obj scale
    // 2. obj rotate
    // 3. obj orbit translate
    // 4. obj orbit rotate
    // 5. obj translate
    // 6. camera orbit rotate
    // 7. camera translate
    // 8. camera rotate

    //setViewMatrix(MV, cam.eye, cam.center, cam.up);
    /*
    MV.setIdentity();

    // 7. camera translate
    cam.translate(MV);

    // 6. camera orbit rotate
    cam.rotate(MV);
    */

    /*
      V = View (inverse of camera)
      T = Translation
      R = Rotation
      S = Scaling
     */
    cam.viewMatrix(MV); // MV = V

    // 5. obj translate
    MV.translate(_center[0], _center[1], _center[2]); // MV = V*T

    // 2. obj rotate
    MV.multiply(_rotation); // MV = V*T*R

    // 1. obj scale
    MV.scale(rescale, rescale, rescale); // MV = V*T*R*S

    gl.uniformMatrix4fv(u_MV, false, MV.storage);
  }

  void draw(GameLoopHtml gameLoop, ShaderProgram prog, Camera cam) {

    RenderingContext gl = prog.gl;

    modelView(gl, prog.u_MV, cam, scale); // set up MV matrix

    gl.bindBuffer(RenderingContext.ARRAY_BUFFER, model.vertexPositionBuffer);
    gl.vertexAttribPointer(prog.a_Position, model.vertexPositionBufferItemSize,
        RenderingContext.FLOAT, false, 0, 0);

    gl.bindBuffer(RenderingContext.ELEMENT_ARRAY_BUFFER, model.vertexIndexBuffer
        );

    model.pieceList.forEach((piece) {
      gl.drawElements(RenderingContext.TRIANGLES, piece.vertexIndexLength,
          RenderingContext.UNSIGNED_SHORT, piece.vertexIndexOffset *
          model.vertexIndexBufferItemSize);
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
  final int vertexIndexBufferItemSize = 2; // size of Uint16Array

  bool modelReady = false; // buffers
  bool piecesReady = false; // multiple OBJ pieces

  String _URL;
  String get modelName => _URL;

  Vector3 _front = new Vector3(1.0, 0.0, 0.0);
  Vector3 _up = new Vector3(0.0, 1.0, 0.0);

  List<Piece> pieceList = new List<Piece>();
  List<Instance> instanceList = new List<Instance>();

  void _createBuffers(RenderingContext gl, List<int> indices, List<double>
      vertCoord, List<double> textCoord, List<double> normCoord) {

    vertexPositionBuffer = gl.createBuffer();
    gl.bindBuffer(RenderingContext.ARRAY_BUFFER, vertexPositionBuffer);
    gl.bufferDataTyped(RenderingContext.ARRAY_BUFFER, new Float32List.fromList(
        vertCoord), RenderingContext.STATIC_DRAW);

    vertexIndexBuffer = gl.createBuffer();
    gl.bindBuffer(RenderingContext.ELEMENT_ARRAY_BUFFER, vertexIndexBuffer);
    gl.bufferDataTyped(RenderingContext.ELEMENT_ARRAY_BUFFER,
        new Uint16List.fromList(indices), RenderingContext.STATIC_DRAW);

    // clean-up
    gl.bindBuffer(RenderingContext.ARRAY_BUFFER, null);
    gl.bindBuffer(RenderingContext.ELEMENT_ARRAY_BUFFER, null);

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

  Model.fromJson(RenderingContext gl, String URL, bool reverse) {

    /*
    // load JSON from URL
    var req = new HttpRequest();
    req.open("GET", URL);
    req.onLoad.listen((ProgressEvent e) {
      String response = req.responseText;
      if (req.status != 200) {
        print("Model.fromURL: error: [$response]");
        return;
      }
      print("Model.fromURL: loaded: json=[$response]");
      Map m = parse(response);
      List<num> vertCoord = m["vertCoord"];
      List<int> vertInd = m["vertInd"];
      _createBuffers(gl, vertCoord, vertInd);
    });
    req.onError.listen((e) {
      print("Model.fromURL: error: [$e]");
    });
    req.send();
    */

    void handleResponse(String response) {

      Map m;
      try {
        m = JSON.decode(response);
      } catch (e) {
        print("Model.fromJson: failure parsing JSON: $e");
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
      print("Model.fromJson: failure fetching JSON from URL: $URL: $err");
    }

    HttpRequest.getString(URL).then(handleResponse).catchError(handleError);
  }

  void printObjStats(Obj o) {

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

    print(
        "model=$_URL indices=${o.indices.length} parts=${o.partList.length} ($min_x,$min_y,$min_z)..($max_x,$max_y,$max_z)=[$size_x,$size_y,$size_z]"
        );
  }

  void loadObj(RenderingContext gl, Obj o) {
    assert(!piecesReady);

    o.partList.forEach((Part pa) {
      Piece pi = addPiece(pa.indexFirst, pa.indexListSize);
      //print("Model.fromOBJ: added part ${pa.name} into piece: offset=${pi.vertexIndexOffset} length=${pi.vertexIndexLength}");
    });

    piecesReady = true;
  }

  frontUpCallbackFunc frontUpCallback;

  void callWhenFrontUpDone(frontUpCallbackFunc callback) {
    //assert(frontUpCallback == null);
    frontUpCallback = callback;
    assert(frontUpCallback != null);
  }

  Model.fromOBJ(RenderingContext gl, this._URL, Vector3 front, Vector3 up) {

    void handleResponse(String response) {
      //print("Model.fromOBJ: fetched OBJ from URL: $URL");

      _front = front.clone();
      _up = up.clone();

      debug("model=$_URL front=$_front up=$_up");

      if (frontUpCallback != null) {
        frontUpCallback();
      }

      Obj obj = new Obj.fromString(_URL, response);

      printObjStats(obj);

      loadObj(gl, obj);

      _createBuffers(gl, obj.indices, obj.vertCoord, obj.textCoord,
          obj.normCoord);
    }

    void handleError(Object err) {
      print("Model.fromOBJ: failure fetching OBJ from URL=$_URL: $err");
    }

    HttpRequest.getString(_URL).then(handleResponse).catchError(handleError);
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

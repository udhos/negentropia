part of shader;

class Instance {
  
  Model model;
  Vector3 center;
  double scale;
  Float32List pickColor;

  Matrix4 MV = new Matrix4.identity(); // model-view matrix
  
  Instance(this.model, this.center, this.scale, [this.pickColor=null]);
  
  void update(GameLoopHtml gameLoop) {
  }
  
  void draw(GameLoopHtml gameLoop, ShaderProgram prog, Camera cam) {
    
    setViewMatrix(MV, cam.eye, cam.center, cam.up);
    
    MV.translate(center[0], center[1], center[2]);
    
    MV.scale(scale, scale, scale);
    
    RenderingContext gl = prog.gl;

    gl.uniformMatrix4fv(prog.u_MV, false, MV.storage);

    gl.bindBuffer(RenderingContext.ARRAY_BUFFER, model.vertexPositionBuffer);
    gl.vertexAttribPointer(prog.a_Position, model.vertexPositionBufferItemSize, RenderingContext.FLOAT, false, 0, 0);
  
    gl.bindBuffer(RenderingContext.ELEMENT_ARRAY_BUFFER, model.vertexIndexBuffer);
    
    model.pieceList.forEach((Piece piece) {
        gl.drawElements(RenderingContext.TRIANGLES, piece.vertexIndexLength, RenderingContext.UNSIGNED_SHORT,
        piece.vertexIndexOffset * model.vertexIndexBufferItemSize);
      });
  }
}

class Piece {
  int vertexIndexOffset;
  int vertexIndexLength;  
  
  Piece(this.vertexIndexOffset, this.vertexIndexLength);
}

class Model {
    
  Buffer vertexPositionBuffer;
  Buffer vertexIndexBuffer;
  final int vertexPositionBufferItemSize = 3; // coord x,y,z
  final int vertexIndexBufferItemSize = 2; // size of Uint16Array
  
  bool modelReady = false;  // buffers
  bool piecesReady = false; // multiple OBJ pieces
  
  String _name;
  String get modelName => _name;
  
  List<Piece> pieceList = new List<Piece>();
  List<Instance> instanceList = new List<Instance>();
  
  void _createBuffers(RenderingContext gl, List<int> indices, List<double> vertCoord, List<double> textCoord, List<double> normCoord) {
        
    vertexPositionBuffer = gl.createBuffer();
    gl.bindBuffer(RenderingContext.ARRAY_BUFFER, vertexPositionBuffer);
    gl.bufferDataTyped(RenderingContext.ARRAY_BUFFER, new Float32List.fromList(vertCoord), RenderingContext.STATIC_DRAW);
    
    vertexIndexBuffer = gl.createBuffer();
    gl.bindBuffer(RenderingContext.ELEMENT_ARRAY_BUFFER, vertexIndexBuffer);
    gl.bufferDataTyped(RenderingContext.ELEMENT_ARRAY_BUFFER, new Uint16List.fromList(indices), RenderingContext.STATIC_DRAW);
        
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
      }
      catch (e) {
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

    HttpRequest.getString(URL)
    .then(handleResponse)
    .catchError(handleError);
  }
  
  void loadObj(RenderingContext gl, Obj o) {
    assert(!piecesReady);
    
    o.partList.forEach((Part pa) {
      Piece pi = addPiece(pa.indexFirst, pa.indexListSize);
      //print("Model.fromOBJ: added part ${pa.name} into piece: offset=${pi.vertexIndexOffset} length=${pi.vertexIndexLength}");
    });
    
    piecesReady = true;
  }
  
  Model.fromOBJ(RenderingContext gl, String URL) {

    void handleResponse(String response) {
      //print("Model.fromOBJ: fetched OBJ from URL: $URL");
      
      Obj obj = new Obj.fromString(URL, response);
      
      _name = URL;
      
      loadObj(gl, obj);
      
      _createBuffers(gl, obj.indices, obj.vertCoord, obj.textCoord, obj.normCoord);
    }

    void handleError(Object err) {
      print("Model.fromOBJ: failure fetching OBJ from URL: $URL: $err");
    }

    HttpRequest.getString(URL)
    .then(handleResponse)
    .catchError(handleError);    
  }

  void addInstance(Instance i) {
    this.instanceList.add(i);
  }
 
  void drawInstances(GameLoopHtml gameLoop, ShaderProgram program, Camera cam) {
    if (!modelReady || !piecesReady) {
      return;
    }
    
    this.instanceList.forEach((Instance i) => i.draw(gameLoop, program, cam));
  }  

  void update(GameLoopHtml gameLoop) {
    this.instanceList.forEach((Instance i) => i.update(gameLoop));
  }  

}

part of shader;

class TexShaderProgram extends ShaderProgram {

  int a_TextureCoord;
  UniformLocation u_Sampler;

  TexShaderProgram(RenderingContext gl, String programName)
      : super(gl, programName);

  /*
  void initContext(RenderingContext gl, Map<String,Texture> textureTable) {
    modelList.forEach((TexModel m) => m.initContext(gl, textureTable));
  }
  */

  void getLocations() {
    super.getLocations();

    a_TextureCoord = gl.getAttribLocation(program, "a_TextureCoord");
    u_Sampler = gl.getUniformLocation(program, "u_Sampler");
  }

  void drawModels(GameLoopHtml gameLoop, Camera cam, Matrix4 pMatrix) {

    if (!shaderReady) {
      return;
    }

    gl.useProgram(program);
    gl.enableVertexAttribArray(a_Position);
    gl.enableVertexAttribArray(a_TextureCoord);

    // send perspective projection matrix uniform
    gl.uniformMatrix4fv(u_P, false, pMatrix.storage);

    modelList.forEach((m) => m.drawInstances(gameLoop, this, cam));

    //gl.disableVertexAttribArray(a_Position); // needed ??
  }

}

class TexPiece extends Piece {

  TextureInfo texInfo;

  TexPiece(int indexOffset, int indexLength) : super(indexOffset, indexLength);
}

class TexModel extends Model {

  Buffer textureCoordBuffer;
  final int textureCoordBufferItemSize = 2; // coord s,t
  Asset asset;
  Map<String, Texture> textureTable;
  final int textureUnit = 1;
  int globeIndexSize;

  /*
  void initContext(RenderingContext gl, Map<String,Texture> textureTable) {
    
    //textureInfoList.forEach((TextureInfo ti) => ti.forceCreateTexture(gl, textureTable));
    pieceList.forEach((TexPiece pi) => pi.texInfo.loadTexture(gl, textureTable));
    
  }
  */

  // redefine _createBuffers() used by Model's constructor
  void _createBuffers(RenderingContext gl, List<int> indices,
      List<double> vertCoord, List<double> textCoord, List<double> normCoord) {

    assert(!modelReady);
    
    log("TexModel._createBuffers model=$modelName");

    textureCoordBuffer = gl.createBuffer();
    gl.bindBuffer(RenderingContext.ARRAY_BUFFER, textureCoordBuffer);
    gl.bufferDataTyped(
        RenderingContext.ARRAY_BUFFER,
        new Float32List.fromList(textCoord),
        RenderingContext.STATIC_DRAW);

    super._createBuffers(gl, indices, vertCoord, textCoord, normCoord);

    assert(modelReady);
  }

  void loadObj(RenderingContext gl, Obj obj) {

    assert(obj != null);

    String mtlURL = "${asset.mtl}/${obj.mtllib}";

    void onMtlLibLoaded(String materialResponse) {

      assert(!piecesReady);

      Map<String, Material> lib = mtllib_parse(materialResponse, mtlURL);
      assert(lib != null);

      int i = 0;

      obj.partList.forEach((pa) {

        String usemtl = pa.usemtl;

        Material mtl = lib[usemtl];
        if (mtl == null) {
          err(
              "loadObj $i: material usemtl=$usemtl NOT FOUND on mtllib=$mtlURL");
        }

        int r = (mtl.Kd[0] * 255.0).round();
        int g = (mtl.Kd[1] * 255.0).round();
        int b = (mtl.Kd[2] * 255.0).round();
        List<int> temporaryColor = [r, g, b, 255];

        String texFile = mtl.map_Kd;

        String textureURL;
        if (texFile != null) {
          textureURL = "${asset.texture}/$texFile";
        }

        TextureInfo texInfo =
            new TextureInfo(gl, textureTable, textureURL, temporaryColor, textureUnit);

        addTexture(pa.indexFirst, pa.indexListSize, texInfo);

        ++i;
      });

      piecesReady = true;
    }

    if (obj.mtllib == null) {
      err("loadObj: model=$modelName undefined OBJ mtllib URL");
      return;
    }

    HttpRequest.getString(mtlURL).then(onMtlLibLoaded).catchError((e) {
      err("loadObj: failure fetching mtllib: $mtlURL: $e");
    });
  }

  void saveIndexSize(int indexSize) {
    globeIndexSize = indexSize; // saves the index size
  }

  TexModel.fromOBJ(RenderingContext gl, String name, String URL, Vector3 front,
      Vector3 up, this.textureTable, this.asset)
      : super.fromOBJ(gl, name, URL, front, up);

  TexModel.fromGlobe(RenderingContext gl, String name, double radius,
      String textureURL, Vector3 front, Vector3 up, this.textureTable, this.asset)
      : super.fromGlobe(gl, name, radius, front, up) {

    log(
        "TexModel.fromGlobe: model=$modelName tex=$textureURL front=$_front up=$_up");   

    assert(!piecesReady);
    assert(pieceList.length == 0);
    
    List<int> temporaryColor = [127, 127, 127, 255];

    TextureInfo texInfo =
        new TextureInfo(gl, textureTable, textureURL, temporaryColor, textureUnit);

    assert(globeIndexSize != null);
    assert(globeIndexSize > 0);
    addTexture(0, globeIndexSize, texInfo);
    
    assert(pieceList.length == 1);
    assert(pieceList.first is TexPiece);
    assert((pieceList.first as TexPiece).texInfo != null);
    assert((pieceList.first as TexPiece).texInfo == texInfo);
    
    piecesReady = true;

    //DEBUG:
    TexPiece tp = pieceList.first as TexPiece;
    log("TexModel.fromGlobe vertexPositionBufferItemSize=$vertexPositionBufferItemSize");
    log("TexModel.fromGlobe textureCoordBufferItemSize=$textureCoordBufferItemSize");
    log("TexModel.fromGlobe vertexIndexBufferItemSize=$vertexIndexBufferItemSize");
    log("TexModel.fromGlobe piece vertexIndexOffset=${tp.vertexIndexOffset}");
    log("TexModel.fromGlobe piece vertexIndexLength=${tp.vertexIndexLength}");
  }

  Piece addPiece(int offset, int length) {
    Piece pi = new TexPiece(offset, length);
    pieceList.add(pi);
    return pi;
  }

  void addTexture(int indexOffset, int indexLength, TextureInfo tex) {
    TexPiece pi = addPiece(indexOffset, indexLength) as TexPiece;
    pi.texInfo = tex;
  }

  void drawInstances(GameLoopHtml gameLoop, ShaderProgram program, Camera cam) {
    if (!modelReady || !piecesReady) {
      return;
    }

    RenderingContext gl = program.gl;

    // vertex coord
    gl.bindBuffer(RenderingContext.ARRAY_BUFFER, vertexPositionBuffer);
    gl.vertexAttribPointer(
        program.a_Position,
        vertexPositionBufferItemSize,
        RenderingContext.FLOAT,
        false,
        0,
        0);

    // texture coord
    gl.bindBuffer(RenderingContext.ARRAY_BUFFER, textureCoordBuffer);
    gl.vertexAttribPointer(
        (program as TexShaderProgram).a_TextureCoord,
        textureCoordBufferItemSize,
        RenderingContext.FLOAT,
        false,
        0,
        0);

    gl.bindBuffer(RenderingContext.ELEMENT_ARRAY_BUFFER, vertexIndexBuffer);

    instanceList.forEach((i) => i.draw(gameLoop, program, cam));
  }

}

class TexInstance extends Instance {

  TexInstance(id, TexModel model, Vector3 center, double scale,
      [Float32List pick = null])
      : super(id, model, center, scale, pick);

  void draw(GameLoopHtml gameLoop, ShaderProgram prog, Camera cam) {

    RenderingContext gl = prog.gl;

    uploadModelView(gl, prog.u_MV, cam, scale); // set up MV matrix

    (model as TexModel).pieceList.forEach((pi) {

      TexPiece tp = pi as TexPiece;
      TextureInfo ti = tp.texInfo;

      // set texture sampler
      int unit = (model as TexModel).textureUnit;
      gl.activeTexture(RenderingContext.TEXTURE0 + unit);
      gl.bindTexture(RenderingContext.TEXTURE_2D, ti.texture);
      gl.uniform1i((prog as TexShaderProgram).u_Sampler, unit);

      gl.drawElements(
          RenderingContext.TRIANGLES,
          tp.vertexIndexLength,
          RenderingContext.UNSIGNED_SHORT,
          tp.vertexIndexOffset * model.vertexIndexBufferItemSize);
    });

  }

}

part of shader;

class TexShaderProgram extends ShaderProgram {
  
  int a_TextureCoord;
  UniformLocation u_Sampler;
   
  TexShaderProgram(RenderingContext gl, String programName) : super(gl, programName);
  
  /*
  void initContext(RenderingContext gl, Map<String,Texture> textureTable) {
    modelList.forEach((TexModel m) => m.initContext(gl, textureTable));
  }
  */

  void getLocations() {
    super.getLocations();

    a_TextureCoord = gl.getAttribLocation(program, "a_TextureCoord");
    u_Sampler      = gl.getUniformLocation(program, "u_Sampler");
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
    
    modelList.forEach((TexModel m) => m.drawInstances(gameLoop, this, cam));

    // clean up
    gl.bindBuffer(RenderingContext.ARRAY_BUFFER, null);
    gl.bindBuffer(RenderingContext.ELEMENT_ARRAY_BUFFER, null);
    gl.bindTexture(RenderingContext.TEXTURE_2D, null);
    
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
  Map<String,Texture> textureTable;

  /*
  void initContext(RenderingContext gl, Map<String,Texture> textureTable) {
    
    //textureInfoList.forEach((TextureInfo ti) => ti.forceCreateTexture(gl, textureTable));
    pieceList.forEach((TexPiece pi) => pi.texInfo.loadTexture(gl, textureTable));
    
  }
  */
  
  // redefine _createBuffers() used by Model's constructor
  void _createBuffers(RenderingContext gl, List<int> indices, List<double> vertCoord, List<double> textCoord, List<double> normCoord) {

    assert(!modelReady);
    
    textureCoordBuffer = gl.createBuffer();
    gl.bindBuffer(RenderingContext.ARRAY_BUFFER, textureCoordBuffer);
    gl.bufferDataTyped(RenderingContext.ARRAY_BUFFER, new Float32List.fromList(textCoord), RenderingContext.STATIC_DRAW);

    super._createBuffers(gl, indices, vertCoord, textCoord, normCoord);
    
    assert(modelReady);
}

  void loadObj(RenderingContext gl, Obj obj) {
    
    String mtlURL = "${asset.mtl}/${obj.mtllib}";

    void onMtlLibLoaded(String materialResponse) {
      
      assert(!piecesReady);

      Map<String,Material> lib = mtllib_parse(materialResponse, mtlURL);
      assert(lib != null);
      
      int i = 0;

      obj.partList.forEach((Part pa) {

        String usemtl = pa.usemtl;      
        
        Material mtl = lib[usemtl];
        if (mtl == null) {
          print("loadObj $i: material usemtl=$usemtl NOT FOUND on mtllib=$mtlURL");
          return;
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
        
        TextureInfo texInfo = new TextureInfo(gl, textureTable, textureURL, temporaryColor);
        
        addTexture(pa.indexFirst, pa.indexListSize, texInfo);
        
        ++i;
      });
      
      piecesReady = true;
      
      //print("loadObj: ${obj.partList.length} parts fed into ${pieceList.length} pieces");
    }

    HttpRequest.getString(mtlURL)
    .then(onMtlLibLoaded)
    .catchError((err) { print("loadObj: failure fetching mtllib: $mtlURL: $err"); });    
  }
  
  TexModel.fromOBJ(RenderingContext gl, String URL,
      this.textureTable, this.asset): super.fromOBJ(gl, URL);

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
    gl.vertexAttribPointer(program.a_Position, vertexPositionBufferItemSize, RenderingContext.FLOAT, false, 0, 0);
    
    // texture coord
    gl.bindBuffer(RenderingContext.ARRAY_BUFFER, textureCoordBuffer);
    gl.vertexAttribPointer((program as TexShaderProgram).a_TextureCoord, textureCoordBufferItemSize, RenderingContext.FLOAT, false, 0, 0);
    
    gl.bindBuffer(RenderingContext.ELEMENT_ARRAY_BUFFER, vertexIndexBuffer);
    
    instanceList.forEach((Instance i) => i.draw(gameLoop, program, cam));
  }  
  
}

class TexInstance extends Instance {
    
  TexInstance(TexModel model, Vector3 center, double scale, [Float32List pick=null]) : super(model, center, scale, pick);

  void draw(GameLoopHtml gameLoop, ShaderProgram prog, Camera cam) {

    setViewMatrix(MV, cam.eye, cam.center, cam.up);
    
    MV.translate(center[0], center[1], center[2]);
    
    MV.scale(scale, scale, scale);
    
    RenderingContext gl = prog.gl;

    gl.uniformMatrix4fv(prog.u_MV, false, MV.storage);
    
    (model as TexModel).pieceList.forEach((Piece pi) {
      
      TexPiece tp = pi as TexPiece;
      TextureInfo ti = tp.texInfo;
      
      // set texture sampler
      int unit = 0;
      gl.activeTexture(RenderingContext.TEXTURE0 + unit);
      gl.bindTexture(RenderingContext.TEXTURE_2D, ti.texture);
      gl.uniform1i((prog as TexShaderProgram).u_Sampler, unit);
      
      gl.drawElements(RenderingContext.TRIANGLES, tp.vertexIndexLength, RenderingContext.UNSIGNED_SHORT,
        tp.vertexIndexOffset * model.vertexIndexBufferItemSize);
    });
    
  }

}
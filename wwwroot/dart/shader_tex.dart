part of shader;

class TexShaderProgram extends ShaderProgram {
  
  int a_TextureCoord;
  UniformLocation u_Sampler;
  //UniformLocation u_Color;  
   
  TexShaderProgram(RenderingContext gl) : super(gl);
  
  void initContext(RenderingContext gl, Map<String,Texture> textureTable) {
    modelList.forEach((TexModel m) => m.initContext(gl, textureTable));
  }

  void getLocations() {
    super.getLocations();

    a_TextureCoord = gl.getAttribLocation(program, "a_TextureCoord");
    u_Sampler      = gl.getUniformLocation(program, "u_Sampler");
    //u_Color        = gl.getUniformLocation(program, "u_Color");
    
    print("TexShaderProgram: locations ready");      
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
    
    /*
    // fallback solid color for textured objects
    List<double> white = [1.0, 1.0, 1.0, 1.0]; // neutral color in multiplication
    gl.uniform4fv(u_Color, new Float32List.fromList(white));
    */    

    modelList.forEach((TexModel m) => m.drawInstances(gameLoop, cam));

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
  int textureCoordBufferItemSize;
  Asset asset;
  Map<String,Texture> textureTable;

  //List<TextureInfo> textureInfoList = new List<TextureInfo>();
  
  void initContext(RenderingContext gl, Map<String,Texture> textureTable) {
    
    //textureInfoList.forEach((TextureInfo ti) => ti.forceCreateTexture(gl, textureTable));
    pieceList.forEach((TexPiece pi) => pi.texInfo.forceCreateTexture(gl, textureTable));
    
  }
  
  // redefine _createBuffers() used by Model's constructor
  void _createBuffers(RenderingContext gl, List<int> indices, List<double> vertCoord, List<double> textCoord, List<double> normCoord) {
            
    textureCoordBuffer = gl.createBuffer();
    gl.bindBuffer(RenderingContext.ARRAY_BUFFER, textureCoordBuffer);
    gl.bufferData(RenderingContext.ARRAY_BUFFER, new Float32List.fromList(textCoord), RenderingContext.STATIC_DRAW);
    textureCoordBufferItemSize = 2; // coord s,t

    super._createBuffers(gl, indices, vertCoord, textCoord, normCoord);
}

  void loadObj(RenderingContext gl, Obj obj) {
    
    String mtlURL = "${asset.mtl}/${obj.mtllib}";

    void onMtlLibLoaded(String materialResponse) {

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
        //print("loadObj $i: usemtl=$usemtl map_Kd=$texFile textureURL=$textureURL");
        
        TextureInfo texInfo = new TextureInfo(gl, textureTable, textureURL, temporaryColor);
        
        addTexture(pa.indexFirst, pa.indexListSize, texInfo);
        
        ++i;
      });
      
      print("loadObj: ${obj.partList.length} parts fed into ${pieceList.length} pieces");
    }

    HttpRequest.getString(mtlURL)
    .then(onMtlLibLoaded)
    .catchError((err) { print("loadObj: failure fetching mtllib: $mtlURL: $err"); });    
  }
  
  TexModel.fromOBJ(RenderingContext gl, ShaderProgram program, String URL,
      this.textureTable, this.asset,
      [void onDone(RenderingContext gl, TexModel m, Obj o, String u)]):
        super.fromOBJ(gl, program, URL, onDone);

  Piece addPiece(int offset, int length) {
    Piece pi = new TexPiece(offset, length);
    pieceList.add(pi);
    return pi;
  }
  
  void addTexture(int indexOffset, int indexLength, TextureInfo tex) {
    TexPiece pi = addPiece(indexOffset, indexLength) as TexPiece;
    pi.texInfo = tex;
    //print("addTexture: offset=${pi.vertexIndexOffset} length=${pi.vertexIndexLength}");
  }

  void drawInstances(GameLoopHtml gameLoop, Camera cam) {
    
    RenderingContext gl = program.gl;

    // vertex coord
    gl.bindBuffer(RenderingContext.ARRAY_BUFFER, vertexPositionBuffer);
    gl.vertexAttribPointer(program.a_Position, vertexPositionBufferItemSize, RenderingContext.FLOAT, false, 0, 0);
    
    // texture coord
    gl.bindBuffer(RenderingContext.ARRAY_BUFFER, textureCoordBuffer);
    gl.vertexAttribPointer((program as TexShaderProgram).a_TextureCoord, textureCoordBufferItemSize, RenderingContext.FLOAT, false, 0, 0);
    
    gl.bindBuffer(RenderingContext.ELEMENT_ARRAY_BUFFER, vertexIndexBuffer);
    
    instanceList.forEach((Instance i) => i.draw(gameLoop, cam));
  }  
  
}

class TexInstance extends Instance {
    
  TexInstance(TexModel model, Vector3 center, double scale) : super(model, center, scale);

  void draw(GameLoopHtml gameLoop, Camera cam) {

    setViewMatrix(MV, cam.eye, cam.center, cam.up);
    
    MV.translate(center[0], center[1], center[2]);
    
    MV.scale(scale, scale, scale);
    
    ShaderProgram prog = model.program;
    RenderingContext gl = prog.gl;

    gl.uniformMatrix4fv(prog.u_MV, false, MV.storage);
    
    /*
    (model as TexModel).textureInfoList.forEach((TextureInfo ti) {
      
      // set texture sampler
      int unit = 1;
      gl.activeTexture(RenderingContext.TEXTURE0 + unit);
      gl.bindTexture(RenderingContext.TEXTURE_2D, ti.texture);
      gl.uniform1i((prog as TexShaderProgram).u_Sampler, unit);
      
      gl.drawElements(RenderingContext.TRIANGLES, ti.indexNumber, RenderingContext.UNSIGNED_SHORT,
        ti.indexOffset * model.vertexIndexBufferItemSize);
    });
    */
    (model as TexModel).pieceList.forEach((Piece pi) {
      
      TexPiece tp = pi as TexPiece;
      TextureInfo ti = tp.texInfo;
      
      // set texture sampler
      int unit = 1;
      gl.activeTexture(RenderingContext.TEXTURE0 + unit);
      gl.bindTexture(RenderingContext.TEXTURE_2D, ti.texture);
      gl.uniform1i((prog as TexShaderProgram).u_Sampler, unit);
      
      gl.drawElements(RenderingContext.TRIANGLES, tp.vertexIndexLength, RenderingContext.UNSIGNED_SHORT,
        tp.vertexIndexOffset * model.vertexIndexBufferItemSize);
    });
    
  }

}
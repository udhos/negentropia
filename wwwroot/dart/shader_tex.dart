part of shader;

class TexShaderProgram extends ShaderProgram {
  
  int a_TextureCoord;
  UniformLocation u_Sampler;
  UniformLocation u_Color;  
   
  TexShaderProgram(RenderingContext gl) : super(gl);
  
  void initContext(RenderingContext gl, Map<String,Texture> textureTable) {
    modelList.forEach((TexModel m) => m.initContext(gl, textureTable));
  }

  void getLocations() {
    super.getLocations();

    a_TextureCoord = gl.getAttribLocation(program, "a_TextureCoord");
    u_Sampler      = gl.getUniformLocation(program, "u_Sampler");
    u_Color        = gl.getUniformLocation(program, "u_Color");
    
    print("TexShaderProgram: locations ready");      
  }

  void drawModels(GameLoopHtml gameLoop, Camera cam, mat4 pMatrix) {
    
    gl.useProgram(program);
    gl.enableVertexAttribArray(a_Position);
    gl.enableVertexAttribArray(a_TextureCoord);

    // send perspective projection matrix uniform
    gl.uniformMatrix4fv(u_P, false, pMatrix.storage);
    
    // fallback solid color for textured objects
    List<double> white = [1.0, 1.0, 1.0, 1.0]; // neutral color in multiplication
    gl.uniform4fv(u_Color, new Float32List.fromList(white));    

    modelList.forEach((TexModel m) => m.drawInstances(gameLoop, cam));

    // clean up
    gl.bindBuffer(RenderingContext.ARRAY_BUFFER, null);
    gl.bindBuffer(RenderingContext.ELEMENT_ARRAY_BUFFER, null);
    gl.bindTexture(RenderingContext.TEXTURE_2D, null);
    
    //gl.disableVertexAttribArray(a_Position); // needed ??
  }
  
}

class TexModel extends Model {
  
  Buffer textureCoordBuffer;
  int textureCoordBufferItemSize;  

  List<TextureInfo> textureInfoList = new List<TextureInfo>();
  
  void initContext(RenderingContext gl, Map<String,Texture> textureTable) {
    
    textureInfoList.forEach((TextureInfo ti) => ti.forceCreateTexture(gl, textureTable));
    
  }
  
  void _createBuffers(RenderingContext gl, List<int> indices, List<double> vertCoord, List<double> textCoord, List<double> normCoord) {
        
    vertexPositionBuffer = gl.createBuffer();
    gl.bindBuffer(RenderingContext.ARRAY_BUFFER, vertexPositionBuffer);
    gl.bufferData(RenderingContext.ARRAY_BUFFER, new Float32List.fromList(vertCoord), RenderingContext.STATIC_DRAW);
    vertexPositionBufferItemSize = 3; // coord x,y,z
    
    textureCoordBuffer = gl.createBuffer();
    gl.bindBuffer(RenderingContext.ARRAY_BUFFER, textureCoordBuffer);
    gl.bufferData(RenderingContext.ARRAY_BUFFER, new Float32List.fromList(textCoord), RenderingContext.STATIC_DRAW);
    textureCoordBufferItemSize = 2; // coord s,t
    
    vertexIndexBuffer = gl.createBuffer();
    gl.bindBuffer(RenderingContext.ELEMENT_ARRAY_BUFFER, vertexIndexBuffer);
    gl.bufferData(RenderingContext.ELEMENT_ARRAY_BUFFER, new Uint16List.fromList(indices), RenderingContext.STATIC_DRAW);
    vertexIndexBufferItemSize = 2; // size of Uint16Array
    
    vertexIndexLength = indices.length;
    
    print("TexModel._createBuffers: vertex index length: ${vertexIndexLength}");
    
    // clean-up
    gl.bindBuffer(RenderingContext.ARRAY_BUFFER, null);
    gl.bindBuffer(RenderingContext.ELEMENT_ARRAY_BUFFER, null);
  }

  /*
  TexModel.fromOBJ(RenderingContext gl, ShaderProgram program, String URL) :
    super.fromOBJ(gl, program, URL);
  */  

  TexModel.fromOBJ(RenderingContext gl, ShaderProgram program, String URL) {

    this.program = program;
    
    void handleResponse(String response) {
      print("TexModel.fromOBJ: fetched OBJ from URL: $URL");
      
      Obj obj = new Obj.fromString(URL, response);
      
      _createBuffers(gl, obj.indices, obj.vertCoord, obj.textCoord, obj.normCoord);
    }

    void handleError(Object err) {
      print("TexModel.fromOBJ: failure fetching OBJ from URL: $URL: $err");
    }

    HttpRequest.getString(URL)
    .then(handleResponse)
    .catchError(handleError);    
  }

  void addTexture(TextureInfo tex) {
    textureInfoList.add(tex); 
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
    
  TexInstance(TexModel model, vec3 center, double scale) : super(model, center, scale);

  void draw(GameLoopHtml gameLoop, Camera cam) {

    setViewMatrix(MV, cam.eye, cam.center, cam.up);
    
    MV.translate(center[0], center[1], center[2]);
    
    MV.scale(scale, scale, scale);
    
    ShaderProgram prog = model.program;
    RenderingContext gl = prog.gl;

    gl.uniformMatrix4fv(prog.u_MV, false, MV.storage);
    
    (model as TexModel).textureInfoList.forEach((TextureInfo ti) {
      
      // set texture sampler
      int unit = 1;
      gl.activeTexture(RenderingContext.TEXTURE0 + unit);
      gl.bindTexture(RenderingContext.TEXTURE_2D, ti.texture);
      gl.uniform1i((prog as TexShaderProgram).u_Sampler, unit);
      
      gl.drawElements(RenderingContext.TRIANGLES, ti.indexNumber, RenderingContext.UNSIGNED_SHORT,
        ti.indexOffset * model.vertexIndexBufferItemSize);
    });
    
  }

}
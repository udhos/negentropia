
part of shader;

class TexShaderProgram extends ShaderProgram {
   
  TexShaderProgram(RenderingContext gl) : super(gl);

}

class TexModel extends Model {

  List<TextureInfo> textureList = new List<TextureInfo>();
  
  TexModel.fromOBJ(RenderingContext gl, ShaderProgram program, String URL) :
    super.fromOBJ(gl, program, URL);

  void addTexture(TextureInfo tex) {
    textureList.add(tex); 
  }

}

class TexInstance extends Instance {
    
  TexInstance(TexModel model, vec3 center, double scale) : super(model, center, scale);

}
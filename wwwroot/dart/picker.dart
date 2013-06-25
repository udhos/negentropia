part of shader;

class PickerShader extends ShaderProgram {

  UniformLocation u_Color;
  List<ShaderProgram> programList;

  PickerShader(RenderingContext gl, this.programList) : super(gl) {
    /*
    // move to draw: scan programs/models/instances
    programList.forEach((p) {
      if (identical(p, this)) {
        return; // skip myself
      }
      if (p is! TexShaderProgram) {
        return;
      }
      if (p.modelList.isEmpty) {
        return;
      }
      p.modelList.forEach((m) {
        if (m is! TexModel) {
          return;
        }
        if (m.instanceList.isEmpty) {
          return;
        }
        p.modelList.add(new PickerModel(m));
        m.instanceList.forEach((i) {
          if (i is! TexInstance) {
            return;
          }
          m.instanceList.add(new PickerInstance(i));
        });
      });
    });
    */
  }
  
  void getLocations() {
    super.getLocations();

    u_Color = gl.getUniformLocation(program, "u_Color");
    
    print("PickerShader: locations ready");      
  }
  
}
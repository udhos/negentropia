part of shader;

class PickerShader extends ShaderProgram {

  UniformLocation u_Color;  

  PickerShader(RenderingContext gl) : super(gl);
  
  void getLocations() {
    super.getLocations();

    u_Color = gl.getUniformLocation(program, "u_Color");
    
    print("PickerShader: locations ready");      
  }
  
}
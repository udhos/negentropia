part of shader;

class PickerInstance extends Instance {
  
  PickerInstance(Instance i): super(i.model, i.center, i.scale, i.pickColor);

  // the whole purpose of this class is to redefine the draw() method
  void draw(GameLoopHtml gameLoop, ShaderProgram prog, Camera cam) {
    RenderingContext gl = prog.gl;
    gl.uniform4fv((prog as PickerShader).u_Color, pickColor);
    super.draw(gameLoop, prog, cam);
  }
}

/*
// New Model not needed 
class PickerModel extends Model {
  PickerModel(Model m) : super.copy(m);
}
*/

class PickerShader extends ShaderProgram {

  UniformLocation u_Color;
  List<ShaderProgram> programList;
  List<PickerInstance> instanceList = new List<PickerInstance>();

  PickerShader(RenderingContext gl, this.programList) : super(gl) {

    programList.forEach((p) {
      p.modelList.forEach((m) {
        m.instanceList.where((i) => i.pickColor != null).forEach((ii) {
          PickerInstance pi = new PickerInstance(ii);
          instanceList.add(pi);
        });
      });
    });
    
  }
  
  void getLocations() {
    super.getLocations();

    u_Color = gl.getUniformLocation(program, "u_Color");
    
    print("PickerShader: locations ready");      
  }
  
  void drawModels(GameLoopHtml gameLoop, Camera cam, Matrix4 pMatrix) {
        
    gl.useProgram(program);
    gl.enableVertexAttribArray(a_Position);

    // send perspective projection matrix uniform
    gl.uniformMatrix4fv(this.u_P, false, pMatrix.storage);

    instanceList.forEach((i) => i.draw(gameLoop, this, cam));

    // clean up
    gl.bindBuffer(RenderingContext.ARRAY_BUFFER, null);
    gl.bindBuffer(RenderingContext.ELEMENT_ARRAY_BUFFER, null);
    
    //gl.disableVertexAttribArray(a_Position); // needed ??
  }
  
}
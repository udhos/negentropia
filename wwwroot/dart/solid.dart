library solid;

import 'dart:web_gl';

import 'shader.dart';

class SolidShader extends ShaderProgram {
  
  UniformLocation u_Color;
  
  SolidShader(RenderingContext gl, String programName) : super(gl, programName);

  void getLocations() {
    super.getLocations();

    u_Color = gl.getUniformLocation(program, "u_Color");
  }

}
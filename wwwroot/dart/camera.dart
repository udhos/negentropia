library camera;

import 'dart:math' as math;

import 'package:vector_math/vector_math.dart';
import 'package:game_loop/game_loop_html.dart';

class Camera {
  final double degreesPerSec = 30.0;
  final double camOrbitRadius = 10.0;
  vec3 eye, center, up;
  double oldAngle, angle;
  
  Camera(this.eye, this.center, this.up);
  
  double get rad => _getRad(1.0);
  
  double _getRad(double interpolation) {
    double deg = interpolation * angle + (1 - interpolation) * oldAngle;
    //double deg = interpolation * (angle - oldAngle) + oldAngle;
    double r = deg * math.PI / 180.0;
    return r;
  }
  
  void update(GameLoopHtml gameLoop) {
    oldAngle = angle;
    angle = gameLoop.gameTime * this.degreesPerSec % 360.0;
  }
    
  void render(GameLoopHtml gameLoop) {
    
    double r = _getRad(gameLoop.renderInterpolationFactor);
    
    eye[0] = camOrbitRadius * math.sin(r);
    eye[2] = camOrbitRadius * math.cos(r);
  }
}

library camera;

import 'dart:math' as math;

import 'package:vector_math/vector_math.dart';
import 'package:game_loop/game_loop_html.dart';

class Camera {
  final double degreesPerSec = 30.0;
  final double camOrbitRadius = 15.0;
  vec3 eye, center, up;
  double oldAngle, angle;
  
  Camera(this.eye, this.center, this.up);
  
  //double get rad => _getRad(1.0);
  double get rad => getRad(0.0);
  
  double getRad(double interpolation) {
    //double deg = interpolation * angle + (1 - interpolation) * oldAngle;
    double deg;
    if (angle > oldAngle) {
      deg = interpolation * (angle       - oldAngle) + oldAngle;
    } else {
      // undo modulo 360 for correct interpolation
      deg = interpolation * (angle + 360 - oldAngle) + oldAngle;
    }
    double r = deg * math.PI / 180.0;
    return r;
  }
  
  void update(GameLoopHtml gameLoop) {
    oldAngle = angle;
    angle = gameLoop.gameTime * this.degreesPerSec % 360;
  }
    
  void render(GameLoopHtml gameLoop) {
    
    double r = getRad(gameLoop.renderInterpolationFactor);
    
    eye[0] = camOrbitRadius * math.sin(r);
    eye[2] = camOrbitRadius * math.cos(r);
  }
}

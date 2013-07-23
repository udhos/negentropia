library camera;

import 'dart:math' as math;

import 'package:vector_math/vector_math.dart';
import 'package:game_loop/game_loop_html.dart';

class Camera {
  final double degreesPerSec = 45.0;
  final double camOrbitRadius = 15.0;
  Vector3 eye, center, up;
  double oldAngle = 0.0;
  double angle = 0.0;
  
  Camera(this.eye, this.center, this.up) {
    _update(0.0);
  }
  
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
  
  double _update(double gameTime) {
    oldAngle = angle;
    angle = gameTime * this.degreesPerSec % 360;    
  }
   
  void update(GameLoopHtml gameLoop) {
    _update(gameLoop.gameTime);
  }
    
  void render(GameLoopHtml gameLoop) {
    
    double r = getRad(gameLoop.renderInterpolationFactor);
    
    eye[0] = camOrbitRadius * math.sin(r);
    eye[2] = camOrbitRadius * math.cos(r);
  }
}

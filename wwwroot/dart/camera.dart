library camera;

import 'dart:math' as math;
import 'interpolate.dart';

import 'package:vector_math/vector_math.dart';

class Camera {
  final double degreesPerSec = 20.0;
  //final double camOrbitRadius = 15.0;
  
  Quaternion _orientation = new Quaternion.identity();
  Vector3 _position = new Vector3.zero();
  
  double _oldAngle = 0.0;
  double _angle = 0.0;
  
  Camera(this._position) {
    update(0.0);
  }
  
  void rotate(Matrix4 MV) {
    MV.setRotation(_orientation.asRotationMatrix());
  }
  
  void translate(Matrix4 MV) {
    MV.translate(- _position);
  }
  
  void update(double gameTime) {
    _oldAngle = _angle;
    _angle = gameTime * this.degreesPerSec % 360.0;    
  }  
  
  /*
  double _getRad(double interpolation) {
    double deg;
    if (_angle > _oldAngle) {
      deg = interpolation * (_angle         - _oldAngle) + _oldAngle;
    } else {
      // undo modulo 360 for correct interpolation
      deg = interpolation * (_angle + 360.0 - _oldAngle) + _oldAngle;
    }
    double r = deg * math.PI / 180.0;
    return r;
  }
  */

  static final Vector3 Y = new Vector3(0.0, 1.0, 0.0);

  void render(double renderInterpolationFactor) {
    
    //double r = _getRad(renderInterpolationFactor);
    double rad = interpolateDegree(_angle, _oldAngle, renderInterpolationFactor) * math.PI / 180.0;

    // FIXME: should apply a rotation quaternion
    _orientation = new Quaternion.axisAngle(Y, rad).conjugated();
  }
}

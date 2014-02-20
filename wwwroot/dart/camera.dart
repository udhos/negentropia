library camera;

//import 'dart:math' as math;

import 'package:vector_math/vector_math.dart';

//import 'interpolate.dart';
import 'logg.dart';
import 'vec.dart';

bool _cameraTracking = false;
bool get cameraTracking => _cameraTracking;

bool _trackWasDown = false;

void trackKey(bool isDown) {
  if (_trackWasDown == isDown) {
    return; // no status change
  }

  // status has switched:

  if (isDown) {
    // only on UP->DOWN: toogle pause on/off
    _cameraTracking = !_cameraTracking;
    String t = cameraTracking ? "on" : "off";
    debug("Camera tracking: $t");
  }

  _trackWasDown = isDown; // update status
}

class Camera {

  Vector3 _position = new Vector3(0.0, 0.0, 10.0);
  Vector3 _focusPosition = new Vector3(0.0, 0.0, 0.0);
  Vector3 _upDirection = new Vector3(0.0, 1.0, 0.0);

  Vector3 get frontDirection => (_focusPosition - _position).normalize();

  String toString() {
    return 'pos=$_position focus=$_focusPosition up=$_upDirection';
  }

  //final double degreesPerSec = 20.0;
  //final double camOrbitRadius = 15.0;

  /*
  Quaternion _orientation = new Quaternion.identity();
  Vector3 _position = new Vector3.zero();

  double _oldAngle = 0.0;
  double _angle = 0.0;
  */

  Camera(Vector3 coord) {
    //update(0.0);
    moveTo(coord);
    _focusPosition = new Vector3(0.0, 0.0, -1.0);
    _upDirection = new Vector3(0.0, 1.0, 0.0);

    debug("new camera: $this");

    if (!vector3Orthogonal(_upDirection, frontDirection)) {
      String fail =
          "new camera: NOT ORTHOGONAL: up=$_upDirection x front=$frontDirection";
      err(fail);
      throw fail;
    }

    assert(vector3Orthogonal(_upDirection, frontDirection));
  }

  void moveTo(Vector3 coord) {
    assert(coord != null);
    _position.setFrom(coord);
  }

  void focusAt(Vector3 coord) {
    assert(coord != null);
    assert(vector3Orthogonal(_upDirection, frontDirection));

    if (coord[0] != _focusPosition[0] || coord[1] != _focusPosition[1] ||
        coord[2] != _focusPosition[2]) {
      double x = _focusPosition[0];
      double y = _focusPosition[1];
      double z = _focusPosition[2];
      _focusPosition.setFrom(coord);
      debug("camera focusAt: from=$_focusPosition to=$coord");
      if (!vector3Orthogonal(_upDirection, frontDirection)) {
        String fail =
            "camera focusAt: NOT ORTHOGONAL: up=$_upDirection x front=$frontDirection";
        err(fail);
        throw fail;
      }
    }

    assert(vector3Orthogonal(_upDirection, frontDirection));
  }

  /*
  void rotate(Matrix4 MV) {
    //MV.setRotation(_orientation.asRotationMatrix());
  }

  void translate(Matrix4 MV) {
    MV.translate(-_position);
  }
  */

  void copyViewMatrix(Matrix4 vm) {
    Matrix4 m = makeViewMatrix(_position, _focusPosition, _upDirection);
    m.copyInto(vm);
  }

  /*
  void update(double gameTime) {
    _oldAngle = _angle;
    _angle = gameTime * this.degreesPerSec % 360.0;
  }
  */

  //static final Vector3 Y = new Vector3(0.0, 1.0, 0.0);

  /*
  void render(double renderInterpolationFactor) {

    double deg = interpolateDegree(_angle, _oldAngle, renderInterpolationFactor
        );
    double rad = deg * math.PI / 180.0;

    // FIXME: should apply a rotation quaternion
    _orientation = new Quaternion.axisAngle(Y, rad).conjugated();
  }
  */
}

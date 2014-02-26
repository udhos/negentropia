library camera;

import 'package:vector_math/vector_math.dart';

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

  static final Vector3 Y = new Vector3(0.0, 1.0, 0.0);

  Vector3 _position = new Vector3(0.0, 0.0, 10.0);
  Vector3 _focusPosition = new Vector3(0.0, 0.0, 0.0);
  Vector3 _upDirection = new Vector3(0.0, 1.0, 0.0);

  Vector3 get frontDirection => (_focusPosition - _position).normalize();
  Vector3 get rightDirection => frontDirection.cross(_upDirection);

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
          "new camera: NOT ORTHOGONAL: up=$_upDirection x front=$frontDirection: dot=${_upDirection.dot(frontDirection)}";
      err(fail);
      throw fail;
    }

    if (!vector3Unit(_upDirection)) {
      String fail =
          "new camera: NOT UNIT: up=$_upDirection length=${_upDirection.length}";
      err(fail);
      throw fail;
    }

    if (!vector3Unit(frontDirection)) {
      String fail =
          "new camera: NOT UNIT: front=$frontDirection length=${frontDirection.length}";
      err(fail);
      throw fail;
    }

    if (!vector3Unit(rightDirection)) {
      String fail =
          "new camera: NOT UNIT: right=$rightDirection length=${rightDirection.length}";
      err(fail);
      throw fail;
    }

    assert(vector3Orthogonal(_upDirection, frontDirection));
    assert(vector3Unit(_upDirection));
    assert(vector3Unit(frontDirection));
    assert(vector3Unit(rightDirection));
  }

  void moveTo(Vector3 coord) {
    assert(coord != null);
    _position.setFrom(coord);
  }

  void focusAt(Vector3 coord) {
    assert(coord != null);
    assert(vector3Orthogonal(_upDirection, frontDirection));
    assert(vector3Unit(_upDirection));
    assert(vector3Unit(frontDirection));
    assert(vector3Unit(rightDirection));

    if (coord[0] == _focusPosition[0] && coord[1] == _focusPosition[1] &&
        coord[2] == _focusPosition[2]) return;

    debug("camera focusAt: from=$_focusPosition to=$coord");

    /*
      Vector3 oldRightDirection = rightDirection; // saves old right direction
      _focusPosition.setFrom(coord); // changes front direction
      _upDirection = oldRightDirection.cross(frontDirection).normalize();
      // new up direction
       */

    _focusPosition.setFrom(coord); // changes front direction
    Vector3 newRightDirection = frontDirection.cross(Y).normalize();
    _upDirection = newRightDirection.cross(frontDirection);

    if (!vector3Orthogonal(_upDirection, frontDirection)) {
      String fail =
          "camera focusAt: NOT ORTHOGONAL: up=$_upDirection x front=$frontDirection: dot=${_upDirection.dot(frontDirection)}";
      err(fail);
      throw fail;
    }

    if (!vector3Unit(_upDirection)) {
      String fail =
          "camera focusAt: NOT UNIT: up=$_upDirection length=${_upDirection.length}";
      err(fail);
      throw fail;
    }

    if (!vector3Unit(frontDirection)) {
      String fail =
          "camera focusAt: NOT UNIT: front=$frontDirection length=${frontDirection.length}";
      err(fail);
      throw fail;
    }

    if (!vector3Unit(rightDirection)) {
      String fail =
          "camera focusAt: NOT UNIT: right=$rightDirection length=${rightDirection.length}";
      err(fail);
      throw fail;
    }

    debug("camera focusAt: $this");

    assert(vector3Orthogonal(_upDirection, frontDirection));
    assert(vector3Unit(_upDirection));
    assert(vector3Unit(frontDirection));
    assert(vector3Unit(rightDirection));
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

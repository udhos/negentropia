library camera;

import 'package:vector_math/vector_math.dart';

import 'logg.dart';
import 'vec.dart';
import 'skybox.dart';

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

  Vector3 get focus => _focusPosition.clone();

  Vector3 get frontVector => _focusPosition - _position;
  Vector3 get frontDirection => frontVector.normalized();
  Vector3 get rightDirection => frontDirection.cross(_upDirection).normalized();

  SkyboxInstance _skybox;

  double get skyboxHalfEdge {
    if (_skybox == null) {
      return null;
    }

    return _skybox.halfEdge;
  }

  void set skybox(SkyboxInstance box) {
    _skybox = box;
  }

  void _skyboxFollowPosition() {
    if (_skybox != null) {
      _skybox.center = _position; // skybox always centered at camera
      //debug("skybox centered at camera: ${_skybox.center}");
    }
  }

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

    _sanity("Camera()");
  }

  void moveForward(double len) {
    moveTo(_position.addScaled(frontDirection, len));

    _sanity("moveForward");
  }

  void _sanity(String label) {
    if (!vector3Unit(_upDirection)) {
      String fail =
          "camera $label: NOT UNIT: up=$_upDirection length=${_upDirection.length}";
      err(fail);
      throw fail;
    }

    if (!vector3Unit(frontDirection)) {
      String fail =
          "camera $label: NOT UNIT: front=$frontDirection length=${frontDirection.length}";
      err(fail);
      throw fail;
    }

    if (!vector3Unit(rightDirection)) {
      String fail =
          "camera $label: NOT UNIT: right=$rightDirection length=${rightDirection.length}";
      err(fail);
      throw fail;
    }

    if (!vector3Orthogonal(_upDirection, frontDirection)) {
      String fail =
          "camera $label: NOT ORTHOGONAL: up=$_upDirection x front=$frontDirection: dot=${_upDirection.dot(frontDirection)}";
      err(fail);
      throw fail;
    }
  }

  void rotateAroundFocusVertical(double radAngle) {
    _position.sub(_focusPosition); // ----: translate focus to origin

    Quaternion q = new Quaternion.axisAngle(_upDirection, radAngle);
    q.rotate(_position);

    _position.add(_focusPosition); // undo: translate focus to origin

    _upDirection = rightDirection.cross(frontDirection).normalized();
    // FIXME why is this needed?

    _sanity("rotateAroundFocusVertical");

    _skyboxFollowPosition();
  }

  void rotateAroundFocusHorizontal(double radAngle) {
    _position.sub(_focusPosition); // ----: translate focus to origin

    Quaternion q = new Quaternion.axisAngle(rightDirection, radAngle);
    q.rotate(_position);

    _position.add(_focusPosition); // undo: translate focus to origin

    _upDirection = rightDirection.cross(frontDirection).normalized();

    _sanity("rotateAroundFocusHorizontal");

    _skyboxFollowPosition();
  }


  void moveTo(Vector3 coord) {
    assert(coord != null);
    _position.setFrom(coord);
    _skyboxFollowPosition();
  }

  void focusAt(Vector3 coord) {

    if (coord[0] == _focusPosition[0] && coord[1] == _focusPosition[1] &&
        coord[2] == _focusPosition[2]) return;

    //debug("camera focusAt: from=$_focusPosition to=$coord");

    _focusPosition.setFrom(coord); // changes front direction
    Vector3 newRightDirection = frontDirection.cross(Y).normalized();
    _upDirection = newRightDirection.cross(frontDirection).normalized();

    _sanity("focusAt");
  }

  /*
  void rotate(Matrix4 MV) {
    //MV.setRotation(_orientation.asRotationMatrix());
  }

  void translate(Matrix4 MV) {
    MV.translate(-_position);
  }
  */

  void viewMatrix(Matrix4 vm) {
    setViewMatrix(vm, _position, _focusPosition, _upDirection);
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

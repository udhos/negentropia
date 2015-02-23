library camera;

import 'package:vector_math/vector_math.dart';

import 'logg.dart';
import 'vec.dart';
import 'skybox.dart';

bool _cameraTracking = false;
bool get cameraTracking => _cameraTracking;
Vector3 cameraFocusTemp = new Vector3.zero();
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
    _skyboxFollowPosition();
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

  Camera(Vector3 coord) {
    moveTo(coord);
    _focusPosition = new Vector3(0.0, 0.0, -1.0);
    _upDirection = new Vector3(0.0, 1.0, 0.0);

    debug("new camera: $this");

    _sanity("Camera()");
  }

  void moveForward(double len) {
    Vector3 newPosition = _position.clone();
    newPosition.addScaled(frontDirection, len);
    moveTo(newPosition);

    _sanity("moveForward");
  }

  void setForwardDistance(double len) {
    Vector3 newPosition = _focusPosition.clone();
    newPosition.addScaled(frontDirection, -len);
    moveTo(newPosition);

    _sanity("setFordwardDistance");
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

  void fixUp() {
    _upDirection = rightDirection.cross(frontDirection).normalized();
  }

  void rotateAroundFocusVertical(double radAngle) {
    _position.sub(_focusPosition); // ----: translate focus to origin

    Quaternion q = new Quaternion.axisAngle(_upDirection, radAngle);
    q.rotate(_position);

    _position.add(_focusPosition); // undo: translate focus to origin

    fixUp(); // FIXME why is this needed?

    _sanity("rotateAroundFocusVertical");

    _skyboxFollowPosition();
  }

  void rotateAroundFocusHorizontal(double radAngle) {
    _position.sub(_focusPosition); // ----: translate focus to origin

    Quaternion q = new Quaternion.axisAngle(rightDirection, radAngle);
    q.rotate(_position);

    _position.add(_focusPosition); // undo: translate focus to origin

    fixUp();

    _sanity("rotateAroundFocusHorizontal");

    _skyboxFollowPosition();
  }

  void moveTo(Vector3 coord) {
    _position.setFrom(coord);

    fixUp(); // FIXME why is this needed?

    _sanity("moveTo");

    _skyboxFollowPosition();
  }

  void focusAt(Vector3 coord) {
    if (coord[0] == _focusPosition[0] &&
        coord[1] == _focusPosition[1] &&
        coord[2] == _focusPosition[2]) return;

    log("camera focusAt: from=$_focusPosition to=$coord");

    _focusPosition.setFrom(coord); // changes front direction
    Vector3 newRightDirection = frontDirection.cross(Y).normalized();
    _upDirection = newRightDirection.cross(frontDirection).normalized();

    _sanity("focusAt");
  }

  void alignHorizontal() {
    _position.y = _focusPosition.y;
    _upDirection.setFrom(Y);

    _sanity("alignHorizontal");

    _skyboxFollowPosition();
  }

  /**
   * Constructs an OpenGL view matrix for this camera into [viewMatrix].
   */
  void loadViewMatrixInto(Matrix4 viewMatrix) {
    setViewMatrix(viewMatrix, _position, _focusPosition, _upDirection);
  }
}

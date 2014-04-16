library camera_control;

import 'dart:math' as math;

//import 'package:vector_math/vector_math.dart';

//import 'logg.dart';
import 'camera.dart';
import 'selection.dart';
import 'message.dart';

double wheelToDistance(int wheelDy) {
  // 100 points => 10.0 meters
  return wheelDy.toDouble() / 10.0;
}

const double DEG_TO_RAD = math.PI / 180.0;

double mouseToRadians(int mouse) {
  // 1 pixel => 1 degree
  return mouse.toDouble() * DEG_TO_RAD;
}

class CameraControl {

  int orbitFocusDx = 0;
  int orbitFocusDy = 0;
  int forwardDy = 0;

  void update(double dt, Camera cam) {

    //
    // Consume pending rotation
    //
    if (orbitFocusDx != 0) {
      cam.rotateAroundFocusVertical(mouseToRadians(orbitFocusDx));
      orbitFocusDx = 0;
    }
    if (orbitFocusDy != 0) {
      cam.rotateAroundFocusHorizontal(mouseToRadians(orbitFocusDy));
      orbitFocusDy = 0;
    }

    if (forwardDy != 0) {
      cam.moveForward(wheelToDistance(forwardDy));
      forwardDy = 0;
    }

  }

  void orbitFocus(int dx, dy) {
    orbitFocusDx += dx;
    orbitFocusDy += dy;
  }

  void moveForward(Camera cam, int dy) {
    if (dy > 0) {
      // getting closer
      double boundingRadius = getSelectionBoundingRadius();
      if (boundingRadius == null) {
        boundingRadius = 1.0;
      }

      double currDistance = cam.frontVector.length;
      if (currDistance - wheelToDistance(dy) < boundingRadius) {
        messageUser("camera: minimum distance reached");
        return;
      }

    }
    forwardDy += dy;
  }

}

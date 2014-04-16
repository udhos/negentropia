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

  double getBoundingRadius() {
    double boundingRadius = getSelectionBoundingRadius();
    if (boundingRadius == null) {
      // ugh: no selected object
      boundingRadius = 1.0;
    }
    return boundingRadius;
  }

  void moveForward(Camera cam, int dy) {
    if (dy > 0) {
      // getting close - closest distance is bounding radius

      double boundingRadius = getBoundingRadius();
      double currDistance = cam.frontVector.length;
      if (currDistance - wheelToDistance(dy) < boundingRadius) {
        messageUser("camera: minimum distance reached: $boundingRadius");
        return;
      }
    } else {
      // getting away - farthest distance is skybox half edge (minus bounding diameter)

      double halfEdge = cam.skyboxHalfEdge;
      if (halfEdge != null) {
        double maxDistance = halfEdge - 2.0 * getBoundingRadius();
        double currDistance = cam.frontVector.length;
        if (currDistance + wheelToDistance(dy) > maxDistance) {
          messageUser("camera: maximum distance reached: $maxDistance");
          return;
        }
      }
    }

    forwardDy += dy;
  }

}

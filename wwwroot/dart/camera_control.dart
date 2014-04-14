library camera_control;

import 'dart:math' as math;

//import 'package:vector_math/vector_math.dart';

import 'logg.dart';
import 'camera.dart';

class CameraControl {

  int orbitFocusDx = 0;
  int orbitFocusDy = 0;
  int forwardDy = 0;

  void update(double dt, Camera cam) {

    //
    // Consume pending rotation
    //
    if (orbitFocusDx != 0) {
      // 1 pixel = 1 degree
      cam.rotateAroundFocusVertical(orbitFocusDx * math.PI / 180.0);
      orbitFocusDx = 0;
    }
    if (orbitFocusDy != 0) {
      cam.rotateAroundFocusHorizontal(orbitFocusDy * math.PI / 180.0);
      orbitFocusDy = 0;
    }

    if (forwardDy != 0) {
      debug("camera forward: $forwardDy");
      cam.moveForward(forwardDy.toDouble());
      forwardDy = 0;
    }

  }

  void orbitFocus(int dx, dy) {
    debug("orbitFocus: dx=$dx dy=$dy");
    orbitFocusDx += dx;
    orbitFocusDy += dy;
  }

  void moveForward(int dy) {
    forwardDy += dy;
  }

}

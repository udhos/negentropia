library camera;

import 'package:vector_math/vector_math.dart';

class Camera {
  vec3 eye, center, up;
  Camera(this.eye, this.center, this.up);
}

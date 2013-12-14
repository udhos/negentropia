library interpolate;

bool _paused = true;

void pause(bool on) {
  _paused = on;
}

bool paused() {
  return _paused;
}

double interpolateDegree(double angleNew, double angleOld, double interpolationFactor) {
  
  if (paused()) {
    return angleNew;  
  }
  
  double deg;
  if (angleNew > angleOld) {
    deg = interpolationFactor * (angleNew         - angleOld) + angleOld;
  } else {
    // undo modulo 360 for correct interpolation
    deg = interpolationFactor * (angleNew + 360.0 - angleOld) + angleOld;
  }
  return deg;  
}
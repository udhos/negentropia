library interpolate;

bool _paused = true;

bool paused() {
  return _paused;
}

bool _wasDown;

void pauseKey(bool isDown) {
  if (_wasDown == isDown) {
    return; // no status change
  }

  // status has switched:

  if (isDown) {
    // only on UP->DOWN: toogle pause on/off
    _paused = !_paused;
  }

  _wasDown = isDown; // update status
}

double interpolateDegree(double angleNew, double angleOld,
    double interpolationFactor) {

  if (paused()) {
    return angleNew;
  }

  double deg;
  if (angleNew > angleOld) {
    deg = interpolationFactor * (angleNew - angleOld) + angleOld;
  } else {
    // undo modulo 360 for correct interpolation
    deg = interpolationFactor * (angleNew + 360.0 - angleOld) + angleOld;
  }
  return deg;
}

library interpolate;

double interpolateDegree(double angleNew, double angleOld, double interpolationFactor) {
  double deg;
  if (angleNew > angleOld) {
    deg = interpolationFactor * (angleNew         - angleOld) + angleOld;
  } else {
    // undo modulo 360 for correct interpolation
    deg = interpolationFactor * (angleNew + 360.0 - angleOld) + angleOld;
  }
  return deg;  
}
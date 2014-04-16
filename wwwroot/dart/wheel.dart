library wheel;

import 'logg.dart';

int normalizeWheel(int dy) {
  if (dy.abs() < 100) {
    // Firefox: Nx3
    return dy * 100 ~/ 3;
  }

  if (dy % 120 == 0) {
    // IE: Nx120
    return dy * 100 ~/ 120;
  }

  if (dy % 100 == 0) {
    // Chrome, Opera: Nx100
    return dy;
  }

  // Unknown delta

  err("normalizeWheel: unexpected dy=$dy");

  return dy;
}

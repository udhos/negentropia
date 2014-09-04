package util

import (
	"math"
)

/*
	http://floating-point-gui.de/errors/comparison/

	public static boolean nearlyEqual(float a, float b, float epsilon) {
		final float absA = Math.abs(a);
		final float absB = Math.abs(b);
		final float diff = Math.abs(a - b);

		if (a == b) { // shortcut, handles infinities
			return true;
		} else if (a == 0 || b == 0 || diff < Float.MIN_NORMAL) {
			// a or b is zero or both are extremely close to it
			// relative error is less meaningful here
			return diff < (epsilon * Float.MIN_NORMAL);
		} else { // use relative error
			return diff / (absA + absB) < epsilon;
		}
	}
*/

func NearlyEqualEpsilon(a, b, epsilon float64) bool {
	absA := math.Abs(a)
	absB := math.Abs(b)
	diff := math.Abs(a - b)

	if a == b {
		// shortcut, handles infinities
		return true
	}

	if a == 0 || b == 0 || diff < math.SmallestNonzeroFloat64 {
		// a or b is zero or both are extremely close to it
		// relative error is less meaningful here
		return diff < (epsilon * math.SmallestNonzeroFloat64)
	}

	// use relative error
	return diff/(absA+absB) < epsilon
}

func NearlyEqual(a, b float64) bool {
	return NearlyEqualEpsilon(a, b, 0.000001)
}

func CloseToZero(f float64) bool {
	return NearlyEqual(f, 0.0)
}

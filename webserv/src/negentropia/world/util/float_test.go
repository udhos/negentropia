package util

import (
	"math"
	"testing"
)

func expect(t *testing.T, a, b float64, expectedResult bool) {
	t.Logf("expectEqual(%v): %v=%v", expectedResult, a, b)
	
	if result := NearlyEqual(a, b); result != expectedResult {
		t.Errorf("Unexpected result: NearlyEqual(%v,%v)=%v", a, b, result)
		return
	}
}

func expectEpsilon(t *testing.T, a, b, epsilon float64, expectedResult bool) {
	t.Logf("expectEqualEpsilon(%v): %v=%v", expectedResult, a, b)
	
	if result := NearlyEqualEpsilon(a, b, epsilon); result != expectedResult {
		t.Errorf("Unexpected: NearlyEqualEpsilon(%v,%v,%v)=%v", a, b, epsilon, result)
		return
	}
}

func TestNearlyEqual(t *testing.T) {

	// big
	expect(t, 1000000, 1000001, true);
	expect(t, 1000001, 1000000, true);
	expect(t, 10000, 10001, false);
	expect(t, 10001, 10000, false);
	
	// big negative
	expect(t, -1000000, -1000001, true);
	expect(t, -1000001, -1000000, true);
	expect(t, -10000, -10001, false);
	expect(t, -10001, -10000, false);
	
	// around 1
	expect(t, 1.0000001, 1.0000002, true);
	expect(t, 1.0000002, 1.0000001, true);
	expect(t, 1.0001, 1.0002, false);
	expect(t, 1.0002, 1.0001, false);
	
	// around -1
	expect(t, -1.000001, -1.000002, true);
	expect(t, -1.000002, -1.000001, true);
	expect(t, -1.0001, -1.0002, false);
	expect(t, -1.0002, -1.0001, false);
	
	// 0 .. 1
	expect(t, 0.000000001000001, 0.000000001000002, true);
	expect(t, 0.000000001000002, 0.000000001000001, true);
	expect(t, 0.000000000001002, 0.000000000001001, false);
	expect(t, 0.000000000001001, 0.000000000001002, false);
	
	// -1 .. 0
	expect(t, -0.000000001000001, -0.000000001000002, true);
	expect(t, -0.000000001000002, -0.000000001000001, true);
	expect(t, -0.000000000001002, -0.000000000001001, false);
	expect(t, -0.000000000001001, -0.000000000001002, false);
	
	// zero
	expect(t, 0.0, 0.0, true);
	expect(t, 0.0, -0.0, true);
	expect(t, -0.0, 0.0, true);
	expect(t, -0.0, -0.0, true);
	
	expect(t, 0.00000001, 0.0, false);
	expect(t, 0.0, 0.00000001, false);
	expect(t, -0.00000001, 0.0, false);
	expect(t, 0.0, -0.00000001, false);
	
	expectEpsilon(t, 0.0, 1e-40, 0.01, true);
	expectEpsilon(t, 1e-40, 0.0, 0.01, true);
	expectEpsilon(t, 0.0, 1e-40, 0.000001, false);
	expectEpsilon(t, 1e-40, 0.0, 0.000001, false);
	
	expectEpsilon(t, 0.0, -1e-40, 0.1, true);
	expectEpsilon(t, -1e-40, 0.0, 0.1, true);
	expectEpsilon(t, 0.0, -1e-40, 0.00000001, false);
	expectEpsilon(t, -1e-40, 0.0, 0.00000001, false);
	
	// infinities
	positiveInfinity := math.Inf(1)
	negativeInfinity := math.Inf(-1)
	expect(t, positiveInfinity, positiveInfinity, true);
	expect(t, negativeInfinity, negativeInfinity, true);
	expect(t, positiveInfinity, negativeInfinity, false);
	expect(t, negativeInfinity, positiveInfinity, false);
	expect(t, positiveInfinity, math.MaxFloat64, false);
	expect(t, negativeInfinity, -math.MaxFloat64, false);
	
	// NaN
	nan := math.NaN()
	expect(t, nan, nan, false);
	expect(t, nan, 0.0, false);
	expect(t, 0.0, nan, false);
	expect(t, nan, -0.0, false);
	expect(t, -0.0, nan, false);
	expect(t, nan, positiveInfinity, false);
	expect(t, positiveInfinity, nan, false);
	expect(t, nan, negativeInfinity, false);
	expect(t, negativeInfinity, nan, false);
	expect(t, nan, math.MaxFloat64, false);
	expect(t, math.MaxFloat64, nan, false);
	expect(t, nan, -math.MaxFloat64, false);
	expect(t, -math.MaxFloat64, nan, false);
	expect(t, nan, math.SmallestNonzeroFloat64, false);
	expect(t, math.SmallestNonzeroFloat64, nan, false);
	expect(t, nan, -math.SmallestNonzeroFloat64, false);
	expect(t, -math.SmallestNonzeroFloat64, nan, false);
	
	// opposite sides of 0
	expect(t, 1.000000001, -1.0, false);
	expect(t, -1.0, 1.000000001, false);
	expect(t, -1.000000001, 1.0, false);
	expect(t, 1.0, -1.000000001, false);
	expect(t, 10 * math.SmallestNonzeroFloat64, 10 * -math.SmallestNonzeroFloat64, true);
	expect(t, 10000 * math.SmallestNonzeroFloat64, 10000 * -math.SmallestNonzeroFloat64, false);
	
	// very close to zero
	expect(t, math.SmallestNonzeroFloat64, -math.SmallestNonzeroFloat64, true);
	expect(t, -math.SmallestNonzeroFloat64, math.SmallestNonzeroFloat64, true);
	expect(t, math.SmallestNonzeroFloat64, 0, true);
	expect(t, 0, math.SmallestNonzeroFloat64, true);
	expect(t, -math.SmallestNonzeroFloat64, 0, true);
	expect(t, 0, -math.SmallestNonzeroFloat64, true);
	
	expect(t, 0.000000001, math.SmallestNonzeroFloat64, false);
	expect(t, math.SmallestNonzeroFloat64, 0.000000001, false);
	expect(t, 0.000000001, -math.SmallestNonzeroFloat64, false);
	expect(t, -math.SmallestNonzeroFloat64, 0.000000001, false);
}

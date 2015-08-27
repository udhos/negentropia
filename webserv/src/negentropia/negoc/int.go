package main

import (
	"math"
)

func intAbs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

func round(a float64) float64 {
	var r float64
	if a < 0 {
		r = math.Ceil(a - 0.5)
	} else {
		r = math.Floor(a + 0.5)
	}
	return r
}

func roundToInt(a float64) int {
	return int(roundToInt(a))
}

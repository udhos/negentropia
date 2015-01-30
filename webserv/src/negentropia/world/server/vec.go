package server

import (
	"fmt"
	//"log"
	"math"
	"strconv"
	//"strings"
	"unicode"

	"github.com/udhos/vectormath"

	"negentropia/world/parser"
	"negentropia/world/util"
)

func quatString(q vectormath.Quat) string {
	var f byte = 'f'
	prec := -1
	bitSize := 64
	return fmt.Sprintf("%s,%s,%s,%s",
		strconv.FormatFloat(float64(q.X), f, prec, bitSize),
		strconv.FormatFloat(float64(q.Y), f, prec, bitSize),
		strconv.FormatFloat(float64(q.Z), f, prec, bitSize),
		strconv.FormatFloat(float64(q.W), f, prec, bitSize))
}

func vector3String(v vectormath.Vector3) string {
	//return fmt.Sprintf("%f,%f,%f", v.X, v.Y, v.Z)
	var f byte = 'f'
	prec := -1
	bitSize := 32
	return fmt.Sprintf("%s,%s,%s",
		strconv.FormatFloat(float64(v.X), f, prec, bitSize),
		strconv.FormatFloat(float64(v.Y), f, prec, bitSize),
		strconv.FormatFloat(float64(v.Z), f, prec, bitSize))
}

func parseVector3Func(result *vectormath.Vector3, text string, f func(rune) bool) error {
	floatSlice, err := parser.ParseFloatVectorFunc(text, 3, f)
	if err != nil {
		return fmt.Errorf("parseVector3Func: error: %v", err)
	}
	vectormath.V3MakeFromElems(result, float32(floatSlice[0]), float32(floatSlice[1]), float32(floatSlice[2]))
	return nil
}

func parseVector3Space(result *vectormath.Vector3, text string) error {
	return parseVector3Func(result, text, unicode.IsSpace)
}

func isComma(r rune) bool {
	return r == ','
}

func parseVector3(result *vectormath.Vector3, text string) error {
	return parseVector3Func(result, text, isComma)
}

/*
func vector3Orthogonal(v1, v2 vectormath.Vector3) bool {
	dot := float64(vectormath.V3Dot(&v1, &v2))
	return util.CloseToZero(dot)
}

func vector3Unit(v vectormath.Vector3) bool {
	length := float64(v.Length())
	return util.CloseToZero(length - 1.0)
}
*/

func v3dot(x1, y1, z1, x2, y2, z2 float64) float64 {
	return x1*x2 + y1*y2 + z1*z2
}

func v3len(x1, y1, z1 float64) float64 {
	return math.Sqrt(v3dot(x1, y1, z1, x1, y1, z1))
}

func v3scalarMul(result []float64, scalar float64) {
	result[0] *= scalar
	result[1] *= scalar
	result[2] *= scalar
}

func v3add(result []float64, x1, y1, z1, x2, y2, z2 float64) {
	result[0] = x1 + x2
	result[1] = y1 + y2
	result[2] = z1 + z2
}

func vector3Orthogonal(v1, v2 vectormath.Vector3) bool {
	x1 := float64(v1.X)
	y1 := float64(v1.Y)
	z1 := float64(v1.Z)
	x2 := float64(v2.X)
	y2 := float64(v2.Y)
	z2 := float64(v2.Z)
	return util.CloseToZero(v3dot(x1, y1, z1, x2, y2, z2))
}

func vector3Unit(v vectormath.Vector3) bool {
	x := float64(v.X)
	y := float64(v.Y)
	z := float64(v.Z)
	return util.CloseToZero(v3len(x, y, z) - 1.0)
}

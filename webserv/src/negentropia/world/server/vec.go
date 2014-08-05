package server

import (
	"fmt"
	//"log"
	//"math"
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

func vector3Orthogonal(v1, v2 vectormath.Vector3) bool {
	dot := float64(vectormath.V3Dot(&v1, &v2))
	return util.CloseToZero(dot)
}

func vector3Unit(v vectormath.Vector3) bool {
	length := float64(v.Length())
	return util.CloseToZero(length - 1.0)
}

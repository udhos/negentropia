package server

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"

	"github.com/udhos/vectormath"
)

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

func parseVector3(result *vectormath.Vector3, text string) error {
	var x, y, z float64
	list := strings.Split(text, ",")
	size := len(list)
	if size != 3 {
		e := fmt.Errorf("parseVector3: text=[%s] size=%d != 3", text, size)
		log.Print(e)
		return e
	}
	var err error
	var i string
	i = strings.TrimSpace(list[0])
	if x, err = strconv.ParseFloat(i, 32); err != nil {
		e := fmt.Errorf("parseVector3: text=[%s] parse x=[%s] failure: %s", text, i, err)
		log.Print(e)
		return e
	}
	i = strings.TrimSpace(list[1])
	if y, err = strconv.ParseFloat(i, 32); err != nil {
		e := fmt.Errorf("parseVector3: text=[%s] parse y=[%s] failure: %s", text, i, err)
		log.Print(e)
		return e
	}
	i = strings.TrimSpace(list[2])
	if z, err = strconv.ParseFloat(i, 32); err != nil {
		e := fmt.Errorf("parseVector3: text=[%s] parse z=[%s] failure: %s", text, i, err)
		log.Print(e)
		return e
	}
	vectormath.V3MakeFromElems(result, float32(x), float32(y), float32(z))
	//log.Printf("parseVector3: text=[%s] result: %s", text, result)
	return nil
}

const MAX_CLOSE_TO_ZERO = 1e-6

func closeToZero(f float64) bool {
	return math.Abs(f) < MAX_CLOSE_TO_ZERO
}

func vector3Orthogonal(v1, v2 vectormath.Vector3) bool {
	dot := float64(vectormath.V3Dot(&v1, &v2))
	return closeToZero(dot)
}

func vector3Unit(v vectormath.Vector3) bool {
	length := float64(v.Length())
	return closeToZero(length - 1.0)
}

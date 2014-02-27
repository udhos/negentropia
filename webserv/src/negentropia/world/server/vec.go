package server

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/spate/vectormath"
)

func parseVector3(result *vectormath.Vector3, text string) error {
	var x, y, z float64
	list := strings.Split(text, ",")
	var err error
	if x, err = strconv.ParseFloat(list[0], 32); err != nil {
		e := fmt.Errorf("parseVector3: text=[%s] parse x=[%s] failure: %s", text, list[0], err)
		log.Print(e)
		return e
	}
	if y, err = strconv.ParseFloat(list[1], 32); err != nil {
		e := fmt.Errorf("parseVector3: text=[%s] parse y=[%s] failure: %s", text, list[1], err)
		log.Print(e)
		return e
	}
	if z, err = strconv.ParseFloat(list[2], 32); err != nil {
		e := fmt.Errorf("parseVector3: text=[%s] parse z=[%s] failure: %s", text, list[2], err)
		log.Print(e)
		return e
	}
	vectormath.V3MakeFromElems(result, float32(x), float32(y), float32(z))
	//log.Printf("parseVector3: text=[%s] result: %s", text, result)
	return nil
}

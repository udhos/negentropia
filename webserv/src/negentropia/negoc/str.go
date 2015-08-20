package main

import (
	"fmt"
	"negentropia/world/parser"
	"strconv"
	"strings"
)

func stringIsTrue(s string) bool {

	s = strings.TrimSpace(s)

	if s == "" {
		return false
	}

	var val bool
	var err error

	if val, err = strconv.ParseBool(s); err != nil {
		log(fmt.Sprintf("stringIsTrue(%s): bad value: %v", s, err))
		return false
	}

	return val
}

func stringIsFalse(s string) bool {

	s = strings.TrimSpace(s)

	if s == "" {
		return false
	}

	var val bool
	var err error

	if val, err = strconv.ParseBool(s); err != nil {
		log(fmt.Sprintf("stringIsFalse(%s): bad value: %v", s, err))
		return false
	}

	return !val
}

func parseVector3(s string) ([]float64, error) {
	v, err := parser.ParseFloatVectorComma(s, 3)
	if err != nil {
		return v, fmt.Errorf("parseVector3: error: %v", err)
	}
	return v, nil
}

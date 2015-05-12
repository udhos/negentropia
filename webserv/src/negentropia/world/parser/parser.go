package parser

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

func parseFloatSlice(list []string) ([]float64, error) {
	result := make([]float64, len(list))

	for i, j := range list {
		j = strings.TrimSpace(j)
		var err error
		if result[i], err = strconv.ParseFloat(j, 64); err != nil {
			return nil, fmt.Errorf("parseFloatSlice: list=[%v] elem[%v]=[%s] failure: %v", list, i, j, err)
		}
	}

	return result, nil
}

func ParseFloatSliceFunc(text string, f func(rune) bool) ([]float64, error) {
	return parseFloatSlice(strings.FieldsFunc(text, f))
}

func ParseFloatSliceSpace(text string) ([]float64, error) {
	return ParseFloatSliceFunc(text, unicode.IsSpace)
}

func ParseFloatVectorFunc(text string, size int, f func(rune) bool) ([]float64, error) {
	list := strings.FieldsFunc(text, f)
	if s := len(list); s != size {
		return nil, fmt.Errorf("ParseFloatVectorFunc: text=[%v] size=%v must be %v", text, s, size)
	}

	return parseFloatSlice(list)
}

func ParseFloatVectorSpace(text string, size int) ([]float64, error) {
	return ParseFloatVectorFunc(text, size, unicode.IsSpace)
}

func ParseFloatVectorComma(text string, size int) ([]float64, error) {
	isComma := func(c rune) bool {
		return c == ','
	}

	return ParseFloatVectorFunc(text, size, isComma)
}

func ParseFloatVector3Comma(text string) ([]float64, error) {
	return ParseFloatVectorComma(text, 3)
}

package parser

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

func ParseFloatSliceFunc(text string, f func(rune) bool) ([]float64, error) {
	list := strings.FieldsFunc(text, f)

	result := make([]float64, len(list))

	for i, j := range list {
		j = strings.TrimSpace(j)
		var err error
		if result[i], err = strconv.ParseFloat(j, 64); err != nil {
			return nil, fmt.Errorf("parseFloatSlice: text=[%v] elem[%v]=[%s] failure: %v", text, i, j, err)
		}
	}

	return result, nil
}

func ParseFloatSliceSpace(text string) ([]float64, error) {
	return ParseFloatSliceFunc(text, unicode.IsSpace)
}

func ParseFloatVectorFunc(text string, size int, f func(rune) bool) ([]float64, error) {
	list := strings.FieldsFunc(text, f)
	if s := len(list); s != size {
		return nil, fmt.Errorf("parseFloatSlice: text=[%v] size=%v must be %v", text, s, size)
	}

	result := make([]float64, size)

	for i, j := range list {
		j = strings.TrimSpace(j)
		var err error
		if result[i], err = strconv.ParseFloat(j, 64); err != nil {
			return nil, fmt.Errorf("parseFloatSlice: text=[%v] elem[%v]=[%s] failure: %v", text, i, j, err)
		}
	}

	return result, nil
}

func ParseFloatVectorSpace(text string, size int) ([]float64, error) {
	return ParseFloatVectorFunc(text, size, unicode.IsSpace)
}

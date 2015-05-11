package main

import (
	"fmt"
	//"strings"
)

func stringIsTrue(s string) bool {

	result := s != ""

	log(fmt.Sprintf("FIXME stringIsTrue(%s)=%v", s, result))

	return result
}

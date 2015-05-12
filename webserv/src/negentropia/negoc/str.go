package main

import (
	"fmt"
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

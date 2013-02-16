package util

import (
	"strings"
)

func GetPort(hostPort string) string {
	pair := strings.Split(hostPort, ":")
	if len(pair) < 2 {
		return ""
	}

	return ":" + pair[1]
}

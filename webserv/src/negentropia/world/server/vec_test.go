package server

import (
	"testing"

	"github.com/udhos/vectormath"
)

func TestParseVector3(t *testing.T) {
	var vec3 vectormath.Vector3

	in := "1.0,2.0,3.0x"
	want := "1,2,3"
	if err := parseVector3(&vec3, in); err != nil {
		t.Errorf("parseVector3(%v): err=%v", in, err)
		return
	}

	out := vector3String(vec3)
	if out != want {
		t.Errorf("parseVector3(%v) = %v, want %v", in, out, want)
	}
}

package server

import (
	"testing"

	"github.com/udhos/vectormath"
)

func expectWant(t *testing.T, in, want string) {
	var vec3 vectormath.Vector3

	if err := parseVector3(&vec3, in); err != nil {
		t.Errorf("parseVector3(%v): err=%v", in, err)
		return
	}

	out := vector3String(vec3)
	if out != want {
		t.Errorf("parseVector3(%v) = %v, want %v", in, out, want)
	}
}

func expectErr(t *testing.T, in string) {
	var vec3 vectormath.Vector3

	if err := parseVector3(&vec3, in); err == nil {
		t.Errorf("parseVector3(%v): accepted, but wanted error", in)
	}
}

func TestParseVector3(t *testing.T) {
	expectWant(t, "1,2,3", "1,2,3")
	expectWant(t, "1.0,2.0,3.0", "1,2,3")
	expectWant(t, " 1.0 , 2.0 , 3.0 ", "1,2,3")
	expectErr(t, "1.0x,2.0,3.0")
	expectErr(t, "1.0,2.0x,3.0")
	expectErr(t, "1.0,2.0,3.0x")
	expectErr(t, "x1.0,2.0,3.0")
	expectErr(t, "1.0,x2.0,3.0")
	expectErr(t, "1.0,2.0,x3.0")
}

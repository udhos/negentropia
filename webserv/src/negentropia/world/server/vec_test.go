package server

import (
	"math"
	"testing"

	"github.com/udhos/vectormath"
)

func expectWant(t *testing.T, in, want string) {
	var vec3 vectormath.Vector3

	t.Logf("expectWant: in=[%v] want=[%v]", in, want)

	if err := parseVector3(&vec3, in); err != nil {
		t.Errorf("parseVector3(%v): err=%v", in, err)
		return
	}

	if out := vector3String(vec3); out != want {
		t.Errorf("parseVector3(%v)=[%v], want [%v]", in, out, want)
	}
}

func expectErr(t *testing.T, in string) {
	var vec3 vectormath.Vector3

	t.Logf("expectErr: in=[%v]", in)

	if err := parseVector3(&vec3, in); err == nil {
		t.Errorf("parseVector3(%v): accepted, but wanted error", in)
	}
}

func TestParseVector3(t *testing.T) {
	expectWant(t, "1,2,3", "1,2,3")
	expectWant(t, "1.1,2.2,3.3", "1.1,2.2,3.3")
	expectWant(t, "1.0,2.0,3.0", "1,2,3")
	expectWant(t, " 1.0 , 2.0 , 3.0 ", "1,2,3")
	expectErr(t, "")
	expectErr(t, " ")
	expectErr(t, ",")
	expectErr(t, " , ")
	expectErr(t, ",,")
	expectErr(t, " , , ")
	expectErr(t, ",,,")
	expectErr(t, " , , , ")
	expectErr(t, "1,2,3,4")
	expectErr(t, "1.0x,2.0,3.0")
	expectErr(t, "1.0,2.0x,3.0")
	expectErr(t, "1.0,2.0,3.0x")
	expectErr(t, "x1.0,2.0,3.0")
	expectErr(t, "1.0,x2.0,3.0")
	expectErr(t, "1.0,2.0,x3.0")
}

func TestQuat(t *testing.T) {
	radAngle := math.Pi
    var axis vectormath.Vector3
	var quat vectormath.Quat
	
	vectormath.V3MakeFromElems(&axis, 1.0, 1.0, 1.0)
	vectormath.V3Normalize(&axis, &axis)
	vectormath.QMakeRotationAxis(&quat, float32(radAngle), &axis)
	
	t.Errorf("TestQuat: 90deg around (1,1,1): quat = %q", quat.String())
}

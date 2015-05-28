package obj

import (
	"fmt"
	//"reflect" // for reflect.DeepEqual
	"testing"
)

func expectInt(t *testing.T, label string, want, got int) {
	if want != got {
		t.Errorf("%s: want=%d got=%d", label, want, got)
	}
}

func sliceEqualInt(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}

	for i, v := range a {
		if v != b[i] {
			return false
		}
	}

	return true
}

func sliceEqualFloat(a, b []float32) bool {
	if len(a) != len(b) {
		return false
	}

	for i, v := range a {
		if v != b[i] {
			return false
		}
	}

	return true
}

func TestCube(t *testing.T) {

	o, err := NewObjFromBuf([]byte(cubeObj), func(msg string) { fmt.Printf("TestCube NewObjFromBuf: log: %s\n", msg) })
	if err != nil {
		t.Errorf("TestCube: NewObjFromBuf: %v", err)
		return
	}

	//expectInt(t, "TestCube vertexCount", 23, o.vertexCount())
	//expectInt(t, "TestCube indexCount", 36, o.indexCount())

	if !sliceEqualInt(cubeIndices, o.Indices) {
		t.Errorf("TestCube: indices: want=%v got=%v", cubeIndices, o.Indices)
	}

	if !sliceEqualFloat(cubeCoord, o.Coord) {
		t.Errorf("TestCube: coord: want=%d%v got=%d%v", len(cubeCoord), cubeCoord, len(o.Coord), o.Coord)
	}
}

func TestRelativeIndex(t *testing.T) {
	o, err := NewObjFromBuf([]byte(relativeObj), func(msg string) { fmt.Printf("TestRelativeIndex NewObjFromBuf: log: %s\n", msg) })
	if err != nil {
		t.Errorf("TestRelativeIndex: NewObjFromBuf: %v", err)
		return
	}

	//indices := o.Indices[:len(o.Indices):len(o.Indices)]
	if !sliceEqualInt(relativeIndices, o.Indices) {
		t.Errorf("TestRelativeIndex: indices: want=%v got=%v", relativeIndices, o.Indices)
	}

	//coord := o.Coord[:len(o.Coord):len(o.Coord)]
	if !sliceEqualFloat(relativeCoord, o.Coord) {
		t.Errorf("TestRelativeIndex: coord: want=%v got=%v", relativeCoord, o.Coord)
	}
}

func TestForwardVertex(t *testing.T) {
	o, err := NewObjFromBuf([]byte(forwardObj), func(msg string) { fmt.Printf("TestForwardVertex NewObjFromBuf: log: %s\n", msg) })
	if err != nil {
		t.Errorf("TestForwardVertex: NewObjFromBuf: %v", err)
		return
	}

	if !sliceEqualInt(forwardIndices, o.Indices) {
		t.Errorf("TestForwardVertex: indices: want=%v got=%v", forwardIndices, o.Indices)
	}

	if !sliceEqualFloat(forwardCoord, o.Coord) {
		t.Errorf("TestForwardVertex: coord: want=%v got=%v", forwardCoord, o.Coord)
	}
}

var cubeIndices = []int{0, 1, 2, 2, 3, 0, 4, 5, 6, 6, 7, 4, 8, 9, 10, 10, 11, 8, 12, 13, 14, 14, 15, 12, 16, 17, 18, 18, 19, 16, 20, 21, 16, 16, 22, 20}
var cubeCoord = []float32{1.0, -1.0, 1.0, -1.0, -1.0, 1.0, -1.0, -1.0, -1.0, 1.0, -1.0, -1.0, 1.0, 1.0, -1.0, -1.0, 1.0, -1.0, -1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, -1.0, -1.0, 1.0, 1.0, -1.0, 1.0, 1.0, 1.0, 1.0, -1.0, 1.0, -1.0, -1.0, 1.0, -1.0, 1.0, 1.0, -1.0, 1.0, -1.0, -1.0, -1.0, -1.0, 1.0, 1.0, -1.0, 1.0, 1.0, 1.0, -1.0, 1.0, 1.0, -1.0, -1.0, 1.0, -1.0, -1.0, -1.0, -1.0, 1.0, -1.0, 1.0, -1.0, -1.0}
var cubeText = []float32{0.5, 0.0, 0.5, 0.0, 0.5, 0.0, 0.5, 0.0, 0.5, 0.0, 0.5, 0.0, 0.5, 0.0, 0.5, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 1.0, 0.0, 1.0, 0.0, 1.0, 0.0, 1.0, 0.0, 1.0, 0.0, 1.0, 0.0, 1.0, 0.0}

var relativeIndices = []int{0, 1, 2, 0, 1, 2, 3, 4, 5, 3, 4, 5, 0, 1, 2, 0, 1, 2}
var relativeCoord = []float32{1.0, 1.0, 1.0, 2.0, 2.0, 2.0, 3.0, 3.0, 3.0, 4.0, 4.0, 4.0, 5.0, 5.0, 5.0, 6.0, 6.0, 6.0}

var forwardIndices = []int{0, 1, 2}
var forwardCoord = []float32{1.0, 1.0, 1.0, 2.0, 2.0, 2.0, 3.0, 3.0, 3.0}

var cubeObj = `
# cube.obj

mtllib cube.mtl

o cube

## comment-begin ##

# This is a multiline commented-out section.
# Notice this section is enclosed between "## comment-begin ##" and "## comment-end ##". 
# This section is fully ignored by this specific OBJ parser.
This uncommented line should cause error on common OBJ parsers.

## comment-end ##

# This is a regular section, processed under usual OBJ specification.

# square bottom
v -1 -1 -1
v -1 -1 1
v 1 -1 1
v 1 -1 -1

# square top
v -1 1 -1
v -1 1 1
v 1 1 1
v 1 1 -1

# uv coord

# red
vt 0.0 0.0

# green
vt 0.5 0.0

# blue
vt 1.0 0.0

usemtl cube_material

# face down
f -6/-2 -7/-2 -8/-2
f -8/-2 -5/-2 -6/-2

# face up
f -1/-2 -4/-2 -3/-2
f -3/-2 -2/-2 -1/-2 

# face right
f -5/-3 -1/-3 -2/-3
f -2/-3 -6/-3 -5/-3

# face left
f -7/-3 -3/-3 -4/-3
f -4/-3 -8/-3 -7/-3

# face front
f -1/-1 -2/-1 -3/-1
f -3/-1 -7/-1 -1/-1

# face back
f -8/-1 -4/-1 -1/-1
f -1/-1 -5/-1 -8/-1

## end-of-file ##

# This is an after-eof section.
# Notice this section follows the marker "## end-of-file ##".
# This section is fully ignored by this specific OBJ parser.
This uncommented line should cause error on common OBJ parsers.
`

var relativeObj = `
o relative_test
v 1 1 1
v 2 2 2
v 3 3 3
f 1 2 3
# this line should affect indices, but not vertex array
f -3 -2 -1
v 4 4 4
v 5 5 5
v 6 6 6
f 4 5 6
# this line should affect indices, but not vertex array
f -3 -2 -1
# these lines should affect indices, but not vertex array
f 1 2 3
f -6 -5 -4
`

var forwardObj = `
o forward_vertices_test
# face pointing to forward vertex definitions
# support for this isn't usual in OBJ parsers
# since it requires multiple passes
# but currently we do support this layout
f 1 2 3
v 1 1 1
v 2 2 2
v 3 3 3
`

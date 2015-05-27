package obj

import (
	"fmt"
	"reflect" // for reflect.DeepEqual
	"testing"
)

func expectInt(t *testing.T, label string, want, got int) {
	if want != got {
		t.Errorf("%s: want=%d got=%d", label, want, got)
	}
}

func TestCube(t *testing.T) {

	o, err := NewObjFromBuf([]byte(cubeObj), func(msg string) { fmt.Printf("TestCube NewObjFromBuf: error: %s\n", msg) })
	if err != nil {
		t.Errorf("TestCube: NewObjFromBuf: %v", err)
		return
	}

	//expectInt(t, "TestCube vertexCount", 23, o.vertexCount())
	//expectInt(t, "TestCube indexCount", 36, o.indexCount())

	if !reflect.DeepEqual(cubeCoord, o.Coord) {
		t.Errorf("TestCube: coord: want=%v got=%v", cubeCoord, o.Coord)
	}
}

func TestRelativeIndex(t *testing.T) {
	o, err := NewObjFromBuf([]byte(relativeObj), func(msg string) { fmt.Printf("TestRelativeIndex NewObjFromBuf: error: %s\n", msg) })
	if err != nil {
		t.Errorf("TestRelativeIndex: NewObjFromBuf: %v", err)
		return
	}

	if !reflect.DeepEqual(relativeCoord, o.Coord) {
		t.Errorf("TestRelativeIndex: coord: want=%v got=%v", cubeCoord, o.Coord)
	}
}

func TestForwardVertex(t *testing.T) {
	o, err := NewObjFromBuf([]byte(forwardObj), func(msg string) { fmt.Printf("TestForwardVertex NewObjFromBuf: error: %s\n", msg) })
	if err != nil {
		t.Errorf("TestForwardVertex: NewObjFromBuf: %v", err)
		return
	}

	if !reflect.DeepEqual(forwardCoord, o.Coord) {
		t.Errorf("TestForwardVertex: coord: want=%v got=%v", cubeCoord, o.Coord)
	}
}

var cubeCoord = []float32{1.1}
var relativeCoord = []float32{2.2}
var forwardCoord = []float32{3.3}

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

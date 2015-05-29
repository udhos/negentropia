package obj

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	//"log"
	"strconv"
	"strings"
	//"unicode"

	"negentropia/world/parser"
	"negentropia/world/util"
)

const FATAL = true
const NON_FATAL = false

type Group struct {
	Name       string
	Smooth     bool
	Usemtl     string
	indexBegin int
	indexCount int
}

type Obj struct {
	Indices       []int     // indices
	Coord         []float32 // vertex data pos=(x,y,z) tex=(tx,ty) norm=(nx,ny,nz)
	Mtllib        string
	Groups        []Group
	BigIndexFound bool
}

type objParser struct {
	lineBuf    []string
	lineCount  int
	vertCoord  []float32
	textCoord  []float32
	normCoord  []float32
	currGroup  *Group
	indexTable map[string]int
	indexCount int
	vertLines  int
	textLines  int
	normLines  int
	faceLines  int // stat-only
	triangles  int // stat-only
}

func (o *Obj) newGroup(name, usemtl string, begin int, smooth bool) *Group {
	gr := Group{Name: name, Usemtl: usemtl, indexBegin: begin, Smooth: smooth}
	o.Groups = append(o.Groups, gr)
	return &gr
}

func (o *Obj) Coord64(i int) float64 {
	return float64(o.Coord[i])
}

func (o *Obj) vertexCount() int {
	return -1
}

func (o *Obj) indexCount() int {
	return -1
}

//type lineParser func(p *objParser, o *Obj, rawLine string) (error, bool)

func NewObjFromBuf(buf []byte, logger func(string)) (*Obj, error) {
	return readObj(bytes.NewBuffer(buf), logger)
}

func NewObjFromReader(rd *bufio.Reader, logger func(string)) (*Obj, error) {
	return readObj(rd, logger)
}

type lineReader interface {
	ReadString(delim byte) (string, error)
}

func readObj(reader lineReader, logger func(msg string)) (*Obj, error) {
	p := &objParser{indexTable: make(map[string]int)}
	o := &Obj{}

	// 1. vertex-only parsing
	if err, fatal := readLines(p, o, reader, logger); err != nil {
		if fatal {
			return o, err
		}
	}

	// 2. full parsing
	if err, fatal := scanLines(p, o, reader, logger); err != nil {
		if fatal {
			return o, err
		}
	}

	// 3. output buffers
	/*
		o.Coord = make([]float32, len(p.vertCoord), len(p.vertCoord))
		for i, v := range p.vertCoord {
			o.Coord[i] = float32(v)
		}
	*/

	if logger != nil {
		logger(fmt.Sprintf("readObj: INPUT lines=%v vertLines=%v textLines=%v normLines=%v faceLines=%v triangles=%v",
			p.lineCount, p.vertLines, p.textLines, p.normLines, p.faceLines, p.triangles))

		logger(fmt.Sprintf("readObj: STATS numberOfIndices=%v indicesArraySize=%v", p.indexCount, len(o.Indices)))
	}

	return o, nil
}

func readLines(p *objParser, o *Obj, reader lineReader, logger func(msg string)) (error, bool) {
	p.lineCount = 0

	for {
		p.lineCount++
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			// parse last line
			if e, fatal := parseLineVertex(p, o, line); e != nil {
				if logger != nil {
					logger(fmt.Sprintf("readLines: %v", e))
				}
				return e, fatal
			}
			break // EOF
		}

		if err != nil {
			// unexpected IO error
			return errors.New(fmt.Sprintf("readLines: error: %v", err)), FATAL
		}

		if e, fatal := parseLineVertex(p, o, line); e != nil {
			if logger != nil {
				logger(fmt.Sprintf("readLines: %v", e))
			}
			if fatal {
				return e, fatal
			}
		}
	}

	return nil, NON_FATAL
}

// parse only vertex linux
func parseLineVertex(p *objParser, o *Obj, rawLine string) (error, bool) {
	line := strings.TrimSpace(rawLine)

	p.lineBuf = append(p.lineBuf, line) // save line for 2nd pass

	switch {
	case line == "" || line[0] == '#':
	case strings.HasPrefix(line, "s "):
	case strings.HasPrefix(line, "o "):
	case strings.HasPrefix(line, "g "):
	case strings.HasPrefix(line, "usemtl "):
	case strings.HasPrefix(line, "mtllib "):
	case strings.HasPrefix(line, "vt "):
	case strings.HasPrefix(line, "vn "):
	case strings.HasPrefix(line, "f "):
	case strings.HasPrefix(line, "v "):
		p.faceLines++

		result, err := parser.ParseFloatSliceSpace(line[2:])
		if err != nil {
			return fmt.Errorf("parseLine %v: [%v]: error: %v", p.lineCount, line, err), NON_FATAL
		}
		coordLen := len(result)
		switch coordLen {
		case 3:
			p.vertCoord = append(p.vertCoord, float32(result[0]), float32(result[1]), float32(result[2]))
		case 4:
			w := result[3]
			p.vertCoord = append(p.vertCoord, float32(result[0]/w), float32(result[1]/w), float32(result[2]/w))
		default:
			return fmt.Errorf("parseLine %v: [%v]: bad number of coords: %v", p.lineCount, line, coordLen), NON_FATAL
		}
	default:
		return fmt.Errorf("parseLine %v: [%v]: unexpected", p.lineCount, line), NON_FATAL
	}

	return nil, NON_FATAL
}

func scanLines(p *objParser, o *Obj, reader lineReader, logger func(msg string)) (error, bool) {

	p.currGroup = o.newGroup("", "", 0, false)

	p.lineCount = 0

	for _, line := range p.lineBuf {
		p.lineCount++

		if e, fatal := parseLine(p, o, line, logger); e != nil {
			if logger != nil {
				logger(fmt.Sprintf("scanLines: %v", e))
			}
			if fatal {
				return e, fatal
			}
		}
	}

	return nil, NON_FATAL
}

func solveRelativeIndex(index, size int) int {
	//fmt.Printf("index=%d size=%d\n", index, size)
	if index > 0 {
		return index - 1
	}
	return size + index
}

func splitSlash(s string) []string {
	isSlash := func(c rune) bool {
		return c == '/'
	}

	return strings.FieldsFunc(s, isSlash)

}

func pushIndex(p *objParser, o *Obj, i int) {
	if i > 65535 {
		o.BigIndexFound = true
	}
	o.Indices = append(o.Indices, i)
	p.currGroup.indexCount++
}

func addVertex(p *objParser, o *Obj, index string) error {
	ind := splitSlash(index)
	size := len(ind)
	if size < 1 || size > 3 {
		return fmt.Errorf("addVertex: line=%d bad index=[%s] size=%d", p.lineCount, index, size)
	}

	v, err := strconv.ParseInt(ind[0], 10, 32)
	if err != nil {
		return fmt.Errorf("addVertex: line=%d bad integer 1st index=[%s]: %v", p.lineCount, ind[0], err)
	}
	vi := solveRelativeIndex(int(v), p.vertLines)

	var ti int
	var tIndex string
	if size > 1 {
		t, e := strconv.ParseInt(ind[1], 10, 32)
		if e != nil {
			return fmt.Errorf("addVertex: line=%d bad integer 2nd index=[%s]: %v", p.lineCount, ind[1], e)
		}
		ti = solveRelativeIndex(int(t), p.textLines)
		tIndex = strconv.Itoa(ti)
	}

	var ni int
	var nIndex string
	if size > 2 {
		n, e := strconv.ParseInt(ind[2], 10, 32)
		if e != nil {
			return fmt.Errorf("addVertex: line=%d bad integer 3rd index=[%s]: %v", p.lineCount, ind[2], e)
		}
		ni = solveRelativeIndex(int(n), p.normLines)
		nIndex = strconv.Itoa(ni)
	}

	absIndex := fmt.Sprintf("%d/%s/%s", vi, tIndex, nIndex)

	// known unified index?
	if i, ok := p.indexTable[absIndex]; ok {
		pushIndex(p, o, i)
		return nil
	}

	vOffset := vi * 3
	o.Coord = append(o.Coord, p.vertCoord[vOffset+0]) // x
	o.Coord = append(o.Coord, p.vertCoord[vOffset+1]) // y
	o.Coord = append(o.Coord, p.vertCoord[vOffset+2]) // z

	if tIndex != "" {
		tOffset := ti * 2
		//fmt.Printf("ti=%d tOffset=%d textCoord=%v len=%d\n", ti, tOffset, p.textCoord, len(p.textCoord))
		o.Coord = append(o.Coord, p.textCoord[tOffset+0]) // u
		o.Coord = append(o.Coord, p.textCoord[tOffset+1]) // v
	}

	if nIndex != "" {
		nOffset := ni * 3
		o.Coord = append(o.Coord, p.normCoord[nOffset+0]) // x
		o.Coord = append(o.Coord, p.normCoord[nOffset+1]) // y
		o.Coord = append(o.Coord, p.normCoord[nOffset+2]) // z
	}

	// add unified index
	pushIndex(p, o, p.indexCount)
	//fmt.Printf("absIndex=%s indexCount=%d\n", absIndex, p.indexCount)
	p.indexTable[absIndex] = p.indexCount
	p.indexCount++

	return nil
}

func parseLine(p *objParser, o *Obj, line string, logger func(msg string)) (error, bool) {

	switch {
	case line == "" || line[0] == '#':
	case strings.HasPrefix(line, "s "):
		smooth := line[2:]
		if s, err := strconv.ParseBool(smooth); err == nil {
			if p.currGroup.Smooth != s {
				// create new group
				p.currGroup = o.newGroup(p.currGroup.Name, p.currGroup.Usemtl, len(o.Indices), s)
			}
		} else {
			return fmt.Errorf("parseLine: line=%d bad boolean smooth=[%s]: %v", p.lineCount, smooth, err), NON_FATAL
		}
	case strings.HasPrefix(line, "o ") || strings.HasPrefix(line, "g "):
		name := line[2:]
		p.currGroup = o.newGroup(name, p.currGroup.Usemtl, len(o.Indices), p.currGroup.Smooth)
	case strings.HasPrefix(line, "usemtl "):
		usemtl := line[7:]
		if p.currGroup.Usemtl == "" {
			// only set the missing material name for group
			p.currGroup.Usemtl = usemtl
		} else if p.currGroup.Usemtl != usemtl {
			// create new group for material
			p.currGroup = o.newGroup(p.currGroup.Name, usemtl, len(o.Indices), p.currGroup.Smooth)
		}
	case strings.HasPrefix(line, "mtllib "):
		mtllib := line[7:]
		if o.Mtllib != "" && logger != nil {
			logger(fmt.Sprintf("parseLine: line=%d mtllib redefinition old=%s new=%s", p.lineCount, o.Mtllib, mtllib))
		}
		o.Mtllib = mtllib
	case strings.HasPrefix(line, "vt "):
		p.textLines++

		tex := line[3:]
		t, err := parser.ParseFloatSliceSpace(tex)
		if err != nil {
			return fmt.Errorf("parseLine: line=%d bad vertex texture=[%s]: %v", p.lineCount, tex, err), NON_FATAL
		}
		size := len(t)
		if size < 2 || size > 3 {
			return fmt.Errorf("parseLine: line=%d bad vertex texture=[%s] size=%d", p.lineCount, tex, size), NON_FATAL
		}
		if size > 2 {
			if w := t[2]; !util.CloseToZero(w) {
				logger(fmt.Sprintf("parseLine: line=%d non-zero third texture coordinate w=%f", p.lineCount, w))
			}
		}
		p.textCoord = append(p.textCoord, float32(t[0]), float32(t[1]))
	case strings.HasPrefix(line, "vn "):
		p.normLines++

		norm := line[3:]
		n, err := parser.ParseFloatVector3Space(norm)
		if err != nil {
			return fmt.Errorf("parseLine: line=%d bad vertex normal=[%s]: %v", p.lineCount, norm, err), NON_FATAL
		}
		p.normCoord = append(p.normCoord, float32(n[0]), float32(n[1]), float32(n[2]))
	case strings.HasPrefix(line, "f "):
		face := line[2:]
		f := strings.Fields(face)
		size := len(f)
		if size < 3 || size > 4 {
			return fmt.Errorf("parseLine: line=%d bad face=[%s] size=%d", p.lineCount, face, size), NON_FATAL
		}
		// triangle face: v0 v1 v2
		// quad face:
		// v0 v1 v2 v3 =>
		// v0 v1 v2
		// v2 v3 v0
		p.triangles++
		if err := addVertex(p, o, f[0]); err != nil {
			return fmt.Errorf("parseLine: line=%d bad face=[%s] index_v0=[%s]: %v", p.lineCount, face, f[0], err), NON_FATAL
		}
		if err := addVertex(p, o, f[1]); err != nil {
			return fmt.Errorf("parseLine: line=%d bad face=[%s] index_v1=[%s]: %v", p.lineCount, face, f[1], err), NON_FATAL
		}
		if err := addVertex(p, o, f[2]); err != nil {
			return fmt.Errorf("parseLine: line=%d bad face=[%s] index_v2=[%s]: %v", p.lineCount, face, f[2], err), NON_FATAL
		}
		if size > 3 {
			// quad face
			p.triangles++
			if err := addVertex(p, o, f[2]); err != nil {
				return fmt.Errorf("parseLine: line=%d bad face=[%s] index_v2=[%s]: %v", p.lineCount, face, f[2], err), NON_FATAL
			}
			if err := addVertex(p, o, f[3]); err != nil {
				return fmt.Errorf("parseLine: line=%d bad face=[%s] index_v3=[%s]: %v", p.lineCount, face, f[3], err), NON_FATAL
			}
			if err := addVertex(p, o, f[0]); err != nil {
				return fmt.Errorf("parseLine: line=%d bad face=[%s] index_v0=[%s]: %v", p.lineCount, face, f[0], err), NON_FATAL
			}
		}
	case strings.HasPrefix(line, "v "):
		p.vertLines++
	default:
		return fmt.Errorf("parseLine %v: [%v]: unexpected", p.lineCount, line), NON_FATAL
	}

	return nil, NON_FATAL
}

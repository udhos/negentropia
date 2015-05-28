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
	Indices []int     // indices
	Coord   []float32 // vertex data pos=(x,y,z) tex=(tx,ty) norm=(nx,ny,nz)
	Mtllib  string
	Groups  []Group
}

type objParser struct {
	lineBuf   []string
	lineCount int
	vertCoord []float64
	textCoord []float64
	currGroup *Group
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
	p := &objParser{}
	o := &Obj{}

	// 1. vertex-only parsing
	if err, fatal := readLines(p, o, reader, logger); err != nil {
		if fatal {
			return o, err
		}
	}

	if logger != nil {
		logger(fmt.Sprintf("readObj: found %v lines", p.lineCount))
	}

	// 2. full parsing
	if err, fatal := scanLines(p, o, reader, logger); err != nil {
		if fatal {
			return o, err
		}
	}

	// 3. output buffers

	o.Coord = make([]float32, len(p.vertCoord), len(p.vertCoord))
	for i, v := range p.vertCoord {
		o.Coord[i] = float32(v)
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
		result, err := parser.ParseFloatSliceSpace(line[2:])
		if err != nil {
			return fmt.Errorf("parseLine %v: [%v]: error: %v", p.lineCount, line, err), NON_FATAL
		}
		coordLen := len(result)
		switch coordLen {
		case 3:
			p.vertCoord = append(p.vertCoord, result[0], result[1], result[2])
		case 4:
			w := result[3]
			p.vertCoord = append(p.vertCoord, result[0]/w, result[1]/w, result[2]/w)
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
		p.textCoord = append(p.textCoord, t[0], t[1])
	case strings.HasPrefix(line, "vn "):
	case strings.HasPrefix(line, "f "):
	case strings.HasPrefix(line, "v "):
	default:
		return fmt.Errorf("parseLine %v: [%v]: unexpected", p.lineCount, line), NON_FATAL
	}

	return nil, NON_FATAL
}

package obj

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	//"log"
	//"strconv"
	"strings"
	//"unicode"

	"negentropia/world/parser"
)

const FATAL = true
const NON_FATAL = false

type Obj struct {
	Indices []int     // indices
	Coord   []float32 // vertex data pos=(x,y,z) tex=(tx,ty) norm=(nx,ny,nz)
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

type objParser struct {
	lineBuf   []string
	lineCount int
	vertCoord []float64
}

//type lineParser func(p *objParser, o *Obj, rawLine string) (error, bool)

func NewObjFromBuf(buf []byte, logger func(string)) (*Obj, error) {
	return readObj(bytes.NewBuffer(buf), logger)
}

func NewObjFromReader(rd *bufio.Reader, logger func(string)) (*Obj, error) {
	return readObj(rd, logger)
}

func readObj(reader lineReader, logger func(msg string)) (*Obj, error) {
	p := &objParser{lineCount: 0}
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

	// 3. output buffers

	o.Coord = make([]float32, len(p.vertCoord), len(p.vertCoord))
	for i, v := range p.vertCoord {
		o.Coord[i] = float32(v)
	}

	return o, nil
}

func readLines(p *objParser, o *Obj, reader lineReader, logger func(msg string)) (error, bool) {
	for {
		p.lineCount++
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			// parse last line
			if e, fatal := parseLine(p, o, line); e != nil {
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

		//log.Printf("DEBUG scanLines %v: [%v]\n", p.lineCount, line)

		if e, fatal := parseLine(p, o, line); e != nil {
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

func parseLine(p *objParser, o *Obj, rawLine string) (error, bool) {
	line := strings.TrimSpace(rawLine)

	p.lineBuf = append(p.lineBuf, line) // save line

	//log.Printf("DEBUG parseLine %v: [%v]\n", p.lineCount, line)
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
		//x, y, z := float32(result[0]), float32(result[1]), float32(result[2])
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

type lineReader interface {
	ReadString(delim byte) (string, error)
}

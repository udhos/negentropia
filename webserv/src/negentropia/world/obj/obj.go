package obj

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	//"strconv"
	"strings"
	//"unicode"

	"negentropia/world/parser"
)

type Obj struct {
	Coord []float64 // vertex coordinates
}

func (o *Obj) parseLine(rawLine string, lineNumber int) {
	line := strings.TrimSpace(rawLine)
	//log.Printf("parseLine %v: [%v]\n", lineNumber, line)
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
			log.Printf("parseLine %v: [%v]: error: %v", lineNumber, line, err)
			return
		}
		coordLen := len(result)
		switch coordLen {
		case 3:
			o.Coord = append(o.Coord, result[0], result[1], result[2])
		case 4:
			w := result[3]
			o.Coord = append(o.Coord, result[0]/w, result[1]/w, result[2]/w)
		default:
			log.Printf("parseLine %v: [%v]: bad number of coords: %v", lineNumber, line, coordLen)
		}
	default:
		log.Printf("parseLine %v: [%v]: unexpected", lineNumber, line)
	}
}

func NewObjFromBuf(buf []byte) (*Obj, error) {
	var o Obj

	b := bytes.NewBuffer(buf)
	lineCount := 0
	for {
		lineCount++
		var line string
		var err error
		line, err = b.ReadString('\n')
		if err == io.EOF {
			if line != "" {
				o.parseLine(line, lineCount)
			}
			break
		}
		if err != nil {
			return nil, errors.New(fmt.Sprintf("NewObjFromBuf: error: %v", err))
		}

		o.parseLine(line, lineCount)
	}

	return &o, nil
}

func NewObjFromReader(rd *bufio.Reader) (*Obj, error) {
	var o Obj

	lineCount := 0
	for {
		lineCount++
		var line string
		var err error
		line, err = rd.ReadString('\n')
		if err == io.EOF {
			if line != "" {
				o.parseLine(line, lineCount)
			}
			break
		}
		if err != nil {
			return nil, errors.New(fmt.Sprintf("NewObjFromReader: error: %v", err))
		}

		o.parseLine(line, lineCount)
	}
	log.Printf("NewObjFromReader: %v lines", lineCount)

	return &o, nil
}

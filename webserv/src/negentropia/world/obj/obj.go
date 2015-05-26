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

type Obj struct {
	Coord []float64 // vertex coordinates
}

const FATAL = true
const NON_FATAL = false

func (o *Obj) parseLine(rawLine string, lineNumber int) (error, bool) {
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
			return fmt.Errorf("parseLine %v: [%v]: error: %v", lineNumber, line, err), NON_FATAL
		}
		coordLen := len(result)
		switch coordLen {
		case 3:
			o.Coord = append(o.Coord, result[0], result[1], result[2])
		case 4:
			w := result[3]
			o.Coord = append(o.Coord, result[0]/w, result[1]/w, result[2]/w)
		default:
			return fmt.Errorf("parseLine %v: [%v]: bad number of coords: %v", lineNumber, line, coordLen), NON_FATAL
		}
	default:
		return fmt.Errorf("parseLine %v: [%v]: unexpected", lineNumber, line), NON_FATAL
	}

	return nil, NON_FATAL
}

type lineReader interface {
	ReadString(delim byte) (string, error)
}

func readObj(reader lineReader, logger func(msg string)) (*Obj, error) {
	var o Obj
	lineCount := 0
	for {
		lineCount++
		var line string
		var err error
		line, err = reader.ReadString('\n')
		if err == io.EOF {
			if line != "" {
				// parse last line
				o.parseLine(line, lineCount)
				if e, _ := o.parseLine(line, lineCount); e != nil {
					if logger != nil {
						logger(fmt.Sprintf("readObj: %v", e))
					}
				}
			}
			break
		}
		if err != nil {
			return nil, errors.New(fmt.Sprintf("readObj: error: %v", err))
		}

		if e, fatal := o.parseLine(line, lineCount); e != nil {
			if logger != nil {
				logger(fmt.Sprintf("readObj: %v", e))
			}
			if fatal {
				return &o, e
			}
		}
	}

	if logger != nil {
		logger(fmt.Sprintf("readObj: %v lines", lineCount))
	}

	return &o, nil
}

func NewObjFromBuf(buf []byte, logger func(string)) (*Obj, error) {

	b := bytes.NewBuffer(buf)
	/*
		var o Obj
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

			if e, fatal = o.parseLine(line, lineCount); err != nil {
			}
		}

		if logger != nil {
			logger(fmt.Sprintf("NewObjFromBuf: %v lines", lineCount))
		}

		return &o, nil
	*/

	return readObj(b, logger)
}

func NewObjFromReader(rd *bufio.Reader, logger func(string)) (*Obj, error) {
	/*
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

		if logger != nil {
			logger(fmt.Sprintf("NewObjFromReader: %v lines", lineCount))
		}

		return &o, nil
	*/

	return readObj(rd, logger)
}

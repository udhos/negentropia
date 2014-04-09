package obj

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
)

type Obj struct {
	VertCoord []float64
}

func parseLine(line string, lineCount int) {
	log.Printf("line %v: [%v]\n", lineCount, line)
}

func NewObjFromBuf(buf []byte) (*Obj, error) {
	b := bytes.NewBuffer(buf)
	lineCount := 0
	for {
		lineCount++
		var line string
		var err error
		line, err = b.ReadString('\n')
		if err == io.EOF {
			if line != "" {
				parseLine(line, lineCount)
			}
			break
		}
		if err != nil {
			return nil, errors.New(fmt.Sprintf("NewObjFromBuf: error: %v", err))
		}

		parseLine(line, lineCount)
	}
	log.Printf("NewObjFromBuf: %v lines", lineCount)

	return nil, errors.New("NewObjFromString: FIXME WRITEME")
}

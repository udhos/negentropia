package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func main() {

	if len(os.Args) < 2 {
		basename := filepath.Base(os.Args[0])
		log.Fatalf("usage: %s input-files", basename)
	}

	for _, input := range os.Args[1:] {
		readFrom(input)
	}

}

func readFrom(filename string) {

	var scanner *bufio.Scanner

	if filename == "-" || filename == "--" {
		log.Printf("reading from stdin")
		scanner = bufio.NewScanner(os.Stdin)
	} else {
		log.Printf("readFrom: reading from: %s", filename)
		input, err := os.Open(filename)
		if err != nil {
			log.Printf("readFrom: error: %v", err)
			return
		}
		scanner = bufio.NewScanner(input)
	}

	lineNum := 0

	for scanner.Scan() {
		lineNum++
		fmt.Printf("line %d: [%v]\n", lineNum, scanner.Text())
	}

	log.Printf("readFrom: done")

	if err := scanner.Err(); err != nil {
		log.Printf("readFrom: error: %v", err)
	}
}

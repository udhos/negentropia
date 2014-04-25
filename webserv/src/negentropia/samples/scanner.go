package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	log.Printf("reading from stdin")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		fmt.Printf("line=[%v]\n", scanner.Text())
	}
	log.Printf("reading from stdin - done")

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

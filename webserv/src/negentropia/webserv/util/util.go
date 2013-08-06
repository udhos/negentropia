package util

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"log"
	"math/rand"
	"strings"
	"time"
)

var (
	randCh chan int64 = make(chan int64)
)

func serveRand() {
	log.Printf("util.serveRand: goroutine started")

	randGen := rand.New(rand.NewSource(time.Now().Unix()))

	for {
		randCh <- randGen.Int63()
	}
}

func init() {
	go serveRand()
}

func GetPort(hostPort string) string {
	pair := strings.Split(hostPort, ":")
	if len(pair) < 2 {
		return ""
	}

	return ":" + pair[1]
}

func RandomSuffix() string {

	n := <-randCh

	buf := &bytes.Buffer{} // buf := new(bytes.Buffer) // which is better??
	if err := binary.Write(buf, binary.BigEndian, n); err != nil {
		log.Printf("handler.RandomSuffix: binary.Write: failed: %s", err)
	}

	return "-" + base64.URLEncoding.EncodeToString(buf.Bytes())
}

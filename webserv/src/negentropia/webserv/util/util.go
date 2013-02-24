package util

import (
	"log"
	"time"
	"bytes"
	"strings"
	"math/rand"
	"encoding/base64"
	"encoding/binary"
)

var (
	randGen *rand.Rand
)

func init() {
	randGen = rand.New(rand.NewSource(time.Now().Unix()))
}

func GetPort(hostPort string) string {
	pair := strings.Split(hostPort, ":")
	if len(pair) < 2 {
		return ""
	}

	return ":" + pair[1]
}

func RandomSuffix() string {
	log.Printf("handler.RandomSuffix: FIXME: randGen.int63() is goroutine unsafe")

	n := randGen.Int63()
	
	// buf := &bytes.Buffer{} // which is better??
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, n); err != nil {
		log.Printf("handler.RandomSuffix: binary.Write: failed: %s", err)
	}
			
	return "-" + base64.StdEncoding.EncodeToString(buf.Bytes())
}

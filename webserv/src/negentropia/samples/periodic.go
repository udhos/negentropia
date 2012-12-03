// Dustin post on golang-nuts

package main

import (
	"log"
	"time"
)

func doSomething(s string) {
	log.Printf("doing something: %v", s)
}

func startPolling() {
	for _ = range time.Tick(2 * time.Second) {
		doSomething("awesome")
	}
}

func main() {
	go startPolling()

	select {}
}

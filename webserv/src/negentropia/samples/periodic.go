// Dustin post on golang-nuts
// See also tickers.go
// See also https://gobyexample.com/tickers

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

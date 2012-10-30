package main

import "fmt"

func generateIds(ids chan int) {
	id := 1
	for {
		ids <- id
		id++
	}
}

func main() {
	ids := make(chan int)
	go generateIds(ids)

	done := make(chan bool)

	// now anyone can simply receive on ids to get the next numeric id
	for i := 0; i < 10; i++ {
		go func(i int) {
			fmt.Printf("%d got counter: %d\n", i, <-ids)
			
			//Inform the caller that we're done :o)
			done <- true
		}(i)
	}

	for i := 0; i < 10; i++ {
		<-done
	}
}

package main

import "time"

func main() {

	done := make(chan bool)
	go loop(done)
	blah(done)
	time.Sleep(10 * time.Second)
	close(done)
	time.Sleep(1 * time.Second)
	println("Finished Test")
}

func blah(done chan<- bool) {
	time.Sleep(5 * time.Second)
	done <- true
}

func loop(done <-chan bool) {
	for {
		select {
		case <-done:
			println("Finished")
			return
		default:
			println("work")
			time.Sleep(1 * time.Second)
		}
	}
}

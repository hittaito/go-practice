package main

import (
	"log"
	"time"
)

func main() {
	doneCh := make(chan struct{})
	for i := 0; i < 10; i++ {
		i := i
		go do(i, doneCh)
	}
	time.Sleep(200 * time.Microsecond)
	close(doneCh)
	time.Sleep(300 * time.Microsecond)
}
func do(n int, doneCh <-chan struct{}) {
	for {
		select {
		case <-doneCh:
			log.Printf("finidhed %d", n)
			return
		default:
			log.Printf("wait %d", n)
			time.Sleep((100 * time.Microsecond))
		}
	}
}

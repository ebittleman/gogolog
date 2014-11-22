package main

import (
	"log"
	"os"
	"time"

	"github.com/ebittleman/gogolog"
)

func main() {
	logger := gogolog.NewLogger(
		gogolog.DEBUG,
		"example ",
		log.Lmicroseconds,

		gogolog.NewWriter(gogolog.DEBUG, os.Stderr),
		gogolog.NewWriter(gogolog.CRIT, gogolog.NewFlushingWriter(os.Stderr, 2, time.Second)),
	)

	a, b := make(chan int), make(chan int)

	go func() {
		a <- Add(1, 2)
	}()

	go func() {
		b <- Add(1, 2)
	}()

	logger.Emergf("My Answer: %d", Add(<-a, <-b))
	<-time.After(time.Second * 2)

}

func Add(a, b int) int {
	return a + b
}

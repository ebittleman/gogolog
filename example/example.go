package main

import (
	"log"
	"os"
	"time"

	"github.com/ebittleman/gogolog"
	"github.com/rcrowley/go-metrics"
)

func main() {
	file, err := os.OpenFile("./metrics", os.O_APPEND|os.O_CREATE|os.O_RDWR, 775)
	if err != nil {
		panic(err)
	}
	flushingWriter := gogolog.NewFlushingWriter(file, 2, time.Second)
	defer flushingWriter.Flush()
	defer file.Close()

	writers := []*gogolog.Writer{
		gogolog.NewWriter(gogolog.INFO, os.Stderr),
		gogolog.NewWriter(gogolog.ALERT, flushingWriter),
	}

	logger := gogolog.New(
		gogolog.DEBUG,
		"example ",
		log.Lmicroseconds,
		writers...,
	)

	a, b := make(chan int), make(chan int)

	go func() {
		a <- Add(1, 2)
	}()

	go func() {
		b <- Add(1, 2)
	}()

	logger.Infof("My Answer: %d", Add(<-a, <-b))
	doMetrics(logger.GetLogger(gogolog.ALERT))
	logger.Info("Done")
}

func Add(a, b int) int {
	return a + b
}

func doMetrics(logger *log.Logger) {
	c := metrics.NewCounter()
	metrics.Register("foo", c)
	c.Inc(47)

	g := metrics.NewGauge()
	metrics.Register("bar", g)
	g.Update(47)

	s := metrics.NewExpDecaySample(1028, 0.015) // or metrics.NewUniformSample(1028)
	h := metrics.NewHistogram(s)
	metrics.Register("baz", h)
	h.Update(47)

	m := metrics.NewMeter()
	metrics.Register("quux", m)
	m.Mark(47)

	t := metrics.NewTimer()
	metrics.Register("bang", t)
	t.Time(func() {})
	t.Update(47)
	go metrics.Log(metrics.DefaultRegistry, time.Second*2, logger)

	<-time.After(time.Second * 2)

	h.Update(200)

	<-time.After(time.Second * 2)
}

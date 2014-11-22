package gogolog

import (
	"bytes"
	"io"
	"sync"
	"time"
)

type FlushingWriter struct {
	output           io.Writer
	flushOnNumWrites int
	flushTimeout     time.Duration

	buf      *bytes.Buffer
	lock     *sync.Mutex
	flushing chan bool
	flushed  chan bool
	timeout  <-chan time.Time
	done     chan struct{}

	numWrites int
}

func NewFlushingWriter(
	output io.Writer,
	flushOnNumWrites int,
	flushTimeout time.Duration,
) io.Writer {
	f := &FlushingWriter{
		output:           output,
		flushOnNumWrites: flushOnNumWrites,
		flushTimeout:     flushTimeout,

		buf:      &bytes.Buffer{},
		lock:     &sync.Mutex{},
		flushing: make(chan bool),
		flushed:  make(chan bool),
		timeout:  nil,
		done:     make(chan struct{}),

		numWrites: 0,
	}

	go f.loop()

	return f
}

func (f *FlushingWriter) loop() {
	f.timeout = time.After(f.flushTimeout)
	for {
		select {
		case <-f.flushed:
			f.timeout = time.After(f.flushTimeout)
		case <-f.flushing:
			f.timeout = nil
		case <-f.timeout:
			f.Flush()
			f.timeout = time.After(f.flushTimeout)
		case <-f.done:
			f.Flush()
			return
		}
	}
}

func (f *FlushingWriter) Done() {
	close(f.done)
}

func (f *FlushingWriter) Write(bytes []byte) (int, error) {
	f.lock.Lock()
	n, err := f.buf.Write(bytes)

	if f.numWrites++; err == nil && f.numWrites >= f.flushOnNumWrites {
		f.flushing <- true
		f.flush()
		f.flushed <- true
	}

	f.lock.Unlock()
	return n, err
}

func (f *FlushingWriter) Flush() (int, error) {
	f.lock.Lock()
	n, err := f.flush()
	f.lock.Unlock()

	return n, err
}

func (f *FlushingWriter) flush() (int, error) {
	n, err := f.output.Write(f.buf.Bytes())

	if err != nil {
		return 0, err
	}

	f.numWrites = 0
	f.buf.Reset()

	return n, nil
}

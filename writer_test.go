package gogolog

import (
	"bytes"
	"testing"
	"time"
)

func TestFlushingWriter_ShouldTimeout(t *testing.T) {
	write := "Here is something in the buffer"
	expected := write
	runTimeoutTest(t, write, expected, 2, time.Millisecond, time.After(time.Millisecond*10))
}

func TestFlushingWriter_ShouldNotTimeout(t *testing.T) {
	write := "Here is something in the buffer"
	expected := ""
	runTimeoutTest(t, write, expected, 2, time.Second, time.After(0))
}

func TestFlushingWriter_ShouldFlush(t *testing.T) {
	write := "Here is something in the buffer"
	expected := write
	runTimeoutTest(t, write, expected, 1, time.Second, time.After(0))
}

func runTimeoutTest(
	t *testing.T,
	write string,
	expected string,
	numWrites int,
	flushTimeout time.Duration,
	waitFor <-chan time.Time,
) {
	buf := &bytes.Buffer{}
	w := NewFlushingWriter(buf, numWrites, flushTimeout)

	w.Write([]byte(write))

	<-waitFor

	actual := buf.String()

	if actual != expected {
		t.Fatalf("expected: `%s`, actual: `%s`", expected, actual)
	}
}

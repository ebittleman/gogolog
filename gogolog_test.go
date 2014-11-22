package gogolog

import (
	"bytes"
	"log"
	"testing"
)

func TestLogger_ImplementsLogger(t *testing.T) {
	var l interface{}
	l = &logger{}

	_, ok := l.(Logger)

	if !ok {
		t.Fatal("struct gogolog.logger does not implement interface gogolog.Logger")
	}

}

func TestLogger_SetLevel(t *testing.T) {
	l := &logger{}

	l.SetLevel(DEBUG)

	if l.level != DEBUG {
		t.Fatalf("Expected: %d, Got: %d", DEBUG, l.level)
	}
}

func TestLogger_Log_ShouldOutput(t *testing.T) {
	buf := &bytes.Buffer{}

	w := NewWriter(WARN, buf)
	l := New(INFO, "gogolog ", log.LstdFlags, w)

	l.Error("Hello Logger!")

	if buf.Len() < 1 {
		t.Fatal("Should have some log output")
	}
}

func TestLogger_Log_ShouldNotOutput(t *testing.T) {
	buf := &bytes.Buffer{}

	w := NewWriter(CRIT, buf)
	l := New(WARN, "gogolog ", log.LstdFlags, w)

	l.Error("Hello Logger!")

	if buf.Len() > 0 {
		t.Fatal("Should not have logged anything")
	}
}

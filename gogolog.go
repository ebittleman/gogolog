package gogolog

import (
	"fmt"
	"io"
	"log"
	"strings"
)

type Level int

const (
	EMERG Level = iota
	ALERT
	CRIT
	ERROR
	WARN
	INFO
	DEBUG
)

type Logger interface {
	SetLevel(Level)
	AddWriter(*Writer)
	GetLogger(Level) *log.Logger

	Log(Level, string)
	Logf(Level, string, ...interface{})

	Emerg(string)
	Emergf(string, ...interface{})

	Alert(string)
	Alertf(string, ...interface{})

	Crit(string)
	Critf(string, ...interface{})

	Error(string)
	Errorf(string, ...interface{})

	Warn(string)
	Warnf(string, ...interface{})

	Info(string)
	Infof(string, ...interface{})

	Debug(string)
	Debugf(string, ...interface{})

	io.Writer
}

type Writer struct {
	level Level
	io.Writer
	log *log.Logger
}

type logger struct {
	level   Level
	prefix  string
	flag    int
	writers []*Writer
}

func New(
	level Level,
	prefix string,
	flag int,
	writers ...*Writer,
) Logger {

	logger := &logger{
		level:   level,
		prefix:  prefix,
		flag:    flag,
		writers: make([]*Writer, 0, len(writers)),
	}

	for _, w := range writers {
		logger.AddWriter(w)
	}

	return logger
}

func NewWriter(level Level, w io.Writer) *Writer {
	return &Writer{Writer: w, level: level}
}

func (l *logger) SetLevel(level Level) {
	l.level = level
}

func (l *logger) AddWriter(w *Writer) {
	w.log = log.New(w, l.prefix, l.flag)
	l.writers = append(l.writers, w)
}

func (l *logger) String() string {
	return fmt.Sprintf("%d", l.level)
}

func (l *logger) Log(level Level, msg string) {
	if level > l.level {
		return
	}

	for _, w := range l.writers {
		if level > w.level {
			continue
		}

		w.log.Println(msg)
	}
}

func (l *logger) Logf(level Level, format string, args ...interface{}) {
	l.Log(level, fmt.Sprintf(format, args...))
}

func (l *logger) Emerg(msg string) {
	l.Log(EMERG, msg)
}

func (l *logger) Emergf(format string, args ...interface{}) {
	l.Logf(EMERG, format, args)
}

func (l *logger) Alert(msg string) {
	l.Log(ALERT, msg)
}

func (l *logger) Alertf(format string, args ...interface{}) {
	l.Logf(ALERT, format, args)
}

func (l *logger) Crit(msg string) {
	l.Log(CRIT, msg)
}

func (l *logger) Critf(format string, args ...interface{}) {
	l.Logf(CRIT, format, args)
}

func (l *logger) Error(msg string) {
	l.Log(ERROR, msg)
}

func (l *logger) Errorf(format string, args ...interface{}) {
	l.Logf(ERROR, format, args)
}

func (l *logger) Warn(msg string) {
	l.Log(WARN, msg)
}

func (l *logger) Warnf(format string, args ...interface{}) {
	l.Logf(WARN, format, args)
}

func (l *logger) Info(msg string) {
	l.Log(INFO, msg)
}

func (l *logger) Infof(format string, args ...interface{}) {
	l.Logf(INFO, format, args)
}

func (l *logger) Debug(msg string) {
	l.Log(DEBUG, msg)
}

func (l *logger) Debugf(format string, args ...interface{}) {
	l.Logf(DEBUG, format, args)
}

func (l *logger) Write(bytes []byte) (int, error) {
	l.Log(l.level, strings.Trim(string(bytes), "\n"))
	return len(bytes), nil
}

func (l *logger) GetLogger(level Level) *log.Logger {
	return log.New(&WrappedLogger{l, level}, "", 0)
}

type WrappedLogger struct {
	l     *logger
	level Level
}

func (wl *WrappedLogger) Write(bytes []byte) (int, error) {
	wl.l.Log(wl.level, strings.Trim(string(bytes), "\n"))
	return len(bytes), nil
}

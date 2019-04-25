package loggers

import (
	"encoding/json"
	"io"
	"sync"

	"sour.is/x/toolbox/log/event"
	"sour.is/x/toolbox/log/scheme"
)

// StdLogger defines an output handler
type StdLogger struct {
	out    io.Writer
	scheme scheme.Scheme
	level  event.Level
	sync.Mutex
}

// NewStdLogger returns a StdLogger
func NewStdLogger(out io.Writer, scheme scheme.Scheme, level event.Level) event.Logger {
	return &StdLogger{out: out, scheme: scheme, level: level}
}

// WriteEvent ouputs an event to stdlogger
func (l *StdLogger) WriteEvent(e *event.Event) {
	if l.level < e.Level {
		return
	}

	l.Lock()
	defer l.Unlock()

	l.out.Write([]byte(l.scheme.FmtEvent(*e)))
}

// SetVerbose sets the event verbose level
func (l *StdLogger) SetVerbose(level event.Level) {
	l.Lock()
	defer l.Unlock()

	l.level = level
}

// JSONLogger defines an output handler
type JSONLogger struct {
	out   io.Writer
	level event.Level
	sync.Mutex
}

// NewJSONLogger returns a StdLogger
func NewJSONLogger(out io.Writer, level event.Level) event.Logger {
	return &JSONLogger{out: out, level: level}
}

// WriteEvent ouputs an event to JSONlogger
func (l *JSONLogger) WriteEvent(e *event.Event) {
	if l.level < e.Level {
		return
	}

	l.Lock()
	defer l.Unlock()

	json.NewEncoder(l.out).Encode(*e)
}

// SetVerbose sets the event verbose level
func (l *JSONLogger) SetVerbose(level event.Level) {
	l.Lock()
	defer l.Unlock()

	l.level = level
}

// FanLogger outputs to a list of loggers
type FanLogger struct {
	sync.RWMutex
	level   event.Level
	loggers []event.Logger
}

// NewFanLogger returns a new logger.
func NewFanLogger(level event.Level, loggers ...event.Logger) event.Logger {
	return &FanLogger{level: level, loggers: loggers}
}

// WriteEvent ouputs an event to FanLogger
func (l *FanLogger) WriteEvent(e *event.Event) {
	if l.level < e.Level {
		return
	}

	l.RLock()
	defer l.RUnlock()

	for i := range l.loggers {
		l.loggers[i].WriteEvent(e)
	}
}

// SetVerbose sets the event verbose level
func (l *FanLogger) SetVerbose(level event.Level) {
	l.Lock()
	defer l.Unlock()

	l.level = level
}

// Add logger to fan
func (l *FanLogger) Add(logger event.Logger) {
	l.Lock()
	defer l.Unlock()

	l.loggers = append(l.loggers, logger)
}

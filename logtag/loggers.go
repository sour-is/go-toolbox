package logtag

import (
	"encoding/json"
	"io"
	"os"
	"sync"

	"sour.is/x/toolbox/mqtt"
)

// Logger outputs events
type Logger interface {
	WriteEvent(*Event)
}

// StdLogger defines an output handler
type StdLogger struct {
	out    io.Writer
	scheme Scheme
	sync.Mutex
}

// WriteEvent ouputs an event to stdlogger
func (s *StdLogger) WriteEvent(e *Event) {
	s.Lock()
	defer s.Unlock()

	s.out.Write([]byte(s.scheme.FmtEvent(*e)))
}

// Default logger
var Default Logger = &StdLogger{out: os.Stderr, scheme: ColorScheme}

// JSONLogger defines an output handler
type JSONLogger struct {
	out io.Writer
	sync.Mutex
}

// WriteEvent ouputs an event to JSONlogger
func (s *JSONLogger) WriteEvent(e *Event) {
	s.Lock()
	defer s.Unlock()

	json.NewEncoder(s.out).Encode(*e)
}

// FanLogger outputs to a list of loggers
type FanLogger struct {
	sync.RWMutex
	outs []Logger
}

// WriteEvent ouputs an event to FanLogger
func (f *FanLogger) WriteEvent(e *Event) {
	f.RLock()
	defer f.RUnlock()

	for i := range f.outs {
		f.outs[i].WriteEvent(e)
	}
}

// Add logger to fan
func (f *FanLogger) Add(l Logger) {
	f.Lock()
	defer f.Unlock()

	f.outs = append(f.outs, l)
}

// MQTTLogger outputs an event to MQTT
type MQTTLogger struct {
	topic string
	sync.Mutex
}

// WriteEvent ouputs an event to MQTT
func (m *MQTTLogger) WriteEvent(e *Event) {
	s, _ := mqtt.NewMessage(m.topic, *e)
	mqtt.Publish(s)
}

package mqtt

import (
	"sync"

	"sour.is/x/toolbox/log/event"
)

// Logger outputs an event to MQTT
type Logger struct {
	topic string
	level event.Level
	sync.Mutex
}

// NewLogger creates a new MQTT Logger
func NewLogger(topic string, level event.Level) event.Logger {
	return &Logger{topic: topic, level: level}
}

// WriteEvent ouputs an event to MQTT
func (l *Logger) WriteEvent(e *event.Event) {
	if l.level < e.Level {
		return
	}

	s, _ := NewMessage(l.topic, *e)
	Publish(s)
}

// SetVerbose sets the event verbose level
func (l *Logger) SetVerbose(level event.Level) {
	l.Lock()
	defer l.Unlock()

	l.level = level
}

// GetVerbose gets the verbose level
func (l *Logger) GetVerbose() event.Level { return l.level }

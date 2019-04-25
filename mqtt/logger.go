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

// WriteEvent ouputs an event to MQTT
func (m *Logger) WriteEvent(e *event.Event) {
	s, _ := NewMessage(m.topic, *e)
	Publish(s)
}

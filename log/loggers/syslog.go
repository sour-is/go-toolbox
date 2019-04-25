package loggers

import (
	"fmt"
	"log"
	"log/syslog"
	"sync"

	"sour.is/x/toolbox/log/event"
	"sour.is/x/toolbox/log/scheme"
)

// SysLogLogger outputs an event to SysLog
type SysLogLogger struct {
	sync.Mutex
	scheme scheme.Scheme
	level  event.Level
	*syslog.Writer
}

// NewSysLogger dial syslogger and return it
func NewSysLogger(addr, tag string) (event.Logger, error) {
	l := SysLogLogger{scheme: scheme.MonoScheme}
	var err error

	if addr == "local" {
		l.Writer, err = syslog.New(
			syslog.LOG_DEBUG|syslog.LOG_DAEMON, tag)
	} else {
		l.Writer, err = syslog.Dial("tcp", addr,
			syslog.LOG_DEBUG|syslog.LOG_DAEMON, tag)
	}
	if err != nil {
		log.Fatal(err)
	}

	return &l, err
}

// WriteEvent ouputs an event to SysLog
func (l *SysLogLogger) WriteEvent(e *event.Event) {
	_, err := l.Writer.Write([]byte(l.scheme.FmtEvent(*e)))
	if err != nil {
		fmt.Println(err)
	}
}

// SetVerbose sets the event verbose level
func (l *SysLogLogger) SetVerbose(level event.Level) {
	l.Lock()
	defer l.Unlock()

	l.level = level
}

// Close passes Close to underlying object.
func (l *SysLogLogger) Close() (err error) {
	return l.Writer.Close()
}

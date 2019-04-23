package logtag

import (
	"fmt"
	"log"
	"log/syslog"
	"sync"
)

// SysLogLogger outputs an event to SysLog
type SysLogLogger struct {
	sync.Mutex
	scheme Scheme

	*syslog.Writer
}

// NewSysLogger dial syslogger and return it
func NewSysLogger(addr, tag string) (Logger, error) {
	l := SysLogLogger{scheme: MonoScheme}
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
func (s *SysLogLogger) WriteEvent(e *Event) {
	_, err := s.Writer.Write([]byte(s.scheme.FmtEvent(*e)))
	if err != nil {
		fmt.Println(err)
	}
}

// Close passes Close to underlying object.
func (s *SysLogLogger) Close() (err error) {
	return s.Writer.Close()
}

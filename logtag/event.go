package logtag

import (
	"os"
	"runtime"
	"time"
)

// Event is a log unit
type Event struct {
	Level   EventLevel `json:"level"`
	Meta    MetaInfo   `json:"meta"`
	Message string     `json:"msg"`
	Tags    Tags       `json:"tags"`
}

var hostname string
var pid int

func init() {
	hostname = "unknown"
	if s, err := os.Hostname(); err == nil {
		hostname = s
	}

	pid = os.Getpid()
}

// MetaInfo is source information about an event
type MetaInfo struct {
	Host string    `json:"host"`
	PID  int       `json:"pid"`
	Time time.Time `json:"time"`
	File string    `json:"file"`
	Line int       `json:"line"`
	Func string    `json:"func"`
}

// NewMetaInfo generate a set of tags about the runtime.
func NewMetaInfo(calldepth int) (m MetaInfo) {
	m.Host = hostname
	m.PID = pid
	m.Time = time.Now()
	var ok bool
	var pc uintptr
	pc, m.File, m.Line, ok = runtime.Caller(calldepth)
	if !ok {
		m.File = "???"
		m.Line = 0
		pc = 0
	}
	details := runtime.FuncForPC(pc)
	m.Func = details.Name()

	return
}

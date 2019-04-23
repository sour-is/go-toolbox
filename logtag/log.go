package logtag

import (
	"fmt"
	"strings"
)

// Output something out
func Output(log Logger, calldepth int, in ...interface{}) {
	e := Event{Level: Vdebug, Tags: make(Tags), Meta: NewMetaInfo(calldepth)}

	var msg strings.Builder

	for _, v := range in {
		switch value := v.(type) {
		case EventLevel:
			e.Level = value
		case Tags:
			for k, t := range value {
				e.Tags[k] = t
			}
		default:
			msg.WriteString(fmt.Sprintf(" %v", value))
		}
	}
	e.Message = msg.String()

	log.WriteEvent(&e)
}

// Outputs generate event in the form Logger, depth, msg, [key, value], [key, value] ...
func Outputs(log Logger, calldepth int, level EventLevel, msg string, tags ...interface{}) {
	e := Event{Level: level, Message: msg, Tags: readTags(tags), Meta: NewMetaInfo(calldepth)}
	log.WriteEvent(&e)
}

// Debugs structured debug message
func Debugs(msg string, tags ...interface{}) { Outputs(Default, 2, Vdebug, msg, tags...) }

// Infos structured debug message
func Infos(msg string, tags ...interface{}) { Outputs(Default, 2, Vinfo, msg, tags...) }

// Warnings structured debug message
func Warnings(msg string, tags ...interface{}) { Outputs(Default, 2, Vwarning, msg, tags...) }

// Errors structured debug message
func Errors(msg string, tags ...interface{}) { Outputs(Default, 2, Verror, msg, tags...) }

// Criticals structured debug message
func Criticals(msg string, tags ...interface{}) { Outputs(Default, 2, Vcritical, msg, tags...) }

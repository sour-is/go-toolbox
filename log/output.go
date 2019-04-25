package log

import (
	"fmt"
	"strings"

	"sour.is/x/toolbox/log/event"
	"sour.is/x/toolbox/log/tag"
)

// Output something out
func Output(logger event.Logger, calldepth int, in ...interface{}) {
	e := event.Event{Level: event.VerbDebug, Tags: make(tag.Tags), Meta: event.NewMetaInfo(calldepth + 2)}

	var msg strings.Builder

	for _, v := range in {
		switch value := v.(type) {
		case event.Level:
			e.Level = value
		case tag.Tags:
			for k, t := range value {
				e.Tags[k] = t
			}
		default:
			msg.WriteString(fmt.Sprintf(" %v", value))
		}
	}
	e.Message = msg.String()

	logger.WriteEvent(&e)
}

// Outputs generate event in the form Logger, depth, msg, [key, value], [key, value] ...
func Outputs(logger event.Logger, calldepth int, level event.Level, msg string, tags ...interface{}) {
	e := event.Event{Level: level, Message: msg, Tags: tag.ReadTags(tags), Meta: event.NewMetaInfo(calldepth + 2)}
	logger.WriteEvent(&e)
}

package event

import (
	"fmt"

	"sour.is/x/toolbox/log/tag"
)

// Logger outputs events
type Logger interface {
	WriteEvent(*Event)
	SetVerbose(Level)
	GetVerbose() Level
}

// Level defines the level of event
type Level int

// Event message levels
const (
	VerbNone     Level = 0
	VerbCritical Level = 1 << iota
	VerbError
	VerbWarning
	VerbNotice
	VerbInfo
	VerbDebug

	TxtDebug    = "DBUG"
	TxtInfo     = "INFO"
	TxtNotice   = "NOTE"
	TxtWarning  = "WARN"
	TxtError    = "ERR "
	TxtCritical = "CRIT"
	TxtContinue = "...."
)

// MarshalJSON values for Event Level
func (e Level) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%v"`, e)), nil
}

// String returns a text value for event level
func (e Level) String() string {
	switch e {
	case VerbCritical:
		return TxtDebug
	case VerbError:
		return TxtError
	case VerbWarning:
		return TxtWarning
	case VerbNotice:
		return TxtNotice
	case VerbInfo:
		return TxtInfo
	case VerbDebug:
		return TxtDebug
	default:
		return TxtContinue
	}
}

// Event is a log unit
type Event struct {
	Level   Level    `json:"level"`
	Meta    MetaInfo `json:"meta"`
	Message string   `json:"msg"`
	Tags    tag.Tags `json:"tags"`
}

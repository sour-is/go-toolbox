package scheme

import (
	"fmt"
	"sort"
	"strings"

	"sour.is/x/toolbox/log/event"
)

// Constants to define commonly used text values
const (
	ASCreset  = "\x1B[0m"
	ASCgrey   = "\x1B[90m"
	ASCblue   = "\x1B[34m"
	ASCgreen  = "\x1B[32m"
	ASCyellow = "\x1B[93m"
	ASCred    = "\x1B[91m"
	ASCired   = "\x1B[7;91;49m"
)

// Scheme defines the colors to use in printing a log line
type Scheme struct {
	Reset    string
	Prefix   string
	Debug    string
	Info     string
	Notice   string
	Warning  string
	Error    string
	Critical string
	Continue string

	TagPrefix  string
	TagInfix   string
	TagPostfix string

	Timestamp string
	LongFile  bool
}

// FmtEvent format event using scheme
func (s Scheme) FmtEvent(e event.Event) string {
	var b strings.Builder

	lines := strings.Split(e.Message, "\n")

	b.WriteString(s.Prefix)
	b.WriteString(e.Meta.Time.Format(s.Timestamp))
	b.WriteRune(' ')

	switch e.Level {
	case event.VerbCritical:
		b.WriteString(s.Critical)
	case event.VerbError:
		b.WriteString(s.Error)
	case event.VerbWarning:
		b.WriteString(s.Warning)
	case event.VerbNotice:
		b.WriteString(s.Notice)
	case event.VerbInfo:
		b.WriteString(s.Info)
	case event.VerbDebug:
		b.WriteString(s.Debug)
	default:
		b.WriteString(s.Continue)
	}

	file := e.Meta.File
	fn := e.Meta.Func
	if !s.LongFile {
		file = shortFile(e.Meta.File)
		fn = shortFile(e.Meta.Func)
	}

	b.WriteString(fmt.Sprintf(" %s[%s:%d] ", fn, file, e.Meta.Line))
	b.WriteString(s.Reset)
	b.WriteString(strings.TrimSpace(lines[0]))

	if len(e.Tags) > 0 {
		var keys []string
		for k := range e.Tags {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			v := e.Tags[k]

			b.WriteRune(' ')
			b.WriteString(s.TagPrefix)
			b.WriteString(k)
			b.WriteString(s.TagInfix)
			b.WriteRune('=')
			b.WriteString(s.TagPostfix)
			b.WriteString(v.String())
			b.WriteString(s.Reset)
		}
	}
	b.WriteString(s.Reset)
	b.WriteString("\n")

	for _, m := range lines[1:] {
		b.WriteString(s.Continue)
		b.WriteRune(' ')
		b.WriteString(s.Reset)
		b.WriteString(strings.TrimSpace(m))
		b.WriteString(s.Reset)
		b.WriteString("\n")
	}

	return b.String()
}

// ColorScheme is default scheme with colors
var ColorScheme = Scheme{
	Reset:    ASCreset,
	Prefix:   ASCgrey,
	Debug:    ASCgrey + event.TxtDebug + ASCgrey,
	Info:     ASCblue + event.TxtInfo + ASCgrey,
	Notice:   ASCgreen + event.TxtNotice + ASCgrey,
	Warning:  ASCyellow + event.TxtWarning + ASCgrey,
	Error:    ASCred + event.TxtError + ASCgrey,
	Critical: ASCired + event.TxtCritical + ASCgrey,
	Continue: ASCgrey + event.TxtContinue + ASCgrey,

	TagPrefix:  ASCgreen,
	TagInfix:   ASCgrey,
	TagPostfix: ASCblue,

	Timestamp: "2006-01-02 15:04:05",
}

// MonoScheme is a scheme without colors
var MonoScheme = Scheme{
	Debug:     event.TxtDebug,
	Info:      event.TxtInfo,
	Notice:    event.TxtNotice,
	Warning:   event.TxtWarning,
	Error:     event.TxtError,
	Critical:  event.TxtCritical,
	Continue:  event.TxtContinue,
	Timestamp: "2006-01-02 15:04:05",
}

func shortFile(file string) string {
	short := file
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			short = file[i+1:]
			break
		}
	}
	return short
}

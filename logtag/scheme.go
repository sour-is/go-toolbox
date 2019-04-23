package logtag

import (
	"fmt"
	"strings"
)

// Constants to define commonly used text values
const (
	Tdebug    = "DBUG"
	Tinfo     = "INFO"
	Tnotice   = "NOTE"
	Twarning  = "WARN"
	Terror    = "ERR "
	Tcritical = "CRIT"
	Tcontinue = "...."

	ASCreset  = "\x1B[0m"
	ASCgrey   = "\x1B[90m"
	ASCblue   = "\x1B[34m"
	ASCgreen  = "\x1B[32m"
	ASCyellow = "\x1B[93m"
	ASCred    = "\x1B[91m"
	ASCired   = "\x1B[7;91;49m"

	Vnone     EventLevel = 0
	Vcritical EventLevel = 1 << iota
	Verror
	Vwarning
	Vnotice
	Vinfo
	Vdebug
)

// EventLevel defines the level of event
type EventLevel int

// MarshalJSON values for EventLevel
func (e EventLevel) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%v"`, e)), nil
}

// String returns a text value for event level
func (e EventLevel) String() string {
	switch e {
	case Vcritical:
		return Tdebug
	case Verror:
		return Terror
	case Vwarning:
		return Twarning
	case Vnotice:
		return Tnotice
	case Vinfo:
		return Tinfo
	case Vdebug:
		return Tdebug
	default:
		return Tcontinue
	}
}

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
func (s Scheme) FmtEvent(e Event) string {
	var b strings.Builder

	lines := strings.Split(e.Message, "\n")

	b.WriteString(s.Prefix)
	b.WriteString(e.Meta.Time.Format(s.Timestamp))
	b.WriteRune(' ')

	switch e.Level {
	case Vcritical:
		b.WriteString(s.Critical)
	case Verror:
		b.WriteString(s.Error)
	case Vwarning:
		b.WriteString(s.Warning)
	case Vnotice:
		b.WriteString(s.Notice)
	case Vinfo:
		b.WriteString(s.Info)
	case Vdebug:
		b.WriteString(s.Debug)
	default:
		b.WriteString(s.Continue)
	}

	file := e.Meta.File
	fn := e.Meta.Func
	if !s.LongFile {
		file = shortFile(e.Meta.File)
	}

	b.WriteString(fmt.Sprintf(" %s[%s:%d] ", fn, file, e.Meta.Line))
	b.WriteString(s.Reset)
	b.WriteString(strings.TrimSpace(lines[0]))

	if len(e.Tags) > 0 {
		for k, v := range e.Tags {
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
		b.WriteString("\n")
	}

	return b.String()
}

// ColorScheme is default scheme with colors
var ColorScheme = Scheme{
	Reset:    ASCreset,
	Prefix:   ASCgrey,
	Debug:    ASCgrey + Tdebug + ASCgrey,
	Info:     ASCblue + Tinfo + ASCgrey,
	Notice:   ASCgreen + Tnotice + ASCgrey,
	Warning:  ASCyellow + Twarning + ASCgrey,
	Error:    ASCred + Terror + ASCgrey,
	Critical: ASCired + Tcritical + ASCgrey,
	Continue: ASCgrey + Tcontinue + ASCgrey,

	TagPrefix:  ASCgreen,
	TagInfix:   ASCgrey,
	TagPostfix: ASCblue,

	Timestamp: "2006-01-02 15:04:05",
}

// MonoScheme is a scheme without colors
var MonoScheme = Scheme{
	Debug:     Tdebug,
	Info:      Tinfo,
	Notice:    Tnotice,
	Warning:   Twarning,
	Error:     Terror,
	Critical:  Tcritical,
	Continue:  Tcontinue,
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

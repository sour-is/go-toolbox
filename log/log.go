package log // import "sour.is/x/toolbox/log"

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"sync"

	"sour.is/x/toolbox/log/event"
	"sour.is/x/toolbox/log/loggers"
	"sour.is/x/toolbox/log/scheme"
)

var logger event.Logger = loggers.NewStdLogger(os.Stderr, scheme.ColorScheme, event.VerbNotice)
var mu = sync.Mutex{}

// SetOutput sets the output destination for the standard logger.
func SetOutput(l event.Logger) {
	mu.Lock()
	defer mu.Unlock()

	logger = l
}

// GetOutput returns the current logger
func GetOutput() event.Logger {
	return logger
}

// SetVerbose for the current global logger
func SetVerbose(level event.Level) {
	logger.SetVerbose(level)
}

// These functions write to the standard logger.

// Print calls Output to print to the standard logger.
// Arguments are handled in the manner of fmt.Print.
func Print(v ...interface{}) { Outputs(logger, 2, event.VerbInfo, fmt.Sprint(v...)) }

// Printf calls Output to print to the standard logger.
// Arguments are handled in the manner of fmt.Printf.
func Printf(format string, v ...interface{}) {
	Outputs(logger, 2, event.VerbInfo, fmt.Sprintf(format, v...))
}

// Write outputs contents of io.Reader to standard logger.
func Write(r io.Reader) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r)

	Outputs(logger, 2, event.VerbInfo, buf.String())
}

// Tee outputs contents of io.Reader to standard logger. and returns a new io.Reader
func Tee(r io.ReadCloser) (w io.ReadCloser) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r)
	r.Close()

	str := buf.String()
	w = ioutil.NopCloser(strings.NewReader(str))

	Outputs(logger, 2, event.VerbInfo, str)

	return
}

// Fatal is equivalent to Print() followed by a call to os.Exit(1).
func Fatal(v ...interface{}) {
	Outputs(logger, 2, event.VerbCritical, fmt.Sprint(v...))
	os.Exit(1)
}

// Fatalf is equivalent to Printf() followed by a call to os.Exit(1).
func Fatalf(format string, v ...interface{}) {
	Outputs(logger, 2, event.VerbCritical, fmt.Sprintf(format, v...))
	os.Exit(1)
}

// Panic is equivalent to Print() followed by a call to panic().
func Panic(v ...interface{}) {
	s := fmt.Sprint(v...)
	Outputs(logger, 2, event.VerbError, s)
	panic(s)
}

// Panicf is equivalent to Printf() followed by a call to panic().
func Panicf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	Outputs(logger, 2, event.VerbError, s)
	panic(s)
}

// Debug outputs to logger with DEBUG level.
func Debug(v ...interface{}) { Outputs(logger, 2, event.VerbDebug, fmt.Sprint(v...)) }

// Debugf formats output to logger with DEBUG level.
func Debugf(format string, v ...interface{}) {
	Outputs(logger, 2, event.VerbDebug, fmt.Sprintf(format, v...))
}

// Debugw outputs io.Reader to logger with DEBUG level.
func Debugw(r io.Reader) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r)
	Outputs(logger, 2, event.VerbDebug, buf.String())
}

// Info outputs to logger with INFO level.
func Info(v ...interface{}) { Outputs(logger, 2, event.VerbInfo, fmt.Sprint(v...)) }

// Infof formatted outputs to logger with INFO level.
func Infof(format string, v ...interface{}) {
	Outputs(logger, 2, event.VerbInfo, fmt.Sprintf(format, v...))
}

// Infow outputs io.Reader to logger with INFO level.
func Infow(r io.Reader) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r)
	Outputs(logger, 2, event.VerbInfo, buf.String())
}

// Notice outputs to logger with NOTICE level.
func Notice(v ...interface{}) { Outputs(logger, 2, event.VerbNotice, fmt.Sprint(v...)) }

// Noticef formatted outputs to logger with NOTICE level.
func Noticef(format string, v ...interface{}) {
	Outputs(logger, 2, event.VerbNotice, fmt.Sprintf(format, v...))
}

// Noticew outputs io.Reader to logger with NOTICE level.
func Noticew(r io.Reader) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r)
	Outputs(logger, 2, event.VerbNotice, buf.String())

}

// Warning outputs to logger with WARNING level.
func Warning(v ...interface{}) { Outputs(logger, 2, event.VerbWarning, fmt.Sprint(v...)) }

// Warningf formatted outputs to logger with WARNING level.
func Warningf(format string, v ...interface{}) {
	Outputs(logger, 2, event.VerbWarning, fmt.Sprintf(format, v...))
}

// Warningw outputs io.Reader to logger with WARNING level.
func Warningw(r io.Reader) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r)
	Outputs(logger, 2, event.VerbWarning, buf.String())

}

// Error outputs to logger with ERROR level.
func Error(v ...interface{}) { Outputs(logger, 2, event.VerbError, fmt.Sprint(v...)) }

// Errorf formatted outputs to logger with ERROR level.
func Errorf(format string, v ...interface{}) {
	Outputs(logger, 2, event.VerbError, fmt.Sprintf(format, v...))
}

// Errorw outputs io.Reader to logger with ERROR level.
func Errorw(r io.Reader) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r)
	Outputs(logger, 2, event.VerbError, buf.String())
}

// Critical outputs to logger with CRITICAL level.
func Critical(v ...interface{}) { Outputs(logger, 2, event.VerbCritical, fmt.Sprint(v...)) }

// Criticalf formatted outputs to logger with CRITICAL level.
func Criticalf(format string, v ...interface{}) {
	Outputs(logger, 2, event.VerbCritical, fmt.Sprintf(format, v...))
}

// Criticalw outputs io.Reader to logger with CRITICAL level.
func Criticalw(r io.Reader) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r)
	Outputs(logger, 2, event.VerbCritical, buf.String())
}

// Debugs structured debug message
func Debugs(msg string, tags ...interface{}) { Outputs(logger, 2, event.VerbDebug, msg, tags...) }

// Infos structured debug message
func Infos(msg string, tags ...interface{}) { Outputs(logger, 2, event.VerbInfo, msg, tags...) }

// Warnings structured debug message
func Warnings(msg string, tags ...interface{}) { Outputs(logger, 2, event.VerbWarning, msg, tags...) }

// Errors structured debug message
func Errors(msg string, tags ...interface{}) { Outputs(logger, 2, event.VerbError, msg, tags...) }

// Criticals structured debug message
func Criticals(msg string, tags ...interface{}) { Outputs(logger, 2, event.VerbCritical, msg, tags...) }

// These functions do nothing. It makes it easy to comment out
// debug lines without having to remove the import.

// NilPrint does nothing.
func NilPrint(_ ...interface{}) {}

// NilPrintf does nothing.
func NilPrintf(_ string, _ ...interface{}) {}

// NilPrintln does nothing.
func NilPrintln(_ ...interface{}) {}

// NilDebug does nothing.
func NilDebug(_ ...interface{}) {}

// NilDebugf does nothing.
func NilDebugf(_ string, _ ...interface{}) {}

// NilDebugw does nothing.
func NilDebugw(_ io.Reader) {}

// NilInfo does nothing.
func NilInfo(_ ...interface{}) {}

// NilInfof does nothing.
func NilInfof(_ string, _ ...interface{}) {}

// NilInfow does nothing.
func NilInfow(_ io.Reader) {}

// NilNotice does nothing.
func NilNotice(_ ...interface{}) {}

// NilNoticef does nothing.
func NilNoticef(_ string, _ ...interface{}) {}

// NilNoticew does nothing.
func NilNoticew(_ io.Reader) {}

// NilWarning does nothing.
func NilWarning(_ ...interface{}) {}

// NilWarningf does nothing.
func NilWarningf(_ string, _ ...interface{}) {}

// NilWarningw does nothing.
func NilWarningw(_ io.Reader) {}

// NilError does nothing.
func NilError(_ ...interface{}) {}

// NilErrorf does nothing.
func NilErrorf(_ string, _ ...interface{}) {}

// NilErrorw does nothing.
func NilErrorw(_ io.Reader) {}

// NilCritical does nothing.
func NilCritical(_ ...interface{}) {}

// NilCriticalf does nothing.
func NilCriticalf(_ string, _ ...interface{}) {}

// NilCriticalw does nothing.
func NilCriticalw(_ io.Reader) {}

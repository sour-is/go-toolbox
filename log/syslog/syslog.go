package syslog

import (
	"io"
	"log/syslog"
	"os"

	"sour.is/x/toolbox/log"
)

var sysLog StripCtrlWriter

/*
var sessionID = "session-id"
var logFormat = "%s %- 16s    %- 6s %- 30s    %12s %d %s"

var accLog StripCtrlWriter
func doAccessLog(name string, w httpsrv.ResponseWriter, r *http.Request, id ident.Ident) bool {
	if r.Method == "HEAD" || r.Method == "OPTIONS" {
		return true
	}

	user := "-"
	if id.IsActive() {
		user = id.GetAspect() + "/" + id.GetIdentity()
	}
	fmt.Fprintf(
		accLog,
		logFormat,
		r.Header.Get(sessionID),
		user,
		r.Method,
		name,
		w.StopTime(),
		w.GetCode(),
		r.RequestURI,
	)

	return true
}
*/

// Dial sets up sysloger to addr
func Dial(addr, tag string) {

	var err error

	var s *syslog.Writer
	if addr == "local" {
		s, err = syslog.New(
			syslog.LOG_DEBUG|syslog.LOG_DAEMON, tag)
	} else {
		s, err = syslog.Dial("tcp", addr,
			syslog.LOG_DEBUG|syslog.LOG_DAEMON, tag)
	}
	if err != nil {
		log.Fatal(err)
	}

	sysLog = StripCtrlWriter{s}

	if err != nil {
		log.Fatal(err)
	}
	multi := io.MultiWriter(os.Stderr, sysLog)
	log.SetOutput(multi)

	log.Noticef("Setup remote logging to %s", addr)

	/*
		var a *syslog.Writer
		if addr == "local" {
			a, err = syslog.New(
				syslog.LOG_DEBUG|syslog.LOG_DAEMON, tag+"-access")
		} else {
			a, err = syslog.Dial("tcp", addr,
				syslog.LOG_DEBUG|syslog.LOG_DAEMON, tag+"-access")
		}
		if err != nil {
			log.Fatal(err)
		}
		accLog = StripCtrlWriter{a}

		httpsrv.NewMiddleware("access-syslog", doAccessLog).Register(httpsrv.EventComplete)
	*/
}

// StripCtrlWriter is a writer that removes control codes from text.
type StripCtrlWriter struct {
	Writer io.WriteCloser
}

// Write removes out terminal control characters before passing to writer.
func (s StripCtrlWriter) Write(b []byte) (n int, err error) {
	i := 0
	for j := 0; i < len(b) && j < len(b); i++ {

		switch b[j] {
		case 0x07, 0x08, 0x09, 0x0A, 0x0D, 0x0E, 0x0F, 0x7F:
			j++

		case 0x1B:
			// ESC
			j++
			// Initial Character
			if b[j] >= 0x40 && b[j] < 0x60 {
				j++
			}
			fallthrough

		case 0x9B:
			// CSI

			for j < len(b) { // Param bytes
				if b[j] >= 0x30 && b[j] < 0x40 {
					j++
				} else {
					break
				}
			}

			for j < len(b) { // intermediate bytes
				if b[j] >= 0x20 && b[j] < 0x30 {
					j++
				} else {
					break
				}
			}

			// Final byte
			if b[j] >= 0x40 && b[j] < 0x80 {
				j++
			}
		}

		if j >= len(b) {
			break
		}

		b[i] = b[j]
		j++
	}

	return s.Writer.Write(b[:i])
}

// Close passes Close to underlying object.
func (s StripCtrlWriter) Close() (err error) {
	return s.Writer.Close()
}

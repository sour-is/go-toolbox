package httpsrv // import "sour.is/x/toolbox/httpsrv"

import (
	"net/http"
	"os"

	"sour.is/x/toolbox/ident"
	"sour.is/x/toolbox/uuid"

	stdlog "log"
)

var sessionId = "session-id"

func doSessionID(_ string, w ResponseWriter, r *http.Request, _ ident.Ident) bool {
	var seq string
	if seq = r.URL.Query().Get(sessionId); seq == "" {
		if seq = r.Header.Get(sessionId); seq == "" {
			seq = uuid.V4() + "::"
			r.Header.Set(sessionId, seq)
		}
	}
	w.Header().Add(sessionId, seq)

	return true
}

var accessLog = stdlog.New(os.Stdout, "", stdlog.Ldate|stdlog.Ltime|stdlog.LUTC)
var logFormat = "%s %- 16s\t%- 16v\t%- 6s %- 30s\t%12s %d %s"

func doAccessLog(name string, w ResponseWriter, r *http.Request, id ident.Ident) bool {
	user := "-"
	if id.IsActive() {
		user = id.GetAspect() + "/" + id.GetIdentity()
	}
	header := r.Header.Get(sessionId)
	if len(header) < 20 {
		header = uuid.NilUUID
	}

	accessLog.Printf(
		logFormat,
		header[19:],
		r.RemoteAddr,
		user,
		r.Method,
		name,
		w.StopTime(),
		w.GetCode(),
		r.RequestURI,
	)

	return true
}

func init() {
	NewMiddleware("session-id", doSessionID).Register(EventPreAuth)
	NewMiddleware("access-log", doAccessLog).Register(EventComplete)
}

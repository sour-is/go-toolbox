package httpsrv // import "sour.is/x/toolbox/httpsrv"

import (
	"net/http"
	"os"

	"sour.is/x/toolbox/ident"
	"sour.is/x/toolbox/log"
	"sour.is/x/toolbox/uuid"

	stdlog "log"
)

var sessionId = "session-id"

func doSessionID(_ string, w ResponseWriter, r *http.Request, _ ident.Ident) {
	var seq string
	if seq = r.URL.Query().Get(sessionId); seq == "" {
		if seq = r.Header.Get(sessionId); seq == "" {
			seq = uuid.V4() + "::"
			r.Header.Set(sessionId, seq)
		}
	}
	w.Header().Add(sessionId, seq)
}

var accessLog = stdlog.New(os.Stdout, "", log.Ldate|log.Ltime|log.LUTC)
var logFormat = "%s %- 16s\t%- 6s %- 30s\t%12s %d %s"

func doAccessLog(name string, w ResponseWriter, r *http.Request, id ident.Ident) {
	user := "-"
	if id.IsActive() {
		user = id.GetAspect() + "/" + id.GetIdentity()
	}
	accessLog.Printf(
		logFormat,
		r.Header.Get(sessionId),
		user,
		r.Method,
		name,
		w.StopTime(),
		w.GetCode(),
		r.RequestURI,
	)
}

func init() {
	NewMiddleware("session-id", doSessionID).Register(EventPreAuth)
	NewMiddleware("access-log", doAccessLog).Register(EventComplete)
}

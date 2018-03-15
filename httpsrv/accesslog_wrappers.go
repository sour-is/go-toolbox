package httpsrv // import "sour.is/x/toolbox/httpsrv"

import (
	"net/http"
	"os"
	"time"

	"github.com/spf13/viper"
	"sour.is/x/toolbox/ident"
	"sour.is/x/toolbox/log"
	"sour.is/x/toolbox/uuid"

	stdlog "log"
)

var accessLog = stdlog.New(os.Stdout, "", log.Ldate|log.Ltime|log.LUTC)
var logFormat = "%s %- 16s\t%- 6s %- 30s\t%12s %d %s"
var sessionId = "session-id"

func addSessionID(w http.ResponseWriter, r *http.Request) (seq string) {
	if seq = r.URL.Query().Get(sessionId); seq == "" {
		if seq = r.Header.Get(sessionId); seq == "" {
			seq = uuid.V4() + ":0000:0000"
			r.Header.Set(sessionId, seq)
		}
	}
	w.Header().Add(sessionId, seq)

	return
}

func Wrapper(inner http.HandlerFunc, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		var nw = NewResponseWriter(w)
		seq := addSessionID(w, r)

		inner.ServeHTTP(nw, r)

		accessLog.Printf(
			logFormat,
			seq,
			"-",
			r.Method,
			name,
			time.Since(start),
			nw.GetCode(),
			r.RequestURI,
		)
	})
}
func IdentWrapper(inner HandlerFunc, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		var nw = NewResponseWriter(w)
		seq := addSessionID(w, r)

		id := ident.GetIdent(viper.GetString("http.idm"), r)

		inner.ServeHTTP(nw, r, id)

		accessLog.Printf(
			logFormat,
			seq,
			id.GetAspect()+"/"+
				id.GetIdentity(),
			r.Method,
			name,
			time.Since(start),
			nw.GetCode(),
			r.RequestURI,
		)
	})
}
func AssetWrapper(name, prefix string, hdlr http.FileSystem) http.Handler {
	fn := http.StripPrefix(prefix, http.FileServer(hdlr))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		var nw = NewResponseWriter(w)
		seq := addSessionID(w, r)

		fn.ServeHTTP(nw, r)

		accessLog.Printf(
			"%s %- 16s\t%- 6s %- 30s\t%12s %d %s",
			seq,
			"-",
			r.Method,
			name,
			time.Since(start),
			nw.GetCode(),
			r.RequestURI,
		)
	})
}

// Custom ResponseWriter that saves the response code so the access log is
// able to display it. The default ResponseWriter is passed by value so after
// ServeHTTP completes the value remains unchanged.
type responseWriter struct {
	ResponseCode int
	HeadersSent  bool
}

type ResponseWriter struct {
	W http.ResponseWriter
	R *responseWriter
}

func NewResponseWriter(w http.ResponseWriter) (r ResponseWriter) {
	r.W = w
	r.R = new(responseWriter)

	return
}

func (w ResponseWriter) WriteHeader(c int) {
	w.R.ResponseCode = c
	if c == 204 || c >= 300 && c <= 399 || c == 410 {
		w.W.WriteHeader(w.R.ResponseCode)
		w.R.HeadersSent = true
	}
}

func (w ResponseWriter) Header() http.Header {
	return w.W.Header()
}

func (w ResponseWriter) Write(b []byte) (int, error) {
	if w.R.ResponseCode == 0 {
		w.R.ResponseCode = 200
	}

	if w.R.HeadersSent == false {
		w.W.WriteHeader(w.R.ResponseCode)
		w.R.HeadersSent = true
	}

	return w.W.Write(b)
}

func (w ResponseWriter) GetCode() int {
	return w.R.ResponseCode
}

package httpsrv

import (
	"net/http"
	"time"
)

// Custom ResponseWriter that saves the response code so the access log is
// able to display it. The default ResponseWriter is passed by value so after
// ServeHTTP completes the value remains unchanged.
type responseWriter struct {
	ResponseCode int
	HeadersSent  bool
	StartTime    time.Time
	Duration     time.Duration
	BytesOut     int
}

type ResponseWriter struct {
	W http.ResponseWriter
	R *responseWriter
}

func wrapResponseWriter(w http.ResponseWriter) (r ResponseWriter) {
	r.W = w
	r.R = new(responseWriter)
	r.R.StartTime = time.Now()

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

	w.R.BytesOut += len(b)
	return w.W.Write(b)
}

func (w ResponseWriter) GetCode() int {
	return w.R.ResponseCode
}

func (w ResponseWriter) Since() time.Duration {
	return time.Since(w.R.StartTime)
}

func (w ResponseWriter) StopTime() time.Duration {
	if w.R.Duration == 0 {
		w.R.Duration = w.Since()
	}

	return w.R.Duration
}

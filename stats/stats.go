package stats

import (
	"sour.is/x/toolbox/ident"
	"net/http"
	"sour.is/x/toolbox/httpsrv"
	"time"
	"bufio"
	"bytes"
)

var httpPipe chan httpData

func init() {
	appStart = time.Now()

	httpPipe = make(chan httpData)
	go recordStats(httpPipe)

	httpsrv.NewMiddleware("gather-stats", doStats).Register(httpsrv.EventComplete)
	httpsrv.IdentRegister("stats", httpsrv.IdentRoutes{
		{"GetStats", "GET", "/v1/stats", getStats},
	})
}

var appStart time.Time

type httpStatsType struct{
	Requests int `json:"requests"`

	Http2xx int `json:"http_2xx"`
	Http3xx int `json:"http_3xx"`
	Http4xx int `json:"http_4xx"`
	Http5xx int `json:"http_5xx"`

	AnonRequests int `json:"anonymous_requests"`

	HeaderBytesOut int `json:"header_bytes_out"`
	ContentBytesOut int `json:"content_bytes_out"`

	BytesOut int `json:"bytes_out"`
}
var httpStats httpStatsType

var httpSeries struct{
	Request5m int
	Request10m int
	Request25m int
	Request60m int


}

func getStats(w http.ResponseWriter, r *http.Request, id ident.Ident) {

	stats := struct{
		AppStart time.Time `json:"app_start"`
		HttpTotals httpStatsType `json:"http_total"`

	}{
		appStart,
		httpStats,
	}

	httpsrv.WriteObject(w, http.StatusOK, stats)
}

func doStats(_ string, w httpsrv.ResponseWriter, r *http.Request, id ident.Ident) {
	httpPipe <- httpData{w,r,id}
}

type httpData struct{
	W httpsrv.ResponseWriter
	R *http.Request
	ID ident.Ident
}

func recordStats(pipe chan httpData) {
	for {
		select {
			case h := <-pipe:
				httpStats.Requests += 1

				code := h.W.GetCode()
				switch {

				case code >= 200 && code < 300:
					httpStats.Http2xx += 1

				case code >= 300 && code < 400:
					httpStats.Http3xx += 1

				case code >= 400 && code < 500:
					httpStats.Http4xx += 1

				case code >= 500:
					httpStats.Http4xx += 1

				}


				if !h.ID.IsActive() {
					httpStats.AnonRequests += 1
				}

				var b bytes.Buffer
				var w = bufio.NewWriter(&b)
				h.W.W.Header().Write(w)

				httpStats.BytesOut += h.W.R.BytesOut

			case <- httpsrv.SignalShutdown:
				return
		}
	}
}
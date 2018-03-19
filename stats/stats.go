package stats

import (
	"bufio"
	"bytes"
	"net/http"
	"sour.is/x/toolbox/httpsrv"
	"sour.is/x/toolbox/ident"
	"sour.is/x/toolbox/log"
	"time"
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

type httpStatsType struct {
	Requests    int           `json:"requests"`
	RequestTime time.Duration `json:"request_time"`

	Http2xx int `json:"http_2xx"`
	Http3xx int `json:"http_3xx"`
	Http4xx int `json:"http_4xx"`
	Http5xx int `json:"http_5xx"`

	AnonRequests int `json:"anonymous_requests"`

	HeaderBytesOut  int `json:"header_bytes_out"`
	ContentBytesOut int `json:"content_bytes_out"`

	BytesOut int `json:"bytes_out"`
}

var httpStats httpStatsType

var httpSeries struct {
	Request1m  int
	Request5m  int
	Request10m int
	Request25m int
	Request60m int
}

var httpCollect struct {
	Request1m  int
	Request5m  int
	Request10m int
	Request25m int
	Request60m int
}

func getStats(w http.ResponseWriter, r *http.Request, id ident.Ident) {

	avgTime := 0
	if httpStats.Requests > 0 {
		avgTime = int(httpStats.RequestTime) / httpStats.Requests
	}

	stats := struct {
		AppStart   time.Time     `json:"app_start"`
		UpTime     time.Duration `json:"uptime"`
		httpStatsType
		AvgTime    int           `json:"req_avg_time"`

		Last1m     int `json:"reqs_1m_last"`
		Request1m  int `json:"reqs_1m"`
		Last5m     int `json:"reqs_5m_last"`
		Request5m  int `json:"reqs_5m"`
		Last10m    int `json:"reqs_10m_last"`
		Request10m int `json:"reqs_10m"`
		Last25m    int `json:"reqs_25m_last"`
		Request25m int `json:"reqs_25m"`
		Last60m    int `json:"reqs_60m_last"`
		Request60m int `json:"reqs_60m"`
	}{
		appStart,
		time.Since(appStart),
		httpStats,
		avgTime,

		httpSeries.Request1m,
		httpCollect.Request1m,
		httpSeries.Request5m,
		httpCollect.Request5m,
		httpSeries.Request10m,
		httpCollect.Request10m,
		httpSeries.Request25m,
		httpCollect.Request25m,
		httpSeries.Request60m,
		httpCollect.Request60m,
	}

	httpsrv.WriteObject(w, http.StatusOK, stats)
}

func doStats(_ string, w httpsrv.ResponseWriter, r *http.Request, id ident.Ident) {
	httpPipe <- httpData{w, r, id}
}

type httpData struct {
	W  httpsrv.ResponseWriter
	R  *http.Request
	ID ident.Ident
}

func recordStats(pipe chan httpData) {
	log.Notice("Begin collecting HTTP stats")

	ticker1m := time.NewTicker(time.Minute)
	defer ticker1m.Stop()

	ticker5m := time.NewTicker(time.Minute * 5)
	defer ticker5m.Stop()

	ticker10m := time.NewTicker(time.Minute * 10)
	defer ticker10m.Stop()

	ticker25m := time.NewTicker(time.Minute * 25)
	defer ticker25m.Stop()

	ticker60m := time.NewTicker(time.Minute * 60)
	defer ticker60m.Stop()

	for {
		select {
		case <-ticker1m.C:
			log.Debug("Rolling 1m stats")
			httpSeries.Request1m = httpCollect.Request5m
			httpCollect.Request1m = 0

		case <-ticker5m.C:
			log.Debug("Rolling 5m stats")
			httpSeries.Request5m = httpCollect.Request5m
			httpCollect.Request5m = 0

		case <-ticker10m.C:
			log.Debug("Rolling 10m stats")
			httpSeries.Request10m = httpCollect.Request10m
			httpCollect.Request10m = 0

		case <-ticker25m.C:
			log.Debug("Rolling 25m stats")
			httpSeries.Request25m = httpCollect.Request25m
			httpCollect.Request25m = 0

		case <-ticker60m.C:
			log.Debug("Rolling 60m stats")
			httpSeries.Request60m = httpCollect.Request60m
			httpCollect.Request60m = 0

		case h := <-pipe:
			httpStats.Requests += 1
			httpCollect.Request1m += 1
			httpCollect.Request5m += 1
			httpCollect.Request10m += 1
			httpCollect.Request25m += 1
			httpCollect.Request60m += 1

			httpStats.RequestTime = h.W.StopTime()

			code := h.W.GetCode()
			switch {

			case code >= 200 && code < 300:
				httpStats.Http2xx += 1

			case code >= 300 && code < 400:
				httpStats.Http3xx += 1

			case code >= 400 && code < 500:
				httpStats.Http4xx += 1

			case code >= 500:
				httpStats.Http5xx += 1
			}

			if !h.ID.IsActive() {
				httpStats.AnonRequests += 1
			}

			var b bytes.Buffer
			var w = bufio.NewWriter(&b)
			h.W.W.Header().Write(w)

			httpStats.BytesOut += h.W.R.BytesOut

		case <-httpsrv.SignalShutdown:
			log.Debug("Shutting Down Stats")
			return
		}
	}
}

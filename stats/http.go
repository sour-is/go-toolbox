package stats

import (
	"net/http"
	"time"

	"sour.is/x/toolbox/httpsrv"
	"sour.is/x/toolbox/ident"
	"sour.is/x/toolbox/log"
	"sour.is/x/toolbox/stats/exposition"
)

var httpPipe chan httpData

func init() {
	httpPipe = make(chan httpData)
	go recordStats(httpPipe)
	httpsrv.NewMiddleware("gather-stats", doStats).Register(httpsrv.EventComplete)

	Register("http", getHTTPstats)
}

func doStats(_ string, w httpsrv.ResponseWriter, r *http.Request, id ident.Ident) bool {
	httpPipe <- httpData{w, r, id}
	return true
}

type httpStatsType struct {
	Requests    int           `json:"reqs"`
	RequestTime time.Duration `json:"req_time"`

	HTTP2xx int `json:"http_2xx"`
	HTTP3xx int `json:"http_3xx"`
	HTTP4xx int `json:"http_4xx"`
	HTTP5xx int `json:"http_5xx"`

	AnonRequests int `json:"reqs_anon"`

	BytesOut int `json:"bytes_out"`
}

var httpStats httpStatsType

type httpSeriesType struct {
	Request1m  int `json:"01m"`
	Request5m  int `json:"05m"`
	Request10m int `json:"10m"`
	Request25m int `json:"25m"`
	Request60m int `json:"60m"`
}

var httpSeries httpSeriesType
var httpCollect httpSeriesType

type httpReqs struct {
	httpStatsType
	AvgTimeNano int    `json:"req_avg_nano"`
	AvgTime     string `json:"req_avg"`

	CurrentCount httpSeriesType `json:"req_counts"`
	LastCount    httpSeriesType `json:"req_counts_last"`
}

func (s httpReqs) String() string {
	return s.Exposition().String()
}
func (s httpReqs) Exposition() (lis exposition.Expositions) {

	e := exposition.New("http_requests_avg_time", exposition.Gauge)
	e.AddRow(float64(s.AvgTimeNano))
	lis = append(lis, e)

	e = exposition.New("http_requests_by_status", exposition.Counter)
	e.AddRow(float64(s.HTTP2xx)).AddTag("code", "200")
	e.AddRow(float64(s.HTTP3xx)).AddTag("code", "300")
	e.AddRow(float64(s.HTTP4xx)).AddTag("code", "400")
	e.AddRow(float64(s.HTTP5xx)).AddTag("code", "500")
	lis = append(lis, e)

	e = exposition.New("http_requests_by_auth", exposition.Counter)
	e.AddRow(float64(s.AnonRequests)).AddTag("auth", "false")
	e.AddRow(float64(s.Requests-s.AnonRequests)).AddTag("auth", "true")
	lis = append(lis, e)

	e = exposition.New("http_requests_total", exposition.Counter)
	e.AddRow(float64(s.Requests))
	lis = append(lis, e)

	e = exposition.New("http_request_bytes_total", exposition.Counter)
	e.AddRow(float64(s.BytesOut))
	lis = append(lis, e)

	e = exposition.New("http_request_freq_sum", exposition.Summary)

	var c int
	if s.LastCount.Request1m == 0 {
		c = s.CurrentCount.Request1m
	} else {
		c = s.LastCount.Request1m
	}
	e.AddRow(float64(c)).AddTag("window", "01m")

	if s.LastCount.Request5m == 0 {
		c = s.CurrentCount.Request5m
	} else {
		c = s.LastCount.Request5m
	}
	e.AddRow(float64(c)).AddTag("window", "05m")

	if s.LastCount.Request10m == 0 {
		c = s.CurrentCount.Request10m
	} else {
		c = s.LastCount.Request10m
	}
	e.AddRow(float64(c)).AddTag("window", "10m")

	if s.LastCount.Request25m == 0 {
		c = s.CurrentCount.Request25m
	} else {
		c = s.LastCount.Request25m
	}
	e.AddRow(float64(c)).AddTag("window", "25m")

	if s.LastCount.Request60m == 0 {
		c = s.CurrentCount.Request60m
	} else {
		c = s.LastCount.Request60m
	}
	e.AddRow(float64(c)).AddTag("window", "60m")
	lis = append(lis, e)

	return
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

	httpsrv.WaitShutdown.Add(1)

	for {
		select {
		case <-ticker1m.C:
			log.Debug("Rolling 1m stats")
			httpSeries.Request1m, httpCollect.Request1m = httpCollect.Request1m, 0

		case <-ticker5m.C:
			log.Debug("Rolling 5m stats")
			httpSeries.Request5m, httpCollect.Request5m = httpCollect.Request5m, 0

		case <-ticker10m.C:
			log.Debug("Rolling 10m stats")
			httpSeries.Request10m, httpCollect.Request10m = httpCollect.Request10m, 0

		case <-ticker25m.C:
			log.Debug("Rolling 25m stats")
			httpSeries.Request25m, httpCollect.Request25m = httpCollect.Request25m, 0

		case <-ticker60m.C:
			log.Debug("Rolling 60m stats")
			httpSeries.Request60m, httpCollect.Request60m = httpCollect.Request60m, 0

		case h := <-pipe:
			httpStats.Requests++
			httpCollect.Request1m++
			httpCollect.Request5m++
			httpCollect.Request10m++
			httpCollect.Request25m++
			httpCollect.Request60m++

			httpStats.RequestTime += h.W.StopTime()

			code := h.W.GetCode()
			switch {

			case code >= 200 && code < 300:
				httpStats.HTTP2xx++

			case code >= 300 && code < 400:
				httpStats.HTTP3xx++

			case code >= 400 && code < 500:
				httpStats.HTTP4xx++

			case code >= 500:
				httpStats.HTTP5xx++
			}

			if !h.ID.IsActive() {
				httpStats.AnonRequests++
			}

			httpStats.BytesOut += h.W.R.BytesOut

		case <-httpsrv.SignalShutdown:
			log.Debug("Shutting Down Stats")
			httpsrv.WaitShutdown.Done()
			return
		}
	}
}

func getHTTPstats() exposition.Expositioner {
	avgTime := 0
	if httpStats.Requests > 0 {
		avgTime = int(httpStats.RequestTime) / httpStats.Requests
	}

	return httpReqs{
		httpStats,
		avgTime,
		time.Duration(avgTime).String(),
		httpCollect,
		httpSeries,
	}
}

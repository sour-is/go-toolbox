package stats

import (
	"database/sql"
	"github.com/cgilling/dbstats"
	"net/http"
	"runtime"
	"sour.is/x/toolbox/httpsrv"
	"sour.is/x/toolbox/ident"
	"sour.is/x/toolbox/log"
	"time"
)

var httpPipe chan httpData

func init() {
	appStart = time.Now().In(time.UTC)
	httpPipe = make(chan httpData)
	dbHooks = make(map[string]*dbstats.CounterHook)
	go recordStats(httpPipe)

	httpsrv.NewMiddleware("gather-stats", doStats).Register(httpsrv.EventComplete)
	httpsrv.IdentRegister("stats", httpsrv.IdentRoutes{
		{"GetStats", "GET", "/v1/stats", getStats},
	})
}

var appStart time.Time

type httpStatsType struct {
	Requests    int           `json:"reqs"`
	RequestTime time.Duration `json:"req_time"`

	Http2xx int `json:"http_2xx"`
	Http3xx int `json:"http_3xx"`
	Http4xx int `json:"http_4xx"`
	Http5xx int `json:"http_5xx"`

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

type Stats struct {
	AppStart   time.Time     `json:"app_start"`
	UpTimeNano time.Duration `json:"uptime_nano"`
	UpTime     string        `json:"uptime"`
	Http       struct {
		httpStatsType
		AvgTimeNano int    `json:"req_avg_nano"`
		AvgTime     string `json:"req_avg"`

		CurrentCount httpSeriesType `json:"req_counts"`
		LastCount    httpSeriesType `json:"req_counts_last"`
	} `json:"http"`
	Runtime runtimeStats       `json:"runtime"`
	DBstats map[string]dbStats `json:"db"`
}

// swagger:operation GET /v1/stats stats getStats
//
// Get Stats
//
// ---
// produces:
//   - "application/json"
// responses:
//   "200":
//     description: Success
//     schema:
//       type: object
//       properties:
//          items:
func getStats(w http.ResponseWriter, r *http.Request, id ident.Ident) {

	avgTime := 0
	if httpStats.Requests > 0 {
		avgTime = int(httpStats.RequestTime) / httpStats.Requests
	}

	stats := Stats{
		appStart,
		time.Since(appStart),
		time.Since(appStart).String(),
		struct {
			httpStatsType
			AvgTimeNano int    `json:"req_avg_nano"`
			AvgTime     string `json:"req_avg"`

			CurrentCount httpSeriesType `json:"req_counts"`
			LastCount    httpSeriesType `json:"req_counts_last"`
		}{
			httpStats,
			avgTime,
			time.Duration(avgTime).String(),

			httpCollect,
			httpSeries,
		},
		getRuntime(),
		getDBstats(),
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
			httpStats.Requests += 1
			httpCollect.Request1m += 1
			httpCollect.Request5m += 1
			httpCollect.Request10m += 1
			httpCollect.Request25m += 1
			httpCollect.Request60m += 1

			httpStats.RequestTime += h.W.StopTime()

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

			httpStats.BytesOut += h.W.R.BytesOut

		case <-httpsrv.SignalShutdown:
			log.Debug("Shutting Down Stats")
			return
		}
	}
}

type runtimeStats struct {
	NumCPU     int `json:"num_cpu"`
	GoRoutines int `json:"go_routines"`

	Alloc        uint64 `json:"alloc"`
	TotalAlloc   uint64 `json:"total_alloc"`
	Sys          uint64 `json:"sys"`
	Lookups      uint64 `json:"lookups"`
	Mallocs      uint64 `json:"mallocs"`
	Frees        uint64 `json:"frees"`
	HeapAlloc    uint64 `json:"heap_alloc"`
	HeapSys      uint64 `json:"heap_sys"`
	HeapIdle     uint64 `json:"heap_idle"`
	HeapInuse    uint64 `json:"heap_inuse"`
	HeapReleased uint64 `json:"heap_released"`
	HeapObjects  uint64 `json:"heap_objects"`
	StackInuse   uint64 `json:"stack_inuse"`
	StackSys     uint64 `json:"stack_sys"`
	MSpanInuse   uint64 `json:"mspan_inuse"`
	MSpanSys     uint64 `json:"mspan_sys"`
	MCacheInuse  uint64 `json:"mcache_inuse"`
	MCacheSys    uint64 `json:"mcache_sys"`
	BuckHashSys  uint64 `json:"buckhash_sys"`
	GCSys        uint64 `json:"gc_sys"`
	OtherSys     uint64 `json:"other_sys"`
	NextGC       uint64 `json:"gc_next"`
	LastGC       uint64 `json:"gc_last"`
	PauseTotalNs uint64 `json:"gc_pause_total"`

	NumGC         uint32  `json:"gc_num"`
	NumForcedGC   uint32  `json:"gc_forced_num"`
	GCCPUFraction float64 `json:"gc_cpu_frac"`
	EnableGC      bool    `json:"gc_enable"`
	DebugGC       bool    `json:"gc_debug"`
}

func getRuntime() (s runtimeStats) {
	s.NumCPU = runtime.NumCPU()
	s.GoRoutines = runtime.NumGoroutine()

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	s.Alloc = m.Alloc
	s.TotalAlloc = m.TotalAlloc
	s.Sys = m.Sys
	s.Lookups = m.Lookups
	s.Mallocs = m.Mallocs
	s.Frees = m.Frees
	s.HeapAlloc = m.HeapAlloc
	s.HeapSys = m.HeapSys
	s.HeapIdle = m.HeapIdle
	s.HeapInuse = m.HeapInuse
	s.HeapReleased = m.HeapReleased
	s.HeapObjects = m.HeapObjects
	s.StackInuse = m.StackInuse
	s.StackSys = m.StackSys
	s.MSpanInuse = m.MSpanInuse
	s.MSpanSys = m.MSpanSys
	s.MCacheInuse = m.MCacheInuse
	s.MCacheSys = m.MCacheSys
	s.BuckHashSys = m.BuckHashSys
	s.GCSys = m.GCSys
	s.OtherSys = m.OtherSys
	s.NextGC = m.NextGC
	s.LastGC = m.LastGC
	s.PauseTotalNs = m.PauseTotalNs

	return
}

var dbHooks map[string]*dbstats.CounterHook

type dbStats struct {
	OpenConns     int `json:"conns_open"`
	TotalConns    int `json:"conns_total"`
	OpenStmts     int `json:"stmts_open"`
	TotalStmts    int `json:"stmts_total"`
	OpenTxs       int `json:"txs_open"`
	TotalTxs      int `json:"txs_total"`
	CommittedTxs  int `json:"txs_committed"`
	RolledbackTxs int `json:"txs_rolledback"`
	Queries       int `json:"queries"`
	Execs         int `json:"execs"`
	RowsIterated  int `json:"rows_inserted"`

	ConnErrs    int `json:"errs_conn"`
	StmtErrs    int `json:"errs_stmt"`
	TxOpenErrs  int `json:"errs_tx_open"`
	TxCloseErrs int `json:"errs_tx_close"`
	QueryErrs   int `json:"errs_query"`
	ExecErrs    int `json:"errs_exec"`
	RowErrs     int `json:"errs_row"`
}

func WrapDB(name string, fn dbstats.OpenFunc) {
	h := &dbstats.CounterHook{}
	s := dbstats.New(fn)
	s.AddHook(h)
	sql.Register(name, s)
	dbHooks[name] = h
}

func getDBstats() (m map[string]dbStats) {
	m = make(map[string]dbStats)
	for k, v := range dbHooks {
		s := dbStats{}
		s.OpenConns = v.OpenConns()
		s.TotalConns = v.TotalConns()
		s.OpenStmts = v.OpenStmts()
		s.TotalStmts = v.TotalStmts()
		s.OpenTxs = v.OpenTxs()
		s.TotalTxs = v.TotalTxs()
		s.CommittedTxs = v.CommittedTxs()
		s.RolledbackTxs = v.RolledbackTxs()
		s.Queries = v.Queries()
		s.Execs = v.Execs()
		s.RowsIterated = v.RowsIterated()
		s.ConnErrs = v.ConnErrs()
		s.StmtErrs = v.StmtErrs()
		s.TxOpenErrs = v.TxOpenErrs()
		s.TxCloseErrs = v.TxCloseErrs()
		s.QueryErrs = v.QueryErrs()
		s.ExecErrs = v.ExecErrs()
		s.RowErrs = v.RowErrs()
		m[k] = s
	}
	return m
}

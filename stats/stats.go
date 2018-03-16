package stats

import (
	"sour.is/x/toolbox/ident"
	"net/http"
	"sour.is/x/toolbox/httpsrv"
	"time"
)

var httpPipe chan httpData

func init() {
	stats.AppStart = time.Now()

	httpPipe = make(chan httpData)
	go recordStats(httpPipe)

	httpsrv.NewMiddleware("gather-stats", doStats).Register(httpsrv.EventComplete)
	httpsrv.IdentRegister("stats", httpsrv.IdentRoutes{
		{"GetStats", "GET", "/v1/stats", getStats},
	})
}

var stats struct{
	AppStart time.Time
	HttpStats struct{
		TotalRequests int
		TotalErrors int
		TotalAnonRequests int
		TotalBytesWritten int
	}
}

func getStats(w http.ResponseWriter, r *http.Request, id ident.Ident) {
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
				stats.HttpStats.TotalRequests += 1
				if h.W.GetCode() >= 400 {
					stats.HttpStats.TotalErrors += 1
				}
				if !h.ID.IsActive() {
					stats.HttpStats.TotalAnonRequests += 1
				}
				stats.HttpStats.TotalBytesWritten += h.W.R.BytesOut

			case <- httpsrv.SignalShutdown:
				return
		}
	}
}
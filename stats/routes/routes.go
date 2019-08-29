package stats

import (
	"net/http"

	"sour.is/x/toolbox/httpsrv"
	"sour.is/x/toolbox/ident"
	"sour.is/x/toolbox/stats"
)

func init() {
	httpsrv.IdentRegister("stats", httpsrv.IdentRoutes{
		{Name: "get-stats", Method: "GET", Pattern: "/v1/stats", HandlerFunc: getStats},
		{Name: "get-metrics", Method: "GET", Pattern: "/metrics", HandlerFunc: getMetrics},
	})

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
func getStats(w httpsrv.ResponseWriter, _ *http.Request, _ ident.Ident) {

	stats := stats.Get()

	httpsrv.WriteObject(w, http.StatusOK, stats)
}

// swagger:operation GET /metrics metrics getMetrics
//
// Get Prometheus Metrics
//
// ---
// produces:
//   - "text/plain"
// responses:
//   "200":
//     description: Success
//     schema:
//       type: string
func getMetrics(w httpsrv.ResponseWriter, _ *http.Request, _ ident.Ident) {

	stats := stats.Get()

	httpsrv.WriteText(w, http.StatusOK, stats.String())
}

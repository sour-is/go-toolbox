package stats

import (
	"fmt"
	"strings"
	"time"

	"sour.is/x/toolbox/stats/exposition"
)

var appStart time.Time

func init() {
	appStart = time.Now().In(time.UTC)
}

// Stats is a point in time status of various metrics
type Stats map[string]exposition.Expositions

func (s Stats) String() string {
	var out strings.Builder
	for _, v := range s {
		out.WriteString(v.String())
	}
	return out.String()
}

// Get returns a current copy of runtime stats.
func Get() Stats {
	s := make(Stats)
	for k, fn := range statRegistry {
		exp := fn()
		s[k] = exp.Exposition()
	}

	return s
}

// StatFn is a func to build stat values to return
type StatFn func() exposition.Expositioner

var statRegistry map[string]StatFn

// Register a new stats handler
func Register(name string, fn StatFn) {
	fmt.Println("stats register", name)
	statRegistry[name] = fn
}

// GetRegistry returns current registry
func GetRegistry() map[string]StatFn {
	return statRegistry
}

func init() {
	fmt.Println("stats init")
	statRegistry = make(map[string]StatFn)
	Register("app", getStats)
}

type appStats struct {
	AppStart   time.Time `json:"app_start"`
	UpTimeNano int64     `json:"uptime_nano"`
	UpTime     string    `json:"uptime"`
}

func getStats() exposition.Expositioner {
	return appStats{
		appStart,
		time.Since(appStart).Nanoseconds(),
		time.Since(appStart).String(),
	}
}
func (a appStats) Exposition() (lis exposition.Expositions) {
	var e exposition.Exposition
	e = exposition.New("app_start", exposition.Gauge)
	e.AddRow(float64(a.AppStart.Unix()))
	lis = append(lis, e)

	e = exposition.New("app_uptime", exposition.Counter)
	e.AddRow(time.Since(appStart).Seconds())

	return
}

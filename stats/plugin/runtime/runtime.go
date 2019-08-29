package runtime

import (
	"fmt"
	"reflect"
	std_runtime "runtime"

	"sour.is/x/toolbox/stats"
	"sour.is/x/toolbox/stats/exposition"
)

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

func getRuntime() exposition.Expositioner {
	var s runtimeStats
	s.NumCPU = std_runtime.NumCPU()
	s.GoRoutines = std_runtime.NumGoroutine()

	var m std_runtime.MemStats
	std_runtime.ReadMemStats(&m)

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

	return s
}

func (s runtimeStats) Exposition() (lis exposition.Expositions) {
	v := reflect.ValueOf(s)
	for i := 0; i < v.NumField(); i++ {
		var e exposition.Exposition
		tag := v.Type().Field(i).Tag.Get("json")

		switch tag {
		case "total_alloc", "lookups", "mallocs", "frees", "gc_pause_total":
			e = exposition.New(fmt.Sprintf("runtime_%s_totals", tag), exposition.Counter)
		default:
			e = exposition.New(fmt.Sprintf("runtime_%s", tag), exposition.Gauge)
		}

		e.AddRow(exposition.ToFloat(v.Field(i)))
		lis = append(lis, e)
	}

	return
}

func (s runtimeStats) String() string {
	return s.Exposition().String()
}

func init() {
	stats.Register("runtime", getRuntime)
}

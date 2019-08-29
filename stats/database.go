package stats

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"github.com/cgilling/dbstats"
	"sour.is/x/toolbox/stats/exposition"
)

type hook struct {
	Name string
	Hook *dbstats.CounterHook
}

var dbHooks map[string]hook

func init() {
	dbHooks = make(map[string]hook)
}

type dbStats struct {
	Name string `json:"name"`

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

// WrapDB wraps up a db connection to trace stats.
func WrapDB(name string, fn dbstats.OpenFunc) {
	h := hook{
		Name: name,
		Hook: &dbstats.CounterHook{},
	}

	s := dbstats.New(fn)
	s.AddHook(h.Hook)
	sql.Register(name, s)
	Register("db."+name, h.getDBstats)
	dbHooks[name] = h
}

// GetDBlist get a list of wrapped dbs
func GetDBlist() (lis []string) {
	for s := range dbHooks {
		lis = append(lis, s)
	}

	return
}

func (h hook) getDBstats() exposition.Expositioner {
	var s dbStats
	s.Name = h.Name
	s.OpenConns = h.Hook.OpenConns()
	s.TotalConns = h.Hook.TotalConns()
	s.OpenStmts = h.Hook.OpenStmts()
	s.TotalStmts = h.Hook.TotalStmts()
	s.OpenTxs = h.Hook.OpenTxs()
	s.TotalTxs = h.Hook.TotalTxs()
	s.CommittedTxs = h.Hook.CommittedTxs()
	s.RolledbackTxs = h.Hook.RolledbackTxs()
	s.Queries = h.Hook.Queries()
	s.Execs = h.Hook.Execs()
	s.RowsIterated = h.Hook.RowsIterated()
	s.ConnErrs = h.Hook.ConnErrs()
	s.StmtErrs = h.Hook.StmtErrs()
	s.TxOpenErrs = h.Hook.TxOpenErrs()
	s.TxCloseErrs = h.Hook.TxCloseErrs()
	s.QueryErrs = h.Hook.QueryErrs()
	s.ExecErrs = h.Hook.ExecErrs()
	s.RowErrs = h.Hook.RowErrs()

	return s
}

func (s dbStats) String() string {
	var out strings.Builder
	out.WriteString(s.Exposition().String())
	return out.String()
}

func (s dbStats) Exposition() (lis exposition.Expositions) {
	v := reflect.ValueOf(s)
	for i := 0; i < v.NumField(); i++ {
		var e exposition.Exposition
		tag := v.Type().Field(i).Tag.Get("json")

		switch tag {
		case "conns_open", "stmts_open", "txs_open":
			e = exposition.New(fmt.Sprintf("db_%s", tag), exposition.Gauge)
		default:
			e = exposition.New(fmt.Sprintf("db_%s", tag), exposition.Counter)
		}

		e.AddRow(exposition.ToFloat(v.Field(i))).AddTag("name", s.Name)
		lis = append(lis, e)
	}

	return
}

package dbm // import "sour.is/x/toolbox/dbm"

/*
Include the driver in your main package.

```
import (
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)
```

*/

import (
	"database/sql"
	"fmt"
	"regexp"
	"runtime/debug"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/spf13/viper"
	"sour.is/x/toolbox/log"
	"sour.is/x/toolbox/uuid"

	"github.com/jmoiron/sqlx"

	sq "github.com/Masterminds/squirrel"

	"reflect"
)

var db *sql.DB
var dbType string
var placeholder sq.PlaceholderFormat
var returns bool

// GetDB eturns a database connection.
// Depricated: Use Transaction instead.
func GetDB() (*sql.Tx, error) {
	if db == nil {
		return nil, fmt.Errorf("database is not configured")
	}

	return db.Begin()
}

func Config() {
	if viper.IsSet("database") {
		pfx := "db." + viper.GetString("database")
		var err error

		dbType = viper.GetString(pfx + ".type")
		connect := viper.GetString(pfx + ".connect")
		max_conn := viper.GetInt(pfx + ".max_conn")

		placeholder = sq.Question
		if strings.Contains(dbType,"postgres") {
			placeholder = sq.Dollar
			returns = true
		}

		if db, err = sql.Open(dbType, connect); err != nil {
			log.Fatal(err)
		}

		if max_conn != 0 {
			db.SetMaxOpenConns(max_conn)
		}

		if err = db.Ping(); err != nil {
			log.Fatal(err)
		}

		re := regexp.MustCompile(`:.*@`)

		log.Notice("DBM: Database Connected: ", re.ReplaceAllString(connect, ":****@"))
	}
}

type Asset struct {
	File func(string) ([]byte, error)
	Dir  func(string) ([]string, error)
}

func Migrate(a Asset) (err error) {
	if viper.IsSet("database") {
		pfx := "db." + viper.GetString("database")
		if !viper.GetBool(pfx) {
			log.Info("Migration is disabled.")
			return
		}
	}

	err = Transaction(func(tx *Tx) (err error) {

		if _, err = tx.Exec(sqlschema); err != nil {
			return
		}

		version := 0
		if err = tx.QueryRow(`select ifnull(max(version), 0) as version from schema_version`).Scan(&version); err != nil {
			log.Println("DBM: ", err)
			return
		}

		log.Infof("DBM: Current Schema Version: %04d", version)

		var d []string
		d, err = a.Dir("schema")
		if err != nil {
			return
		}

		sort.Strings(d)
		for _, name := range d {
			var v int
			v, err = strconv.Atoi(strings.SplitN(name, "-", 2)[0])
			if err != nil {
				continue
			}

			if version >= v {
				continue
			}

			log.Print("DBM: Migrating ", name)

			var file []byte
			if file, err = a.File("schema/" + name); err != nil {
				return
			}

			if _, err = tx.Exec(string(file)); err != nil {
				return
			}

			if _, err = tx.Exec(`insert into schema_version (version) values (?)`, v); err != nil {
				return
			}

			log.Print("DBM: Finished ", name)
		}

		return
	})

	if err != nil {
		log.Error("DBM: Migration Check Failed")
		return
	}

	log.Notice("DBM: Migration Check Complete")
	return
}

var sqlschema = `
CREATE TABLE IF NOT EXISTS schema_version (
  version     INT(8) NOT NULL,
  updated_on  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (version)
);`

type Tx struct{
	*sql.Tx
	DbType string
	Placeholder sq.PlaceholderFormat
	Returns bool
}

func NewTx(db *sql.DB, dbType string, placeholder sq.PlaceholderFormat, returns bool) (tx *Tx, err error) {
	tx = new(Tx)
	tx.Tx, err = db.Begin()
	tx.Placeholder = placeholder
	tx.DbType = dbType
	tx.Returns = returns

	return
}
// Transaction starts a new database transction and executes the supplied func.
func Transaction(txFunc func(*Tx) error) (err error) {
	tx, err := NewTx(db, dbType, placeholder, returns)

	if err != nil {
		log.Error(err.Error())
		return
	}
	defer func() {
		if p := recover(); p != nil {
			switch p := p.(type) {
			case error:
				err = p
			default:
				err = fmt.Errorf("%s", p)
			}
		}
		if err != nil {
			tx.Rollback()
			log.Error(err.Error())

			debug.PrintStack()
			return
		}
		err = tx.Commit()
	}()
	err = txFunc(tx)
	return err
}
func Transactionx(txFunc func(*sqlx.Tx) error) (err error) {
	dbx := sqlx.NewDb(db, dbType)

	tx, err := dbx.Beginx()
	if err != nil {
		log.Error("dbm.Transaction ", err.Error())
		return
	}
	defer func() {
		if p := recover(); p != nil {
			switch p := p.(type) {
			case error:
				err = p
			default:
				err = fmt.Errorf("%s", p)
			}
		}
		if err != nil {
			tx.Rollback()
			log.Error(err.Error())

			debug.PrintStack()
			return
		}
		err = tx.Commit()
	}()
	err = txFunc(tx)
	return err
}

var txMap = make(map[string]*Tx)
var txMutex = sync.Mutex{}

func txSet(id string, tx *Tx) {
	txMutex.Lock()
	defer txMutex.Unlock()
	txMap[id] = tx
}
func txGet(id string) *Tx {
	txMutex.Lock()
	defer txMutex.Unlock()
	return txMap[id]
}
func txRm(id string) {
	txMutex.Lock()
	defer txMutex.Unlock()
	delete(txMap, id)
}

// TransactionContinue returns a transaction that can be continued by suppling the
// TxID that gets passed into the txFunc.
func TransactionContinue(TxID string, txFunc func(*Tx, string) error) (err error) {
	var tx *Tx

	if TxID == "" {

		TxID = uuid.V4()
		tx, err = NewTx(db, dbType, placeholder, returns)
		if err != nil {
			log.Error(err.Error())
			return
		}
		txSet(TxID, tx)

	} else {
		if tx = txGet(TxID); tx != nil {
			err = txFunc(tx, TxID)
			return err
		}
		return fmt.Errorf("unable to find tx %s", TxID)
	}

	defer func() {
		if p := recover(); p != nil {
			switch p := p.(type) {
			case error:
				err = p
			default:
				err = fmt.Errorf("%s", p)
			}
		}

		if err != nil {
			tx.Rollback()
			log.Error(err.Error())

			debug.PrintStack()
			return
		}

		err = tx.Commit()
		if err != nil {
			log.Error(err.Error())
		}

		txRm(TxID)
	}()

	err = txFunc(tx, TxID)
	return err
}

func Count(tx *Tx, table string, where sq.Eq) (count uint64, err error) {

	s := sq.Select("count(1)")
	s = s.RunWith(tx)
	s = s.PlaceholderFormat(tx.Placeholder)

	s = s.From(table)
	if where != nil {
		s = s.Where(where)
	}

	log.Debug(s.ToSql())

	err = s.QueryRow().Scan(&count)

	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}

		log.Debug(err.Error())
		return
	}

	return count, nil
}

type FetchMap func(row *sql.Rows) error

func Fetch(tx *Tx, table string, cols []string, where sq.Eq, limit, offset uint64, fn FetchMap) (err error) {

	s := sq.Select(cols...)
	s = s.PlaceholderFormat(tx.Placeholder)
	s = s.RunWith(tx)

	s = s.From(table)
	s = s.Limit(limit)
	s = s.Offset(offset)
	if where != nil {
		s = s.Where(where)
	}

	rows, err := s.Query()
	if err != nil {
		return
	}
	defer rows.Close()
	return fn(rows)
}

func Insert(tx *Tx, table string) sq.InsertBuilder {
	s := sq.Insert(table)
	s = s.PlaceholderFormat(tx.Placeholder)
	s = s.RunWith(tx)

	return s
}

func Update(tx *Tx, table string) sq.UpdateBuilder {
	s := sq.Update(table)
	s = s.PlaceholderFormat(tx.Placeholder)
	s = s.RunWith(tx)

	return s
}

type DbInfo struct{
	Table   string
	Cols    []string
	columns map[string]string
}
func (d DbInfo) Col(column string) (col string) {
	var ok bool
	if col, ok = d.columns[column]; !ok {
		panic("Col not found on table: " + column)
	}

	return col
}
func GetDbInfo(o interface{}) (d DbInfo) {
	t := reflect.TypeOf(o)

	d.Table = t.Name()
	d.columns = make(map[string]string)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		table := field.Tag.Get("table")
		if table != "" {
			d.Table = table
		}
		tag := field.Tag.Get("db")

		if tag == "" {
			tag = field.Tag.Get("json")
		}
		if tag == "" || tag == "-" {
			tag = field.Name
		}
		d.columns[field.Name] = tag
		d.Cols = append(d.Cols, tag)
	}

	return d
}

func Replace(
	tx *Tx,
	table string,
	where sq.Eq,
	update sq.UpdateBuilder,
	insert sq.InsertBuilder,
) (found bool, id int64, err error) {

	var num uint64
	if num, err = Count(tx, table, where); err == nil && num == 0 {

		log.Debug(insert.ToSql())

		var result sql.Result
		result, err = insert.Exec()
		if err != nil {
			log.Warning(err.Error())
			return
		}

		var affected int64
		if affected, err = result.RowsAffected(); err != nil {
			return
		}
		if affected == 0 {
			err = fmt.Errorf("update Failed. %d rows affected", num)
		}

	} else if err == nil && num > 0 {

		found = true

		log.Debug(update.ToSql())

		if tx.Returns {
			result := update.QueryRow()
			err = result.Scan(&id)
			if err != nil {
				log.Debug(err.Error())
				return
			}

		} else {
			var result sql.Result
			result, err = update.Exec()
			if err != nil {
				log.Debug(err.Error())
				return
			}

			id, err = result.LastInsertId()
			if err != nil {
				log.Debug(err.Error())
				return
			}
		}
	}

	return
}
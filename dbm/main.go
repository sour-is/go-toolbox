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
)

var db *sql.DB

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

		name := viper.GetString(pfx + ".type")
		connect := viper.GetString(pfx + ".connect")
		max_conn := viper.GetInt(pfx + ".max_conn")

		if db, err = sql.Open(name, connect); err != nil {
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

	err = Transaction(func(tx *sql.Tx) (err error) {

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

// Transaction starts a new database transction and executes the supplied func.
func Transaction(txFunc func(*sql.Tx) error) (err error) {
	tx, err := db.Begin()
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

var txMap = make(map[string]*sql.Tx)
var txMutex = sync.Mutex{}

func txSet(id string, tx *sql.Tx) {
	txMutex.Lock()
	defer txMutex.Unlock()
	txMap[id] = tx
}
func txGet(id string) *sql.Tx {
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
func TransactionContinue(TxID string, txFunc func(*sql.Tx, string) error) (err error) {
	var tx *sql.Tx

	if TxID == "" {

		TxID = uuid.V4()
		tx, err = db.Begin()
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

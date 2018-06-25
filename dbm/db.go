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
	"github.com/spf13/viper"
	"regexp"
	"sour.is/x/toolbox/log"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"time"
)

type DB struct {
	Conn        *sql.DB
	DbType      string
	Placeholder sq.PlaceholderFormat
	Returns     bool
}

// GetDB returns a database connection.
func GetDB(pfx string) (db DB, err error) {

	dbType := viper.GetString(pfx + ".type")
	connect := viper.GetString(pfx + ".connect")
	maxConn := viper.GetInt(pfx + ".max_conn")
	maxLifetime := viper.GetInt(pfx + ".max_lifetime")

	if dbType == "" {
		log.Fatal("Database Type is not set!")
	}
	if connect == "" {
		log.Fatal("Database Connect is not set!")
	}

	var conn *sql.DB
	if conn, err = sql.Open(dbType, connect); err != nil {
		log.Error(err)
		return
	}

	if maxConn != 0 {
		conn.SetMaxOpenConns(maxConn)
	}

	if maxLifetime == 0 {
		maxLifetime = 5
	}
	conn.SetConnMaxLifetime(5 * time.Minute)

	if err = conn.Ping(); err != nil {
		log.Error(err)
		return
	}

	db.Conn = conn
	db.DbType = dbType
	db.Placeholder = sq.Question
	if strings.Contains(db.DbType, "postgres") {
		db.Placeholder = sq.Dollar
		db.Returns = true
	}

	connect = regexp.MustCompile(`:.*@`).ReplaceAllString(connect, ":****@")
	connect = regexp.MustCompile(`password=.[[:graph:]]+`).ReplaceAllString(connect, "password=****")

	log.Notice("DBM: Database Connected: ", connect)

	return
}

var stdDB DB

func Config() {
	if viper.IsSet("database") {
		pfx := "db." + viper.GetString("database")

		var err error

		stdDB, err = GetDB(pfx)
		if err != nil {
			log.Fatal(err.Error())
		}
	}
}

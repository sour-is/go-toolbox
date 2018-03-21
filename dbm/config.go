package dbm

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/spf13/viper"
	"regexp"
	"sour.is/x/toolbox/log"
	"strings"
)

var dbType string
var placeholder sq.PlaceholderFormat
var returns bool

func Config() {
	if viper.IsSet("database") {
		pfx := "db." + viper.GetString("database")
		var err error

		dbType = viper.GetString(pfx + ".type")
		connect := viper.GetString(pfx + ".connect")
		max_conn := viper.GetInt(pfx + ".max_conn")

		placeholder = sq.Question
		if strings.Contains(dbType, "postgres") {
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

		connect = regexp.MustCompile(`:.*@`).ReplaceAllString(connect, ":****@")
		connect = regexp.MustCompile(`password=.[[:graph:]]+`).ReplaceAllString(connect, "password=****")

		log.Notice("DBM: Database Connected: ", connect)
	}
}

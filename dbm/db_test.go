package dbm // import "sour.is/x/toolbox/dbm"

import (
	"bytes"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/spf13/viper"
	"database/sql"
	"fmt"
	"sour.is/x/toolbox/log"
)

var defaultTestConfig = `
database = "test"

[db.test]
type = "sqlmock"
connect = "user@host/db"

[db.pg]
type = "postgres"
connect = """
    host=localhost
    port=5432
    dbname=jonlundy
    sslmode=disable
    """
`

func TestMain(m *testing.M) {
	viper.SetConfigType("toml")
	viper.ReadConfig(bytes.NewBuffer([]byte(defaultTestConfig)))
	fmt.Println("HELLO")
	log.SetVerbose(log.Vdebug)
	log.SetFlags(log.LstdFlags|log.Lshortfile)

	m.Run()
}

func TestConfig(t *testing.T) {
	Convey("Given a mock database", t, func() {
		log.Println("TEST RUN")
		var mockDB *sql.DB
		var err error
		var mock sqlmock.Sqlmock

		mockDB, mock, err = sqlmock.NewWithDSN("user@host/db")
		if err != nil {
			log.Printf("Unable to create database: %s", err)
		}
		defer mockDB.Close()
		mock.ExpectBegin()

		GetDB("db.test")

		// we make sure that all expectations were met
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Logf("there were unfulfilled expectations: %s", err)
		}
		So(err, ShouldBeNil)

	})
}

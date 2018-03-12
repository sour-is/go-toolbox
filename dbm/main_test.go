package dbm // import "sour.is/x/toolbox/dbm"

import (
	"bytes"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/spf13/viper"
)

var defaultTestConfig = `
database = "test"
[db.test]
type = "sqlmock"
connect = "user@host/db"
`

func TestMain(m *testing.M) {
	viper.SetConfigType("toml")
	viper.ReadConfig(bytes.NewBuffer([]byte(defaultTestConfig)))
	m.Run()
}

func TestGetDB(t *testing.T) {
	Convey("Given a mock database", t, func() {

		var err error
		var mock sqlmock.Sqlmock

		db, mock, err = sqlmock.New()
		if err != nil {
			t.Logf("Unable to create database: %s", err)
		}
		So(err, ShouldBeNil)
		defer db.Close()

		mock.ExpectBegin()

		GetDB()

		// we make sure that all expectations were met
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Logf("there were unfulfilled expectations: %s", err)
		}
		So(err, ShouldBeNil)

	})
}

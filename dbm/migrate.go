package dbm

import (
	"github.com/spf13/viper"
	"sort"
	"sour.is/x/toolbox/log"
	"strconv"
	"strings"
)

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

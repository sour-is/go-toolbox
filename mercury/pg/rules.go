package pg

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"sour.is/x/toolbox/dbm"
	"sour.is/x/toolbox/log"
	"sour.is/x/toolbox/mercury"
)

// GetRules get list of rules
func GetRules(user string) (lis mercury.Rules, err error) {
	err = dbm.Transaction(func(tx *dbm.Tx) error {
		return tx.Fetch(
			"mercury_rules_vw",
			[]string{"role", "type", "match"},
			squirrel.Eq{"user_id": user},
			0, 0, nil,
			func(rows *sql.Rows) (err error) {
				var role, typ, match string
				for rows.Next() {
					err = rows.Scan(&role, &typ, &match)
					if err != nil {
						log.Debug(err)
						return
					}
					lis = append(lis, mercury.Rule{Role: role, Type: typ, Match: match})
				}
				return err
			},
		)
	})

	return
}

package pg

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"sour.is/x/toolbox/dbm"
	"sour.is/x/toolbox/log"
)

// GetGroups get list of groups
func GetGroups(user string) (lis []string, err error) {
	err = dbm.Transaction(func(tx *dbm.Tx) error {
		return tx.Fetch(
			"souris_groups_vw",
			[]string{"group_id"},
			squirrel.Eq{"user_id": user},
			0, 0, nil,
			func(rows *sql.Rows) (err error) {
				var s string
				for rows.Next() {
					err = rows.Scan(&s)
					if err != nil {
						log.Debug(err)
						return
					}
					lis = append(lis, s)
				}
				return err
			},
		)
	})

	return
}

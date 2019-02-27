package pg

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"sour.is/x/toolbox/dbm"
	"sour.is/x/toolbox/log"
	"sour.is/x/toolbox/mercury"
)

// Notify stores the attributes for a registry space
type Notify struct {
	Name   string `json:"name" view:"mercury_notify_vw"`
	Match  string `json:"match"`
	Event  string `json:"event"`
	Method string `json:"-" db:"method"`
	URL    string `json:"-" db:"url"`
}

// GetNotify get list of rules
func GetNotify(event string) (lis mercury.ListNotify) {
	err := dbm.Transaction(func(tx *dbm.Tx) error {
		return tx.Fetch(
			"mercury_notify_vw",
			[]string{"name", "match", "event", "method", "url"},
			squirrel.Eq{"event": event},
			0, 0, nil,
			func(rows *sql.Rows) (err error) {
				var name, match, event, method, url string
				for rows.Next() {
					err = rows.Scan(&name, &match, &event, &method, &url)
					if err != nil {
						log.Debug(err)
						return
					}
					log.Debugf("%s %s %s %s %s", name, match, event, method, url)
					lis = append(lis, mercury.Notify{
						Name:   name,
						Match:  match,
						Event:  event,
						Method: method,
						URL:    url,
					})
				}
				return err
			},
		)
	})
	if err != nil {
		log.Error(err)
		return nil
	}

	return
}

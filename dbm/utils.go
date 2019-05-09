package dbm

import (
	"database/sql"
	"fmt"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"sour.is/x/toolbox/log"
)

// Count returns a row count for query
func Count(tx *Tx, table string, where interface{}) (count uint64, err error) {
	return tx.Count(table, where)
}

// FetchMap is passed to fetch for assigning to objects.
type FetchMap func(row *sql.Rows) error

// Fetch fetch rows and pass through supplied function
func Fetch(tx *Tx, table string, cols []string, where interface{}, limit, offset uint64, order []string, fn FetchMap) (err error) {
	return tx.Fetch(table, cols, where, limit, offset, order, fn)
}

// Insert begin new insert statement.
func Insert(tx *Tx, table string) sq.InsertBuilder {
	return tx.Insert(table)
}

// Update begin new update statement.
func Update(tx *Tx, table string) sq.UpdateBuilder {
	return tx.Update(table)
}

// Replace begin new replace statement.
func Replace(
	tx *Tx,
	d DbInfo,
	o interface{},
	where interface{},
	update sq.UpdateBuilder,
	insert sq.InsertBuilder,
) (found bool, err error) {
	return tx.Replace(d, o, where, update, insert)
}

// Count returns a row count for query
func (tx *Tx) Count(table string, where interface{}) (count uint64, err error) {

	s := sq.Select("count(1)")
	s = s.RunWith(tx.Tx)
	s = s.PlaceholderFormat(tx.Placeholder)

	s = s.From(table)
	if where != nil {
		s = s.Where(where)
	}

	ss, args, _ := s.ToSql()
	log.Debug("Count", "sql", ss, "args", args)

	err = s.QueryRowContext(tx.Context).Scan(&count)

	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}

		log.Debug(err.Error())
		return
	}

	return count, nil
}

// Fetch fetch rows and pass through supplied function
func (tx *Tx) Fetch(table string, cols []string, where interface{}, limit, offset uint64, order []string, fn FetchMap) (err error) {
	s := sq.Select(cols...)
	s = s.PlaceholderFormat(tx.Placeholder)
	s = s.RunWith(tx.Tx)

	s = s.From(table)
	if limit > 0 {
		s = s.Limit(limit)
		s = s.Offset(offset)
	}
	if where != nil {
		s = s.Where(where)
	}
	if len(order) > 0 {
		s = s.OrderBy(order...)
	}

	sql, args, _ := s.ToSql()
	log.Debugs("Fetch", "sql", sql, "args", args)

	rows, err := s.QueryContext(tx.Context)
	if err != nil {
		return
	}
	defer rows.Close()
	return fn(rows)
}

// Select begin new select statement.
func (tx *Tx) Select(cols []string, table string) sq.SelectBuilder {
	s := sq.Select(cols...).From(table)
	s = s.PlaceholderFormat(tx.Placeholder)
	s = s.RunWith(tx.Tx)

	return s
}

// Insert begin new insert statement.
func (tx *Tx) Insert(table string) sq.InsertBuilder {
	s := sq.Insert(table)
	s = s.PlaceholderFormat(tx.Placeholder)
	s = s.RunWith(tx.Tx)

	return s
}

// Update begin new update statement.
func (tx *Tx) Update(table string) sq.UpdateBuilder {
	s := sq.Update(table)
	s = s.PlaceholderFormat(tx.Placeholder)
	s = s.RunWith(tx.Tx)

	return s
}

// Delete begin new delete statement.
func (tx *Tx) Delete(table string) sq.DeleteBuilder {
	s := sq.Delete(table)
	s = s.PlaceholderFormat(tx.Placeholder)
	s = s.RunWith(tx.Tx)

	return s
}

// Replace begin new replace statement.
func (tx *Tx) Replace(
	d DbInfo,
	o interface{},
	where interface{},
	update sq.UpdateBuilder,
	insert sq.InsertBuilder,
) (found bool, err error) {
	var num uint64
	auto, dest, err := d.StructMap(o, d.Auto)

	if num, err = tx.Count(d.Table, where); err == nil && num == 0 {
		if tx.Returns {
			if len(d.Auto) > 0 {
				insert = insert.Suffix(`RETURNING "` + strings.Join(auto, `","`) + "\"")

				s, a, _ := insert.ToSql()
				log.Debugs("Replace", "sql", s, "args", a, "return", d.Auto, "return_for", auto)

				row := insert.QueryRowContext(tx.Context)
				err = row.Scan(dest...)
				if err != nil {
					log.Debug(err.Error())
					return
				}
			} else {
				log.Debug(insert.ToSql())
				insert.ExecContext(tx.Context)
			}

		} else {
			s, a, _ := insert.ToSql()
			log.Debugs("Replace", "sql", s, "args", a)

			var result sql.Result
			result, err = insert.ExecContext(tx.Context)
			if err != nil {
				log.Debug(err.Error())
				return
			}

			var lastID int64
			lastID, err = result.LastInsertId()
			if err != nil {
				log.Debug(err.Error())
				return
			}

			row := tx.Select(auto, d.View).Where(sq.Eq{d.Auto[0]: lastID}).QueryRow()
			err = row.Scan(dest...)
			if err != nil {
				log.Debug(err.Error())
				return
			}
		}

	} else if err == nil && num > 0 {
		found = true
		update = update.Where(where)

		if tx.Returns {
			if len(d.Auto) > 0 {
				update = update.Suffix(`RETURNING "` + strings.Join(auto, `","`) + "\"")

				s, a, _ := update.ToSql()
				log.Debugs("Replace", "sql", s, "args", a, "return", d.Auto, "return_for", auto)

				row := update.QueryRowContext(tx.Context)
				err = row.Scan(dest...)
				if err != nil {
					log.Debug(err.Error())
					return
				}
			} else {
				log.Debug(update.ToSql())
				update.ExecContext(tx.Context)
			}
		} else {
			s, a, _ := update.ToSql()
			log.Debugs("Replace", "sql", s, "args", a)

			var result sql.Result

			result, err = update.ExecContext(tx.Context)
			if err != nil {
				log.Warning(err.Error())
				return
			}

			var affected int64
			if affected, err = result.RowsAffected(); err != nil {
				return
			}

			if affected == 0 {
				found = false
				err = fmt.Errorf("update Failed. %d rows affected", num)
			}

			row := tx.Select(auto, d.View).Where(where).QueryRowContext(tx.Context)
			err = row.Scan(dest...)
			if err != nil {
				log.Debug(err.Error())
				return
			}
		}
	}

	return
}

// Ping attempt a connection to database
func Ping() bool {
	err := stdDB.Conn.Ping()
	if err != nil {
		log.Error(err)
		return false
	}
	return true
}

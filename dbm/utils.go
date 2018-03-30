package dbm

import (
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"sour.is/x/toolbox/log"
)

func Count(tx *Tx, table string, where sq.Eq) (count uint64, err error) {
	return tx.Count(table, where)
}

type FetchMap func(row *sql.Rows) error

func Fetch(tx *Tx, table string, cols []string, where sq.Eq, limit, offset uint64, fn FetchMap) (err error) {
	return tx.Fetch(table, cols, where, limit, offset, fn)
}

func Insert(tx *Tx, table string) sq.InsertBuilder {
	return tx.Insert(table)
}

func Update(tx *Tx, table string) sq.UpdateBuilder {
	return tx.Update(table)
}

func Replace(
	tx *Tx,
	table string,
	where sq.Eq,
	update sq.UpdateBuilder,
	insert sq.InsertBuilder,
) (found bool, id uint64, err error) {
	return tx.Replace(table, where, update, insert)
}

func (tx *Tx) Count(table string, where sq.Eq) (count uint64, err error) {

	s := sq.Select("count(1)")
	s = s.RunWith(tx.Tx)
	s = s.PlaceholderFormat(tx.Placeholder)

	s = s.From(table)
	if where != nil {
		s = s.Where(where)
	}

	log.Debug(s.ToSql())

	err = s.QueryRow().Scan(&count)

	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}

		log.Debug(err.Error())
		return
	}

	return count, nil
}

func (tx *Tx) Fetch(table string, cols []string, where sq.Eq, limit, offset uint64, fn FetchMap) (err error) {

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
	log.Debug(s.ToSql())
	rows, err := s.Query()
	if err != nil {
		return
	}
	defer rows.Close()
	return fn(rows)
}

func (tx *Tx) Insert(table string) sq.InsertBuilder {
	s := sq.Insert(table)
	s = s.PlaceholderFormat(tx.Placeholder)
	s = s.RunWith(tx.Tx)

	return s
}

func (tx *Tx) Update(table string) sq.UpdateBuilder {
	s := sq.Update(table)
	s = s.PlaceholderFormat(tx.Placeholder)
	s = s.RunWith(tx.Tx)

	return s
}

func (tx *Tx) Replace(
	table string,
	where sq.Eq,
	update sq.UpdateBuilder,
	insert sq.InsertBuilder,
) (found bool, id uint64, err error) {

	var num uint64
	if num, err = Count(tx, table, where); err == nil && num == 0 {

		log.Debug(insert.ToSql())

		if tx.Returns {
			var result sq.RowScanner
			result = insert.QueryRow()
			err = result.Scan(&id)
			log.Debugf("RETURNING %d", id)
			if err != nil {
				log.Debug(err.Error())
				return
			}

		} else {
			var result sql.Result
			result, err = insert.Exec()
			if err != nil {
				log.Debug(err.Error())
				return
			}
			var lastId int64
			lastId, err = result.LastInsertId()
			if err != nil {
				log.Debug(err.Error())
				return
			}

			id = uint64(lastId)
		}

	} else if err == nil && num > 0 {

		found = true
		update = update.Where(where)

		log.Debug(update.ToSql())
		var result sql.Result

		result, err = update.Exec()
		if err != nil {
			log.Warning(err.Error())
			return
		}

		var affected int64
		if affected, err = result.RowsAffected(); err != nil {
			return
		}

		if affected == 0 {
			err = fmt.Errorf("update Failed. %d rows affected", num)
		}
	}

	return
}

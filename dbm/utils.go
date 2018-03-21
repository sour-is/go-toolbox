package dbm

import (
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"sour.is/x/toolbox/log"
)

func Count(tx *Tx, table string, where sq.Eq) (count uint64, err error) {

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

type FetchMap func(row *sql.Rows) error

func Fetch(tx *Tx, table string, cols []string, where sq.Eq, limit, offset uint64, fn FetchMap) (err error) {

	s := sq.Select(cols...)
	s = s.PlaceholderFormat(tx.Placeholder)
	s = s.RunWith(tx.Tx)

	s = s.From(table)
	s = s.Limit(limit)
	s = s.Offset(offset)
	if where != nil {
		s = s.Where(where)
	}

	rows, err := s.Query()
	if err != nil {
		return
	}
	defer rows.Close()
	return fn(rows)
}

func Insert(tx *Tx, table string) sq.InsertBuilder {
	s := sq.Insert(table)
	s = s.PlaceholderFormat(tx.Placeholder)
	s = s.RunWith(tx.Tx)

	return s
}

func Update(tx *Tx, table string) sq.UpdateBuilder {
	s := sq.Update(table)
	s = s.PlaceholderFormat(tx.Placeholder)
	s = s.RunWith(tx.Tx)

	return s
}

func Replace(
	tx *Tx,
	table string,
	where sq.Eq,
	update sq.UpdateBuilder,
	insert sq.InsertBuilder,
) (found bool, id uint64, err error) {

	var num uint64
	if num, err = Count(tx, table, where); err == nil && num == 0 {

		log.Debug(insert.ToSql())

		var result sql.Result
		result, err = insert.Exec()
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

	} else if err == nil && num > 0 {

		found = true

		log.Debug(update.ToSql())

		if tx.Returns {
			result := update.QueryRow()
			err = result.Scan(&id)
			if err != nil {
				log.Debug(err.Error())
				return
			}

		} else {
			var result sql.Result
			result, err = update.Exec()
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
	}

	return
}
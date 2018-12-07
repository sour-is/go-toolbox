
package ctrl

import (
	"context"
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"

	model "<...model...>"
	"sour.is/x/toolbox/dbm"
	"sour.is/x/toolbox/dbm/qry"
	"sour.is/x/toolbox/log"

	opentracing "github.com/opentracing/opentracing-go"
)

// Mode of insert or update
type Mode int

const (
	// Insert row to table
	Insert Mode = iota
	// Update row in table
	Update
)

// NotFoundError matching row not found
type NotFoundError string

// Error format message as string
func (n NotFoundError) Error() string {
	return string(n)
}

// ParseError failed to parse query
type ParseError string

// Error format message as string
func (n ParseError) Error() string {
	return string(n)
}

{{range .Types}}
// {{.Name}} adds transaction to model
type {{.Name}} struct {
	*model.{{.Name}}
	Where sq.Eq
	Tx *dbm.Tx
}
func get{{.Name}}Tx(tx *dbm.Tx, q qry.Input) (lis []model.{{.Name}}, err error) {
	sp, _ := opentracing.StartSpanFromContext(tx.Context, "ctrl.get{{.Name}}Tx")
	defer sp.Finish()
	
	var o model.{{.Name}}
	if q.DbInfo == nil {
		d := dbm.GetDbInfo(o)
		q.DbInfo = &d
	}
	dcols, dest, err := q.StructMap(&o, nil)
	if err != nil {
		log.Debug(err)
		return
	}
	err = tx.Fetch(q.View, dcols, q.Search, q.Limit, q.Offset, q.Sort,
		func(rows *sql.Rows) (err error) {
			for rows.Next() {
                i := 0
				{{range .Fields}}{{if .Container}}i, _ = q.Index("{{.Name}}"); dest[i] = &model.{{.Container}}{}{{end}}
                {{end}}

				err = rows.Scan(dest...)
				if err != nil {
					log.Debug(err)
					return
				}

				{{range .Fields}}{{if .Container}}i, _ = q.Index("{{.Name}}"); o.{{.Name}} = *dest[i].(*model.{{.Container}}){{end}}
                {{end}}
				lis = append(lis, o)
                _ = i
			}
			return
		})
	if err != nil {
		log.Debug(err)
		return
	}
    if lis == nil {
        lis = make([]model.{{.Name}}, 0)
    }

	return
}
{{if .ROnly}}{{else}}
func delete{{.Name}}Tx(tx *dbm.Tx, where interface{}) (err error) {
	db := dbm.GetDbInfo(model.{{.Name}}{})

	// ----
    _, err = tx.Delete(db.Table).Where(where).Exec()
	if err != nil {
		log.Debug(err)
		return
	}

	return
}

// Delete{{.Name}} will delete using the provided id.
func Delete{{.Name}}ByID(id uint64) (err error) {
    var idx string
	d := dbm.GetDbInfo(model.{{.Name}}{})
    idx, err = d.Col("ID")
    if err != nil {
        return
    }
    err = dbm.Transaction(func(tx *dbm.Tx) (err error) {
        return delete{{.Name}}Tx(tx, sq.Eq{idx: id})
    })

    return
}

// Delete{{.Name}} will delete using the provided where statement.
func Delete{{.Name}}(where interface{}) (err error) {
    err = dbm.Transaction(func(tx *dbm.Tx) (err error) {
        return delete{{.Name}}Tx(tx, where)
    })

    return
}

// Save stores {{.Name}} to DB
func (o {{.Name}}) Save() (err error) {
	op := o.{{.Name}}
	
	// ----
	if op == nil {
		err = fmt.Errorf("uninitialized {{.Name}}")
		return
	}
	d := dbm.GetDbInfo(*op)
	// ----
	
	setMap := map[string]interface{}{
		{{range .Fields}}
		{{if .Auto}}{{else if .ROnly}}{{else}}d.ColPanic("{{.Name}}"): {{if .Container}}model.{{.Container}}(o.{{.Name}}){{else}}o.{{.Name}}{{end}},{{end}}
        {{end}}
    }
	
	// ----
	_, err = o.Tx.Replace(
		d, op, o.Where,
		dbm.Update(o.Tx, d.Table).SetMap(setMap),
		dbm.Insert(o.Tx, d.Table).SetMap(setMap),
	)
	if err != nil {
		log.Error(err)
	}
	return
}
// Put{{.Name}} prepares a transaction to save {{.Name}}
func Put{{.Name}}(id uint64, fn func(Mode, dbm.DbInfo, {{.Name}}) error) (o model.{{.Name}}, err error) {
		d := dbm.GetDbInfo(o)

		co := {{.Name}}{
			{{.Name}}: &o,
			Where: sq.Eq{d.ColPanic("ID"): id},
		}
		gfn := get{{.Name}}Tx

		// ----
		var nerr error // Non aborting error
		err = dbm.Transaction(func(tx *dbm.Tx) (err error) {
			var mode = Insert

			if id != 0 {
				lis, gErr := gfn(tx, qry.Input{&d, co.Where, 1, 0, nil})
				if gErr != nil {
					err = gErr
					return
				}
				if len(lis) > 0 {
					mode = Update
					o = lis[0]
				}
			}

			co.Tx = tx
			err = fn(mode, d, co)
			if err != nil {
				switch err.(type) {
				case NotFoundError:
					nerr = err
					err = nil
				case ParseError:
					nerr = err
					err = nil
				}

				return
			}

			return
		})
		if nerr != nil {
			err = nerr
		}

		return
}
{{end}}

// List{{.Name}} queries a list of {{.Name}}
func List{{.Name}}(where interface{}, limit, offset uint64, order []string) (lis []model.{{.Name}}, err error) {
	return List{{.Name}}Context(context.Background(), where, limit, offset, order)
}
// List{{.Name}}Context queries a list of {{.Name}} with Context
func List{{.Name}}Context(ctx context.Context, where interface{}, limit, offset uint64, order []string) (lis []model.{{.Name}}, err error) {
	d := dbm.GetDbInfo(model.{{.Name}}{})

	return List{{.Name}}Qry(ctx, qry.Input{&d, where, limit, offset, order})
}
// List{{.Name}}Qry queries a list of {{.Name}} with Context
func List{{.Name}}Qry(ctx context.Context, q qry.Input) (lis []model.{{.Name}}, err error) {
		fn := get{{.Name}}Tx
    
        // ----
		err = dbm.QueryContext(ctx, func(tx *dbm.Tx) (err error) {
			lis, err = fn(tx, q)
			return
		})
		return
}
// List{{.Name}}ByID will query a list of {{.Name}} by ids
func List{{.Name}}ByID(ids []uint64) (lis []model.{{.Name}}, count uint64, err error) {
	return List{{.Name}}ByIDContext(context.Background(), ids)
}
// List{{.Name}}ByIDContext will query a list of {{.Name}} by ids with Context
func List{{.Name}}ByIDContext(ctx context.Context, ids []uint64) (lis []model.{{.Name}}, count uint64, err error) {
		d := dbm.GetDbInfo(model.{{.Name}}{})
		fn := List{{.Name}}Qry

        // ----
		where := sq.Eq{d.ColPanic("ID"): ids}
		count = uint64(len(ids))
		lis, err = fn(ctx, qry.Input{&d, where, count, 0, nil})
		return
}
// List{{.Name}}Count will count the number of items returned
func List{{.Name}}Count(where interface{}, limit, offset uint64) (lis []model.{{.Name}}, count uint64, err error) {
	return List{{.Name}}CountContext(context.Background(), where, limit, offset)
}
// List{{.Name}}CountContext will count the number of items returned with Context
func List{{.Name}}CountContext(ctx context.Context, where interface{}, limit, offset uint64) (lis []model.{{.Name}}, count uint64, err error) {
		d := dbm.GetDbInfo(model.{{.Name}}{})
		fn := get{{.Name}}Tx

		// ----
		err = dbm.QueryContext(ctx, func(tx *dbm.Tx) (err error) {
			count, _ = dbm.Count(tx, d.View, where)
			lis, err = fn(tx, qry.Input{&d, where, limit, offset, nil})
			return
		})
		return
}
{{end}}

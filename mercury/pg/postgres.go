package pg

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Masterminds/squirrel"
	opentracing "github.com/opentracing/opentracing-go"
	"sour.is/x/toolbox/dbm"
	"sour.is/x/toolbox/dbm/qry"
	"sour.is/x/toolbox/gql"
	"sour.is/x/toolbox/log"
	"sour.is/x/toolbox/mercury"
)

// Space stores a space value
type Space struct {
	ID    uint64   `json:"id" db:",AUTO" table:"mercury_spaces"`
	Space string   `json:"space"`
	Tags  []string `json:"tags" cont:"ListStrings"`
	Notes []string `json:"notes" cont:"ListStrings"`
	Items []Value  `json:"items" db:"-"`
}

// Value stores the attributes for a registry space
type Value struct {
	ID     uint64   `json:"id" db:",AUTO" table:"mercury_values"`
	Seq    uint64   `json:"seq"`
	Name   string   `json:"name"`
	Values []string `json:"values" cont:"ListStrings"`
	Notes  []string `json:"notes"  cont:"ListStrings"`
	Tags   []string `json:"tags"   cont:"ListStrings"`
}

// Config stores the attributes for a registry space
type Config struct {
	ID     uint64   `json:"id" view:"mercury_registry_vw"`
	Seq    uint64   `json:"seq"`
	Space  string   `json:"space"`
	Name   string   `json:"name"`
	Values []string `json:"values" cont:"ListStrings"`
	Notes  []string `json:"notes" cont:"ListStrings"`
	Tags   []string `json:"tags" cont:"ListStrings"`
}

// ListSpaceContext queries a list of User with Context
func ListSpaceContext(ctx context.Context, where interface{}, limit, offset uint64, order []string) (lis []Space, err error) {
	d := dbm.GetDbInfo(Space{})

	return ListSpaceQry(ctx, qry.Input{
		DbInfo: &d,
		Search: where,
		Limit:  limit,
		Offset: offset,
		Sort:   order,
	})
}

// ListSpace queries a list of Space
func ListSpace(where interface{}, limit, offset uint64, order []string) (lis []Space, err error) {
	return ListSpaceContext(context.Background(), where, limit, offset, order)
}

// ListSpaceQry queries a list of Space with Context
func ListSpaceQry(ctx context.Context, q qry.Input) (lis []Space, err error) {
	fn := getSpaceTx

	// ----
	err = dbm.QueryContext(ctx, func(tx *dbm.Tx) (err error) {
		lis, err = fn(tx, q)
		return
	})
	return
}

// SpaceTx adds transaction to model
type SpaceTx struct {
	*Space
	Where squirrel.Eq
	Tx    *dbm.Tx
}

func getSpaceTx(tx *dbm.Tx, q qry.Input) (lis []Space, err error) {
	sp, _ := opentracing.StartSpanFromContext(tx.Context, "ctrl.getSpaceTx")
	defer sp.Finish()

	var o Space
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
				i, _ = q.Index("Notes")
				dest[i] = &gql.ListStrings{}
				i, _ = q.Index("Tags")
				dest[i] = &gql.ListStrings{}

				err = rows.Scan(dest...)
				if err != nil {
					log.Debug(err)
					return
				}

				i, _ = q.Index("Notes")
				o.Notes = *dest[i].(*gql.ListStrings)
				i, _ = q.Index("Tags")
				o.Tags = *dest[i].(*gql.ListStrings)

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
		lis = make([]Space, 0)
	}

	return
}

// Save stores Space to DB
func (o SpaceTx) Save() (err error) {
	op := o.Space

	// ----
	if op == nil {
		err = fmt.Errorf("uninitialized Space")
		return
	}
	d := dbm.GetDbInfo(*op)
	// ----

	setMap := map[string]interface{}{
		d.ColPanic("Space"): o.Space.Space,
		d.ColPanic("Notes"): gql.ListStrings(o.Notes),
		d.ColPanic("Tags"):  gql.ListStrings(o.Tags),
	}
	log.Debug("REPLACE ", setMap, " WHERE ", o.Where)
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

// PutSpace prepares a transaction to save Space
func PutSpace(id uint64, fn func(Mode, dbm.DbInfo, SpaceTx) error) (o Space, err error) {
	d := dbm.GetDbInfo(o)

	co := SpaceTx{
		Space: &o,
		Where: squirrel.Eq{d.ColPanic("ID"): id},
	}
	gfn := getSpaceTx

	// ----
	var nerr error // Non aborting error
	err = dbm.Transaction(func(tx *dbm.Tx) (err error) {
		var mode = Insert

		if id != 0 {
			lis, gErr := gfn(tx, qry.Input{DbInfo: &d, Search: co.Where, Limit: 1})
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

// ConfigTx adds transaction to model
type ConfigTx struct {
	*Config
	Where squirrel.Eq
	Tx    *dbm.Tx
}

func getConfigTx(tx *dbm.Tx, q qry.Input) (lis []Config, err error) {
	sp, _ := opentracing.StartSpanFromContext(tx.Context, "ctrl.getConfigTx")
	defer sp.Finish()

	var o Config
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
				i, _ = q.Index("Values")
				dest[i] = &gql.ListStrings{}
				i, _ = q.Index("Notes")
				dest[i] = &gql.ListStrings{}
				i, _ = q.Index("Tags")
				dest[i] = &gql.ListStrings{}

				err = rows.Scan(dest...)
				if err != nil {
					log.Debug(err)
					return
				}

				i, _ = q.Index("Values")
				o.Values = *dest[i].(*gql.ListStrings)
				i, _ = q.Index("Notes")
				o.Notes = *dest[i].(*gql.ListStrings)
				i, _ = q.Index("Tags")
				o.Tags = *dest[i].(*gql.ListStrings)

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
		lis = make([]Config, 0)
	}

	return
}

// ListConfigQry queries a list of Config with Context
func ListConfigQry(ctx context.Context, q qry.Input) (lis []Config, err error) {
	fn := getConfigTx

	// ----
	err = dbm.QueryContext(ctx, func(tx *dbm.Tx) (err error) {
		lis, err = fn(tx, q)
		return
	})
	return
}

// ListConfigContext queries a list of User with Context
func ListConfigContext(ctx context.Context, where interface{}, limit, offset uint64, order []string) (lis []Config, err error) {
	d := dbm.GetDbInfo(Config{})

	return ListConfigQry(ctx, qry.Input{
		DbInfo: &d,
		Search: where,
		Limit:  limit,
		Offset: offset,
		Sort:   order,
	})
}

// ListConfig queries a list of Config
func ListConfig(where interface{}, limit, offset uint64, order []string) (lis []Config, err error) {
	return ListConfigContext(context.Background(), where, limit, offset, order)
}

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

// Mode of insert or update
type Mode int

const (
	// Insert row to table
	Insert Mode = iota
	// Update row in table
	Update
)

// WriteConfig writes a config map to database
func WriteConfig(tx *dbm.Tx, config mercury.ArraySpace) (err error) {
	d := dbm.GetDbInfo(Space{})

	// convert to map
	spaceMap := config.ToSpaceMap()

	// get names of each space
	var names []string
	for k := range spaceMap {
		names = append(names, k)
	}

	// get current spaces
	lis, err := ListSpace(squirrel.Eq{d.ColPanic("Space"): names}, 0, 0, nil)
	if err != nil {
		return
	}

	// determine which are being updated
	var updateSpaces []Space
	ids := make(map[string]uint64)
	for _, n := range lis {
		ids[n.Space] = n.ID
		updateSpaces = append(updateSpaces, n)
	}

	// update spaces
	for _, u := range updateSpaces {
		s, ok := spaceMap[u.Space]
		if !ok {
			continue
		}

		_, err = PutSpace(u.ID, func(mode Mode, db dbm.DbInfo, o SpaceTx) (err error) {
			o.Notes = s.Notes
			o.Tags = s.Tags

			return o.Save()
		})
		if err != nil {
			return
		}
		log.Debugf("UPDATED %d SPACES", len(updateSpaces))
	}

	// determine spaces to add
	var newSpaces []Space
	var newNames []string
	var curIDs []uint64
	for _, n := range names {
		if id, ok := ids[n]; !ok {
			newNames = append(newNames, n)
			s := spaceMap[n]
			newSpaces = append(newSpaces, Space{
				Space: s.Space,
				Tags:  s.Tags,
				Notes: s.Notes,
			})
		} else {
			curIDs = append(curIDs, id)
		}
	}

	// if there are any new spaces allocate and assign ids
	if len(newNames) > 0 {
		err = tx.Fetch(
			fmt.Sprintf("generate_series( 1, %d )", len(newNames)),
			[]string{"nextval('mercury_spaces_id_seq')"},
			nil, 0, 0, nil,
			func(row *sql.Rows) error {
				var u uint64
				i := 0
				for row.Next() {

					err := row.Scan(&u)
					if err != nil {
						return err
					}
					newSpaces[i].ID = u
					ids[newNames[i]] = u

					i++
				}
				return nil
			})
		if err != nil {
			return
		}

		// write new spaces
		err = WriteSpaces(tx, newSpaces)
		if err != nil {
			return
		}
		log.Debugf("WROTE %d NEW SPACES", len(newSpaces))
	}

	// extract all values
	var attrs []Value
	for ns, c := range spaceMap {
		nsID := ids[ns]
		for i, a := range c.List {
			attrs = append(attrs, Value{
				ID:     nsID,
				Seq:    uint64(i),
				Name:   a.Name,
				Values: a.Values,
				Tags:   a.Tags,
				Notes:  a.Notes,
			})
		}
	}

	// write all values to db.
	err = WriteValues(tx, squirrel.Eq{dbm.GetDbInfo(Value{}).ColPanic("ID"): curIDs}, attrs)
	log.Debugf("WROTE %d ATTRS", len(attrs))

	return
}

// WriteSpaces writes the spaces to db
func WriteSpaces(tx *dbm.Tx, lis []Space) (err error) {
	d := dbm.GetDbInfo(Space{})

	if len(lis) == 0 {
		return nil
	}

	newInsert := func() squirrel.InsertBuilder {
		return tx.Insert(d.Table).Columns(
			d.ColPanic("ID"),
			d.ColPanic("Space"),
			d.ColPanic("Tags"),
			d.ColPanic("Notes"),
		)
	}
	chunk := int(65000 / 3)
	insert := newInsert()
	for i, s := range lis {
		insert = insert.Values(
			s.ID,
			s.Space,
			gql.ListStrings(s.Tags),
			gql.ListStrings(s.Notes),
		)

		if i > 0 && i%chunk == 0 {
			log.Debugf("inserting %v rows into %v", i%chunk, d.Table)
			log.Debug(insert.ToSql())

			_, err = insert.Exec()
			if err != nil {
				log.Error(err)
				return
			}

			insert = newInsert()
		}
	}
	if len(lis)%chunk > 0 {
		log.Debugf("inserting %v rows into %v", len(lis)%chunk, d.Table)
		log.Debug(insert.ToSql())

		_, err = insert.Exec()
		if err != nil {
			log.Error(err)
			return
		}
	}

	return
}

// WriteValues writes the values to db
func WriteValues(tx *dbm.Tx, delete squirrel.Sqlizer, lis []Value) (err error) {
	d := dbm.GetDbInfo(Value{})

	log.Debug("DELETE ", delete)
	_, err = tx.Delete(d.Table).Where(delete).Exec()
	if err != nil {
		return
	}
	log.Debug(d.Table, len(lis))

	if len(lis) == 0 {
		return nil
	}

	newInsert := func() squirrel.InsertBuilder {
		return tx.Insert(d.Table).Columns(
			d.ColPanic("ID"),
			d.ColPanic("Seq"),
			d.ColPanic("Name"),
			d.ColPanic("Values"),
			d.ColPanic("Notes"),
			d.ColPanic("Tags"),
		)
	}
	chunk := int(65000 / 3)
	insert := newInsert()
	for i, s := range lis {
		insert = insert.Values(
			s.ID,
			s.Seq,
			s.Name,
			gql.ListStrings(s.Values),
			gql.ListStrings(s.Notes),
			gql.ListStrings(s.Tags),
		)
		log.Debug(s.Name)

		if i > 0 && i%chunk == 0 {
			log.Debugf("inserting %v rows into %v", i%chunk, d.Table)
			log.Debug(insert.ToSql())

			_, err = insert.Exec()
			if err != nil {
				log.Error(err)
				return
			}

			insert = newInsert()
		}
	}
	if len(lis)%chunk > 0 {
		log.Debugf("inserting %v rows into %v", len(lis)%chunk, d.Table)
		log.Debug(insert.ToSql())

		_, err = insert.Exec()
		if err != nil {
			log.Error(err)
			return
		}
	}

	return
}

func getWhere(search mercury.NamespaceSearch, d dbm.DbInfo) squirrel.Sqlizer {
	var where squirrel.Or
	for _, m := range search {
		switch m.(type) {
		case mercury.NamespaceNode:
			where = append(where, squirrel.Eq{d.ColPanic("Space"): m.Value()})
		case mercury.NamespaceStar:
			where = append(where, squirrel.Like{d.ColPanic("Space"): m.Value()})
		case mercury.NamespaceTrace:
			e := squirrel.Expr(`? LIKE concat(`+d.ColPanic("Space")+`, '%')`, m.Value())
			where = append(where, e)
		}
	}
	return where
}

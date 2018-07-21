package dbm

import (
	"fmt"
	"reflect"
	"strings"
)

// DbInfo database model metadata
type DbInfo struct {
	Table string
	View  string
	Cols  []string
	SCols []string
	index map[string]int
	Auto  []string
}
// SCol returns the struct column names
func (d DbInfo) SCol(column string) string {
	return d.SCols[d.Index(column)]
}
// Col returns the mapped column names
func (d DbInfo) Col(column string) string {
	return d.Cols[d.Index(column)]
}
// Index returns the column number
func (d DbInfo) Index(column string) (idx int) {
	var ok bool
	if idx, ok = d.index[column]; !ok {
		panic("Col not found on table: " + column)
	}

	return
}
// GetDbInfo builds a metadata struct
func GetDbInfo(o interface{}) (d DbInfo) {
	t := reflect.TypeOf(o)

	d.Table = t.Name()
	d.index = make(map[string]int)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		table := field.Tag.Get("table")
		if table != "" {
			d.Table = table
		}
		view := field.Tag.Get("view")
		if view != "" {
			d.View = view
		}

		sp := append(strings.SplitN(field.Tag.Get("db"), ",", 2), "")

		tag, opt := sp[0], sp[1]

		if opt == "AUTO" {
			d.Auto = append(d.Auto, field.Name)
		}

		if tag == "" {
			tag = field.Tag.Get("json")
		}

		if tag == "-" {
			continue
		}

		if tag == "" {
			tag = field.Name
		}

		d.index[field.Name] = len(d.Cols)
		d.Cols = append(d.Cols, tag)
		d.SCols = append(d.SCols, field.Name)
	}
	if d.View == "" {
		d.View = d.Table
	}

	return d
}

// StructMap accepts a struct and the columns to be set and returns a []interface{} that can be passed to a query scan
func (d DbInfo) StructMap(o interface{}, cols []string) (fields []string, targets []interface{}, err error) {
	if cols == nil {
		cols = d.SCols
	}

	r := reflect.ValueOf(o)
	e := r.Elem()
	if e.Kind() == reflect.Struct {
		// exported field
		for _, field := range cols {
			f := e.FieldByName(field)
			if f.IsValid() {
				// A Value can be changed only if it is
				// addressable and was not obtained by
				// the use of unexported struct fields.
				if f.CanSet() && f.CanAddr() {
					fields = append(fields, d.Col(field))
					targets = append(targets, f.Addr().Interface())
				} else {
					err = fmt.Errorf("field %s cannot be set", field)
					break
				}
			} else {
				err = fmt.Errorf("field %s is not valid", field)
				break
			}
		}
	} else {
		err = fmt.Errorf("object %s is not struct", e.Kind())
	}

	return
}

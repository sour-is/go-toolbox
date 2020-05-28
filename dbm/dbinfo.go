package dbm

import (
	"fmt"
	"reflect"
	"strings"
	"unicode"
)

// DbInfo database model metadata
type DbInfo struct {
	Table string
	View  string
	Cols  []string
	SCols []string
	GCols []string
	index map[string]int
	Auto  []string
	ID    string
	HasID bool
	ROnly []string
}

// SCol returns the struct column names
func (d DbInfo) SCol(column string) (s string, err error) {
	idx, ok := d.Index(column)
	if !ok {
		err = fmt.Errorf("column not found on table: %v", column)
		return
	}
	return d.SCols[idx], err
}

// GCol returns the graphql column names
func (d DbInfo) GCol(column string) (s string, err error) {
	idx, ok := d.Index(column)
	if !ok {
		err = fmt.Errorf("column not found on table: %v", column)
		return
	}
	return d.GCols[idx], err
}

// Col returns the mapped column names
func (d DbInfo) Col(column string) (s string, err error) {
	idx, ok := d.Index(column)
	if !ok {
		err = fmt.Errorf("column not found on table: %v", column)
		return
	}
	return d.Cols[idx], err
}

// ColPanic returns the mapped column names will panic if col does not exist
func (d DbInfo) ColPanic(column string) string {
	idx, ok := d.Index(column)
	if !ok {
		err := fmt.Errorf("column not found on table: %v", column)
		panic(err)
	}
	return d.Cols[idx]
}

// Index returns the column number
func (d DbInfo) Index(column string) (idx int, ok bool) {
	idx, ok = d.index[column]
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

		dbField := field.Tag.Get("db")
		sp := append(strings.SplitN(dbField, ",", 2), "")

		tag := sp[0]

		if strings.Contains(dbField, ",AUTO") {
			d.Auto = append(d.Auto, field.Name)
		}

		if strings.Contains(dbField, ",RO") {
			d.ROnly = append(d.ROnly, field.Name)
		}

		if strings.Contains(dbField, ",ID") {
			d.HasID = true
			d.ID = field.Name
		}

		json := field.Tag.Get("json")
		if tag == "" {
			tag = json
		}

		graphql := field.Tag.Get("graphql")
		if tag == "" {
			tag = graphql
		}

		if tag == "-" {
			continue
		}

		if tag == "" {
			tag = field.Name
		}

		d.index[field.Name] = len(d.Cols)

		if _, ok := d.index[tag]; !ok && tag != "" {
			d.index[tag] = len(d.Cols)
		}
		if _, ok := d.index[json]; !ok && json != "" {
			d.index[json] = len(d.Cols)
		}
		if _, ok := d.index[graphql]; !ok && graphql != "" {
			d.index[graphql] = len(d.Cols)
		} else if !ok && graphql == "" {
			a := []rune(field.Name)
			for i := 0; i < len(a); i++ {
				if unicode.IsLower(a[i]) {
					break
				}
				a[i] = unicode.ToLower(a[i])
			}
			graphql = string(a)
			d.index[graphql] = len(d.Cols)
		}

		d.Cols = append(d.Cols, tag)
		d.SCols = append(d.SCols, field.Name)
		d.GCols = append(d.GCols, graphql)
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
					col, err := d.Col(field)
					if err != nil {
						break
					}
					fields = append(fields, col)
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

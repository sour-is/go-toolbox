package dbm

import (
	"reflect"
	"strings"
)

type DbInfo struct {
	Table   string
	Cols    []string
	columns map[string]string
	Auto    []string
}

func (d DbInfo) Col(column string) (col string) {
	var ok bool
	if col, ok = d.columns[column]; !ok {
		panic("Col not found on table: " + column)
	}

	return col
}
func GetDbInfo(o interface{}) (d DbInfo) {
	t := reflect.TypeOf(o)

	d.Table = t.Name()
	d.columns = make(map[string]string)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		table := field.Tag.Get("table")
		if table != "" {
			d.Table = table
		}
		sp := strings.SplitN(field.Tag.Get("db"),",", 2)

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

		d.columns[field.Name] = tag
		d.Cols = append(d.Cols, tag)
	}

	return d
}

func ApplyUint(o interface{}, auto []string, vals []uint64) {
	r := reflect.ValueOf(o)
	e := r.Elem()
	if e.Kind() == reflect.Struct {
		// exported field
		for i, field := range auto {
			f := e.FieldByName(field)
			if f.IsValid() {
				// A Value can be changed only if it is
				// addressable and was not obtained by
				// the use of unexported struct fields.
				if f.CanSet() {
					// change value of N
					switch f.Kind() {
					case reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8, reflect.Uint:
						if !f.OverflowUint(vals[i]) {
							f.SetUint(vals[i])
						}
					case reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8, reflect.Int:
						if int64(vals[i]) >= 0 && !f.OverflowInt(int64(vals[i])) {
							f.SetInt(int64(vals[i]))
						}
					}
				}
			}
		}
	}
}

func ApplyInt(o interface{}, auto []string, vals []int64) {
	r := reflect.ValueOf(o)
	e := r.Elem()
	if e.Kind() == reflect.Struct {
		// exported field
		for i, field := range auto {
			f := e.FieldByName(field)
			if f.IsValid() {
				// A Value can be changed only if it is
				// addressable and was not obtained by
				// the use of unexported struct fields.
				if f.CanSet() {
					// change value of N
					switch f.Kind() {
					case reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8, reflect.Uint:
						if !f.OverflowUint(uint64(vals[i])) {
							f.SetUint(uint64(vals[i]))
						}
					case reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8, reflect.Int:
						if !f.OverflowInt(vals[i]) {
							f.SetInt(vals[i])
						}
					}
				}
			}
		}
	}
}
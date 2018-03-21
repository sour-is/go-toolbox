package dbm

import "reflect"

type DbInfo struct {
	Table   string
	Cols    []string
	columns map[string]string
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
		tag := field.Tag.Get("db")

		if tag == "" {
			tag = field.Tag.Get("json")
		}
		if tag == "" || tag == "-" {
			tag = field.Name
		}
		d.columns[field.Name] = tag
		d.Cols = append(d.Cols, tag)
	}

	return d
}

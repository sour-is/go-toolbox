package exposition

import (
	"fmt"
	"reflect"
	"strings"

	"sour.is/x/toolbox/mercury"
)

// Tags for counter differentiation
type Tags map[string]string

// Type for metric counter type
type Type string

// Metric types
const (
	Counter Type = "counter"
	Gauge        = "gauge"
	Summary      = "summary"
)

// Exposition for tracking/formatting counters
type Exposition struct {
	Name string
	Type Type
	Rows []Row
}

// Expositions list of Expositions
type Expositions []Exposition

// Row is a line item
type Row struct {
	Tags  Tags
	Value float64
}

// String outputs a string reprisentation
func (row Row) String() string {
	var out strings.Builder
	var tags []string
	for key, val := range row.Tags {
		tags = append(tags, fmt.Sprintf("%s=\"%s\"", key, val))
	}
	if len(tags) > 0 {
		out.WriteString("{")
		out.WriteString(strings.Join(tags, ","))
		out.WriteString("}")

	}
	out.WriteString(fmt.Sprintf(" %v\n", row.Value))
	return out.String()
}

// String outputs a string reprisentation
func (e Exposition) String() string {
	var out strings.Builder
	out.WriteString(fmt.Sprintf("# TYPE %s %s\n", e.Name, e.Type))
	for _, row := range e.Rows {
		out.WriteString(e.Name)
		out.WriteString(row.String())
	}

	return out.String()
}

// String outputs a string reprisentation
func (e Expositions) String() string {
	var out strings.Builder
	for _, exp := range e {
		out.WriteString(exp.String())
	}

	return out.String()
}

// New creates a new Exposition
func New(name string, expType Type) (e Exposition) {
	e.Name, e.Type = name, expType

	return
}

// AddRow creates a row on exposition
func (e *Exposition) AddRow(value float64) *Row {
	var row Row
	row.Value, row.Tags = value, make(Tags)

	e.Rows = append(e.Rows, row)
	return &row
}

// AddTag adds a tag to a row
func (row *Row) AddTag(name, value string) *Row {
	row.Tags[name] = value

	return row
}

// ToFloat converts a reflect.Value to float64
func ToFloat(v reflect.Value) float64 {
	switch v.Type().Name() {
	case "float32", "float64":
		return float64(v.Float())
	case "bool":
		var b int
		if v.Bool() {
			b = 1
		}
		return float64(b)
	case "uint", "uint64", "uint32":
		return float64(v.Uint())
	case "int", "int32", "int64":
		return float64(v.Int())
	}
	return 0.0
}

// Expositioner converts a stat object into a list of expositions.
type Expositioner interface {
	Exposition() Expositions
}

// ToSpaceValues Convert Exposition to Mercury Space
func (e Expositions) ToSpaceValues() (lis []mercury.Value) {
	var seq uint64
	for _, exp := range e {
		for _, row := range exp.Rows {
			seq++

			v := mercury.Value{
				Name:   exp.Name,
				Seq:    seq,
				Values: []string{fmt.Sprintf("%v", row.Value)},
			}

			v.Tags = append(v.Tags, string(exp.Type))
			for key, val := range row.Tags {
				v.Tags = append(v.Tags, fmt.Sprintf("%s/%s", key, val))
			}

			lis = append(lis, v)
		}
	}

	return
}

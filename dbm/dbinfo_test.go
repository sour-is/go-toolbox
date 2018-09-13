package dbm

import "testing"
import (
	. "github.com/smartystreets/goconvey/convey"
)

type testStruct struct {
	Uint64 uint64
	Uint32 uint32
	Uint16 uint16
	Uint8  uint8

	Int64 int64
	Int32 int32
	Int16 int16
	Int8  int8
}

type dbTable struct {
	Col1 int       `table:"table_name" view:"view_name"`
	Col2 string    `db:"col2,AUTO"`
	Col3 []float32 `json:"col3" graphql:"graphCol3"`
	Col4 uint64    `graphql:"graphCol4"`
}

// TestDbInfo test the DbInfo function
func TestDbInfo(t *testing.T) {
	Convey("Givin a table struct", t, func() {
		o := dbTable{}
		d := GetDbInfo(o)

		So(d.Table, ShouldEqual, "table_name")
		So(d.View, ShouldEqual, "view_name")

		tests := []struct {
			input  string
			expect interface{}
		}{
			{"Col1", 0},
			{"Col2", 1},
			{"col2", 1},
			{"Col3", 2},
			{"col3", 2},
			{"graphCol3", 2},
			{"Col4", 3},
			{"graphCol4", 3},
		}
		for _, tt := range tests {
			i, err := d.Index(tt.input)
			So(err, ShouldBeNil)
			So(i, ShouldEqual, tt.expect)
		}

		tests = []struct {
			input  string
			expect interface{}
		}{
			{"Col1", "Col1"},
			{"Col2", "Col2"},
			{"col2", "Col2"},
			{"Col3", "Col3"},
			{"graphCol3", "Col3"},
			{"col3", "Col3"},
			{"Col4", "Col4"},
			{"graphCol4", "Col4"},
		}
		for _, tt := range tests {
			i, err := d.SCol(tt.input)
			So(err, ShouldBeNil)
			So(i, ShouldEqual, tt.expect)
		}

		tests = []struct {
			input  string
			expect interface{}
		}{
			{"Col1", "Col1"},
			{"Col2", "col2"},
			{"col2", "col2"},
			{"Col3", "col3"},
			{"graphCol3", "col3"},
			{"graphCol4", "graphCol4"},
		}
		for _, tt := range tests {
			i, err := d.Col(tt.input)
			So(err, ShouldBeNil)
			So(i, ShouldEqual, tt.expect)
		}

		So(len(d.Auto), ShouldEqual, 1)
		So(d.Auto[0], ShouldContainSubstring, "Col2")
	})
}

func BenchmarkDbInfo(b *testing.B) {
	o := dbTable{}

	for i := 0; i < b.N; i++ {
		GetDbInfo(o)
	}
}

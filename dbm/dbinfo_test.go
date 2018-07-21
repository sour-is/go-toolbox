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
	Col3 []float32 `json:"col3"`
}

// TestDbInfo test the DbInfo function
func TestDbInfo(t *testing.T) {
	Convey("Givin a table struct", t, func() {
		o := dbTable{}
		d := GetDbInfo(o)

		So(d.Table, ShouldEqual, "table_name")
		So(d.View, ShouldEqual, "view_name")

		So(d.Index("Col1"), ShouldEqual, 0)
		So(d.Index("Col2"), ShouldEqual, 1)
		So(d.Index("Col3"), ShouldEqual, 2)

		So(d.SCol("Col1"), ShouldEqual, "Col1")
		So(d.SCol("Col2"), ShouldEqual, "Col2")
		So(d.SCol("Col3"), ShouldEqual, "Col3")

		So(d.Col("Col1"), ShouldEqual, "Col1")
		So(d.Col("Col2"), ShouldEqual, "col2")
		So(d.Col("Col3"), ShouldEqual, "col3")

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

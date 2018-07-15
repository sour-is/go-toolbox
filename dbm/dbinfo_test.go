package dbm

import "testing"
import . "github.com/smartystreets/goconvey/convey"

type TestStruct struct {
	Uint64 uint64
	Uint32 uint32
	Uint16 uint16
	Uint8  uint8

	Int64 int64
	Int32 int32
	Int16 int16
	Int8  int8
}

const MaxUint = ^uint64(0)
const MinUint = 0
const MaxInt = int64(MaxUint >> 1)
const MinInt = -MaxInt - 1

func TestApplyUint(t *testing.T) {
	o := TestStruct{}

	Convey("Given a set of Uints", t, func() {
		ApplyUint(&o, []string{"Uint64", "Uint32", "Uint16", "Uint8"}, []uint64{MaxUint, MaxUint, 256, 256})

		So(o.Uint64, ShouldEqual, MaxUint)
		So(o.Uint32, ShouldEqual, 0)
		So(o.Uint16, ShouldEqual, 256)
		So(o.Uint8, ShouldEqual, 0)
	})

	Convey("Given a set of Ints", t, func() {
		ApplyInt(&o, []string{"Int64", "Int32", "Int16", "Int8"}, []int64{MaxInt, MaxInt, 128, 128})

		So(o.Int64, ShouldEqual, MaxInt)
		So(o.Int32, ShouldEqual, 0)
		So(o.Int16, ShouldEqual, 128)
		So(o.Int8, ShouldEqual, 0)
	})
}

type dbTable struct {
	Col1 int       `table:"table_name" view:"view_name"`
	Col2 string    `db:"col2,AUTO"`
	Col3 []float32 `json:"col3"`
}

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

	GetDbInfo(o)
}
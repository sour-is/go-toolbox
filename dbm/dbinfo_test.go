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

	Convey("Given a set of Uints", t, func(){
		ApplyUint(&o, []string{"Uint64","Uint32", "Uint16", "Uint8"}, []uint64{MaxUint, MaxUint, 256, 256})

		So(o.Uint64, ShouldEqual, MaxUint)
		So(o.Uint32, ShouldEqual, 0)
		So(o.Uint16, ShouldEqual, 256)
		So(o.Uint8, ShouldEqual, 0)
	})

	Convey("Given a set of Ints", t, func(){
		ApplyInt(&o, []string{"Int64","Int32","Int16","Int8"}, []int64{MaxInt, MaxInt, 128, 128})

		So(o.Int64, ShouldEqual, MaxInt)
		So(o.Int32, ShouldEqual, 0)
		So(o.Int16, ShouldEqual, 128)
		So(o.Int8,  ShouldEqual, 0)
	})
}


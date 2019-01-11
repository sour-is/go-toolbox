package gql

import "testing"
import . "github.com/smartystreets/goconvey/convey"

func TestUnmarshalUint64(t *testing.T) {
	type test struct {
		In  interface{}
		Out uint64
		Err bool
	}
	tt := []test{
		{"", 0, false},
		{"0", 0, false},
		{"1", 1, false},
		{`""`, 0, false},
		{"xxx", 0, true},
		{3, 3, false},
	}

	Convey("Tests for UnmarshalUint64", t, func() {
		for _, v := range tt {
			r, err := UnmarshalUint64(v.In)

			if v.Err {
				So(err, ShouldNotBeNil)
			} else {
				So(err, ShouldBeNil)
			}

			So(r, ShouldEqual, v.Out)
		}
	})
}

func TestUnmarshalUint32(t *testing.T) {
	type test struct {
		In  interface{}
		Out uint32
		Err bool
	}
	tt := []test{
		{"", 0, false},
		{"0", 0, false},
		{"1", 1, false},
		{`""`, 0, false},
		{"xxx", 0, true},
		{3, 3, false},
	}

	Convey("Tests for UnmarshalUint32", t, func() {
		for _, v := range tt {
			r, err := UnmarshalUint32(v.In)

			if v.Err {
				So(err, ShouldNotBeNil)
			} else {
				So(err, ShouldBeNil)
			}

			So(r, ShouldEqual, v.Out)
		}
	})
}

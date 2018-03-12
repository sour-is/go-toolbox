package uuid // import "sour.is/x/toolbox/uuid"

import (
	"crypto/rand"
	"fmt"
	"strings"
	"testing"

	"github.com/bouk/monkey"

	. "github.com/smartystreets/goconvey/convey"
)

func TestV4(t *testing.T) {
	Convey("Generate a V4 UUID", t, func() {
		v4 := V4()
		VERSION := []uint8{'4'}
		CHECK := []uint8{'8', '9', 'a', 'b'}

		So(v4[14], ShouldBeIn, VERSION)
		So(v4[19], ShouldBeIn, CHECK)
	})

	Convey("Given randomness runs out", t, func() {
		monkey.Patch(rand.Read, func(p []byte) (n int, err error) {
			return 0, fmt.Errorf("Something Broke!")
		})
		v4 := V4()
		So(v4, ShouldEqual, NilUuid)
	})
}

func TestV5(t *testing.T) {
	Convey("Generate a V5 UUID", t, func() {
		v4 := V5("test", NilUuid)
		VERSION := []uint8{'5'}
		CHECK := []uint8{'8', '9', 'a', 'b'}

		So(v4[14], ShouldBeIn, VERSION)
		So(v4[19], ShouldBeIn, CHECK)
	})
}

func TestV6(t *testing.T) {
	Convey("Generate a V6 UUID", t, func() {
		v4 := V6(NilUuid, false)
		VERSION := []uint8{'6'}
		CHECK := []uint8{'8', '9', 'a', 'b'}

		So(v4[14], ShouldBeIn, VERSION)
		So(v4[19], ShouldBeIn, CHECK)
	})
}

func TestBytes(t *testing.T) {
	Convey("Convert UUID string to []byte", t, func() {
		b := Bytes(NilUuid)

		So(len(b), ShouldEqual, 16)
	})
}

func TestFromHexChar(t *testing.T) {
	Convey("Test conversion of hex values to bytes", t, func() {
		for i, c := range hextable {
			d, ok := fromHexChar(c)
			So(d, ShouldEqual, i)
			So(ok, ShouldBeTrue)
		}

		for i, c := range strings.ToUpper(hextable) {
			d, ok := fromHexChar(c)
			So(d, ShouldEqual, i)
			So(ok, ShouldBeTrue)
		}

		c, ok := fromHexChar('-')
		So(c, ShouldEqual, 0)
		So(ok, ShouldBeFalse)
	})
}

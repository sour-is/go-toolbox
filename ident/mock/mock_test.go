package mock // import "sour.is/x/toolbox/ident/mock"

import (
	"net/http"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"sour.is/x/toolbox/ident"
)

func TestMock(t *testing.T) {

	Convey("On init mock should register itself.", t, func() {
		So(ident.IdentSet, ShouldContainKey, "mock")
	})

	Convey("Given a valid request with config options", t, func() {
		ident.RegisterConfig("mock", map[string]string{
			"ident":  "ident",
			"aspect": "aspect",
			"name":   "name",
		})

		req := http.Request{}
		req.RemoteAddr = "127.0.0.1"

		Convey("A request is received", func() {
			u := ident.GetIdent("mock", &req)
			So(u.GetIdentity(), ShouldEqual, "ident")
			So(u.GetAspect(), ShouldEqual, "aspect")
			So(u.GetDisplay(), ShouldContainSubstring, "name")
			So(u.GetDisplay(), ShouldContainSubstring, req.RemoteAddr)
			So(u.IsActive(), ShouldBeTrue)
			So(u.HasRole("any"), ShouldBeTrue)
			So(u.HasGroup("any"), ShouldBeTrue)
		})
	})
}

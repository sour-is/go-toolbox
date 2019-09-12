package header // import "sour.is/x/toolbox/ident/header"

import (
	"net/http"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"sour.is/x/toolbox/ident"
)

func TestHeaderIdent(t *testing.T) {
	Convey("On init header should register itself.", t, func() {
		So(ident.GetHandlers(), ShouldContainKey, "header")
	})

	Convey("Given a valid request", t, func() {
		req := http.Request{}
		req.Header = make(http.Header)
		req.Header.Set("user_ident", "ident")
		req.Header.Set("user_aspect", "aspect")
		req.Header.Set("user_name", "name")

		Convey("Header returns a valid Ident", func() {
			u := ident.GetIdent("header", &req)
			So(u.GetIdentity(), ShouldEqual, "ident")
			So(u.GetAspect(), ShouldEqual, "aspect")
			So(u.GetDisplay(), ShouldEqual, "name")
			So(u.IsActive(), ShouldBeTrue)
			So(u.HasRole("any"), ShouldBeTrue)
			So(u.HasGroup("any"), ShouldBeTrue)
		})
	})

	Convey("Given a invalid request", t, func() {
		req := http.Request{}
		req.Header = make(http.Header)
		req.Header.Set("user_ident", "")
		req.Header.Set("user_aspect", "")
		req.Header.Set("user_name", "")

		Convey("Header returns a valid Ident", func() {
			u := ident.GetIdent("header", &req)
			So(u.GetIdentity(), ShouldEqual, "anon")
			So(u.GetAspect(), ShouldEqual, "none")
			So(u.GetDisplay(), ShouldEqual, "Guest User")
			So(u.IsActive(), ShouldBeFalse)
			So(u.HasRole("any"), ShouldBeFalse)
			So(u.HasGroup("any"), ShouldBeFalse)
		})
	})

	Convey("Given a invalid request with config options", t, func() {
		ident.RegisterConfig("header", map[string]string{
			"ident":  "ident",
			"aspect": "aspect",
			"name":   "name",
		})

		req := http.Request{}
		req.Header = make(http.Header)
		req.Header.Set("user_ident", "anon")
		req.Header.Set("user_aspect", "default")
		req.Header.Set("user_name", "Guest User")

		Convey("Header returns a valid Ident", func() {
			u := ident.GetIdent("header", &req)
			So(u.GetIdentity(), ShouldEqual, "anon")
			So(u.GetAspect(), ShouldEqual, "none")
			So(u.GetDisplay(), ShouldEqual, "Guest User")
			So(u.IsActive(), ShouldBeFalse)
			So(u.HasRole("any"), ShouldBeFalse)
			So(u.HasGroup("any"), ShouldBeFalse)
		})
	})

}

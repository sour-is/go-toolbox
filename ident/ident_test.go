package ident // import "sour.is/x/toolbox/ident"

import (
	"net/http"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNullUser(t *testing.T) {
	Convey("Given a logged in user", t, func() {
		u := NewNullUser("ident", "aspect", "name", true)
		So(u.IsActive(), ShouldBeTrue)
		So(u.HasRole("any"), ShouldBeTrue)
		So(u.HasGroup("any"), ShouldBeTrue)
		So(u.GetAspect(), ShouldEqual, "aspect")
		So(u.GetDisplay(), ShouldEqual, "name")
		So(u.GetIdentity(), ShouldEqual, "ident")
	})
	Convey("Given a logged out user", t, func() {
		u := NewNullUser("ident", "aspect", "name", false)
		So(u.IsActive(), ShouldBeFalse)
		So(u.HasRole("any"), ShouldBeFalse)
		So(u.HasGroup("any"), ShouldBeFalse)
	})
}

func TestGetIdent(t *testing.T) {

	Register("null", func(r *http.Request) Ident { return NullUser{} })
	RegisterConfig("null", map[string]string{"foo": "bar"})

	req := new(http.Request)

	a := NewNullUser("ident", "aspect", "name", true)
	Register("active", a.MakeHandlerFunc())

	i := NewNullUser("ident", "aspect", "name", false)
	Register("inactive", i.MakeHandlerFunc())

	Convey("Given a request to decode", t, func() {

		Convey("If valid return an active Ident", func() {
			u := GetIdent("active", req)
			So(u.IsActive(), ShouldBeTrue)
		})

		Convey("If invalid return an inactive Ident", func() {
			u := GetIdent("inactive", req)
			So(u.IsActive(), ShouldBeFalse)
		})

		Convey("With no registered handler", func() {
			So(func() { GetIdent("none", req) }, ShouldPanicWith, "GetIdentity Plugin does not exist!")
		})

		Convey("Test multiple handlers in list (empty)", func() {
			u := GetIdent("", req)
			So(u.IsActive(), ShouldBeFalse)
		})

		Convey("Test multiple handlers in list", func() {
			u := GetIdent("inactive active", req)
			So(u.IsActive(), ShouldBeTrue)
		})

		Convey("Test multiple handlers in list (reversed)", func() {
			u := GetIdent("active inactive", req)
			So(u.IsActive(), ShouldBeTrue)
		})
	})

}

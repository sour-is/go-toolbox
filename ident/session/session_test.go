package session

import (
	"io"
	"net/http"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"sour.is/go/ident"
)

func TestSession(t *testing.T) {
	req := http.Request{}
	req.Header = make(http.Header)

	Convey("On session should register itself.", t, func() {
		So(ident.IdentSet, ShouldContainKey, "session")
	})

	sess := NewSession("ident", "aspect", "display name")

	Convey("Given a valid session", t, func() {
		Convey("When the authorization header is not set", func() {
			u := ident.GetIdent("session", &req)

			So(u.GetIdentity(), ShouldEqual, "anon")
			So(u.GetAspect(), ShouldEqual, "none")
			So(u.GetDisplay(), ShouldEqual, "Guest User")
			So(u.IsActive(), ShouldBeFalse)
			So(u.HasRole("any"), ShouldBeFalse)
			So(u.HasGroup("any"), ShouldBeFalse)
		})

		Convey("If the authorization header is not for session", func() {
			req.Header.Set("authorization", "basic 12345")

			u := ident.GetIdent("session", &req)
			So(u.IsActive(), ShouldBeFalse)
		})

		Convey("If the authorization header is set for non-existant session", func() {
			req.Header.Set("authorization", "session 12345")

			u := ident.GetIdent("session", &req)
			So(u.IsActive(), ShouldBeFalse)
		})

		Convey("When the authorization header is set", func() {
			req.Header.Set("authorization", "session "+sess)

			u := ident.GetIdent("session", &req)

			So(u.GetIdentity(), ShouldEqual, "ident")
			So(u.GetAspect(), ShouldEqual, "aspect")
			So(u.GetDisplay(), ShouldEqual, "display name")
			So(u.IsActive(), ShouldBeTrue)
			So(u.HasRole("any"), ShouldBeFalse)
			So(u.HasGroup("any"), ShouldBeFalse)
		})

		Convey("If the session has been deleted", func() {
			DeleteSession(sess)

			u := ident.GetIdent("session", &req)

			So(u.IsActive(), ShouldBeFalse)
		})
	})

}

type nopCloser struct {
	io.Reader
}

func (nopCloser) Close() error { return nil }

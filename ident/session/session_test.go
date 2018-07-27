package session // import "sour.is/x/toolbox/ident/session"

import (
		"net/http/httptest"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"sour.is/x/toolbox/ident"
	"github.com/spf13/viper"
	"strings"
)

func TestSession(t *testing.T) {
	config := `
	[idm.session]
    [idm.session.user-groups]
      "ident" = [ "admin" ]
    [idm.session.group-roles]
      "admin" = [ "reader", "writer" ]
	`

	viper.SetConfigType("toml")
	viper.ReadConfig(strings.NewReader(config))

	t.Log(viper.AllSettings())
	Config()

	req := httptest.NewRequest("GET", "/some-url", nil)

	Convey("On session should register itself.", t, func() {
		So(ident.IdentSet, ShouldContainKey, "session")
	})

	sess := NewSession("ident", "aspect", "display name", nil, nil)

	Convey("Given a valid session", t, func() {
		Convey("When the authorization header is not set", func() {
			u := ident.GetIdent("session", req)

			So(u.GetIdentity(), ShouldEqual, "anon")
			So(u.GetAspect(), ShouldEqual, "none")
			So(u.GetDisplay(), ShouldEqual, "Guest User")
			So(u.IsActive(), ShouldBeFalse)
			So(u.HasRole("any"), ShouldBeFalse)
			So(u.HasGroup("any"), ShouldBeFalse)
		})

		Convey("If the authorization header is not for session", func() {
			req.Header.Set("authorization", "basic 12345")

			u := ident.GetIdent("session", req)
			So(u.IsActive(), ShouldBeFalse)
		})

		Convey("If the authorization header is set for non-existant session", func() {
			req.Header.Set("authorization", "session 12345")

			u := ident.GetIdent("session", req)
			So(u.IsActive(), ShouldBeFalse)
		})

		Convey("When the authorization header is set", func() {
			req.Header.Set("authorization", "session " + sess.(User).Session)

			u := ident.GetIdent("session", req)

			So(u.GetIdentity(), ShouldEqual, "ident")
			So(u.GetAspect(), ShouldEqual, "aspect")
			So(u.GetDisplay(), ShouldEqual, "display name")
			So(u.IsActive(), ShouldBeTrue)
			So(u.HasRole("any"), ShouldBeFalse)
			So(u.HasGroup("any"), ShouldBeFalse)

			t.Log(u)

			So(u.HasGroup("admin"), ShouldBeTrue)
			So(u.HasRole("writer"), ShouldBeTrue)
			So(u.HasRole("reader"), ShouldBeTrue)

		})

		Convey("If the session has been deleted", func() {
			DeleteSession(sess.(User).Session)

			u := ident.GetIdent("session", req)

			So(u.IsActive(), ShouldBeFalse)
		})
	})
}
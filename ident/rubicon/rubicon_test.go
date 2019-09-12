package rubicon // import "sour.is/x/toolbox/ident/rubicon"

import (
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"github.com/bouk/monkey"
	. "github.com/smartystreets/goconvey/convey"
	"sour.is/x/toolbox/ident"
)

// TestRubicon tests the rubicon ident methods
func TestRubicon(t *testing.T) {
	req := http.Request{}
	req.URL = new(url.URL)

	res := new(http.Response)
	var client *http.Client
	clientCalled := false
	clientURL := ""

	monkey.PatchInstanceMethod(reflect.TypeOf(client), "Get", func(cl *http.Client, url string) (*http.Response, error) {
		clientCalled = true
		clientURL = url
		return res, nil
	})
	defer monkey.UnpatchAll()

	Convey("On rubicon mock should register itself.", t, func() {
		So(ident.GetHandlers(), ShouldContainKey, "rubicon")
	})

	Convey("Given a valid request with config options", t, func() {
		ident.RegisterConfig("rubicon", map[string]string{
			"idm": "https://request",
		})

		Convey("When the access_token is not set", func() {
			u := ident.GetIdent("rubicon", &req)
			So(clientCalled, ShouldBeFalse)
			//So(cacheIdent, ShouldBeNil)

			So(u.GetIdentity(), ShouldEqual, "anon")
			So(u.GetAspect(), ShouldEqual, "none")
			So(u.GetDisplay(), ShouldEqual, "Guest User")
			So(u.IsActive(), ShouldBeFalse)
			So(u.HasRole("any"), ShouldBeFalse)
			So(u.HasGroup("any"), ShouldBeFalse)
		})

		Convey("When the remote service returns an empty response", func() {
			req.URL.RawQuery = "access_token=asdf"
			res.Body = nopCloser{strings.NewReader("{}")}
			clientCalled = false
			clientURL = ""

			u := ident.GetIdent("rubicon", &req)
			So(clientCalled, ShouldBeTrue)
			So(clientURL, ShouldContainSubstring, "https://request")
			So(clientURL, ShouldContainSubstring, "?user_token=asdf")

			So(u.GetIdentity(), ShouldEqual, "anon")
			So(u.GetAspect(), ShouldEqual, "none")
			So(u.GetDisplay(), ShouldEqual, "Guest User")
			So(u.IsActive(), ShouldBeFalse)
			So(u.HasRole("any"), ShouldBeFalse)
			So(u.HasGroup("any"), ShouldBeFalse)
		})

		Convey("When the remote service returns an valid response", func() {
			req.URL.RawQuery = "access_token=asdf"
			res.Body = nopCloser{strings.NewReader(`{"userId":1,"userName":"ident","email":"foo@bar.tld","firstName":"first","lastName":"last"}`)}
			clientCalled = false
			clientURL = ""

			u := ident.GetIdent("rubicon", &req)
			So(clientCalled, ShouldBeTrue)

			So(u.GetIdentity(), ShouldEqual, "ident")
			So(u.GetAspect(), ShouldEqual, "rubicon")
			So(u.GetDisplay(), ShouldContainSubstring, "first last")
			So(u.IsActive(), ShouldBeTrue)
			So(u.HasRole("any"), ShouldBeTrue)
			So(u.HasGroup("any"), ShouldBeTrue)

			Convey("After successful authenticaton it should be stored in cache.", func() {
				_, ok := store.Get("asdf")
				So(ok, ShouldBeTrue)
			})
		})

		Convey("When the remote service returns an valid response (cached)", func() {
			req.URL.RawQuery = "access_token=asdf"
			res.Body = nopCloser{strings.NewReader(`{"userId":1,"userName":"ident","email":"foo@bar.tld","firstName":"first","lastName":"last"}`)}
			clientCalled = false
			clientURL = ""

			u := ident.GetIdent("rubicon", &req)
			So(clientCalled, ShouldBeFalse)

			So(u.GetIdentity(), ShouldEqual, "ident")
			So(u.GetAspect(), ShouldEqual, "rubicon")
			So(u.GetDisplay(), ShouldContainSubstring, "first last")
			So(u.IsActive(), ShouldBeTrue)
			So(u.HasRole("any"), ShouldBeTrue)
			So(u.HasGroup("any"), ShouldBeTrue)
		})

	})
}

type nopCloser struct {
	io.Reader
}

func (nopCloser) Close() error { return nil }

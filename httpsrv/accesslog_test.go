package httpsrv

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"net/http/httptest"
	"sour.is/x/toolbox/ident"
	"bytes"
	"bufio"
)

func TestDoSession(t *testing.T) {
	Convey("Given a HTTP Request validate session-id sets header", t, func(){
		r := httptest.NewRequest("GET", "/", nil)
		w := wrapResponseWriter(httptest.NewRecorder())

		doSessionID("", w, r, ident.NullUser{})

		So(r.Header.Get(sessionId), ShouldNotBeEmpty)
		So(w.Header().Get(sessionId), ShouldNotBeEmpty)
	})

	Convey("Given a HTTP Request validate session-id passes received header", t, func(){
		r := httptest.NewRequest("GET", "/", nil)
		r.Header.Set(sessionId, "SOMEVALUE")

		w := wrapResponseWriter(httptest.NewRecorder())

		doSessionID("", w, r, ident.NullUser{})

		So(r.Header.Get(sessionId), ShouldNotBeEmpty)
		So(r.Header.Get(sessionId), ShouldEqual, "SOMEVALUE")

		So(w.Header().Get(sessionId), ShouldNotBeEmpty)
		So(w.Header().Get(sessionId), ShouldEqual, "SOMEVALUE")
	})
}


func TestDoAccessLog(t *testing.T) {
	var b bytes.Buffer
	var l = bufio.NewWriter(&b)

	accessLog.SetOutput(l)
	Convey("Given a request access log writes to stdout", t, func(){
		r := httptest.NewRequest("GET", "/some-url", nil)
		w := wrapResponseWriter(httptest.NewRecorder())
		w.WriteHeader(200)

		doAccessLog("NAME", w, r, ident.NullUser{"IDENT", "ASPECT", "name", true})

		l.Flush()
		str := b.String()
		b.Reset()

		So(str, ShouldNotBeBlank)
		So(str, ShouldContainSubstring, "NAME")
		So(str, ShouldContainSubstring, "ASPECT")
		So(str, ShouldContainSubstring, "IDENT")
		So(str, ShouldContainSubstring, "/some-url")
	})
}
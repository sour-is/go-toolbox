package httpsrv

import (
	"net/http"
	"testing"

	"github.com/bouk/monkey"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/spf13/viper"
	"net/http/httptest"
	"sour.is/x/toolbox/ident"
)

func TestIdentWrapper(t *testing.T) {

	Convey("Given a HTTP Request validate handler completes", t, func() {
		r := httptest.NewRequest("GET", "/", nil)
		w := WrapResponseWriter(httptest.NewRecorder())

		fn := identWrapper("TEST", stdWrapper(doNothing))
		fn.ServeHTTP(w, r)

		So(r.Header.Get(sessionId), ShouldNotBeEmpty)
	})

	monkey.Patch(viper.GetString, func(str string) string {
		t.Log("GOT " + str)
		return "something"
	})

	monkey.Patch(ident.GetIdent, func(str string, r *http.Request) ident.Ident {
		t.Log("GOT " + str)

		return ident.NullUser{"IDENT", "ASPECT", "NAME", true}
	})

	defer monkey.UnpatchAll()

	Convey("Given a HTTP Request validate ident handler completes", t, func() {
		r := httptest.NewRequest("GET", "/", nil)
		ww := httptest.NewRecorder()
		w := WrapResponseWriter(ww)

		fn := identWrapper("TEST", doHello)
		fn.ServeHTTP(w, r)

		So(r.Header.Get(sessionId), ShouldNotBeEmpty)
		So(w.Header().Get("x-user-ident"), ShouldEqual, "IDENT")
		So(w.Header().Get("x-user-aspect"), ShouldEqual, "ASPECT")
		So(w.Header().Get("x-user-name"), ShouldEqual, "NAME")

		So(ww.Body.String(), ShouldEqual, "HELLO NAME")
	})
}

func doNothing(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("NOTHING"))
}

func doHello(w ResponseWriter, r *http.Request, id ident.Ident) {
	w.Write([]byte("HELLO " + id.GetDisplay()))
}

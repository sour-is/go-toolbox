package routes // import "sour.is/x/toolbox/httpsrv/routes"

import (
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"sour.is/x/toolbox/httpsrv"
	"sour.is/x/toolbox/ident"
)

func TestDefaultRoutes(t *testing.T) {
	Convey("Given a request for index", t, func() {
		r, _ := http.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()
		i := ident.NewNullUser("user", "default", "name", true)

		W := httpsrv.WrapResponseWriter(w)
		Index(W, r, i)
		txt := w.Body.String()

		So(txt, ShouldContainSubstring, "Welcome")
		So(txt, ShouldContainSubstring, "user")
		So(txt, ShouldContainSubstring, "name")
	})
}

package routes // sour.is/go/uuid/routes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"sour.is/go/uuid"
)

func TestRoutes(t *testing.T) {
	Convey("Given a request for v4 uuid", t, func() {
		r, _ := http.NewRequest("GET", "/v1/uuid/v4", nil)
		w := httptest.NewRecorder()

		v4(w, r)
		u := w.Body.String()

		So(len(u), ShouldEqual, 36)
		So(u, ShouldContainSubstring, "-")
		b := uuid.Bytes(u)
		So(len(b), ShouldEqual, 16)
	})

	Convey("Given a request for v5 uuid", t, func() {
		r, _ := http.NewRequest("GET", "/v1/uuid/v5", nil)
		w := httptest.NewRecorder()

		v5(w, r)
		u := w.Body.String()

		So(len(u), ShouldEqual, 36)
		So(u, ShouldContainSubstring, "-")
		b := uuid.Bytes(u)
		So(len(b), ShouldEqual, 16)
	})

	Convey("Given a request for v6 uuid", t, func() {
		r, _ := http.NewRequest("GET", "/v1/uuid/v6", nil)
		w := httptest.NewRecorder()

		v6(w, r)
		u := w.Body.String()

		So(len(u), ShouldEqual, 36)
		So(u, ShouldContainSubstring, "-")
		b := uuid.Bytes(u)
		So(len(b), ShouldEqual, 16)
	})

	Convey("Given a request for v6 uuid origin", t, func() {
		r, _ := http.NewRequest("GET", "/v1/uuid/v6?origin=true", nil)
		w := httptest.NewRecorder()

		v6(w, r)
		u := w.Body.String()

		So(u, ShouldEqual, uuid.V6(uuid.NilUuid, true))
		So(len(u), ShouldEqual, 36)
		So(u, ShouldContainSubstring, "-")
		b := uuid.Bytes(u)
		So(len(b), ShouldEqual, 16)
	})

}

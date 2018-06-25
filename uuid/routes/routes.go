package routes // import "sour.is/x/toolbox/uuid/routes"

import (
	"fmt"
	"net/http"

	"sour.is/x/toolbox/httpsrv"
	"sour.is/x/toolbox/uuid"
)

func init() {
	httpsrv.HttpRegister("uuid", httpsrv.HttpRoutes{
		{
			Name:        "v4",
			Method:      "GET",
			Pattern:     "/v1/uuid/v4",
			HandlerFunc: v4,
		},
		{
			Name:        "v5",
			Method:      "GET",
			Pattern:     "/v1/uuid/v5",
			HandlerFunc: v5,
		},
		{
			Name:        "v6",
			Method:      "GET",
			Pattern:     "/v1/uuid/v6",
			HandlerFunc: v6,
		},
	})
}

func v4(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, uuid.V4())
}

func v5(w http.ResponseWriter, r *http.Request) {
	var ns string
	var name string

	if ns = r.URL.Query().Get("ns"); ns == "" {
		ns = "00000000-0000-0000-0000-000000000000"
	}

	if name = r.URL.Query().Get("name"); name == "" {
		name = "text"
	}

	fmt.Fprint(w, uuid.V5(name, ns))
}

func v6(w http.ResponseWriter, r *http.Request) {
	var ns string
	var origin bool

	if ns = r.URL.Query().Get("ns"); ns == "" {
		ns = "00000000-0000-0000-0000-000000000000"
	}

	origin = false
	if ok := r.URL.Query().Get("origin"); ok == "true" {
		origin = true
	}

	fmt.Fprint(w, uuid.V6(ns, origin))
}
